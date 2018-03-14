package commands

import (
	"bufio"
	"io/ioutil"
	"os"

	"github.com/statecrafthq/borg/utils"
	"github.com/urfave/cli"
)

func merge(c *cli.Context) error {
	dst := c.String("dst")
	exist := utils.FileExists(dst)
	if dst == "" {
		return cli.NewExitError("Destination file is not provided", 1)
	}
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

	isFirst := true
	for _, file := range c.StringSlice("src") {
		body, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		err = utils.IterateFeaturesRaw(body, func(value []byte) error {
			if isFirst {
				isFirst = false
			} else {
				_, err = w.WriteString(",\n")
				if err != nil {
					return err
				}
			}
			_, err := w.Write(value)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
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

func CreateMergeCommands() []cli.Command {
	return []cli.Command{
		{
			Name:  "merge",
			Usage: "Merge GeoJSON files",
			Flags: []cli.Flag{
				cli.StringSliceFlag{
					Name:  "source, src",
					Usage: "Path to dataset",
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
				return merge(c)
			},
		},
	}
}
