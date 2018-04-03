package commands

import (
	"fmt"

	"github.com/statecrafthq/borg/commands/ops"

	"github.com/statecrafthq/borg/utils"
	"github.com/urfave/cli"
)

func doFinalize(c *cli.Context) error {
	src := c.String("src")
	dst := c.String("dst")
	if src == "" {
		return cli.NewExitError("You should provide source file", 1)
	}
	if dst == "" {
		return cli.NewExitError("You should provide destination file", 1)
	}

	//
	// Tear Up
	//

	e := utils.AssumeNotExists(dst, c.Bool("force"))
	if e != nil {
		return e
	}

	//
	// Main Cycle
	//

	e = ops.RecordTransformer(src, dst, func(row map[string]interface{}) (map[string]interface{}, error) {
		if geom, ok := row["geometry"]; ok {

			// Check if already optimized
			if _, ok := row["$geometry_src"]; ok {
				return row, nil
			}

			// Convert types
			src := geom
			coords := utils.ParseFloat4(geom.([]interface{}))

			// Repair
			repaired, e := utils.PolygonRepair(coords)
			if e != nil {
				fmt.Println(row)
				fmt.Println(e)
			} else {
				coords = repaired

				// Measure area
				// areas = append(areas, utils.MeasureArea(coords))

				// Optimize
				coords = ops.OptimizePolygon(coords)

				// Repair again
				repairedAgain, e := utils.PolygonRepair(coords)
				if e != nil {
					fmt.Println(row)
					fmt.Println(repaired)
					fmt.Println(coords)
					fmt.Println(e)
				} else {
					coords = repairedAgain
				}
			}

			// Save updated geometry
			row["geometry"] = coords
			row["$geometry_src"] = src
		}
		return row, nil
	})
	if e != nil {
		return e
	}

	return nil
}

func CreateFinalizeCommands() []cli.Command {
	return []cli.Command{
		{
			Name:  "finalize",
			Usage: "Finalize Dataset",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "src",
					Usage: "Source dataset",
				},
				cli.StringFlag{
					Name:  "dst",
					Usage: "Destination dataset",
				},
				cli.BoolFlag{
					Name:  "force, f",
					Usage: "Overwrite file if exists",
				},
			},
			Action: func(c *cli.Context) error {
				return doFinalize(c)
			},
		},
	}
}
