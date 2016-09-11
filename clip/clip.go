// Package clip provides geometric clipping primitives.
//
// Slices geometries by axis parallel lines. Based on the implementation in
// mapbox's geojson-vt project:
// https://raw.githubusercontent.com/mapbox/geojson-vt-cpp/b9df0caa3a5c06838528c1989d11a697ece6e79b/include/mapbox/geojsonvt/clip.hpp
package clip

import (
	"errors"
	"time"

	"github.com/twpayne/go-geom"

	"github.com/chrismcbride/go-vector-tiler/metrics"
	"github.com/chrismcbride/go-vector-tiler/planar"
)

// Errors
var (
	ErrEmptyResult     = errors.New("Empty geometry from clip")
	ErrUnsupportedType = errors.New("Unsupported geometry type")
)

// ByRectangle clips a geometry by a rectangle. Returns an error on empty clip.
func ByRectangle(g geom.T, xmin, xmax, ymin, ymax float64) (geom.T, error) {
	defer metrics.LogElapsedTime(time.Now(), "clip.ByRectangle")

	clippedX, err := ByAxis(g, planar.XAxis, xmin, xmax)
	if err != nil {
		return nil, err
	}
	return ByAxis(clippedX, planar.YAxis, ymin, ymax)
}

// ByAxis clips a geometry between two bounds along an axis
func ByAxis(g geom.T, axis planar.Axis, min, max float64) (geom.T, error) {
	switch geometry := g.(type) {
	case *geom.MultiPolygon:
		return MultiPolygonByAxis(geometry, axis, min, max)
	case *geom.Polygon:
		return PolygonByAxis(geometry, axis, min, max)
	case *geom.LinearRing:
		return LinearRingByAxis(geometry, axis, min, max)
	default:
		return nil, ErrUnsupportedType
	}
}

// MultiPolygonByAxis clips a multipolygon along axis bounds
func MultiPolygonByAxis(
	mp *geom.MultiPolygon, axis planar.Axis,
	min, max float64) (*geom.MultiPolygon, error) {

	result := geom.NewMultiPolygon(geom.XY)
	numPolys := mp.NumPolygons()
	for i := 0; i < numPolys; i++ {
		clipped, err := PolygonByAxis(mp.Polygon(i), axis, min, max)
		if err != ErrEmptyResult {
			if err = result.Push(clipped); err != nil {
				// This error is only possible if the clipped polygon has a different
				// number of dimensions than the multipolygon. Since that would break
				// an invariant this library maintains, we're panicing here.
				panic(err)
			}
		}
	}
	if result.NumPolygons() > 0 {
		return result, nil
	}
	return nil, ErrEmptyResult
}

// PolygonByAxis clips a polygon along axis bounds
func PolygonByAxis(
	p *geom.Polygon, axis planar.Axis, min, max float64) (*geom.Polygon, error) {

	result := geom.NewPolygon(geom.XY)
	numRings := p.NumLinearRings()
	for i := 0; i < numRings; i++ {
		clipped, err := LinearRingByAxis(p.LinearRing(i), axis, min, max)
		if i == 0 && err == ErrEmptyResult {
			// If the shell is empty, no sense clipping the holes
			return nil, ErrEmptyResult
		}
		if err != ErrEmptyResult {
			if err = result.Push(clipped); err != nil {
				panic(err)
			}
		}
	}
	if result.NumLinearRings() > 0 {
		return result, nil
	}
	return nil, ErrEmptyResult
}

// LinearRingByAxis clips a linear ring along axis bounds
func LinearRingByAxis(
	lr *geom.LinearRing, axis planar.Axis,
	min, max float64) (*geom.LinearRing, error) {

	var resultCoords []float64
	addCoord := func(c planar.Coord) {
		resultCoords = append(resultCoords, c.X(), c.Y())
	}
	flatCoords := lr.FlatCoords()
	endIndex := len(flatCoords) - 2
	minLine := axis.Line(min)
	maxLine := axis.Line(max)
	for i := 0; i < endIndex; i += 2 {
		aCoord := planar.Coord(flatCoords[i:(i + 2)])
		bCoord := planar.Coord(flatCoords[(i + 2):(i + 4)])
		a := aCoord.ValueAtAxis(axis)
		b := bCoord.ValueAtAxis(axis)

		if a < min {
			if b >= min {
				// ---|-->  |
				addCoord(minLine.Intersection(aCoord, bCoord))
				if b > max { // ---|-----|-->
					addCoord(maxLine.Intersection(aCoord, bCoord))
				} else if i == endIndex {
					// At the last point and B is in bounds. Include it
					addCoord(bCoord)
				}
			}
		} else if a > max {
			if b <= max {
				// |  <--|---
				addCoord(maxLine.Intersection(aCoord, bCoord))
				if b < min { // <--|----|---
					addCoord(minLine.Intersection(aCoord, bCoord))
				} else if i == endIndex {
					// last point
					addCoord(bCoord)
				}
			}
		} else {
			addCoord(aCoord)
			if b <= min { // <--|---  |
				addCoord(minLine.Intersection(aCoord, bCoord))
			} else if b >= max { // |  ---|-->
				addCoord(maxLine.Intersection(aCoord, bCoord))
			}
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
