package commands

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/statecrafthq/borg/utils"

	"github.com/urfave/cli"

	"encoding/json"

	"cloud.google.com/go/storage"
)

type CurrentSyncStatus struct {
	Hash   string `json:"hash"`
	Latest string `json:"latest"`
}

func readStatus(ctx context.Context, bucket *storage.BucketHandle, name string) (*CurrentSyncStatus, error) {
	currentName := "imports/" + name + "/CURRENT"
	reader, err := bucket.Object(currentName).NewReader(ctx)
	if err != nil {
		if err.Error() != "storage: object doesn't exist" {
			return nil, err
		}
		return nil, nil
	}
	defer reader.Close()
	ex, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	res := &CurrentSyncStatus{}
	err = json.Unmarshal(ex, res)
	if err == nil {
		return res, nil
	}
	return nil, err
}

func writeStatus(ctx context.Context, bucket *storage.BucketHandle, name string, hash string, fileName string) error {
	currentName := "imports/" + name + "/CURRENT"
	writer := bucket.Object(currentName).NewWriter(ctx)
	state := &CurrentSyncStatus{Hash: hash, Latest: fileName}
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	_, err = writer.Write(data)
	if err != nil {
		return err
	}
	err = writer.Close()
	if err != nil {
		return err
	}
	return nil
}

func uploadFile(ctx context.Context, bucket *storage.BucketHandle, name string, fileName string, src string) error {
	writer := bucket.Object("imports/" + name + "/" + fileName).NewWriter(ctx)
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(writer, file)
	if err != nil {
		return err
	}
	return writer.Close()
}

func sync(c *cli.Context) error {
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

	// Calculate HASH of a current file
	hash, err := utils.SHA256File(file)
	if err != nil {
		return err
	}

	// Init Bucket API
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	bucket := client.Bucket("data.openland.com")

	// Loading latest state
	status, err := readStatus(ctx, bucket, name)
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
		return nil
	}

	// Upload new version
	log.Println("Changed!")
	ext := filepath.Ext(file)
	fname := name + "_" + (time.Now().Format("2006_01_02_150405")) + ext
	log.Printf(fname)
	err = uploadFile(ctx, bucket, name, fname, file)
	if err != nil {
		return err
	}

	// Persisting state
	err = writeStatus(ctx, bucket, name, hash, fname)
	if err != nil {
		return err
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
	}
}
