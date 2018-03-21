package commands

import (
	"bufio"
	"encoding/json"
	"os"

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

	wktOutput, e := os.Create(dst)
	if e != nil {
		return e
	}
	writer := bufio.NewWriter(wktOutput)

	e = ops.RecordReader(src, func(row map[string]interface{}) error {
		if geom, ok := row["geometry"]; ok {

			// Convert types
			coords := utils.ParseFloat4(geom.([]interface{}))

			// Repair
			coords, e := utils.PolygonRepair(coords)
			if e != nil {
				return e
			}

			// Measure area
			// areas = append(areas, utils.MeasureArea(coords))

			// Optimize
			coords = ops.OptimizePolygon(coords)

			// Repair again
			coords, e = utils.PolygonRepair(coords)
			if e != nil {
				return e
			}

			// Save updated geometry
			row["geometry"] = coords
		}

		// Write result

		b, e := json.Marshal(row)
		if e != nil {
			return e
		}

		_, e = writer.Write(b)
		if e != nil {
			return e
		}

		_, e = writer.WriteString("\n")
		if e != nil {
			return e
		}

		return nil
	})
	if e != nil {
		return e
	}

	//
	// Tear Down
	//

	e = writer.Flush()
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
