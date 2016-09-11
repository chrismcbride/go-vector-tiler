package planar

import (
	"testing"

	"github.com/chrismcbride/go-vector-tiler/bounds"
)

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

func TestXAxisLine(t *testing.T) {
	xLine := XAxis.Line(10)
	for _, tc := range []struct {
		start, end Coord
		result     Coord
	}{
		{
			start:  Coord{0, 0}, // basic case
			end:    Coord{20, 20},
			result: Coord{10, 10},
		},
		{
			start:  Coord{5, 10}, // slope > 1
			end:    Coord{15, 25},
			result: Coord{10, 17.5},
		},
		{
			start:  Coord{0, 10}, // point ends on axis line
			end:    Coord{10, 10},
			result: Coord{10, 10},
		},
		{
			start:  Coord{25, 45}, // start > end
			end:    Coord{5, 5},
			result: Coord{10, 15},
		},
	} {
		i := xLine.Intersection(tc.start, tc.end)
		if !i.Equals(tc.result) {
			t.Errorf("Expected %v, got %v. For context: %v", tc.result, i, tc)
		}
	}
}

func TestYAxisLine(t *testing.T) {
	yLine := YAxis.Line(10)
	for _, tc := range []struct {
		start, end Coord
		result     Coord
	}{
		{
			start:  Coord{0, 0}, // basic case
			end:    Coord{20, 20},
			result: Coord{10, 10},
		},
		{
			start:  Coord{10, 5}, // slope > 1
			end:    Coord{25, 15},
			result: Coord{17.5, 10},
		},
		{
			start:  Coord{10, 0}, // point ends on axis line
			end:    Coord{10, 10},
			result: Coord{10, 10},
		},
		{
			start:  Coord{45, 25}, // start > end
			end:    Coord{5, 5},
			result: Coord{15, 10},
		},
	} {
		i := yLine.Intersection(tc.start, tc.end)
		if !i.Equals(tc.result) {
			t.Errorf("Expected %v, got %v. For context: %v", tc.result, i, tc)
		}
	}
}

func TestInclusiveAxisBounds(t *testing.T) {
	xBounds := NewInclusiveAxisBounds(XAxis, 10, 20)
	for _, tc := range []struct {
		c      Coord
		result bounds.Result
	}{
		{
			c:      Coord{0, 0},
			result: bounds.LessThan,
		},
		{
			c:      Coord{0, 15}, // Should ignore Y Axis
			result: bounds.LessThan,
		},
		{
			c:      Coord{10, 15}, // On the line counts
			result: bounds.Inside,
		},
		{
			c:      Coord{15, 15},
			result: bounds.Inside,
		},
		{
			c:      Coord{20, 15},
			result: bounds.Inside,
		},
		{
			c:      Coord{21, 15},
			result: bounds.GreaterThan,
		},
	} {
		cmp := xBounds.CompareCoord(tc.c)
		if cmp != tc.result {
			t.Errorf("Expected %v, got %v. For context: %v", tc.result, cmp, tc)
		}
	}
}
