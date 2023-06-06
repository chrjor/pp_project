// milestones.go
// Christian Jordan
// MileStone object for path planning

package config

import (
	"math"
	"sync"
)

// MileStone is a node-like struct that represents a point in the path plan
type MileStone struct {
	point        *Point       // Point in the path plan
	parent       *MileStone   // Parent of the milestone
	children     []*MileStone // Children of the milestone
	lock         sync.Mutex   // Lock for updating cost value (float32)
	Cost         float32      // Cost of the milestone (lengthance from start)
	OccupiedFlag int32        // Milestone is occupied by a thread
}

// Create a new MileStone, assume point is feasible
func NewMileStone(pt *Point) *MileStone {
	return &MileStone{
		point:    pt,
		parent:   nil,
		children: make([]*MileStone, 0),
		lock:     sync.Mutex{},
	}
}

// Set the point of a milestone
func (ms *MileStone) GetPoint() *Point { return ms.point }

// Add a child to a milestone
func (ms *MileStone) SetChild(c *MileStone) {
	for idx, child := range ms.children {
		if child == nil {
			ms.children[idx] = c
			return
		}
	}
	ms.children = append(ms.children, c)
}

// Remove a child from a milestone
func (ms *MileStone) RemoveChild(c *MileStone) {
	for idx, child := range ms.children {
		if child == c {
			ms.children[idx] = nil
			return
		}
	}
}

// Set the parent of a milestone
func (ms *MileStone) SetParent(p *MileStone) {
	if ms.parent != nil {
		ms.parent.RemoveChild(ms)
	}
	ms.parent = p
}

func (ms *MileStone) SetCost(cost float32) { ms.Cost = cost }

// Update the cost of a milestone and all of its children
func (ms *MileStone) UpdateCost(diff float32) {
	ms.lock.Lock()
	ms.Cost += diff
	ms.lock.Unlock()
	for _, child := range ms.children {
		if child == nil {
			continue
		}
		child.UpdateCost(diff)
	}
}

// Set the milestone's new point location as min(delta, length) lengthance from
// its nearest neighbor towards the direction of its current (sampled) point
func (ms *MileStone) ShortenPathToNearest(nearest *MileStone, delta float32) {
	length := CalcDistance(ms.point, nearest.point)
	delta = float32(math.Min(float64(delta), float64(length)))
	ms.point.X = nearest.point.X + (ms.point.X-nearest.point.X)*delta/length
	ms.point.Y = nearest.point.Y + (ms.point.Y-nearest.point.Y)*delta/length
}
