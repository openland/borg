package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/buger/jsonparser"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
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
		return nil, errors.New("Unsupported poing type")
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

func IterateFeatures(data []byte, strict bool, cb func(feature *geojson.Feature) error) error {
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
				fmt.Println("Error in record:")
				fmt.Println(string(value))
				fmt.Println(err)
			}
		}()
		bar.Set(offset)

		// TODO: Handle errors!
		feature := &geojson.Feature{}
		err = json.Unmarshal(value, &feature)
		if err != nil {
			log.Panic(err)
		}

		// If failed ignore all other
		existingError = cb(feature)
		if !strict && existingError != nil {
			panic(existingError)
		}
	}, "features")

	if err != nil {
		return err
	}
	if existingError != nil {
		return existingError
	}
	return nil
}
