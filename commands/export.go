package commands

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/statecrafthq/borg/commands/ops"
	"github.com/statecrafthq/borg/utils"
	"github.com/urfave/cli"
)

func doExportParcels(c *cli.Context) error {
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

	file, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer file.Close()
	w := bufio.NewWriter(file)

	// Header of file
	_, err = w.WriteString(`{"type":"FeatureCollection", "features": [` + "\n")
	if err != nil {
		return err
	}

	// Body
	isFirst := true
	err = ops.RecordReader(src, func(a map[string]interface{}) error {
		if isFirst {
			isFirst = false
		} else {
			_, err = w.WriteString(",\n")
			if err != nil {
				return err
			}
		}

		record := "{\"type\": \"Feature\""

		// Properties
		record = record + ", \"properties\": {"
		record = record + "\"id\":\"" + a["id"].(string) + "\""
		record = record + "}"

		// Geomertry
		record = record + ", \"geometry\": {"
		record = record + "\"type\":\"MultiPolygon\""
		record = record + ",\"coordinates\":"
		g, err := json.Marshal(a["geometry"])
		if err != nil {
			return err
		}
		record = record + string(g)
		record = record + "}"
		// a["geometry"]

		// End
		record = record + "}"

		_, err = w.WriteString(record)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Footer of file
	_, err = w.WriteString("\n]}")
	if err != nil {
		return err
	}

	err = w.Flush()
	if err != nil {
		return err
	}

	return nil
}

func CreateExportCommands() []cli.Command {
	return []cli.Command{
		{
			Name:  "export",
			Usage: "Export Dataset",
			Subcommands: []cli.Command{
				{
					Name: "parcels",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "source, src",
							Usage: "Path to source file",
						},
						cli.StringFlag{
							Name:  "dest, dst",
							Usage: "Path to destination file",
						},
						cli.BoolFlag{
							Name:  "force, f",
							Usage: "Overwrite file if exists",
						},
					},
					Action: func(c *cli.Context) error {
						return doExportParcels(c)
					},
				},
			},
		},
	}
}
