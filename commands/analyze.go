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
	totalCount := 0
	emptyCount := 0
	multiCount := 0
	withHolesCount := 0
	notAnalyzed := 0
	notConvex := 0
	trianglesCount := 0
	rectangleCount := 0
	topIds := make([]string, 0)
	ops.RecordReader(src, func(row map[string]interface{}) error {
		totalCount++
		if geometry, ok := row["geometry"]; ok {
			coords := utils.ParseFloat4(geometry.([]interface{}))

			// Multipolygons
			if len(coords) > 1 {
				multiCount++
			}

			// Holes
			if len(coords[0]) > 1 {
				for _, poly := range coords {
					if len(poly) > 1 {
						withHolesCount++
						break
					}
				}
			}

			// Without holes and single polygon
			if len(coords) == 1 && len(coords[0]) == 1 {
				isConvex := true
				line := coords[0][0]
				wasPositive := false
				wasNegative := false
				for i := 0; i < len(line)-2; i++ {

					dx1 := line[i+1][1] - line[i][1]
					dy1 := line[i+1][0] - line[i][0]
					dx2 := line[i+2][1] - line[i+1][1]
					dy2 := line[i+2][0] - line[i+1][0]

					cross := dx1*dy2 - dy1*dx2

					if cross > 0 {
						if wasNegative {
							isConvex = false
							break
						}
						wasPositive = true
					} else if cross < 0 {
						if wasPositive {
							isConvex = false
							break
						}
						wasNegative = true
					}
				}
				if !isConvex {
					notConvex++
					notAnalyzed++
				} else {
					// If triangles
					if len(line) == 4 {
						trianglesCount++
					}
					if len(line) == 5 {
						rectangleCount++
						topIds = append(topIds, row["id"].(string))
						if len(topIds) == 1000 {
							fmt.Println(line)
							for i := 0; i < len(line)-2; i++ {

								dx1 := line[i+1][1] - line[i][1]
								dy1 := line[i+1][0] - line[i][0]
								dx2 := line[i+2][1] - line[i+1][1]
								dy2 := line[i+2][0] - line[i+1][0]

								fmt.Println(math.Atan2(dy2-dy1, dx2-dx1))
							}
						}
					}
				}
				// line := coords[0][0]
			} else {
				notAnalyzed++
			}
			//
		} else {
			emptyCount++
		}
		return nil
	})
	emoji.Printf(":bar_chart: Stats:\n")
	fmt.Printf("-- Total: %d\n", totalCount)
	fmt.Printf("-- Empty: %d\n", emptyCount)
	fmt.Printf("-- Triangles: %d\n", trianglesCount)
	fmt.Printf("-- Rectangle: %d\n", rectangleCount)
	fmt.Printf("-- Multi Poly: %d\n", multiCount)
	fmt.Printf("-- With Holes: %d\n", withHolesCount)
	fmt.Printf("-- Complex: %d\n", notConvex)
	fmt.Printf("-- Total Not Analyzed: %d\n", notAnalyzed)
	fmt.Println(topIds[1000])
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
			},
			Action: func(c *cli.Context) error {
				return analyzeDataset(c)
			},
		},
	}
}
