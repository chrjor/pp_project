//
// task.go
// Christian Jordan
// Task interface for concurrent package
//

package concurrent

import (
	"pp_project/config"
	"pp_project/pathfind"
)

// UpdateTask updates the config space with a new sample point. It implemnents
// the Callable interface
type UpdateTask struct {
	ctx *config.ConfigSpace // Config space to update
}

// NewUpdateTask creates a new UpdateTask
func NewUpdateTask(ctx *config.ConfigSpace) *UpdateTask {
	return &UpdateTask{ctx: ctx}
}

func (t *UpdateTask) GetDistToGoal() float32 {
	return t.ctx.Path.GetDistToGoal()
}

// Run the task
func (t *UpdateTask) Run() {
	var sample *config.MileStone
	point := pathfind.SamplePoint(t.ctx)
	if t.ctx.Feasible(point) {
		// Create new MileStone
		sample = config.NewMileStone(point)
		// Run RRT* algorithm
		pathfind.RRTstar(sample, t.ctx)
	}
}

// TaskFuture implements the Future interface
type TaskFuture struct {
	task   *UpdateTask
	result chan interface{}
}

func (f *TaskFuture) Get() interface{} {
	return f.task.GetDistToGoal()
}

func NewTaskFuture(task interface{}) Future {
	return &TaskFuture{
		task:   task.(*UpdateTask),
		result: make(chan interface{}, 1),
	}
}
