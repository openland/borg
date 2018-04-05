package commands

import (
	"fmt"
	"math"

	"github.com/statecrafthq/borg/commands/ops"
	"github.com/statecrafthq/borg/geometry"
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
		if g, ok := row["geometry"]; ok {
			multipoly := geometry.NewGeoMultipolygon(utils.ParseFloat4(g.([]interface{})))
			proj := geometry.NewProjection(multipoly.Center())
			projected := multipoly.Project(proj)

			// Classificator
			if len(projected.Polygons) > 1 {
				multiCount++
				notConvex++
				notAnalyzed++

				extras.AppendString("shape_type", "miltipolygon")
			} else {
				poly := projected.Polygons[0]
				t := poly.Classify()
				if t == geometry.TypeComplexPolygon {
					notConvex++

					extras.AppendString("shape_type", "complex")
				} else if t == geometry.TypePolygonWithHoles {
					withHolesCount++
					notConvex++

					extras.AppendString("shape_type", "complex")
				} else if t == geometry.TypeTriangle {
					trianglesCount++

					// Upgrade field data
					sides := poly.Edges()
					extras.AppendString("shape_type", "triangle")
					extras.AppendFloat("side1", sides[0])
					extras.AppendFloat("side2", sides[1])
					extras.AppendFloat("side3", sides[2])
				} else if t == geometry.TypeRectangle {
					rectangleCount++
					fourPointCount++
					// Upgrade field data
					sides := poly.Edges()
					small := math.Min((sides[0]+sides[2])/2, (sides[1]+sides[3])/2)
					large := math.Max((sides[0]+sides[2])/2, (sides[1]+sides[3])/2)
					extras.AppendString("shape_type", "rectangle")
					extras.AppendFloat("side1", large)
					extras.AppendFloat("side2", small)
				} else if t == geometry.TypeQuadriliteral {
					fourPointCount++
					// Upgrade field data
					sides := poly.Edges()
					extras.AppendString("shape_type", "quadriliteral")
					extras.AppendFloat("side1", sides[0])
					extras.AppendFloat("side2", sides[1])
					extras.AppendFloat("side3", sides[2])
					extras.AppendFloat("side4", sides[3])
				} else if t == geometry.TypeConvexPolygon {
					extras.AppendString("shape_type", "convex")
				} else if t == geometry.TypeBroken {
					emptyCount++
					extras.AppendString("shape_type", "broken")
				}

				//
				// Building Fitting
				//

				// Kassita-1: 12ft x 35ft (3.6576 x 10.668)
				// Kassita-2: 10ft x 35ft (3.048  x 12.192)
				kassita1 := ops.LayoutRectangle(poly, 3.6576, 10.668)
				kassita2 := ops.LayoutRectangle(poly, 3.048, 12.192)

				if kassita1.Analyzed || kassita2.Analyzed {
					extras.AppendString("analyzed", "true")
				} else {
					extras.AppendString("analyzed", "false")
					notAnalyzed++
				}

				if kassita1.Analyzed && kassita1.Fits {
					extras.AppendString("project_kassita1", "true")
					if kassita1.HasLocation {
						extras.AppendFloat("project_kassita1_angle", kassita1.Angle)
						center := kassita1.Center.Unproject(proj)
						extras.AppendFloat("project_kassita1_lon", center.Longitude)
						extras.AppendFloat("project_kassita1_lat", center.Latitude)
					}
				} else {
					extras.AppendString("project_kassita1", "false")
				}

				if kassita2.Analyzed && kassita2.Fits {
					extras.AppendString("project_kassita2", "true")
					if kassita2.HasLocation {
						extras.AppendFloat("project_kassita2_angle", kassita2.Angle)
						center := kassita2.Center.Unproject(proj)
						extras.AppendFloat("project_kassita2_lon", center.Longitude)
						extras.AppendFloat("project_kassita2_lat", center.Latitude)
					}
				} else {
					extras.AppendString("project_kassita2", "false")
				}
			}
		} else {
			emptyCount++
			notAnalyzed++

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
