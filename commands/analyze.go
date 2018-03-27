package commands

import (
	"fmt"
	"math"

	"github.com/statecrafthq/borg/commands/ops"
	"github.com/statecrafthq/borg/utils"
	"github.com/urfave/cli"
	"gopkg.in/kyokomi/emoji.v1"
)

func analyzeDataset(c *cli.Context) error {
	src := c.String("src")
	if src == "" {
		return cli.NewExitError("You should provide current file", 1)
	}
	dst := c.String("dst")
	if dst == "" {
		return cli.NewExitError("You should provide destination file", 1)
	}
	e := utils.AssumeNotExists(dst, c.Bool("force"))
	if e != nil {
		return e
	}

	//
	// Stats counter
	//

	totalCount := 0
	emptyCount := 0
	multiCount := 0
	withHolesCount := 0
	notAnalyzed := 0
	notConvex := 0
	trianglesCount := 0
	rectangleCount := 0
	fourPointCount := 0

	e = ops.RecordTransformer(src, dst, func(row map[string]interface{}) (map[string]interface{}, error) {
		totalCount++
		extras, e := ops.LoadExtras(row["extras"])
		if e != nil {
			return nil, e
		}
		if geometry, ok := row["geometry"]; ok {
			coords := utils.ParseFloat4(geometry.([]interface{}))

			// Project and convert to meters
			projected := utils.ProjectToPlane(coords)

			// Classificator
			t := ops.ClassifyParcelGeometry(projected)
			if t == ops.TypeMultipolygon {
				multiCount++
				notConvex++
				notAnalyzed++

				extras.AppendString("shape_type", "miltipolygon")
				extras.AppendString("analyzed", "false")
			} else if t == ops.TypeComplexPolygon {
				notConvex++
				notAnalyzed++

				extras.AppendString("shape_type", "complex")
				extras.AppendString("analyzed", "false")
			} else if t == ops.TypePolygonWithHoles {
				withHolesCount++
				notConvex++
				notAnalyzed++

				extras.AppendString("shape_type", "convex")
				extras.AppendString("analyzed", "false")
			} else if t == ops.TypeTriangle {
				trianglesCount++

				// Upgrade field data
				sides := utils.GetSides(projected[0][0])
				extras.AppendString("shape_type", "triangle")
				extras.AppendString("analyzed", "false")
				extras.AppendFloat("side1", sides[0])
				extras.AppendFloat("side2", sides[1])
				extras.AppendFloat("side3", sides[2])
			} else if t == ops.TypeRectangle {
				rectangleCount++
				fourPointCount++

				// Upgrade field data
				sides := utils.GetSides(projected[0][0])
				small := math.Min((sides[0]+sides[2])/2, (sides[1]+sides[3])/2)
				large := math.Max((sides[0]+sides[2])/2, (sides[1]+sides[3])/2)
				extras.AppendString("shape_type", "rectangle")
				extras.AppendString("analyzed", "true")
				extras.AppendFloat("side1", large)
				extras.AppendFloat("side2", small)

				// Check Kassitas
				// Element1: 12ft x 35ft (3.6576 x 10.668)
				// Element2: 10ft x 35ft (3.048  x 12.192)

				// Element 1
				if small > 3.6576 && large > 10.668 {
					extras.AppendString("project_kassita1", "true")
				} else {
					extras.AppendString("project_kassita1", "false")
				}

				// Element 2
				if small > 3.048 && large > 12.192 {
					extras.AppendString("project_kassita2", "true")
				} else {
					extras.AppendString("project_kassita2", "false")
				}

			} else if t == ops.TypeQuadriliteral {
				fourPointCount++

				// Upgrade field data
				sides := utils.GetSides(projected[0][0])
				extras.AppendString("shape_type", "quadriliteral")
				extras.AppendString("analyzed", "false")
				extras.AppendFloat("side1", sides[0])
				extras.AppendFloat("side2", sides[1])
				extras.AppendFloat("side3", sides[2])
				extras.AppendFloat("side4", sides[3])
			} else if t == ops.TypeConvexPolygon {
				notAnalyzed++

				extras.AppendString("shape_type", "convex")
				extras.AppendString("analyzed", "false")
			} else if t == ops.TypeBroken {
				emptyCount++
				extras.AppendString("shape_type", "broken")
				extras.AppendString("analyzed", "false")
			}
		} else {
			emptyCount++

			extras.AppendString("shape_type", "broken")
			extras.AppendString("analyzed", "false")
		}

		row["extras"] = extras
		return row, nil
	})
	if e != nil {
		return e
	}
	emoji.Printf(":bar_chart: Stats:\n")
	fmt.Printf("-- Total: %d\n", totalCount)
	fmt.Printf("-- Empty: %d\n", emptyCount)
	fmt.Printf("-- Triangles: %d\n", trianglesCount)
	fmt.Printf("-- Rectangle: %d\n", rectangleCount)
	fmt.Printf("-- Four point: %d\n", fourPointCount)
	fmt.Printf("-- Multi Poly: %d\n", multiCount)
	fmt.Printf("-- With Holes: %d\n", withHolesCount)
	fmt.Printf("-- Complex: %d\n", notConvex)
	fmt.Printf("-- Total Not Analyzed: %d\n", notAnalyzed)
	return nil
}

func CreateAnalyzeCommands() []cli.Command {
	return []cli.Command{
		{
			Name:  "analyze",
			Usage: "Analyze Datasets",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "source, src",
					Usage: "path to source file",
				},
				cli.StringFlag{
					Name:  "destination, dst",
					Usage: "path to destination file",
				},
				cli.BoolFlag{
					Name:  "force, f",
					Usage: "Overwrite file if exists",
				},
			},
			Action: func(c *cli.Context) error {
				return analyzeDataset(c)
			},
		},
	}
}
