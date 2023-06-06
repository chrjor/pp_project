//
// stealing.go
// Christian Jordan
// Work stealing algorithm
//

package concurrent

import (
	"math/rand"
	"sync"
)

// NewWorkStealingExecutor returns an Executor that is implemented using the
// work-stealing algorithm.
// @param capacity - The number of goroutines in the pool
// @param threshold - The number of items that a goroutine in the pool can
// grab from the executor in one time period.
func NewWorkStealingExecutor(capacity, threshold int) ExecutorService {
	var workers []Worker
	for i := 0; i < capacity; i++ {
		workers = append(workers, Worker{workerQueue: NewUnBoundedDEQueue()})
	}
	executor := &Executor{
		workers:        workers,
		globalQueue:    NewUnBoundedDEQueue(),
		shutdown:       make(chan interface{}),
		thresholdQueue: threshold,
		wg:             sync.WaitGroup{},
	}
	for worker := 0; worker < capacity; worker++ {
		go executor.runWorkStealer(worker)
	}
	return executor
}

// runWorkStealer is the main worker instructions for the worker stealing routine
func (e *Executor) runWorkStealer(me int) {
	for {
		select {
		case <-e.shutdown:
			return
		default:
			// Attempt to pop and call a task from local queue
			successfulCall := e.popCallTask(me)

			// If local queue is empty, steal from another worker
			if !successfulCall {
				var randSteal int
				for {
					randSteal = rand.Intn(len(e.workers))
					if randSteal == me {
						continue
					}
					break
				}
				// Steal from random worker, if they have greater than 500 tasks
				steal := 1000
				if e.workers[randSteal].workerQueue.Size() > int32(steal) {
					for i := 0; i < steal; i++ {
						task := e.workers[randSteal].workerQueue.PopTop()
						if task != nil {
							e.workers[me].workerQueue.PushBottom(task)
						}
					}
				}
			}
		}
	}
}
