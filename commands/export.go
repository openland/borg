package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/statecrafthq/borg/geometry"

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
	exportRetired := c.Bool("export-retired")

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
	err = ops.RecordReader(src, func(row map[string]interface{}) error {

		// Check retired
		if !exportRetired {
			if ret, ok := row["retired"]; ok {
				if ret.(bool) {
					return nil
				}
			}
		}

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
		record = record + "\"id\":\"" + row["id"].(string) + "\""
		bounds := geometry.NewGeoMultipolygon(utils.ParseFloat4(row["geometry"].([]interface{}))).Bounds()
		record = record + ",\"max_lat\":" + fmt.Sprintf("%f", bounds.MaxLatitude)
		record = record + ",\"max_lon\":" + fmt.Sprintf("%f", bounds.MaxLongitude)
		record = record + ",\"min_lat\":" + fmt.Sprintf("%f", bounds.MinLatitude)
		record = record + ",\"min_lon\":" + fmt.Sprintf("%f", bounds.MinLongitude)
		record = record + "}"

		// Geomertry
		record = record + ", \"geometry\": {"
		record = record + "\"type\":\"MultiPolygon\""
		record = record + ",\"coordinates\":"

		g, err := json.Marshal(row["geometry"])
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
					Name: "geometry",
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
						cli.BoolFlag{
							Name:  "export-retired",
							Usage: "Export retired records too",
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
