// pathplan.go
// Christian Jordan
// Path plan data tree object

package config

import (
	"container/heap"
	"sync/atomic"
)

// PathPlan is a struct used for path planning
type PathPlan struct {
	pathHead   *MileStone // Root of the path tree
	Goal       *MileStone // Goal point
	distToGoal float32    // Distance to goal
	Radius     float32    // Visibility radius
	DeltaDist  float32    // Max distance from branch to new milestone
}

// Create a new PathPlan
func NewPathPlan(delta float32,
	radius float32,
	goal *Point,
	start *Point,
) *PathPlan {

	return &PathPlan{
		pathHead:   NewMileStone(start),
		Goal:       NewMileStone(goal),
		distToGoal: 0,
		Radius:     radius,
		DeltaDist:  delta,
	}
}

// Get the path head
func (path *PathPlan) GetDistToGoal() float32 {
	return path.distToGoal
}

// Get all neighbors in visibilty of a MileStone, assuming MileStone is feasible
func (path *PathPlan) GetNN(new_ms *MileStone) NeighborHeap {

	// Use DFS to find all neighbors within radius and return as heap
	var neighborhood NeighborHeap
	var DFS func(*MileStone)

	DFS = func(ms *MileStone) {
		dist := CalcDistance(ms.point, new_ms.point)
		if dist <= path.Radius {
			neighbor := NewNeighborItem(ms, dist)
			heap.Push(&neighborhood, neighbor)
		}
		for _, child := range ms.children {
			if child != nil && child.OccupiedFlag == 1 {
				DFS(child)
			}
		}
	}
	DFS(path.pathHead)
	return neighborhood
}

// Set neighbors as occupied, return nil and unset all if any are occupied
func SetNNOccupied(neighbors NeighborHeap) NeighborHeap {
	for idx, nItem := range neighbors {
		if !atomic.CompareAndSwapInt32(&nItem.Neighbor.OccupiedFlag, 0, 1) {
			SetNNUnoccupied(neighbors[:idx+1])
			return nil
		}
	}
	return neighbors
}

// Set neighbors as unoccupied
func SetNNUnoccupied(neighbors NeighborHeap) {
	for _, nItem := range neighbors {
		atomic.CompareAndSwapInt32(&nItem.Neighbor.OccupiedFlag, 1, 0)
	}
}

// Draw the path plan
func (path *PathPlan) Draw() {}
