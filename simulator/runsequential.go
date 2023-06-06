package main

import (
	"pp_project/concurrent"
	"pp_project/config"
)

// RunSequential runs the pathfinding algorithm sequentially
func RunSequential(input string,
	sample_size int,
) []float32 {
	var progress []float32

	// Read the configuration space from the input file
	configSpace := config.NewConfigSpace(input)

	for i := 0; i < sample_size; i++ {
		task := concurrent.NewUpdateTask(configSpace)
		task.Run()
		progress = append(progress, task.GetDistToGoal())
	}
	return progress
}
