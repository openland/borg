package commands

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/statecrafthq/borg/commands/ops"

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

func mergeOls(c *cli.Context) error {

	dst := c.String("dst")
	latest := c.String("latest")
	previous := c.String("previous")
	out := c.String("out")
	if latest == "" {
		return cli.NewExitError("You should provide latest file", 1)
	}
	if previous == "" {
		return cli.NewExitError("You should provide previous file", 1)
	}
	if out == "" {
		return cli.NewExitError("You should provide output file", 1)
	}

	// Destination
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

	//
	// Preflight operations
	//

	dstFile, e := os.Create(out)
	if e != nil {
		return e
	}
	defer dstFile.Close()
	writer := bufio.NewWriter(dstFile)

	//
	// Applying
	//

	e = ops.DiffReader(previous, latest, func(a *map[string]interface{}, b *map[string]interface{}) error {
		if a != nil && b != nil {
			// Merging two records
			merged, e := ops.Merge(*a, *b)

			// Writing to file
			bytes, e := json.Marshal(merged)
			if e != nil {
				return e
			}
			_, e = writer.Write(bytes)
			if e != nil {
				return e
			}
			_, e = writer.WriteString("\n")
			if e != nil {
				return e
			}
		} else if a != nil {
			bytes, e := json.Marshal(*a)
			if e != nil {
				return e
			}
			_, e = writer.Write(bytes)
			if e != nil {
				return e
			}
			_, e = writer.WriteString("\n")
			if e != nil {
				return e
			}
		} else if b != nil {
			bytes, e := json.Marshal(*a)
			if e != nil {
				return e
			}
			_, e = writer.Write(bytes)
			if e != nil {
				return e
			}
			_, e = writer.WriteString("\n")
			if e != nil {
				return e
			}
		}
		return nil
	})
	if e != nil {
		return e
	}

	e = writer.Flush()
	if e != nil {
		return e
	}

	return nil
}

func CreateMergeCommands() []cli.Command {
	return []cli.Command{
		{
			Name:  "merge",
			Usage: "Merge GeoJSON files",
			Subcommands: []cli.Command{
				{
					Name: "geojson",
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
				{
					Name: "ols",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "new, latest, updated",
							Usage: "Latest version of OLS file",
						},
						cli.StringFlag{
							Name:  "old, previous",
							Usage: "Path to previous dataset",
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
						return mergeOls(c)
					},
				},
			},
		},
	}
}
