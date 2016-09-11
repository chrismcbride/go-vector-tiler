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

	"github.com/chrismcbride/go-vector-tiler/bounds"
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
		aCmp := axisBounds.CompareCoord(aCoord)
		bCmp := axisBounds.CompareCoord(bCoord)

		switch {
		case aCmp == bounds.LessThan && bCmp != bounds.LessThan:
			// ---|-->  |
			addCoord(axisBounds.IntersectMin(aCoord, bCoord))
			if bCmp == bounds.GreaterThan {
				// ---|-----|-->
				addCoord(axisBounds.IntersectMax(aCoord, bCoord))
			}
		case aCmp == bounds.GreaterThan && bCmp != bounds.GreaterThan:
			// |  <--|---
			addCoord(axisBounds.IntersectMax(aCoord, bCoord))
			if bCmp == bounds.LessThan {
				// <--|----|---
				addCoord(axisBounds.IntersectMin(aCoord, bCoord))
			}
		case aCmp == bounds.Inside:
			addCoord(aCoord)
			if bCmp == bounds.LessThan {
				// <--|---  |
				addCoord(axisBounds.IntersectMin(aCoord, bCoord))
			} else if bCmp == bounds.GreaterThan {
				// |  ---|-->
				addCoord(axisBounds.IntersectMax(aCoord, bCoord))
			}
		}

		if i == endIndex && bCmp == bounds.Inside {
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
