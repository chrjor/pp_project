# Python script that processes benchmark data and produce speedup graph
# This script is called by process_speedup.sh, which runs the actual
# tests.


import matplotlib.pyplot as plt
import pandas as pd
import numpy as np
from pandas.api.types import CategoricalDtype
import seaborn as sns

np.float = float

# Create dataframe
data = pd.read_table("benchmark/output.txt", sep='\s+')
data["avg"] = data.iloc[:,3:8].mean(axis=1)

# Calculate speedup
sizes = data["sample_size"].unique()
num_sizes = sizes.size
data_seq = data.iloc[:num_sizes,[0,1,2,8]]
data_par = data.iloc[num_sizes:,[0,1,2,8]]
threads = data_par["threads"].unique()

# Re-index data_seq
data_seq.set_index("sample_size", inplace=True)

# Calculate speedup for each row
def speedup(row):
    return data_seq.loc[row["sample_size"]]["avg"] / row["avg"]

data_par["speedup"] = data_par.apply(speedup, axis=1)

# Plot speedup graphs
for model in data_par["model"].unique():
    if model != "s":
        data_plot = data_par[data_par["model"] == model]
        sns.set_style("whitegrid")
        sns.lineplot(
            x="threads", 
            y="speedup", 
            hue="sample_size", 
            data=data_plot
            ).set(
                title=f"{model} Speedup Graph", 
                xlabel="Threads",
                ylabel="Speedup")
        sns.despine()

    # Create PNG
    plt.savefig(f"benchmark/{model}_graph.png")
    plt.clf()