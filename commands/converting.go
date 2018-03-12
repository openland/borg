package commands

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/statecrafthq/borg/commands/formats"
	"github.com/statecrafthq/borg/utils"
	"github.com/twpayne/go-geom/encoding/geojson"
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
	// Check if file in path
	//

	_, e := exec.LookPath("ogr2ogr")
	if e != nil {
		return e
	}

	//
	// Starting conversion
	//

	command := exec.Command("ogr2ogr", "-f", "GeoJSON", "-t_srs", "crs:84", dst, src)
	var out bytes.Buffer
	command.Stdout = &out
	err := command.Run()
	if err != nil {
		fmt.Printf("Execution failed: \n%s\n", out.String())
		log.Fatal(err)
	}
	return nil
}

func converGeoJson(c *cli.Context) error {
	src := c.String("src")
	dst := c.String("dst")
	if src == "" {
		return cli.NewExitError("Source file is not provided", 1)
	}
	if dst == "" {
		return cli.NewExitError("Destination file is not provided", 1)
	}

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
	geometry := &geojson.FeatureCollection{}
	err = json.Unmarshal(body, &geometry)
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

	for _, element := range geometry.Features {

		// Parsing IDs
		idValue, err := formats.NewYorkId(element)
		if err != nil {
			continue
		}

		// Parsing Coordinates
		coordinates, err := utils.ConvertGeometry(element.Geometry)
		if err != nil {
			return err
		}

		// Preparing Bundle
		fields := make(map[string]interface{})
		fields["id"] = idValue
		fields["geometry"] = coordinates
		log.Printf(idValue)

		// Writing
		marshaled, err := json.Marshal(fields)
		if err != nil {
			return err
		}
		fmt.Fprintln(w, string(marshaled))
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
							Name:  "format",
							Usage: "ny_blocks, ny_parcels, sf_blocks, sf_parcels",
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
