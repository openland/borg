package commands

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/statecrafthq/borg/commands/ops"
	"github.com/statecrafthq/borg/utils"
	"github.com/urfave/cli"

	"encoding/json"
)

func writeRecord(writer *bufio.Writer, new map[string]interface{}) error {
	bytes, err := json.Marshal(new)
	if err != nil {
		return err
	}
	_, err = writer.Write(bytes)
	if err != nil {
		return err
	}
	_, err = writer.WriteString("\n")
	if err != nil {
		return err
	}
	return nil
}

func doDiff(src string, updated string, out string, ignoreRemoved bool) error {
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
	// Diffing
	//

	err := ops.DiffReader(src, updated, func(srcLine *map[string]interface{}, updLine *map[string]interface{}) error {
		if srcLine != nil && updLine != nil {
			changed, e := utils.IsChanged(*srcLine, *updLine)
			if e != nil {
				return e
			}
			if changed {
				// Record changed
				e = writeRecord(writer, *updLine)
				if e != nil {
					return e
				}
			} else {
				// Record is same
			}
		} else if srcLine != nil {
			// Throw if there are missing record
			fmt.Println("Record was removed!")
			fmt.Println(srcLine)
			if !ignoreRemoved {
				return cli.NewExitError("Record was removed!", 1)
			}
		} else if updLine != nil {
			e = writeRecord(writer, *updLine)
			if e != nil {
				return e
			}
		} else {
			return errors.New("Internal inconsistency")
		}
		return nil
	})

	if err != nil {
		return err
	}

	e = writer.Flush()
	if e != nil {
		return e
	}

	return nil
}

func diff(c *cli.Context) error {
	// Validate argumens
	src := c.String("current")
	updated := c.String("updated")
	out := c.String("out")
	ignoreRemoved := c.Bool("ignore-removed")
	if src == "" {
		return cli.NewExitError("You should provide current file", 1)
	}
	if updated == "" {
		return cli.NewExitError("You should provide updated file", 1)
	}
	if out == "" {
		return cli.NewExitError("You should provide output file", 1)
	}

	return doDiff(src, updated, out, ignoreRemoved)
}

func CreateDiffCommands() []cli.Command {
	return []cli.Command{
		{
			Name:  "diff",
			Usage: "Get changed lines from ols file",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "current",
					Usage: "Path to old dataset",
				},
				cli.StringFlag{
					Name:  "updated",
					Usage: "Path to updated dataset",
				},
				cli.StringFlag{
					Name:  "out",
					Usage: "Path to differenced dataset",
				},
				cli.BoolFlag{
					Name:  "ignore-removed",
					Usage: "Ignore removed records",
				},
			},
			Action: func(c *cli.Context) error {
				return diff(c)
			},
		},
	}
}
