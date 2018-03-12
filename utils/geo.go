package utils

import (
	"errors"

	"github.com/twpayne/go-geom"
)

func convertCoord(coord geom.Coord) []float64 {
	return []float64{coord[0], coord[1]}
}

func convertCoordArray1(coords []geom.Coord) [][]float64 {
	res := make([][]float64, 0)
	for _, e := range coords {
		res = append(res, convertCoord(e))
	}
	return res
}

func convertCoordArray2(coords [][]geom.Coord) [][][]float64 {
	res := make([][][]float64, 0)
	for _, e := range coords {
		res = append(res, convertCoordArray1(e))
	}
	return res
}
func convertCoordArray3(coords [][][]geom.Coord) [][][][]float64 {
	res := make([][][][]float64, 0)
	for _, e := range coords {
		res = append(res, convertCoordArray2(e))
	}
	return res
}

func ConvertGeometry(g geom.T) ([][][][]float64, error) {
	switch g := g.(type) {
	case *geom.Polygon:
		return [][][][]float64{convertCoordArray2(g.Coords())}, nil
	case *geom.MultiPolygon:
		return convertCoordArray3(g.Coords()), nil
	default:
		return nil, errors.New("Unsupported poing type")
	}
}
