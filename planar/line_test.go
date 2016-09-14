package planar

import "testing"

func TestXAxisLineIntersection(t *testing.T) {
	// All lines intersect where X=10
	for _, tc := range []struct {
		line   *Line
		result Coord
	}{
		{ // basic case
			line:   NewLine(Coord{0, 0}, Coord{20, 20}),
			result: Coord{10, 10},
		},
		{ // slope > 1
			line:   NewLine(Coord{5, 10}, Coord{15, 25}),
			result: Coord{10, 17.5},
		},
		{ // point ends on axis line
			line:   NewLine(Coord{0, 10}, Coord{10, 20}),
			result: Coord{10, 20},
		},
		{ // start > end
			line:   NewLine(Coord{25, 45}, Coord{5, 5}),
			result: Coord{10, 15},
		},
	} {
		i := tc.line.IntersectWithAxis(XAxis, 10)
		if !i.Equals(tc.result) {
			t.Errorf("Expected %v, got %v. For context: %v", tc.result, i, tc)
		}
	}
}

func TestYAxisLineIntersection(t *testing.T) {
	// All lines intersect where Y=10
	for _, tc := range []struct {
		line   *Line
		result Coord
	}{
		{ // basic case
			line:   NewLine(Coord{0, 0}, Coord{20, 20}),
			result: Coord{10, 10},
		},
		{ // slope > 1
			line:   NewLine(Coord{10, 5}, Coord{25, 15}),
			result: Coord{17.5, 10},
		},
		{ // point ends on axis line
			line:   NewLine(Coord{10, 0}, Coord{20, 10}),
			result: Coord{20, 10},
		},
		{ // start > end
			line:   NewLine(Coord{45, 25}, Coord{5, 5}),
			result: Coord{15, 10},
		},
	} {
		i := tc.line.IntersectWithAxis(YAxis, 10)
		if !i.Equals(tc.result) {
			t.Errorf("Expected %v, got %v. For context: %v", tc.result, i, tc)
		}
	}
}
