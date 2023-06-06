// sampling.go
// Christian Jordan
// Sampling helper functions and objects

package pathfind

import (
	"math/rand"
	"pp_project/config"
)

// SamplePoint samples a random point in the configuration space
func SamplePoint(space *config.ConfigSpace) *config.Point {
	randX := rand.Float32() * float32(space.WinWidth)
	randY := rand.Float32() * float32(space.WinHeight)
	return config.NewPoint(randX, randY)
}
