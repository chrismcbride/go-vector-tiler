package planar

import "testing"

func TestCoord(t *testing.T) {
	c1 := Coord{5, 10}
	if c1.X() != 5 {
		t.Errorf("Expected %d, got %f", 5, c1.X())
	}
	if c1.Y() != 10 {
		t.Errorf("Expected %d got %f", 10, c1.Y())
	}
	if c1.ValueAtAxis(XAxis) != c1.X() {
		t.Errorf("Expected %f, got %f", c1.X(), c1.ValueAtAxis(XAxis))
	}
	if c1.ValueAtAxis(YAxis) != c1.Y() {
		t.Errorf("Expected %f, got %f", c1.Y(), c1.ValueAtAxis(YAxis))
	}
	if !c1.Equals(Coord{5, 10}) {
		t.Error("Equals check failed")
	}
	if c1.Equals(Coord{5.0001, 10}) {
		t.Error("Not equala check failed")
	}
}
