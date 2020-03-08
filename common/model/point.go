package model

// Math operations on a point
type Math interface {
	Translate(point2 *Point) *Point
}

// Point represents a grid point projected on a floor surface
type Point struct {
	X int
	Y int
}

// Translate two points to a new point
func (point *Point) Translate(point1 *Point) *Point {
	return &Point{point.X + point1.X, point.Y + point1.X}
}
