package planar

// Coord represents A 2D coordinate as a slice of 64bit floats
type Coord []float64

// X return the X value of the coord. Assumes that is position 0.
func (c Coord) X() float64 {
	return c.ValueAtAxis(XAxis)
}

// Y returns the Y value of the coord. Assumes that is position 1.
func (c Coord) Y() float64 {
	return c.ValueAtAxis(YAxis)
}

// ValueAtAxis returns the coord value at a given axis
func (c Coord) ValueAtAxis(axis Axis) float64 {
	return c[int(axis)]
}

// Equals compares two coords
func (c Coord) Equals(other Coord) bool {
	return (c.X() == other.X()) && (c.Y() == other.Y())
}
