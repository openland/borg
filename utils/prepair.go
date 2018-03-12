package utils

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func invokeRepair(src string) (string, error) {
	_, e := exec.LookPath("prepair")
	if e != nil {
		return "", e
	}
	command := exec.Command("prepair", "--wkt", src)
	var out bytes.Buffer
	command.Stdout = &out
	err := command.Run()
	if err != nil {
		fmt.Printf("Execution failed: \n%s\n", out.String())
		return "", err
	}

	return out.String(), nil
}

func convertToWkt(src [][][][]float64) string {
	res := "MULTIPOLYGON ("
	isFirstPoly := true
	for _, poly := range src {
		p := "("
		isFirstCircle := true
		for _, circle := range poly {
			c := "("
			isFirst := true
			for _, p := range circle {
				if !isFirst {
					c = c + ", "
				} else {
					isFirst = false
				}
				c = c + fmt.Sprintf("%f %f", p[0], p[1])
			}
			c = c + ")"
			if !isFirstCircle {
				p = p + ", "
			} else {
				isFirstCircle = false
			}
			p = p + c
		}
		p = p + ")"
		if !isFirstPoly {
			res = res + ", "
		} else {
			isFirstPoly = false
		}
		res = res + p
	}
	return res + ")"
}

func splitBrackets(src string) []string {
	opened := 1
	start := 0
	w := strings.Trim(src, " ")
	res := make([]string, 0)
	for i, r := range w {
		if i <= start {
			continue
		}
		c := string(r)
		if c == "(" {
			opened = opened + 1
		} else if c == ")" {
			opened = opened - 1
			if opened < 0 {
				// Handle error!
			} else if opened == 0 {
				res = append(res, strings.Trim(w[start+1:i], " "))
				start = i + 2 // There should be comma between
			}
		}
	}
	if start < len(w) {
		res = append(res, strings.Trim(w[start+1:len(w)-1], " "))
	}
	return res
}

func parseWkt(src string) ([][][][]float64, error) {
	if !strings.HasPrefix(src, "MULTIPOLYGON") {
		return [][][][]float64{}, errors.New("String should start from MULTIPOLYGON")
	}
	w := strings.TrimPrefix(src, "MULTIPOLYGON")
	res := make([][][][]float64, 0)
	body := splitBrackets(w)[0]
	for _, poly := range splitBrackets(body) {
		polyParsed := make([][][]float64, 0)
		for _, circle := range splitBrackets(poly) {
			circleParsed := make([][]float64, 0)
			for _, point := range strings.Split(circle, ",") {
				pointSplited := strings.Split(strings.Trim(point, " "), " ")
				lat, err := strconv.ParseFloat(strings.Trim(pointSplited[0], " "), 64)
				if err != nil {
					return [][][][]float64{}, err
				}
				lon, err := strconv.ParseFloat(strings.Trim(pointSplited[1], " "), 64)
				if err != nil {
					return [][][][]float64{}, err
				}
				circleParsed = append(circleParsed, []float64{lat, lon})
			}

			polyParsed = append(polyParsed, circleParsed)
		}
		res = append(res, polyParsed)
	}
	return res, nil
}

func PolygonRepair(src [][][][]float64) ([][][][]float64, error) {
	wkt := convertToWkt(src)
	res, err := invokeRepair(wkt)
	if err != nil {
		return [][][][]float64{}, err
	}
	parsed, err := parseWkt(res)
	if err != nil {
		return [][][][]float64{}, err
	}
	return parsed, nil
}
