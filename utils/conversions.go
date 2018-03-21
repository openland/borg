package utils

func ParseFloat1(coord []interface{}) []float64 {
	return []float64{coord[0].(float64), coord[1].(float64)}
}

func ParseFloat2(coord []interface{}) [][]float64 {
	res := make([][]float64, 0)
	for _, e := range coord {
		res = append(res, ParseFloat1(e.([]interface{})))
	}
	return res
}

func ParseFloat3(coord []interface{}) [][][]float64 {
	res := make([][][]float64, 0)
	for _, e := range coord {
		res = append(res, ParseFloat2(e.([]interface{})))
	}
	return res
}

func ParseFloat4(coord []interface{}) [][][][]float64 {
	res := make([][][][]float64, 0)
	for _, e := range coord {
		res = append(res, ParseFloat3(e.([]interface{})))
	}
	return res
}
