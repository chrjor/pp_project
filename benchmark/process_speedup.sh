#!/bin/bash

# Bash script that runs benchmark sample_sizes and ccreate speed-up graphs
# 
# Example usage:
#     ./process_speedup.sh input.txt
#     ./process_speedup.sh --help


# Display help information
Usage()
{
    echo "Usage: process_speedup.sh [input_file]|[options]" 
    echo
    echo "options:"
    echo "     --help       display this help"
    echo
}

# Check command line args passed to script
if [[ "$1" == "--help" ]]; then
    Usage
    exit 0
elif [ $# -ne 1 ]; then
    Usage
    exit 1
fi

# Set benchmark parameters
input="$1"
sample_sizes=( "40000" "80000" "160000" "640000" )

# Create/prep output the file
out=$(mktemp)
header=( "model" "sample_size" "threads" ) 
header+=( "test1" "test2"  "test3" "test4" "test5" )
echo "${header[@]}" >> "$out"

# Run sequential operations
for size in "${sample_sizes[@]}"; do
    sequential=( "s" "$size" "1" )
    for test in {1..5}; do
        sequential+=($(go run pp_project/simulator b "$size" "$input"))
    done
    echo "${sequential[@]}" >> "$out"
done

# Average threaded speedup
for strategy in "wb" "ws"; do
    for size in "${sample_sizes[@]}"; do
        for threads in 2 4 6 8 12; do
            threaded=( "$strategy" "$size" "$threads" )
            for test in {1..5}; do 
                threaded+=($(go run pp_project/simulator b "$size" "$input" "$strategy" "$threads"))
            done
            echo "${threaded[@]}" >> "$out"
        done
    done
done

# Export to python
python3 benchmark/process_speedup.py

# Clean up
rm "$out"
