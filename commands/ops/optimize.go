package ops

import "github.com/statecrafthq/borg/commands/ops/simplify"

func OptimizeLine(line [][]float64) [][]float64 {
	return simplify.Simplify(line, 0.0001, false)
}
