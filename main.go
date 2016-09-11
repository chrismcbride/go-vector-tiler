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
	clipped, err := clip.ByRectangle(goGeom, -75.20, -75.10, 39.90, 40)
	if err == clip.ErrEmptyResult {
		fmt.Printf("Empty clip")
		os.Exit(1)
	} else if err != nil {
		panic(err)
	}
	if err := writeGeojson(clipped, "clipped.geojson"); err != nil {
		panic(err)
	}
}

func readGeojson(in io.Reader) (geom.T, error) {
	geojsonBytes, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, err
	}
	var goGeom geom.T
	return goGeom, geojson.Unmarshal(geojsonBytes, &goGeom)
}

func writeGeojson(g geom.T, filename string) error {
	geojsonBytes, err := geojson.Marshal(g)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, geojsonBytes, 0644)
}
