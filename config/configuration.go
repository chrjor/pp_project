// configuration.go
// Christian Jordan
// Configuration space data structure

package config

import (
	"math"
	"pp_project/utils"
	"strconv"
	"strings"
)

// ConfigSpace is a struct used for path planning
type ConfigSpace struct {
	Path       *PathPlan  // Root of the tree
	Obstacles  []Obstacle // Obstacles in the configuration space
	WinHeight  float32    // Window height
	WinWidth   float32    // Window width
	ConfigPath string     // Path to config file
}

// Obstacle is an interface for objects in the configuration space
type Obstacle interface {
	Draw()
	Collision(*Point) bool
}

// Create a new configuration space from a config file
func NewConfigSpace(configPath string) *ConfigSpace {
	// Initialize variables
	var winWidth, winHeight float64
	var radius, delta float64
	var start *Point
	var goal *Point
	var obstacles []Obstacle

	// Parse config file
	config := utils.ReadFile(configPath)

	for _, line := range config {
		line := strings.Split(line, ",")
		if line[0] == "window" {
			winHeight, _ = strconv.ParseFloat(line[1], 32)
			winWidth, _ = strconv.ParseFloat(line[2], 32)
		} else if line[0] == "radius" {
			radius, _ = strconv.ParseFloat(line[1], 32)
		} else if line[0] == "delta" {
			delta, _ = strconv.ParseFloat(line[1], 32)
		} else if line[0] == "start" {
			x, _ := strconv.ParseFloat(line[1], 32)
			y, _ := strconv.ParseFloat(line[2], 32)
			start = NewPoint(float32(x), float32(y))
		} else if line[0] == "goal" {
			x, _ := strconv.ParseFloat(line[1], 32)
			y, _ := strconv.ParseFloat(line[2], 32)
			goal = NewPoint(float32(x), float32(y))
		} else if line[0] == "rectangle" {
			x, _ := strconv.ParseFloat(line[1], 32)
			y, _ := strconv.ParseFloat(line[2], 32)
			w, _ := strconv.ParseFloat(line[3], 32)
			h, _ := strconv.ParseFloat(line[4], 32)
			obstacles = append(obstacles,
				NewRectangle(float32(x), float32(y), float32(w), float32(h)))
		} else if line[0] == "circle" {
			x, _ := strconv.ParseFloat(line[1], 32)
			y, _ := strconv.ParseFloat(line[2], 32)
			r, _ := strconv.ParseFloat(line[3], 32)
			obstacles = append(obstacles,
				NewCircle(float32(x), float32(y), float32(r)))
		}
	}

	return &ConfigSpace{
		Path: NewPathPlan(float32(delta),
			float32(radius),
			goal,
			start,
		),
		Obstacles:  obstacles,
		WinHeight:  float32(winHeight),
		WinWidth:   float32(winWidth),
		ConfigPath: "/data/",
	}
}

// Add an obstacle to the configuration space
func (c *ConfigSpace) AddObstacle(o Obstacle) {
	c.Obstacles = append(c.Obstacles, o)
}

// Check if a point is feasible in the configuration space
func (c *ConfigSpace) Feasible(pt *Point) bool {
	for _, o := range c.Obstacles {
		if o.Collision(pt) {
			return false
		}
	}
	return true
}

// Draw the configuration space
func (c *ConfigSpace) Draw() {
	c.Path.Draw()
	for _, o := range c.Obstacles {
		o.Draw()
	}
}

// Calculate the distance between two points in the configuration space
func CalcDistance(pt1 *Point, pt2 *Point) float32 {
	return float32(math.Sqrt(math.Pow(float64(pt1.X-pt2.X), 2) +
		math.Pow(float64(pt1.Y-pt2.Y), 2)))
}
