package commands

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"

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
			},
		},
	}
}
