package commands

import (
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

			// Convert types
			coords := utils.ParseFloat4(geom.([]interface{}))

			// Repair
			coords, e := utils.PolygonRepair(coords)
			if e != nil {
				return nil, e
			}

			// Measure area
			// areas = append(areas, utils.MeasureArea(coords))

			// Optimize
			coords = ops.OptimizePolygon(coords)

			// Repair again
			coords, e = utils.PolygonRepair(coords)
			if e != nil {
				return nil, e
			}

			// Save updated geometry
			row["geometry"] = coords
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
