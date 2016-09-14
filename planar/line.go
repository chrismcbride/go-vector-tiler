package planar

// Line represents a 2D line
type Line struct {
	start Coord
	end   Coord
}

// NewLine constructs a Line
func NewLine(start, end Coord) *Line {
	return &Line{start, end}
}

// IntersectWithAxis returns the intersection point between a line and a
// position on an axis
func (l *Line) IntersectWithAxis(axis Axis, position float64) Coord {
	axisDistance := position - l.start.ValueAtAxis(axis)
	oppositeAxis := axis.Invert()
	riseOverRun :=
		(l.end.ValueAtAxis(oppositeAxis) - l.start.ValueAtAxis(oppositeAxis)) /
			(l.end.ValueAtAxis(axis) - l.start.ValueAtAxis(axis))
	newValue := axisDistance*riseOverRun + l.start.ValueAtAxis(oppositeAxis)
	if axis == XAxis {
		return Coord{position, newValue}
	}
	return Coord{newValue, position}
}
