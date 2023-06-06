//
// unbounded.go
// Christian Jordan
// Unbounded double ended lock free queue
//

package concurrent

import (
	"sync/atomic"
	"unsafe"
)

type Task interface{}

type DEQueue interface {
	PushBottom(task Task)
	IsEmpty() bool //returns whether the queue is empty
	PopTop() Task
	PopBottom() Task
	Size() int32
}

type TaskNode struct {
	task      Task
	next      *TaskNode
	prev      *TaskNode
	takenFlag int32
}

type UnBoundedDEQueue struct {
	head  *TaskNode
	tail  *TaskNode
	count int32
}

// NewUnBoundedDEQueue returns an empty DEQueue
func NewUnBoundedDEQueue() *UnBoundedDEQueue {
	dummyHead := TaskNode{task: nil, next: nil, prev: nil, takenFlag: 0}
	dummyTail := TaskNode{task: nil, next: nil, prev: &dummyHead, takenFlag: 0}
	dummyHead.next = &dummyTail

	return &UnBoundedDEQueue{
		head:  &dummyHead,
		tail:  &dummyTail,
		count: 0,
	}
}

// PushBottom adds a task to the bottom of the deque
func (q *UnBoundedDEQueue) PushBottom(task Task) {
	for {
		// Get dummy tail node pointer
		tail := (*TaskNode)(atomic.LoadPointer(
			(*unsafe.Pointer)(unsafe.Pointer(&q.tail))))
		if !atomic.CompareAndSwapInt32(&tail.takenFlag, 0, 1) {
			continue
		}

		// Get prev node pointer
		prev := (*TaskNode)(atomic.LoadPointer(
			(*unsafe.Pointer)(unsafe.Pointer(&tail.prev))))
		if !atomic.CompareAndSwapInt32(&prev.takenFlag, 0, 1) {
			atomic.StoreInt32(&tail.takenFlag, 0)
			continue
		}

		// Check if last node is valid and push new task
		if prev == (*TaskNode)(atomic.LoadPointer(
			(*unsafe.Pointer)(unsafe.Pointer(&q.tail.prev)))) {
			newTask := &TaskNode{
				task:      task,
				next:      tail,
				prev:      prev,
				takenFlag: 1}

			if atomic.CompareAndSwapPointer(
				(*unsafe.Pointer)(unsafe.Pointer(&prev.next)),
				unsafe.Pointer(tail),
				unsafe.Pointer(newTask)) {

				atomic.StorePointer(
					(*unsafe.Pointer)(unsafe.Pointer(&q.tail.prev)),
					unsafe.Pointer(newTask))
				atomic.StoreInt32(&newTask.takenFlag, 0)
				atomic.StoreInt32(&prev.takenFlag, 0)
				atomic.StoreInt32(&tail.takenFlag, 0)
				atomic.AddInt32(&q.count, 1)
				return
			}
		}
		atomic.StoreInt32(&tail.takenFlag, 0)
		atomic.StoreInt32(&prev.takenFlag, 0)
	}
}

func (q *UnBoundedDEQueue) PopTop() Task {
	// Get first node pointer (if deque is not empty)
	firstNode := (*TaskNode)(atomic.LoadPointer(
		(*unsafe.Pointer)(unsafe.Pointer(&q.head.next))))
	if firstNode == (*TaskNode)(atomic.LoadPointer(
		(*unsafe.Pointer)(unsafe.Pointer(&q.tail)))) ||
		!atomic.CompareAndSwapInt32(&firstNode.takenFlag, 0, 1) {
		return nil
	}

	// Get pointer to next node
	next := (*TaskNode)(atomic.LoadPointer(
		(*unsafe.Pointer)(unsafe.Pointer(&firstNode.next))))
	if !atomic.CompareAndSwapInt32(&next.takenFlag, 0, 1) {
		atomic.StoreInt32(&firstNode.takenFlag, 0)
		return nil
	}

	// Pop first node pointer
	if firstNode == (*TaskNode)(atomic.LoadPointer(
		(*unsafe.Pointer)(unsafe.Pointer(&q.head.next)))) {

		atomic.AddInt32(&q.count, -1)
		atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&next.prev)),
			unsafe.Pointer(q.head))
		atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&q.head.next)),
			unsafe.Pointer(next))
		atomic.StoreInt32(&next.takenFlag, 0)
		return firstNode.task
	}
	atomic.StoreInt32(&next.takenFlag, 0)
	atomic.StoreInt32(&firstNode.takenFlag, 0)
	return nil
}

func (q *UnBoundedDEQueue) PopBottom() Task {
	// Get lastNode pointer if deque is not empty
	lastNode := (*TaskNode)(atomic.LoadPointer(
		(*unsafe.Pointer)(unsafe.Pointer(&q.tail.prev))))
	if atomic.LoadInt32(&q.tail.takenFlag) == 1 ||
		lastNode == (*TaskNode)(atomic.LoadPointer(
			(*unsafe.Pointer)(unsafe.Pointer(&q.head)))) ||
		!atomic.CompareAndSwapInt32(&lastNode.takenFlag, 0, 1) {
		return nil
	}

	// Get prev node pointer
	prev := (*TaskNode)(atomic.LoadPointer(
		(*unsafe.Pointer)(unsafe.Pointer(&lastNode.prev))))
	if !atomic.CompareAndSwapInt32(&prev.takenFlag, 0, 1) {
		atomic.StoreInt32(&lastNode.takenFlag, 0)
		return nil
	}

	// Pop tail node pointer if lastNode is valid
	if lastNode != (*TaskNode)(atomic.LoadPointer(
		(*unsafe.Pointer)(unsafe.Pointer(&q.head)))) &&
		lastNode == (*TaskNode)(atomic.LoadPointer(
			(*unsafe.Pointer)(unsafe.Pointer(&q.tail.prev)))) {

		atomic.AddInt32(&q.count, -1)
		atomic.StorePointer(
			(*unsafe.Pointer)(unsafe.Pointer(&prev.next)),
			unsafe.Pointer(q.tail))
		atomic.StorePointer(
			(*unsafe.Pointer)(unsafe.Pointer(&q.tail.prev)),
			unsafe.Pointer(prev))
		atomic.StoreInt32(&prev.takenFlag, 0)
		return lastNode.task
	}
	atomic.StoreInt32(&prev.takenFlag, 0)
	atomic.StoreInt32(&lastNode.takenFlag, 0)
	return nil
}

func (q *UnBoundedDEQueue) Size() int32 { return atomic.LoadInt32(&q.count) }

func (q *UnBoundedDEQueue) IsEmpty() bool { return q.Size() == 0 }
