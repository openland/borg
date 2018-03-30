package commands

import (
	"log"
	"path/filepath"
	"regexp"
	"time"

	"github.com/statecrafthq/borg/commands/ops"
	"github.com/statecrafthq/borg/utils"

	"github.com/urfave/cli"
)

func sync(c *cli.Context) error {
	file := c.String("file")
	name := c.String("name")
	if file == "" {
		return cli.NewExitError("Source file is not provided", 1)
	}
	if name == "" {
		return cli.NewExitError("Dataset name is not provided", 1)
	}
	statusPath := "imports/" + name + "/CURRENT"
	var validID = regexp.MustCompile(`^[a-z0-9_]+$`)
	if !validID.MatchString(name) {
		return cli.NewExitError("Invalid name", 1)
	}

	// Calculate HASH of a current file
	hash, err := utils.SHA256File(file)
	if err != nil {
		return err
	}

	// Loading latest state
	status, err := ops.ReadStatus(statusPath)
	if err != nil {
		return err
	}

	// Checking latest hash
	changed := true
	if status != nil {
		if status.Hash == hash {
			changed = false
		}
	}
	if !changed {
		log.Println("Dataset wasn't changed")
		return nil
	}

	// Upload new version
	log.Println("Dataset was changed")
	ext := filepath.Ext(file)
	fname := name + "_" + (time.Now().Format("2006_01_02_150405")) + ext
	err = ops.UploadFile(name, fname, file)
	if err != nil {
		return err
	}

	// Persisting state
	err = ops.WriteStatus(statusPath, hash, fname)
	if err != nil {
		return err
	}
	return nil
}

func download(c *cli.Context) error {
	file := c.String("file")
	name := c.String("name")
	if file == "" {
		return cli.NewExitError("Source file is not provided", 1)
	}
	if name == "" {
		return cli.NewExitError("Dataset name is not provided", 1)
	}
	var validID = regexp.MustCompile(`^[a-z0-9_]+$`)
	if !validID.MatchString(name) {
		return cli.NewExitError("Invalid name", 1)
	}
	statusPath := "imports/" + name + "/CURRENT"

	// Loading latest state
	var status *ops.CurrentSyncStatus
	var err error
	keyFile := c.String("key")
	if keyFile != "" {
		status, err = ops.ReadStatusFromFile(keyFile)
		if err != nil {
			return err
		}
	} else {
		status, err = ops.ReadStatus(statusPath)
		if err != nil {
			return err
		}
	}
	if status == nil {
		return cli.NewExitError("Unable to find dataset", 1)
	}

	// Downloading
	err = ops.DownloadFile(name, *status, file)
	if err != nil {
		return err
	}

	// Exporting key
	exportKey := c.String("export-key")
	if exportKey != "" {
		ops.WriteStatusToFile(exportKey, status.Hash, status.Latest)
	}

	return nil
}

func CreateSyncCommands() []cli.Command {
	return []cli.Command{
		{
			Name:  "sync",
			Usage: "Sync Dataset",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file",
					Usage: "Path to dataset",
				},
				cli.StringFlag{
					Name:  "name",
					Usage: "Unique name of dataset",
				},
			},
			Action: func(c *cli.Context) error {
				return sync(c)
			},
		},
		{
			Name:  "download",
			Usage: "Download Dataset",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file, out",
					Usage: "Path to dataset",
				},
				cli.StringFlag{
					Name:  "name",
					Usage: "Unique name of dataset",
				},
				cli.StringFlag{
					Name:  "key",
					Usage: "Explicitly download speicific version instead of latest one",
				},
				cli.StringFlag{
					Name:  "export-key",
					Usage: "Export key during download",
				},
			},
			Action: func(c *cli.Context) error {
				return download(c)
			},
		},
	}
}
