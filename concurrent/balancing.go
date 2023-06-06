//
// balancing.go
// Christian Jordan
// Work balancing algorithm
//

package concurrent

import (
	"math"
	"math/rand"
	"sync"
)

// NewWorkBalancingExecutor returns an ExecutorService that is implemented using
// the work-balancing algorithm.
// @param capacity - The number of goroutines in the pool
// @param threshold - The number of items that a goroutine in the pool can
// grab from the executor in one time period.
// @param thresholdBalance - The threshold used to know when to perform
// balancing.
func NewWorkBalancingExecutor(capacity, thresholdQueue, thresholdBalance int,
) ExecutorService {
	var workers []Worker
	for i := 0; i < capacity; i++ {
		workers = append(workers, Worker{workerQueue: NewUnBoundedDEQueue()})
	}
	executor := &Executor{
		workers:          workers,
		globalQueue:      NewUnBoundedDEQueue(),
		shutdown:         make(chan interface{}),
		thresholdBalance: thresholdBalance,
		thresholdQueue:   thresholdQueue,
		wg:               sync.WaitGroup{},
	}
	for worker := 0; worker < capacity; worker++ {
		go executor.runWorkBalancer(worker)
	}
	return executor
}

// runWorkBalancer is the main worker instructions for the worker balancing routine
func (e *Executor) runWorkBalancer(me int) {
	for {
		select {
		case <-e.shutdown:
			return
		default:
			// Attempt to pop and call a task from local queue
			e.popCallTask(me)

			// Balance queues if needed
			if e.globalQueue.IsEmpty() {
				size := e.workers[me].workerQueue.Size()
				// Randomly select victim
				if rand.Intn(int(size+1)) == int(size) {
					victim := rand.Intn(len(e.workers))
					diff := math.Abs(float64(e.workers[victim].workerQueue.Size() -
						e.workers[me].workerQueue.Size()))
					// If difference is greater than threshold, balance
					if diff >= float64(e.thresholdBalance) {
						e.balance(victim, me)
					}
				}
			}
		}
	}
}

// Balances the local queues of two workers
func (e *Executor) balance(victim, me int) {
	if e.workers[victim].workerQueue.Size() > e.workers[me].workerQueue.Size() {
		// If victim queue is larger, grab from victim and place into me
		for i := 0; i < e.thresholdBalance/2; i++ {
			task := e.workers[victim].workerQueue.PopTop()
			if task != nil {
				e.workers[me].workerQueue.PushBottom(task)
			}
		}
	} else {
		// If victim queue is smaller, grab from me and place into victim
		for i := 0; i < e.thresholdBalance/2; i++ {
			task := e.workers[me].workerQueue.PopTop()
			if task != nil {
				e.workers[victim].workerQueue.PushBottom(task)
			}
		}
	}
}
