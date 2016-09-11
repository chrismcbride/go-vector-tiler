// Clips geometries by axis parallel lines. Based on the implementation in
// mapbox's geojson-vt project:
// https://raw.githubusercontent.com/mapbox/geojson-vt-cpp/b9df0caa3a5c06838528c1989d11a697ece6e79b/include/mapbox/geojsonvt/clip.hpp
package clip

import (
	"fmt"
  "time"

	"github.com/twpayne/go-geom"

	"github.com/chrismcbride/go-vector-tiler/metrics"
	"github.com/chrismcbride/go-vector-tiler/planar"
)

func ClipByRectangle(g geom.T, xmin, xmax, ymin, ymax float64) (geom.T, bool) {
	defer metrics.LogElapsedTime(time.Now(), "ClipByRectangle")

	if clippedX, empty := ClipByAxis(g, planar.XAxis, xmin, xmax); empty {
		return nil, true
	} else {
		return ClipByAxis(clippedX, planar.YAxis, ymin, ymax)
	}
}

func ClipByAxis(g geom.T, axis planar.Axis, min, max float64) (geom.T, bool) {
	var clippedGeom geom.T
	var isEmpty bool
	switch geometry := g.(type) {
	case *geom.MultiPolygon:
		mp := ClipMultiPolygonByAxis(geometry, axis, min, max)
		isEmpty = mp == nil
		clippedGeom = mp
	case *geom.Polygon:
		poly := ClipPolygonByAxis(geometry, axis, min, max)
		isEmpty = poly == nil
		clippedGeom = poly
	case *geom.LinearRing:
		lr := ClipLinearRingByAxis(geometry, axis, min, max)
		isEmpty = lr == nil
		clippedGeom = lr
	default:
		panic(fmt.Sprintf("Unsupported geom type: %T", geometry))
	}
	return clippedGeom, isEmpty
}

func ClipMultiPolygonByAxis(
	mp *geom.MultiPolygon, axis planar.Axis,
	min, max float64) *geom.MultiPolygon {

	result := geom.NewMultiPolygon(geom.XY)
	numPolys := mp.NumPolygons()
	for i := 0; i < numPolys; i++ {
		clipped := ClipPolygonByAxis(mp.Polygon(i), axis, min, max)
		if clipped != nil {
			if err := result.Push(clipped); err != nil {
				// This error is only possible if the clipped polygon has a different
				// number of dimensions than the multipolygon. Since that would break
				// an invariant this library maintains, we're panicing here.
				panic(err)
			}
		}
	}
	if result.NumPolygons() > 0 {
		return result
	}
	return nil
}

func ClipPolygonByAxis(
	p *geom.Polygon, axis planar.Axis, min, max float64) *geom.Polygon {

	result := geom.NewPolygon(geom.XY)
	numRings := p.NumLinearRings()
	for i := 0; i < numRings; i++ {
		clipped := ClipLinearRingByAxis(p.LinearRing(i), axis, min, max)
		if i == 0 && clipped == nil {
			// If the shell is empty, no sense clipping the holes
			return nil
		}
		if clipped != nil {
			if err := result.Push(clipped); err != nil {
				panic(err)
			}
		}
	}
	if result.NumLinearRings() > 0 {
		return result
	}
	return nil
}

func ClipLinearRingByAxis(
	lr *geom.LinearRing, axis planar.Axis, min, max float64) *geom.LinearRing {

	var resultCoords []float64
	addCoord := func(c planar.Coord) {
		resultCoords = append(resultCoords, c.X(), c.Y())
	}
	flatCoords := lr.FlatCoords()
	length := len(flatCoords)
	minLine := axis.Line(min)
	maxLine := axis.Line(max)
	for i := 0; i < length; i += 2 {
		aCoord := planar.Coord(flatCoords[i:(i + 2)])
		bCoord := planar.Coord(flatCoords[(i + 2):(i + 4)])
		a := aCoord.AtAxis(axis)
		b := bCoord.AtAxis(axis)

		if a < min {
			if b >= min {
				// ---|-->  |
				addCoord(minLine.Intersection(aCoord, bCoord))
				if b > max { // ---|-----|-->
					addCoord(maxLine.Intersection(aCoord, bCoord))
				} else if i == (length - 2) {
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
				} else if i == (length - 2) {
					// last point
					addCoord(bCoord)
				}
			}
		} else {
			addCoord(aCoord)
			if b < min { // <--|---  |
				addCoord(minLine.Intersection(aCoord, bCoord))
			} else if b > max { // |  ---|-->
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
		return geom.NewLinearRingFlat(geom.XY, resultCoords)
	}
	return nil
}
