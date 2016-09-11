// Package planar provides primitives for planar geometry
package planar

import "github.com/chrismcbride/go-vector-tiler/bounds"

// An Axis in euclidean space
type Axis int

// Axis definitions
const (
	XAxis Axis = 0
	YAxis      = 1
)

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
type AxisLine interface {
	// Given two coordinates representing a line, returns the coordiante where
	// the line meets the infinite axis line
	Intersection(start, end Coord) Coord
}

// Line creates a line on this axis at a given value
func (a Axis) Line(axisValue float64) AxisLine {
	switch a {
	case XAxis:
		return xAxisLine(axisValue)
	case YAxis:
		return yAxisLine(axisValue)
	default:
		panic("Invalid axis")
	}
}

type xAxisLine float64

func (x xAxisLine) Intersection(start, end Coord) Coord {
	xDistance := float64(x) - start.X()
	riseOverRun := (end.Y() - start.Y()) / (end.X() - start.X())
	newY := xDistance*riseOverRun + start.Y()
	return Coord{float64(x), newY}
}

type yAxisLine float64

func (y yAxisLine) Intersection(start, end Coord) Coord {
	yDistance := float64(y) - start.Y()
	runOverRise := (end.X() - start.X()) / (end.Y() - start.Y())
	newX := yDistance*runOverRise + start.X()
	return Coord{newX, float64(y)}
}

// InclusiveAxisBounds inclusively compares geometries against axis boundaries
type InclusiveAxisBounds struct {
	axis Axis
	min  float64
	max  float64
}

// NewInclusiveAxisBounds creates InclusiveAxisBounds
func NewInclusiveAxisBounds(a Axis, min, max float64) *InclusiveAxisBounds {
	return &InclusiveAxisBounds{a, min, max}
}

// CompareCoord checks a coordinate against the axis boundary
func (ab *InclusiveAxisBounds) CompareCoord(c Coord) bounds.Result {
	axisValue := c.ValueAtAxis(ab.axis)
	if axisValue < ab.min {
		return bounds.LessThan
	} else if axisValue > ab.max {
		return bounds.GreaterThan
	} else {
		return bounds.Inside
	}
}
