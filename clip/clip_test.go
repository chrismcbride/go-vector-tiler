package clip

import (
	"testing"

	geom "github.com/twpayne/go-geom"

	"github.com/chrismcbride/go-vector-tiler/planar"
)

func TestByXAxis(t *testing.T) {
	lr := geom.NewLinearRingFlat(geom.XY, []float64{
		10, 10,
		30, 10,
		30, 30,
		10, 30,
		10, 10,
	})

	clipped, err := ByAxis(lr, planar.NewAxisBounds(planar.XAxis, 15, 25))
	if err != nil {
		t.Error(err)
	}
	out := clipped.(*geom.LinearRing)
	coords := out.FlatCoords()
	expected := []float64{
		15, 10,
		25, 10,
		25, 30,
		15, 30,
		15, 10,
	}
	if !compareFloatSlice(expected, coords) {
		t.Errorf("Expected %v, got %v", expected, coords)
	}
}

func TestByYAxis(t *testing.T) {
	lr := geom.NewLinearRingFlat(geom.XY, []float64{
		10, 10,
		30, 10,
		30, 30,
		10, 30,
		10, 10,
	})

	clipped, err := ByAxis(lr, planar.NewAxisBounds(planar.YAxis, 15, 35))
	if err != nil {
		t.Error(err)
	}
	out := clipped.(*geom.LinearRing)
	coords := out.FlatCoords()
	// Note that this output is rotated by 1 compared to the input. This is
	// because the first line segment is out of bounds and skipped initially.
	expected := []float64{
		30, 15,
		30, 30,
		10, 30,
		10, 15,
		30, 15,
	}
	if !compareFloatSlice(expected, coords) {
		t.Errorf("Expected %v, got %v", expected, coords)
	}
}

func TestByRectangle(t *testing.T) {
	lr := geom.NewLinearRingFlat(geom.XY, []float64{
		10, 10,
		30, 10,
		30, 30,
		10, 30,
		10, 10,
	})

	clipped, err := ByRectangle(lr, 15, 25, 15, 35)
	if err != nil {
		t.Error(err)
	}
	out := clipped.(*geom.LinearRing)
	coords := out.FlatCoords()
	expected := []float64{
		25, 15,
		25, 30,
		15, 30,
		15, 15,
		25, 15,
	}
	if !compareFloatSlice(expected, coords) {
		t.Errorf("Expected %v, got %v", expected, coords)
	}
}

func compareFloatSlice(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
