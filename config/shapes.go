// shapes.go
// Christian Jordan
// Shape data structures

package config

// Point is a general struct used for points
type Point struct {
	X float32
	Y float32
}

// Rectangle is a obstacle rectangle
type Rectangle struct {
	pt *Point
	w  float32
	h  float32
}

// Circle is a obstacle circle
type Circle struct {
	pt *Point
	r  float32
}

// NewPoint creates a new Point
func NewPoint(x, y float32) *Point {
	return &Point{x, y}
}

// NewRectangle creates a new Rectangle
func NewRectangle(x, y, w, h float32) *Rectangle {
	return &Rectangle{NewPoint(x, y), w, h}
}

// NewCircle creates a new Circle
func NewCircle(x float32, y float32, r float32) *Circle {
	return &Circle{NewPoint(x, y), r}
}

// Checks if a point collides with a Rectangle
func (r *Rectangle) Collision(pt *Point) bool {
	if pt.X >= r.pt.X && pt.X <= r.pt.X+r.w &&
		pt.Y >= r.pt.Y && pt.Y <= r.pt.Y+r.h {
		return true
	}
	return false
}

// Checks if a point collides with a Circle
func (c *Circle) Collision(pt *Point) bool {
	if CalcDistance(c.pt, pt) <= c.r {
		return true
	}
	return false
}

func (r *Rectangle) Draw() {}
func (c *Circle) Draw()    {}
