package main

import (
	"pp_project/concurrent"
	"pp_project/config"
)

// RunParallel runs the pathfinding algorithm in parallel
func RunParallel(input string,
	sample_size int,
	threads int,
	strategy string,
) []concurrent.Future {
	var executor concurrent.ExecutorService
	var progress []concurrent.Future

	// Read the configuration space from the input file
	configSpace := config.NewConfigSpace(input)

	if strategy == "wb" {
		// Run the work balancing executor
		executor = concurrent.NewWorkBalancingExecutor(
			threads,
			sample_size/(threads),
			sample_size/(threads*threads),
		)
	} else if strategy == "ws" {
		// Run the work stealing executor
		executor = concurrent.NewWorkStealingExecutor(
			threads,
			sample_size/(threads),
		)
	}

	for i := 0; i < sample_size; i++ {
		task := concurrent.NewUpdateTask(configSpace)
		f := executor.Submit(task)
		progress = append(progress, f)
	}
	executor.Shutdown() // Shutdown the executor

	return progress
}
