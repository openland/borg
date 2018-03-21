package ops

import "github.com/statecrafthq/borg/commands/ops/simplify"

func OptimizeLine(line [][]float64) [][]float64 {
	res := simplify.Simplify(line, 0.00001, false)
	// Avoid too simplifyed lines
	if len(res) <= 3 {
		return line
	}
	return res
}

func OptimizePolygon(multipoly [][][][]float64) [][][][]float64 {
	res := make([][][][]float64, 0)
	for _, poly := range multipoly {
		p := make([][][]float64, 0)
		for _, ring := range poly {
			p = append(p, OptimizeLine(ring))
		}
		res = append(res, p)
	}
	return res
}
