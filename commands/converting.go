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
	"strings"

	"github.com/buger/jsonparser"
	"github.com/statecrafthq/borg/commands/formats"
	"github.com/statecrafthq/borg/utils"
	"github.com/twpayne/go-geom/encoding/geojson"
	"github.com/urfave/cli"
	"gopkg.in/cheggaaa/pb.v1"
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
	formatID := c.String("format")
	if src == "" {
		return cli.NewExitError("Source file is not provided", 1)
	}
	if dst == "" {
		return cli.NewExitError("Destination file is not provided", 1)
	}
	if formatID == "" {
		return cli.NewExitError("Format is not provided", 1)
	}

	allFormats := formats.Formats()
	if _, ok := allFormats[strings.ToLower(formatID)]; !ok {
		return cli.NewExitError("Unable to find required format", 1)
	}
	format := allFormats[strings.ToLower(formatID)]

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

	// geometry := &geojson.FeatureCollection{}
	// err = json.Unmarshal(body, &geometry)
	// if err != nil {
	// 	return err
	// }

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
	bar := pb.StartNew(len(body))
	_, err = jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		defer func() {
			// recover from panic if one occured. Set err to nil otherwise.
			if err := recover(); err != nil {
				// 	err = errors.New("array index out of bounds")
				fmt.Println("Error in record:")
				fmt.Println(string(value))
			}
		}()

		bar.Set(offset)

		// jsonparser.Get(body, "geometry")

		// TODO: Handle errors!
		feature := &geojson.Feature{}
		err = json.Unmarshal(value, &feature)
		if err != nil {
			log.Panic(err)
		}

		// Parsing IDs
		idValue, err := format.ID(feature)
		if err != nil {
			log.Panic(err)
		}

		// Parsing Coordinates
		coordinates, err := utils.ConvertGeometry(feature.Geometry)
		if err != nil {
			log.Panic(err)
		}

		// Preparing Bundle
		fields := make(map[string]interface{})
		fields["id"] = idValue
		fields["geometry"] = coordinates

		// Writing
		marshaled, err := json.Marshal(fields)
		if err != nil {
			log.Panic(err)
		}
		fmt.Fprintln(w, string(marshaled))
	}, "features")

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
							Name:  "format",
							Usage: "ny_blocks, ny_parcels, sf_blocks, sf_parcels",
						},
						cli.BoolFlag{
							Name:  "force, f",
							Usage: "Overwrite file if exists",
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
