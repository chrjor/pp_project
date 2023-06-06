package main

import (
	"fmt"
	"os"
	"pp_project/concurrent"
	"strconv"
	"time"
)

const usage = "Usage: simulator [mode] [sample_size] input_file [parallelization] [number of threads] \n" +
	"[mode] = (b) run benchmark mode, (d) run image draw mode\n" +
	"[sample_size] = The number of samples to be generated\n" +
	"input_file = The file used to set up the configuration space\n" +
	"[parallelization] = (wb) work balancing, (ws) work stealing\n" +
	"[number of threads] = Runs parallel version of the program with the specified number of threads,\n" +
	"                      if not specified, runs the sequential version of the program.\n"

func main() {

	if len(os.Args) < 4 || len(os.Args) > 6 {
		fmt.Println(usage)
		return
	}

	// Parse command line arguments
	var strategy string
	threads := 1
	mode := os.Args[1]
	sample_size, _ := strconv.Atoi(os.Args[2])
	input := os.Args[3]
	if len(os.Args) == 6 {
		strategy = os.Args[4]
		threads, _ = strconv.Atoi(os.Args[5])
	}

	// Run benchmark mode
	var start time.Time
	var end float64
	if mode == "b" {
		start = time.Now()
	}

	// Run the simulation
	var pathOutput interface{}
	if threads == 1 {
		pathOutput = RunSequential(input, sample_size)
	} else {
		pathOutput = RunParallel(input, sample_size, threads, strategy)
	}

	// Print run-time or draw the configuration space
	if mode == "b" {
		end = time.Since(start).Seconds()
		fmt.Printf("%.2f\n", end)
	} else if mode == "d" {
		var dist float32
		if threads == 1 {
			out := pathOutput.([]float32)
			dist = out[len(out)-1]
			fmt.Println("Distance after", sample_size, "iterations: ", dist)
		} else {
			out := pathOutput.([]concurrent.Future)
			dist = out[len(out)-1].Get().(float32)
			fmt.Println("Distance after", sample_size, "iterations: ", dist)
		}
		if dist == 0 {
			fmt.Println("No Goal!")
		} else {
			fmt.Println("Goal!")
		}
	}
}
