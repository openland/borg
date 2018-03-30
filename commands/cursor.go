package commands

import (
	"os"

	"gopkg.in/kyokomi/emoji.v1"

	"github.com/statecrafthq/borg/commands/ops"
	"github.com/urfave/cli"
)

func cursorGet(c *cli.Context) error {
	out := c.String("out")
	key := c.String("key")
	name := c.String("name")
	dataset := c.String("dataset")
	reset := c.Bool("reset")
	if out == "" {
		return cli.NewExitError("Output is not provided", 1)
	}
	if key == "" {
		return cli.NewExitError("Key file name is not provided", 1)
	}
	if dataset == "" {
		return cli.NewExitError("Dataset name is not provided", 1)
	}
	if name == "" {
		return cli.NewExitError("Cursor name is not provided", 1)
	}

	var cursor *ops.CurrentSyncStatus
	var err error

	// Latest cursor
	latestCursor, err := ops.ReadStatus("imports/" + dataset + "/CURRENT")
	if err != nil {
		return err
	}
	if latestCursor == nil {
		return cli.NewExitError("Unable to find dataset", 1)
	}

	// Reading cursor
	if !reset {
		cursor, err = ops.ReadStatus("cursors/" + name + "/CURSOR")
		if err != nil {
			return err
		}
	}

	if cursor == nil {
		emoji.Println(":file_cabinet: (Reset) Downloading latest dataset")
		// Loading latest if there are no cursors or reset
		err = ops.DownloadFile(dataset, *latestCursor, out)
		if err != nil {
			return err
		}
		cursor = latestCursor
	} else {
		// Loading latest and cursor'ed one
		if latestCursor.Hash == cursor.Hash {
			emoji.Println(":file_cabinet: Dataset not changed")
			// Create empty
			dstFile, e := os.Create(out)
			if e != nil {
				return e
			}
			defer dstFile.Close()
		} else {
			// Download latest
			emoji.Println(":file_cabinet: Downloading latest dataset")
			err = ops.DownloadFile(dataset, *latestCursor, "_latest.ols")
			if err != nil {
				return err
			}
			defer os.Remove("_latest.ols")

			// Download cursor
			emoji.Println(":file_cabinet: Downloading cursored dataset")
			err = ops.DownloadFile(dataset, *cursor, "_processed.ols")
			if err != nil {
				return err
			}
			defer os.Remove("_processed.ols")

			// Build diff
			emoji.Println(":file_cabinet: Diffing datasets")
			err = doDiff("_latest.ols", "_processed.ols", out)
			if err != nil {
				return err
			}
		}
	}

	// Exporting key
	err = ops.WriteStatusToFile(key, cursor.Hash, cursor.Latest)
	if err != nil {
		return err
	}

	return nil
}

func cursorSet(c *cli.Context) error {
	key := c.String("key")
	name := c.String("name")
	if key == "" {
		return cli.NewExitError("Key file name is not provided", 1)
	}
	if name == "" {
		return cli.NewExitError("Cursor name is not provided", 1)
	}

	cursor, err := ops.ReadStatusFromFile(key)
	if err != nil {
		return err
	}

	err = ops.WriteStatus("cursors/"+key+"/CURSOR", cursor.Hash, cursor.Latest)
	if err != nil {
		return err
	}

	return nil
}

func CreateCursorCommands() []cli.Command {
	return []cli.Command{
		{
			Name:  "cursor",
			Usage: "Cursor operations",
			Subcommands: []cli.Command{
				{
					Name:  "get",
					Usage: "Get current cursor",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name",
							Usage: "Unique name of cursor",
						},
						cli.StringFlag{
							Name:  "key",
							Usage: "Path to key",
							Value: "cursor.json",
						},
						cli.StringFlag{
							Name:  "dst,out",
							Usage: "Path to changed records",
						},
						cli.StringFlag{
							Name:  "dataset",
							Usage: "Unique name of dataset",
						},
						cli.BoolFlag{
							Name:  "reset",
							Usage: "Resetting cursor",
						},
					},
					Action: func(c *cli.Context) error {
						return cursorGet(c)
					},
				},
				{
					Name:  "set",
					Usage: "Set current cursor",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name",
							Usage: "Unique name of cursor",
						},
						cli.StringFlag{
							Name:  "key",
							Usage: "Path to key",
							Value: "cursor.json",
						},
					},
					Action: func(c *cli.Context) error {
						return cursorSet(c)
					},
				},
			},
		},
	}
}
