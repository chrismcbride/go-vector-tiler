package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"

	"github.com/chrismcbride/go-vector-tiler/clip"
)

func main() {
	goGeom, err := readGeojson(os.Stdin)
	if err != nil {
		panic(err)
	}
	if goGeom.Layout() != geom.XY {
		panic("Only 2D geometries supported")
	}
	switch goGeom.(type) {
	case *geom.MultiPolygon, *geom.Polygon, *geom.LinearRing:
		clipped, empty := clip.ClipByRectangle(goGeom, -75.20, -75.10, 39.90, 40)
		if empty {
			fmt.Printf("Empty clip")
			os.Exit(1)
		}
		if err := writeGeojson(clipped, "clipped.geojson"); err != nil {
			panic(err)
		}
		os.Exit(0)
	default:
		panic(fmt.Sprintf("Unsupported geom type: %T", goGeom))
	}
}

func readGeojson(in io.Reader) (geom.T, error) {
	if geojsonBytes, err := ioutil.ReadAll(in); err != nil {
		return nil, err
	} else {
		var goGeom geom.T
		err := geojson.Unmarshal(geojsonBytes, &goGeom)
		return goGeom, err
	}
}

func writeGeojson(g geom.T, filename string) error {
	if geojsonBytes, err := geojson.Marshal(g); err != nil {
		return err
	} else {
		return ioutil.WriteFile(filename, geojsonBytes, 0644)
	}
}
