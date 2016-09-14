// Package planar provides primitives for planar geometry
package planar

// An Axis in euclidean space
type Axis int

// Axis definitions
const (
	XAxis Axis = 0
	YAxis Axis = 1
)

// Invert returns the opposite axis
func (a Axis) Invert() Axis {
	if a == XAxis {
		return YAxis
	}
	return XAxis
}

// AxisBounds contains bounding lines for a given axis
type AxisBounds struct {
	Axis Axis
	Min  float64
	Max  float64
}

// NewAxisBounds creates an AxisBounds
func NewAxisBounds(a Axis, min, max float64) *AxisBounds {
	return &AxisBounds{
		Axis: a,
		Min:  min,
		Max:  max,
	}
}

// IntersectMin returns the intersection point with the min axis boundary
func (ab *AxisBounds) IntersectMin(start, end Coord) Coord {
	return NewLine(start, end).IntersectWithAxis(ab.Axis, ab.Min)
}

// IntersectMax returns the intersection point with the max axis boundary
func (ab *AxisBounds) IntersectMax(start, end Coord) Coord {
	return NewLine(start, end).IntersectWithAxis(ab.Axis, ab.Max)
}
