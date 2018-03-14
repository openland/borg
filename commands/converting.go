package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/statecrafthq/borg/commands/drivers"
	"github.com/statecrafthq/borg/utils"
	"github.com/urfave/cli"
)

func convertShapefile(c *cli.Context) error {
	src := c.String("src")
	dst := c.String("dst")
	if src == "" {
		return cli.NewExitError("Source file is not provided", 1)
	}
	if dst == "" {
		return cli.NewExitError("Destination file is not provided", 1)
	}

	//
	// Check if exists
	//

	exist := utils.FileExists(dst)
	if exist {
		if c.Bool("force") {
			e := os.Remove(dst)
			if e != nil {
				return e
			}
		} else {
			return cli.NewExitError("File already exists. Use --force for overwriting.", 1)
		}
	}

	//
	// Starting conversion
	//

	return utils.ShapefileToGeoJson(src, dst)
}

func converGeoJson(c *cli.Context) error {
	src := c.String("src")
	dst := c.String("dst")
	driverID := c.String("driver")
	strict := c.Bool("strict")
	if src == "" {
		return cli.NewExitError("Source file is not provided", 1)
	}
	if dst == "" {
		return cli.NewExitError("Destination file is not provided", 1)
	}
	if driverID == "" {
		return cli.NewExitError("driver is not provided", 1)
	}

	allDrivers := drivers.Drivers()
	if _, ok := allDrivers[strings.ToLower(driverID)]; !ok {
		return cli.NewExitError("Unable to find required driver", 1)
	}
	driver := allDrivers[strings.ToLower(driverID)]

	//
	// Existing file
	//

	exist := utils.FileExists(dst)
	if exist {
		if c.Bool("force") {
			e := os.Remove(dst)
			if e != nil {
				return e
			}
		} else {
			return cli.NewExitError("File already exists. Use --force for overwriting.", 1)
		}
	}

	//
	// Decoding Geometry
	//

	body, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	//
	// Generating of JSVC
	//

	file, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer file.Close()
	w := bufio.NewWriter(file)

	//
	// Iterating each feature
	//
	err = utils.IterateFeatures(body, strict, func(feature *utils.Feature) error {

		// Loading ID
		idValue, err := driver.ID(feature)
		if err != nil {
			return err
		}

		// Parsing Coordinates
		// Ignore if geometry missing
		if feature.Geometry == nil {
			return nil
		}
		coordinates, err := utils.SerializeGeometry(*feature.Geometry)
		if err != nil {
			return err
		}

		// Fixing invalid polygons
		err = utils.ValidateGeometry(coordinates)
		if err != nil {
			coordinates, err = utils.PolygonRepair(coordinates)
			if err != nil {
				return err
			}
			err = utils.ValidateGeometry(coordinates)
			if err != nil {
				return err
			}
		}

		// Loading Extras
		extras := drivers.NewExtras()
		err = driver.Extras(feature, &extras)
		if err != nil {
			return err
		}

		// Preparing Bundle
		fields := make(map[string]interface{})
		fields["id"] = idValue[0]
		if len(idValue) > 1 {
			fields["displayId"] = idValue[1:]
		}
		fields["geometry"] = coordinates
		fields["extras"] = extras

		// Writing
		marshaled, err := json.Marshal(fields)
		if err != nil {
			return nil
		}
		fmt.Fprintln(w, string(marshaled))

		return nil
	})
	if err != nil {
		return err
	}

	return w.Flush()
}

func CreateConvertingCommands() []cli.Command {
	return []cli.Command{
		{
			Name:  "convert",
			Usage: "Convert Datasets",
			Subcommands: []cli.Command{
				{
					Name:  "shapefile",
					Usage: "Converting Shapefile to GeoJSON",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "source, src",
							Usage: "path to source file",
						},
						cli.StringFlag{
							Name:  "dest, dst",
							Usage: "path to destination file",
						},
						cli.BoolFlag{
							Name:  "force, f",
							Usage: "Overwrite file if exists",
						},
					},
					Action: func(c *cli.Context) error {
						return convertShapefile(c)
					},
				},
				{
					Name:  "geojson",
					Usage: "Converting GeoJSON to jsvc file for Openland importing",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "source, src",
							Usage: "path to source file",
						},
						cli.StringFlag{
							Name:  "dest, dst",
							Usage: "path to destination file",
						},
						cli.StringFlag{
							Name:  "format,driver",
							Usage: "ny_blocks, ny_parcels, sf_blocks, sf_parcels",
						},
						cli.BoolFlag{
							Name:  "force, f",
							Usage: "Overwrite file if exists",
						},
						cli.BoolFlag{
							Name:  "strict",
							Usage: "Crash on invalid record",
						},
					},
					Action: func(c *cli.Context) error {
						return converGeoJson(c)
					},
				},
			},
		},
	}
}
