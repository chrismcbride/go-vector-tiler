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

func TestPolygonByAxis(t *testing.T) {
	lr := geom.NewLinearRingFlat(geom.XY, []float64{
		0, 0,
		50, 0,
		50, 10,
		20, 10,
		20, 20,
		30, 20,
		30, 30,
		50, 30,
		50, 40,
		25, 40,
		25, 50,
		0, 50,
		0, 60,
		25, 60,
		0, 0,
	})
	poly := geom.NewPolygon(geom.XY)
	poly.Push(lr)

	clipped, err := PolygonByAxis(
		poly, planar.NewAxisBounds(planar.XAxis, 10, 40))
	if err != nil {
		t.Error(err)
	}
	coords := clipped.FlatCoords()
	expected := []float64{
		10, 0,
		40, 0,
		40, 10,
		20, 10,
		20, 20,
		30, 20,
		30, 30,
		40, 30,
		40, 40,
		25, 40,
		25, 50,
		10, 50,
		10, 60,
		25, 60,
		10, 24,
		10, 0,
	}
	if !compareFloatSlice(expected, coords) {
		t.Errorf("Expected %v, got %v", expected, coords)
	}
}

func TestLineStringByAxis(t *testing.T) {
	ls := geom.NewLineStringFlat(geom.XY, []float64{
		0, 0,
		50, 0,
		50, 10,
		20, 10,
		20, 20,
		30, 20,
		30, 30,
		50, 30,
		50, 40,
		25, 40,
		25, 50,
		0, 50,
		0, 60,
		25, 60,
		30, 60,
	})
	clipped, err := LineStringByAxis(
		ls, planar.NewAxisBounds(planar.XAxis, 10, 40))
	if err != nil {
		t.Error(err)
	}
	out := clipped.(*geom.MultiLineString)
	if out.NumLineStrings() != 4 {
		t.Errorf("Expected 4 line strings, got %d", out.NumLineStrings())
	}
	expected := [][]float64{
		{10, 0, 40, 0},
		{40, 10, 20, 10, 20, 20, 30, 20, 30, 30, 40, 30},
		{40, 40, 25, 40, 25, 50, 10, 50},
		{10, 60, 25, 60, 30, 60},
	}
	for i, coords := range expected {
		actual := out.LineString(i).FlatCoords()
		if !compareFloatSlice(coords, actual) {
			t.Errorf("Expected %v, got %v", coords, actual)
		}
	}
}

func TestMultiPointByAxis(t *testing.T) {
	mp := geom.NewMultiPointFlat(geom.XY, []float64{
		0, 0,
		50, 0,
		50, 10,
		20, 10,
		20, 20,
		30, 20,
		30, 30,
		50, 30,
		50, 40,
		25, 40,
		25, 50,
		0, 50,
		0, 60,
		25, 60,
	})
	clipped, err := MultiPointByAxis(
		mp, planar.NewAxisBounds(planar.XAxis, 10, 40))
	if err != nil {
		t.Error(err)
	}
	expected := []float64{
		20, 10,
		20, 20,
		30, 20,
		30, 30,
		25, 40,
		25, 50,
		25, 60,
	}
	coords := clipped.FlatCoords()
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
