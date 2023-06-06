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
	point    *Point            // Point in the path plan
	parent   *MileStone        // Parent of the milestone
	children MileStoneChildren // Children of the milestone
	lock     sync.Mutex        // Lock for milestone
	Cost     float32           // Cost of the milestone
	ParDist  float32           // Distance from parent
}

// Create a new MileStone, assume point is feasible
func NewMileStone(pt *Point) *MileStone {
	return &MileStone{
		point:    pt,
		parent:   nil,
		children: NewChildrenList(),
		lock:     sync.Mutex{},
	}
}

// Set the point of a milestone
func (ms *MileStone) GetPoint() *Point { return ms.point }

// Add a child to a milestone
func (ms *MileStone) SetChild(c *MileStone) {
	ms.children.Add(c)
}

// Remove a child from a milestone
func (ms *MileStone) RemoveChild(c *MileStone) {
	ms.children.Remove(c)
}

// Set the parent of a milestone
func (ms *MileStone) SetParent(p *MileStone, dist float32) {
	if ms.parent != nil {
		ms.parent.RemoveChild(ms)
	}
	ms.ParDist = dist
	ms.parent = p
}

func (ms *MileStone) SetCost(cost float32) {
	ms.lock.Lock()
	defer ms.lock.Unlock()
	ms.Cost = cost
}

// Update the cost of a milestone and all of its children
func (ms *MileStone) UpdateCost(diff float32) {
	ms.SetCost(ms.Cost + diff)
	children := ms.children.GetChildren()
	for _, child := range children {
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
