// Primitives for planar geometry
package planar

type Axis int

const (
	XAxis Axis = 0
	YAxis      = 1
)

type Coord []float64

func (c Coord) X() float64 {
	return c[int(XAxis)]
}

func (c Coord) Y() float64 {
	return c[int(YAxis)]
}

func (c Coord) AtAxis(axis Axis) float64 {
	return c[int(axis)]
}

func (c Coord) Equals(other Coord) bool {
	return (c.X() == other.X()) && (c.Y() == other.Y())
}

// An infinite line at a given value on an axis
type AxisLine interface {
	// Given two coordinates representing a line, returns the coordiante where
	// the line meets the infinite axis line
	Intersection(start, end Coord) Coord
}

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
