// RRTstar.go
// Christian Jordan
// RRT* implementation
// Algorithm used:
// http://motion.cs.illinois.edu/RoboticSystems/MotionPlanningHigherDimensions.html

package pathfind

import (
	"pp_project/config"
	"sync/atomic"
)

// RRT* algorithm. Assumes samplePt is feasible
func RRTstar(ms *config.MileStone, space *config.ConfigSpace) float32 {
	// Find all neighbors of new MileStone within visibility radius
	neighborhood := space.Path.GetNN(ms)
	neighborhood = config.SetNNOccupied(neighborhood)
	atomic.SwapInt32(&ms.OccupiedFlag, 1)

	// Extend the path from the given point to the nearest point in the tree
	neighborhood = ExtendPath(neighborhood, space, ms)

	if neighborhood != nil && space.Path.Goal.Cost == 0 && IsGoalVisible(ms, space) {
		// Prioritize first connection to goal point when visible
		goalDist := config.CalcDistance(ms.GetPoint(), space.Path.Goal.GetPoint())
		space.Path.Goal.SetParent(ms)
		ms.SetChild(space.Path.Goal)
		space.Path.Goal.SetCost(ms.Cost + goalDist)
	} else {
		// Rewire the tree to account for the new MileStone
		Rewire(ms, neighborhood, space)
		config.SetNNUnoccupied(neighborhood)
		atomic.SwapInt32(&ms.OccupiedFlag, 0)
	}
	return space.Path.Goal.Cost
}

// Extend the path from the given point to the nearest point in the tree
func ExtendPath(neighborhood config.NeighborHeap, space *config.ConfigSpace,
	mileStone *config.MileStone,
) config.NeighborHeap {
	// Update new MileStone to be a child of nearest neighbor
	if neighborhood != nil {
		nearest := neighborhood[0].Neighbor

		// Shorten path to nearest neighbor according to delta
		mileStone.ShortenPathToNearest(nearest, space.Path.DeltaDist)
		if !space.Feasible(mileStone.GetPoint()) {
			return nil
		}

		// Set parent and child
		newDist := config.CalcDistance(nearest.GetPoint(), mileStone.GetPoint())
		mileStone.SetParent(nearest)
		nearest.SetChild(mileStone)
		mileStone.SetCost(nearest.Cost + newDist)
	}
	return neighborhood
}

// Rewire the tree to account for the new MileStone
func Rewire(newMileStone *config.MileStone,
	nHeap config.NeighborHeap,
	space *config.ConfigSpace,
) {
	if nHeap == nil {
		return
	}
	for _, nItem := range nHeap {
		// Calc distance between newMileStone and neighbor
		distBetween := config.CalcDistance(newMileStone.GetPoint(),
			nItem.Neighbor.GetPoint())

		// Calc distances passing through the newMileStone
		newDistThrough := newMileStone.Cost + distBetween
		if newDistThrough < nItem.Neighbor.Cost {
			nItem.Neighbor.SetParent(newMileStone)
			nItem.Neighbor.RemoveChild(newMileStone)
			newMileStone.SetChild(nItem.Neighbor)
			nItem.Neighbor.UpdateCost(newDistThrough - nItem.Neighbor.Cost)
			continue
		}

		// Calc distances passing to the newMileStone
		newDistTo := nItem.Neighbor.Cost + distBetween
		if newDistTo < newMileStone.Cost {
			newMileStone.SetParent(nItem.Neighbor)
			newMileStone.RemoveChild(nItem.Neighbor)
			nItem.Neighbor.SetChild(newMileStone)
			newMileStone.UpdateCost(newDistTo - newMileStone.Cost)
			continue
		}
	}
}

func IsGoalVisible(ms *config.MileStone, space *config.ConfigSpace) bool {
	return config.CalcDistance(ms.GetPoint(),
		space.Path.Goal.GetPoint()) <= space.Path.Radius
}
