// pathplan.go
// Christian Jordan
// Path plan data tree object

package config

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
func (path *PathPlan) GetNN(newMS *MileStone) NeighborHeap {

	// Find all neighbors of new MileStone within visibility radius
	var neighborhood NeighborHeap
	RecurseNN := func(ms *MileStone) {
		dist := CalcDistance(ms.point, newMS.point)
		if dist <= path.Radius {
			neighbor := NewNeighborItem(ms, dist)
			neighborhood.Push(neighbor)
		}
	}
	BranchApply(path.pathHead.children, RecurseNN)

	return neighborhood
}

// Draw the path plan
func (path *PathPlan) Draw() {}
