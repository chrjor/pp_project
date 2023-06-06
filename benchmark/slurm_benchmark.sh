#!/bin/bash
#
#SBATCH --job-name=proj2_benchmark 
#SBATCH --output=./benchmark/%j.%N.stdout
#SBATCH --error=./benchmark/%j.%N.stderr
#SBATCH --chdir=~/mpcs52060/project-3-chrjor/proj3/
#SBATCH --partition=general
#SBATCH --nodes=1
#SBATCH --ntasks=1
#SBATCH --cpus-per-task=16
#SBATCH --mem=1G
#SBATCH --exclusive
#SBATCH --time=4:00:00

module load golang/1.16.2 
./benchmark/process_speedup.sh
