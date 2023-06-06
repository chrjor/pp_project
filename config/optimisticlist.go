//
// lock.go
// Christian Jordan
// Optimistic synchronized lock data structure for concurrent package
//

package config

import (
	"sync"
)

// Children is an interface for a milestone's children
type MileStoneChildren interface {
	Add(*MileStone)
	Remove(*MileStone) bool
	Contains(*MileStone) bool
	GetChildren() []*MileStone
}

// milestoneChildren is a struct that implements the ChildrenList interface
type milestoneChildren struct {
	head *child
	tail *child
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
func NewChildrenList() MileStoneChildren {
	return &milestoneChildren{
		head: newChild(nil, nil),
		tail: newChild(nil, nil),
	}
}

// Add adds a child to the list
func (c *milestoneChildren) Add(body *MileStone) {
	// Create new child
	newChildRef := newChild(body, nil)

	// Create child refs
	prevChild := c.head
	curChild := c.head
	for {
		if curChild.dist > newChildRef.dist || curChild == c.tail {
			prevChild.lock.Lock()
			curChild.lock.Lock()
			defer prevChild.lock.Unlock()
			defer curChild.lock.Unlock()

			// Start over if invalid
			if !c.validate(prevChild, curChild) {
				prevChild.lock.Unlock()
				curChild.lock.Unlock()
				prevChild = c.head
				curChild = c.head
				continue
			}

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
func (c *milestoneChildren) Remove(child *MileStone) bool {
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
			defer prevChild.lock.Unlock()
			defer curChild.lock.Unlock()

			// Start over if invalid
			if !c.validate(prevChild, curChild) {
				prevChild.lock.Unlock()
				curChild.lock.Unlock()
				prevChild = c.head
				curChild = c.head
				continue
			}

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
func (c *milestoneChildren) Contains(child *MileStone) bool {
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
func (c *milestoneChildren) validate(prevChild *child, curChild *child) bool {
	node := c.head
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

// GetChildren returns a list of children
func (c *milestoneChildren) GetChildren() []*MileStone {
	children := make([]*MileStone, 0)
	node := c.head
	for {
		if node.next == nil {
			return children
		} else {
			children = append(children, node.body)
			node = node.next
		}
	}
}
