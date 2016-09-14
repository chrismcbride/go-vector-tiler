// Package clip provides geometric clipping primitives.
//
// Slices geometries by axis parallel lines. Based on the implementation in
// mapbox's geojson-vt project:
// https://raw.githubusercontent.com/mapbox/geojson-vt-cpp/b9df0caa3a5c06838528c1989d11a697ece6e79b/include/mapbox/geojsonvt/clip.hpp
package clip

import (
	"errors"
	"time"

	geom "github.com/twpayne/go-geom"

	"github.com/chrismcbride/go-vector-tiler/metrics"
	"github.com/chrismcbride/go-vector-tiler/planar"
)

// Errors
var (
	ErrEmptyResult     = errors.New("Empty geometry from clip")
	ErrUnsupportedType = errors.New("Unsupported geometry type")
)

// ByRectangle clips a geometry by a rectangle (inclusive).
// Returns an error on empty clip.
func ByRectangle(g geom.T, xmin, xmax, ymin, ymax float64) (geom.T, error) {
	defer metrics.LogElapsedTime(time.Now(), "clip.ByRectangle")

	xAxisBounds := planar.NewAxisBounds(planar.XAxis, xmin, xmax)
	clippedX, err := ByAxis(g, xAxisBounds)
	if err != nil {
		return nil, err
	}
	yAxisBounds := planar.NewAxisBounds(planar.YAxis, ymin, ymax)
	return ByAxis(clippedX, yAxisBounds)
}

// ByAxis inclusively clips a geometry between two bounds along an axis
func ByAxis(g geom.T, axisBounds *planar.AxisBounds) (geom.T, error) {
	switch geometry := g.(type) {
	case *geom.MultiPolygon:
		return MultiPolygonByAxis(geometry, axisBounds)
	case *geom.MultiLineString:
		return MultiLineStringByAxis(geometry, axisBounds)
	case *geom.MultiPoint:
		return MultiPointByAxis(geometry, axisBounds)
	case *geom.LineString:
		return LineStringByAxis(geometry, axisBounds)
	case *geom.Polygon:
		return PolygonByAxis(geometry, axisBounds)
	case *geom.LinearRing:
		return LinearRingByAxis(geometry, axisBounds)
	default:
		return nil, ErrUnsupportedType
	}
}

// MultiPolygonByAxis inclusively clips a multipolygon along axis bounds
func MultiPolygonByAxis(
	mp *geom.MultiPolygon,
	axisBounds *planar.AxisBounds) (*geom.MultiPolygon, error) {

	result := geom.NewMultiPolygon(geom.XY)
	numPolys := mp.NumPolygons()
	for i := 0; i < numPolys; i++ {
		clipped, err := PolygonByAxis(mp.Polygon(i), axisBounds)
		if err != ErrEmptyResult {
			result.Push(clipped)
		}
	}
	if result.NumPolygons() > 0 {
		return result, nil
	}
	return nil, ErrEmptyResult
}

// MultiLineStringByAxis inclusively clips a multiLineString along axis bounds
// Returns either a LineString or MultiLineString
func MultiLineStringByAxis(
	mls *geom.MultiLineString, axisBounds *planar.AxisBounds) (geom.T, error) {

	result := geom.NewMultiLineString(geom.XY)
	numLines := mls.NumLineStrings()
	for i := 0; i < numLines; i++ {
		clipped, err := LineStringByAxis(mls.LineString(i), axisBounds)
		if err != ErrEmptyResult {
			switch geometry := clipped.(type) {
			case *geom.LineString:
				result.Push(geometry)
			case *geom.MultiLineString:
				numLines2 := geometry.NumLineStrings()
				for j := 0; j < numLines2; j++ {
					result.Push(geometry.LineString(j))
				}
			}
		}
	}
	if result.Empty() {
		return nil, ErrEmptyResult
	}
	if result.NumLineStrings() == 1 {
		return result.LineString(0), nil
	}
	return result, nil
}

// MultiPointByAxis inclusively clips a multipoint along axis bounds
func MultiPointByAxis(
	mp *geom.MultiPoint,
	axisBounds *planar.AxisBounds) (*geom.MultiPoint, error) {

	var resultCoords []float64
	flatCoords := mp.FlatCoords()
	endIndex := len(flatCoords)
	for i := 0; i < endIndex; i += 2 {
		coord := planar.Coord(flatCoords[i:(i + 2)])
		value := coord.ValueAtAxis(axisBounds.Axis)

		if value >= axisBounds.Min && value <= axisBounds.Max {
			resultCoords = append(resultCoords, coord.X(), coord.Y())
		}
	}
	if len(resultCoords) > 0 {
		return geom.NewMultiPointFlat(geom.XY, resultCoords), nil
	}
	return nil, ErrEmptyResult
}

// PolygonByAxis inclusively clips a polygon along axis bounds
func PolygonByAxis(
	p *geom.Polygon, axisBounds *planar.AxisBounds) (*geom.Polygon, error) {

	result := geom.NewPolygon(geom.XY)
	numRings := p.NumLinearRings()
	for i := 0; i < numRings; i++ {
		clipped, err := LinearRingByAxis(p.LinearRing(i), axisBounds)
		if i == 0 && err == ErrEmptyResult {
			// If the shell is empty, no sense clipping the holes
			return nil, ErrEmptyResult
		}
		if err != ErrEmptyResult {
			result.Push(clipped)
		}
	}
	if result.NumLinearRings() > 0 {
		return result, nil
	}
	return nil, ErrEmptyResult
}

// LinearRingByAxis inclusively clips a linear ring along axis bounds
func LinearRingByAxis(
	lr *geom.LinearRing,
	axisBounds *planar.AxisBounds) (*geom.LinearRing, error) {

	var resultCoords []float64
	addCoord := func(c planar.Coord) {
		resultCoords = append(resultCoords, c.X(), c.Y())
	}
	flatCoords := lr.FlatCoords()
	endIndex := len(flatCoords) - 2
	for i := 0; i < endIndex; i += 2 {
		aCoord := planar.Coord(flatCoords[i:(i + 2)])
		bCoord := planar.Coord(flatCoords[(i + 2):(i + 4)])
		a := aCoord.ValueAtAxis(axisBounds.Axis)
		b := bCoord.ValueAtAxis(axisBounds.Axis)

		switch {
		case a < axisBounds.Min:
			if b >= axisBounds.Min {
				// ---|-->  |
				addCoord(axisBounds.IntersectMin(aCoord, bCoord))
				if b > axisBounds.Max {
					// ---|-----|-->
					addCoord(axisBounds.IntersectMax(aCoord, bCoord))
				}
			}
		case a > axisBounds.Max:
			if b <= axisBounds.Max {
				// |  <--|---
				addCoord(axisBounds.IntersectMax(aCoord, bCoord))
				if b < axisBounds.Min {
					// <--|----|---
					addCoord(axisBounds.IntersectMin(aCoord, bCoord))
				}
			}
		default:
			addCoord(aCoord)
			if b < axisBounds.Min {
				// <--|---  |
				addCoord(axisBounds.IntersectMin(aCoord, bCoord))
			} else if b > axisBounds.Max {
				// |  ---|-->
				addCoord(axisBounds.IntersectMax(aCoord, bCoord))
			}
		}

		if i == (endIndex-2) && b >= axisBounds.Min && b <= axisBounds.Max {
			// At the last point and B is in bounds. Include it, otherwise it will be
			// part of the next line segment
			addCoord(bCoord)
		}
	}
	// close the ring
	numCoords := len(resultCoords)
	if numCoords > 0 {
		firstCoord := planar.Coord(resultCoords[0:2])
		lastCoord := planar.Coord(resultCoords[numCoords-2 : numCoords])
		if !firstCoord.Equals(lastCoord) {
			addCoord(firstCoord)
		}
		return geom.NewLinearRingFlat(geom.XY, resultCoords), nil
	}
	return nil, ErrEmptyResult
}

// LineStringByAxis inclusively clips a line string along axis bounds. May
// return either a LineString or MultiLineString.
func LineStringByAxis(
	ls *geom.LineString, axisBounds *planar.AxisBounds) (geom.T, error) {

	var currentLineCoords []float64
	addCoord := func(c planar.Coord) {
		currentLineCoords = append(currentLineCoords, c.X(), c.Y())
	}
	multiLine := geom.NewMultiLineString(geom.XY)
	newLineString := func() {
		multiLine.Push(geom.NewLineStringFlat(geom.XY, currentLineCoords))
		currentLineCoords = make([]float64, 0)
	}
	flatCoords := ls.FlatCoords()
	endIndex := len(flatCoords) - 2
	for i := 0; i < endIndex; i += 2 {
		aCoord := planar.Coord(flatCoords[i:(i + 2)])
		bCoord := planar.Coord(flatCoords[(i + 2):(i + 4)])
		a := aCoord.ValueAtAxis(axisBounds.Axis)
		b := bCoord.ValueAtAxis(axisBounds.Axis)

		switch {
		case a < axisBounds.Min:
			if b >= axisBounds.Min {
				// ---|-->  |
				addCoord(axisBounds.IntersectMin(aCoord, bCoord))
				if b > axisBounds.Max {
					// ---|-----|-->
					addCoord(axisBounds.IntersectMax(aCoord, bCoord))
					newLineString()
				}
			}
		case a > axisBounds.Max:
			if b <= axisBounds.Max {
				// |  <--|---
				addCoord(axisBounds.IntersectMax(aCoord, bCoord))
				if b < axisBounds.Min {
					// <--|----|---
					addCoord(axisBounds.IntersectMin(aCoord, bCoord))
					newLineString()
				}
			}
		default:
			addCoord(aCoord)
			if b < axisBounds.Min {
				// <--|---  |
				addCoord(axisBounds.IntersectMin(aCoord, bCoord))
				newLineString()
			} else if b > axisBounds.Max {
				// |  ---|-->
				addCoord(axisBounds.IntersectMax(aCoord, bCoord))
				newLineString()
			}
		}

		if i == (endIndex-2) && b >= axisBounds.Min && b <= axisBounds.Max {
			// At the last point and B is in bounds. Include it, otherwise it will be
			// part of the next line segment
			addCoord(bCoord)
		}
	}

	// add the final line
	if len(currentLineCoords) > 0 {
		newLineString()
	}

	if multiLine.Empty() {
		return nil, ErrEmptyResult
	}
	if multiLine.NumLineStrings() == 1 {
		return multiLine.LineString(0), nil
	}
	return multiLine, nil
}
