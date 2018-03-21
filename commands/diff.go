package commands

import (
	"bufio"
	"errors"
	"os"

	"github.com/statecrafthq/borg/commands/ops"
	"github.com/statecrafthq/borg/utils"
	"github.com/urfave/cli"

	"encoding/json"
)

func writeDiffAdded(writer *bufio.Writer, added map[string]interface{}) error {
	record := make(map[string]interface{})
	record["action"] = "added"
	record["new"] = added
	bytes, err := json.Marshal(record)
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

func writeDiffRemoved(writer *bufio.Writer, removed map[string]interface{}) error {
	record := make(map[string]interface{})
	record["action"] = "removed"
	record["old"] = removed
	bytes, err := json.Marshal(record)
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

func writeDiffUpdated(writer *bufio.Writer, old map[string]interface{}, new map[string]interface{}) error {
	record := make(map[string]interface{})
	record["action"] = "updated"
	record["old"] = old
	record["new"] = new
	bytes, err := json.Marshal(record)
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

func diff(c *cli.Context) error {
	// Validate argumens
	src := c.String("current")
	updated := c.String("updated")
	out := c.String("out")
	if src == "" {
		return cli.NewExitError("You should provide current file", 1)
	}
	if updated == "" {
		return cli.NewExitError("You should provide updated file", 1)
	}
	if out == "" {
		return cli.NewExitError("You should provide output file", 1)
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
				e = writeDiffUpdated(writer, *srcLine, *updLine)
				if e != nil {
					return e
				}
			} else {
				// Record is same
			}
		} else if srcLine != nil {
			// Removed
			e = writeDiffRemoved(writer, *srcLine)
			if e != nil {
				return e
			}
		} else if updLine != nil {
			e = writeDiffAdded(writer, *updLine)
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
			},
			Action: func(c *cli.Context) error {
				return diff(c)
			},
		},
	}
}
