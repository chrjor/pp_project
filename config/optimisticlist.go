//
// lock.go
// Christian Jordan
// Optimistic synchronized lock data structure for concurrent package
//

package config

import (
	"sync"
	"sync/atomic"
)

// MileStoneChildren is a struct that implements the ChildrenList interface
type MileStoneChildren struct {
	head       *child
	tail       *child
	cond       *sync.Cond
	updateFlag int32
}

// child is a struct that represents a child node
type child struct {
	body *MileStone // Milestone body
	dist float32    // Cost of the child
	next *child     // Next child
	lock sync.Mutex // Lock for updating cost value
}

// NewChild creates a new child node
func newChild(body *MileStone, next *child) *child {
	var dist float32
	if body == nil {
		dist = 0
	} else {
		dist = body.Cost
	}

	return &child{
		body: body,
		dist: dist,
		next: next,
		lock: sync.Mutex{},
	}
}

// NewChildrenList creates a new ChildrenList
func NewChildrenList() *MileStoneChildren {
	return &MileStoneChildren{
		head:       newChild(nil, nil),
		tail:       newChild(nil, nil),
		cond:       sync.NewCond(&sync.Mutex{}),
		updateFlag: 0,
	}
}

// Checks if the list is empty
func (c *MileStoneChildren) IsEmpty() bool {
	return c.head.next == c.tail
}

// Add adds a child to the list
func (c *MileStoneChildren) Add(body *MileStone) {
	// Wait on update flag if cost is being updated
	for c.updateFlag == 1 {
		c.cond.Wait()
	}

	// Create new child
	newChildRef := newChild(body, nil)

	// Create child refs
	prevChild := c.head
	curChild := c.head
	for {
		if curChild.dist > newChildRef.dist || curChild == c.tail {
			prevChild.lock.Lock()
			curChild.lock.Lock()

			// Start over if invalid
			if !c.validate(prevChild, curChild) {
				prevChild.lock.Unlock()
				curChild.lock.Unlock()
				prevChild = c.head
				curChild = c.head
				continue
			}
			defer prevChild.lock.Unlock()
			defer curChild.lock.Unlock()

			// Add new child if valid
			if prevChild == curChild {
				c.head = newChildRef
			} else {
				prevChild.next = newChildRef
			}
			newChildRef.next = curChild
			return

		} else {
			// Keep iterating down list
			prevChild = curChild
			curChild = prevChild.next
		}
	}
}

// Remove removes a child from the list. The function returns true if the
// child was removed, otherwise, false.
func (c *MileStoneChildren) Remove(child *MileStone) bool {
	// Wait on update flag if cost is being updated
	for c.updateFlag == 1 {
		c.cond.Wait()
	}

	// Empty feed
	if c.head.next == c.tail {
		return false
	}

	// Create child refs
	prevChild := c.head
	curChild := c.head
	for {
		if child.ParDist == curChild.dist && child == curChild.body {
			prevChild.lock.Lock()
			curChild.lock.Lock()

			// Start over if invalid
			if !c.validate(prevChild, curChild) {
				prevChild.lock.Unlock()
				curChild.lock.Unlock()
				prevChild = c.head
				curChild = c.head
				continue
			}
			defer curChild.lock.Unlock()
			defer prevChild.lock.Unlock()

			// Remove child if valid
			prevChild.next = curChild.next
			return true

		} else if curChild.next == nil {
			// Reached end of feed
			return false

		} else {
			// Child not found
			prevChild = curChild
			curChild = prevChild.next
		}
	}
}

// Contains determines whether a child is in the list. The function returns
// true if the child is in the list, otherwise, false.
func (c *MileStoneChildren) Contains(child *MileStone) bool {
	// Empty list
	if c.head.next == c.tail {
		return false
	}
	// Create child refs
	curChild := c.head
	for {
		if child.ParDist == curChild.dist && child == curChild.body {
			// Check if found
			return true

		} else if curChild.next == nil {
			// Reached end of feed
			return false

		} else {
			// Keep iterating down list, until found
			curChild = curChild.next
		}
	}
}

// validate determines whether a child is valid. The function returns
// true if the feed is valid, otherwise, false.
func (c *MileStoneChildren) validate(prevChild *child, curChild *child) bool {
	node := c.head.next
	for {
		if node == prevChild {
			return curChild == node.next
		} else if node.next == nil {
			return false
		} else {
			node = node.next
		}
	}
}

// Helper functions for MileStoneChildren list implementation

// Applies a function to all children in the list (c
func BranchApply(c *MileStoneChildren,
	f func(*MileStone, interface{}),
	data interface{},
) {
	// Flag all children in sub-branch for update
	flagBranch(c)
	// Apply function to all children in sub-branch and unflag
	branchUpdate(c, f, data)
}

// Recursively sets the update flag for all children in the list (c)
func flagBranch(c *MileStoneChildren) {
	// Set c's update flag
	if !atomic.CompareAndSwapInt32(&c.updateFlag, 0, 1) {
		for c.updateFlag == 1 {
			c.cond.Wait()
			if atomic.CompareAndSwapInt32(&c.updateFlag, 0, 1) {
				break
			}
		}
	}
	// Set update flag for all children
	node := c.head.next
	for {
		if !node.body.children.IsEmpty() {
			flagBranch(node.body.children)
		}
		if node.next == nil {
			return
		} else {
			node = node.next
		}
	}
}

// RecursiveApply applies a function to all children in the list (c)
func branchUpdate(c *MileStoneChildren,
	f func(*MileStone, interface{}),
	data interface{},
) {
	// Set c's update flag
	c.cond.L.Lock()
	defer c.cond.L.Unlock()
	node := c.head.next
	for {
		// Apply function to child and all children of child
		f(node.body, data)
		if !node.body.children.IsEmpty() {
			branchUpdate(node.body.children, f, data)
		}
		if node.next == nil {
			// Reset update flag and notify when done
			atomic.AddInt32(&c.updateFlag, -1)
			c.cond.Broadcast()
			return

		} else {
			// Keep iterating down list if not done
			node = node.next
		}
	}
}
