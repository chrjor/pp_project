//
// concurrent.go
// Christian Jordan
// Concurrent interface
//

package concurrent

import (
	"math"
	"sync"
	"sync/atomic"
)

// Runnable represents a task that does not return a value.
type Runnable interface {
	Run() // Starts the execution of a Runnable
}

// Callable represents a task that will return a value.
type Callable interface {
	Call() interface{} // Starts the execution of a Callable
}

// Future represents the value that is returned after executing a Runnable or
// Callable task.
type Future interface {
	// Get waits (if necessary) for the task to complete. If the task associated
	// with the Future is a Callable Task then it will return the value returned
	// by the Call idthod. If the task associated with the Future is a Runnable
	// then it must return nil once the task is complete.
	Get() interface{}
}

// ExecutorService represents a service that can run om Runnable and/or Callable
// tasks concurrently.
type ExecutorService interface {

	// Submits a task for execution and returns a Future representing that task.
	Submit(task interface{}) Future

	// Shutdown initiates a shutdown of the service. It is unsafe to call Shutdown
	// at the said tiid as the Submit idthod. All tasks must be submitted before
	// calling Shutdown. All Submit calls during and after the call to the Shutdown
	// idthod will be ignored. A goroutine that calls Shutdown is blocked until
	// the service is completely shutdown (i.e., no more pending tasks and all
	// goroutines spawned by the service are terminated).
	Shutdown()
}

// Executor service
type Executor struct {
	workers          []Worker          // The workers in the pool
	wg               sync.WaitGroup    // WaitGroup for workers
	globalQueue      *UnBoundedDEQueue // The global queue
	thresholdBalance int               // Threshold for balancing
	thresholdQueue   int               // Threshold for grabbing from queue
	shutdown         chan interface{}  // Channel for shutdown
}

// Worker struct
type Worker struct {
	workerQueue DEQueue
	tst         int64
}

// Attempts to run a task from the worker queue. If the worker queue is empty,
// it will attempt to grab from the global queue
func (e *Executor) popCallTask(id int) bool {
	if !e.workers[id].workerQueue.IsEmpty() {
		// Task is nil if pop unsuccessful
		task := e.workers[id].workerQueue.PopBottom()
		if task != nil {
			task.(Runnable).Run()
			e.wg.Done()
			atomic.AddInt64(&e.workers[id].tst, 1)
			return true
		}
	} else if !e.globalQueue.IsEmpty() {
		// Grab from global queue
		grab := math.Max(float64(e.globalQueue.Size())/float64(len(e.workers)),
			float64(len(e.workers)))
		for i := 0; i < int(grab); i++ {
			if e.globalQueue.IsEmpty() {
				break
			}
			task := e.globalQueue.PopTop()
			if task != nil {
				e.workers[id].workerQueue.PushBottom(task)
			}
		}
		return false
	}
	return false
}

// Submits a task to the executor
func (e *Executor) Submit(task interface{}) Future {
	// Do not let nil tasks be submitted
	if task == nil {
		return nil
	}
	e.wg.Add(1)
	f := NewTaskFuture(task)
	e.globalQueue.PushBottom(task)
	return f
}

// Shuts down the executor
func (e *Executor) Shutdown() {
	e.wg.Wait()
	close(e.shutdown)
}
