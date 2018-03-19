package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"

	"github.com/buger/jsonparser"
	"github.com/twpayne/go-geom"
	enc "github.com/twpayne/go-geom/encoding/geojson"
	"gopkg.in/cheggaaa/pb.v1"
)

func serializeCoord(coord geom.Coord) []float64 {
	return []float64{coord[0], coord[1]}
}

func parseCoord(coord []float64) geom.Coord {
	return geom.Coord{coord[0], coord[1]}
}

func serializeCoordArray1(coords []geom.Coord) [][]float64 {
	res := make([][]float64, 0)
	for _, e := range coords {
		res = append(res, serializeCoord(e))
	}
	return res
}

func parseCoordArray1(coords [][]float64) []geom.Coord {
	res := make([]geom.Coord, 0)
	for _, e := range coords {
		res = append(res, parseCoord(e))
	}
	return res
}

func serializeCoordArray2(coords [][]geom.Coord) [][][]float64 {
	res := make([][][]float64, 0)
	for _, e := range coords {
		res = append(res, serializeCoordArray1(e))
	}
	return res
}

func parseCoordArray2(coords [][][]float64) [][]geom.Coord {
	res := make([][]geom.Coord, 0)
	for _, e := range coords {
		res = append(res, parseCoordArray1(e))
	}
	return res
}

func serializeCoordArray3(coords [][][]geom.Coord) [][][][]float64 {
	res := make([][][][]float64, 0)
	for _, e := range coords {
		res = append(res, serializeCoordArray2(e))
	}
	return res
}

func parseCoordArray3(coords [][][][]float64) [][][]geom.Coord {
	res := make([][][]geom.Coord, 0)
	for _, e := range coords {
		res = append(res, parseCoordArray2(e))
	}
	return res
}

func SerializeGeometry(g geom.T) ([][][][]float64, error) {
	switch g := g.(type) {
	case *geom.Polygon:
		return [][][][]float64{serializeCoordArray2(g.Coords())}, nil
	case *geom.MultiPolygon:
		return serializeCoordArray3(g.Coords()), nil
	default:
		log.Println(g)
		return nil, errors.New("Unsupported gometry type")
	}
}

func ParseGeometry(polys [][][][]float64) geom.T {
	return geom.NewMultiPolygon(geom.XY).MustSetCoords(parseCoordArray3(polys))
}

func ValidateGeometry(polygons [][][][]float64) error {
	for _, poly := range polygons {
		for _, cirlce := range poly {
			if len(cirlce) < 4 {
				return errors.New("too short circle length")
			}
			for i, pointI := range cirlce {
				for j, pointJ := range cirlce {
					if i == j || j == len(cirlce)-1 || i == len(cirlce)-1 {
						continue
					}
					if pointI[0] == pointJ[0] && pointI[1] == pointJ[1] {
						return errors.New("self touching polygon circle")
					}
				}
			}
		}
	}
	return nil
}

type Feature struct {
	Geometry   *[][][][]float64
	Properties map[string]interface{}
}

func IterateFeaturesRaw(data []byte, cb func(feature []byte) error) error {
	bar := pb.StartNew(len(data))
	var existingError error
	_, err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if existingError != nil {
			return
		}
		bar.Set(offset)
		existingError = cb(value)
	}, "features")
	if err != nil {
		return err
	}
	if existingError != nil {
		return existingError
	}
	return nil
}

func IterateFeatures(data []byte, strict bool, displayErrors bool, cb func(feature *Feature) error) error {
	bar := pb.StartNew(len(data))
	var existingError error
	_, err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if strict && existingError != nil {
			return
		}
		defer func() {
			// recover from panic if one occured. Set err to nil otherwise.
			if err := recover(); err != nil {
				// 	err = errors.New("array index out of bounds")
				if displayErrors {
					fmt.Println("Error in record:")
					fmt.Println(string(value))
					fmt.Println(err)
				}
			}
		}()
		bar.Set(offset)

		// Parsing Properties
		v, t, _, err := jsonparser.Get(value, "properties")
		if err != nil {
			log.Panic(err)
		}
		properties := make(map[string]interface{})
		if t != jsonparser.NotExist {
			err = json.Unmarshal(v, &properties)
			if err != nil {
				log.Panic(err)
			}
		}

		// Parsing Geometry
		v, t, _, err = jsonparser.Get(value, "geometry")
		if err != nil {
			log.Panic(err)
		}
		var geometry geom.T
		hasGeometry := false
		if t != jsonparser.NotExist && t != jsonparser.Null {
			err = enc.Unmarshal(v, &geometry)

			if err != nil {
				log.Panic(err)
			} else {
				hasGeometry = true
			}
		}

		// Building Feature
		var geometryRef *[][][][]float64
		if hasGeometry {
			res, err := SerializeGeometry(geometry)
			if err != nil {
				log.Panic(err)
			}
			geometryRef = &res
		}
		feature := &Feature{Properties: properties, Geometry: geometryRef}

		// If failed ignore all subsequent
		// TODO: How we can handle this better?
		existingError = cb(feature)
		if !strict && existingError != nil {
			panic(existingError)
		}
	}, "features")
	bar.Finish()

	if err != nil {
		return err
	}
	if existingError != nil {
		return existingError
	}
	return nil
}

// Ported from https://github.com/mapbox/geojson-area/blob/master/index.js

func rad(src float64) float64 {
	return src * math.Pi / 180
}

func measureRingArea(coords [][]float64) float64 {
	var p1, p2, p3 []float64
	var lowerIndex, middleIndex, upperIndex, i int
	var area float64
	area = 0
	coordsLength := len(coords)

	if coordsLength > 2 {
		for i = 0; i < coordsLength; i++ {
			if i == coordsLength-2 { // i = N-2
				lowerIndex = coordsLength - 2
				middleIndex = coordsLength - 1
				upperIndex = 0
			} else if i == coordsLength-1 { // i = N-1
				lowerIndex = coordsLength - 1
				middleIndex = 0
				upperIndex = 1
			} else { // i = 0 to N-3
				lowerIndex = i
				middleIndex = i + 1
				upperIndex = i + 2
			}
			p1 = coords[lowerIndex]
			p2 = coords[middleIndex]
			p3 = coords[upperIndex]
			area += (rad(p3[0]) - rad(p1[0])) * math.Sin(rad(p2[1]))
		}

		// wgs84.RADIUS = 6378137
		area = area * 6378137 * 6378137 / 2
	}

	return area
}

func MeasureArea(src [][][][]float64) float64 {
	var res float64 = 0
	for _, poly := range src {
		if len(poly) > 0 {
			res += math.Abs(measureRingArea(poly[0]))
			for i := 1; i < len(poly); i++ {
				res += math.Abs(measureRingArea(poly[i]))
			}
		}
	}
	return res
}
