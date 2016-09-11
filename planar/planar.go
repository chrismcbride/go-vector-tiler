// Package planar provides primitives for planar geometry
package planar

import "github.com/chrismcbride/go-vector-tiler/bounds"

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

// Line creates a line on this axis at a given value
func (a Axis) Line(axisValue float64) AxisLine {
	return AxisLine{a, axisValue}
}

// Coord represents A 2D coordinate as a slice of 64bit floats
type Coord []float64

// X return the X value of the coord. Assumes that is position 0.
func (c Coord) X() float64 {
	return c[int(XAxis)]
}

// Y returns the Y value of the coord. Assumes that is position 1.
func (c Coord) Y() float64 {
	return c[int(YAxis)]
}

// ValueAtAxis returns the coord value at a given axis
func (c Coord) ValueAtAxis(axis Axis) float64 {
	return c[int(axis)]
}

// Equals compares two coords
func (c Coord) Equals(other Coord) bool {
	return (c.X() == other.X()) && (c.Y() == other.Y())
}

// AxisLine represents an infinite line at a given value on an axis
type AxisLine struct {
	axis  Axis
	value float64
}

// Intersection takes two coordinates representing a line, and returns the
// coordinate where the line meets the infinite axis line
func (al AxisLine) Intersection(start, end Coord) Coord {
	axisDist := float64(al.value) - start.ValueAtAxis(al.axis)
	other := al.axis.Invert()
	riseOverRun := (end.ValueAtAxis(other) - start.ValueAtAxis(other)) /
		(end.ValueAtAxis(al.axis) - start.ValueAtAxis(al.axis))
	newValue := axisDist*riseOverRun + start.ValueAtAxis(other)
	if al.axis == XAxis {
		return Coord{float64(al.value), newValue}
	}
	return Coord{newValue, float64(al.value)}
}

// AxisBounds inclusively compares geometries against axis boundaries
type AxisBounds struct {
	axis    Axis
	min     float64
	max     float64
	minLine AxisLine
	maxLine AxisLine
}

// NewAxisBounds creates an AxisBounds
func NewAxisBounds(a Axis, min, max float64) *AxisBounds {
	return &AxisBounds{
		axis:    a,
		min:     min,
		max:     max,
		minLine: a.Line(min),
		maxLine: a.Line(max),
	}
}

// CompareCoord checks a coordinate against the axis boundary
func (ab *AxisBounds) CompareCoord(c Coord) bounds.Result {
	axisValue := c.ValueAtAxis(ab.axis)
	if axisValue < ab.min {
		return bounds.LessThan
	} else if axisValue > ab.max {
		return bounds.GreaterThan
	}
	return bounds.Inside
}

// IntersectMin returns the intersection point with the min axis boundary
func (ab *AxisBounds) IntersectMin(start, end Coord) Coord {
	return ab.minLine.Intersection(start, end)
}

// IntersectMax returns the intersection point with the max axis boundary
func (ab *AxisBounds) IntersectMax(start, end Coord) Coord {
	return ab.maxLine.Intersection(start, end)
}
