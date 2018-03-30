package ops

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"

	"cloud.google.com/go/storage"
	"github.com/statecrafthq/borg/utils"
)

type CurrentSyncStatus struct {
	Hash   string `json:"hash"`
	Latest string `json:"latest"`
}

func ReadStatus(ctx context.Context, bucket *storage.BucketHandle, name string) (*CurrentSyncStatus, error) {
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

func ReadStatusFromFile(fileName string) (*CurrentSyncStatus, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	ex, err := ioutil.ReadAll(file)
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

func WriteStatus(ctx context.Context, bucket *storage.BucketHandle, name string, hash string, fileName string) error {
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

func WriteStatusToFile(outFileName string, hash string, fileName string) error {
	file, err := os.Create(outFileName)
	if err != nil {
		return err
	}
	defer file.Close()
	state := &CurrentSyncStatus{Hash: hash, Latest: fileName}
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func UploadFile(ctx context.Context, bucket *storage.BucketHandle, name string, fileName string, src string) error {
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

func DownloadFile(ctx context.Context, bucket *storage.BucketHandle, name string, status CurrentSyncStatus, dst string) error {
	reader, err := bucket.Object("imports/" + name + "/" + status.Latest).NewReader(ctx)
	if err != nil {
		return err
	}
	defer reader.Close()
	file, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, reader)
	if err != nil {
		return err
	}
	hash, err := utils.SHA256File(dst)
	if err != nil {
		return err
	}
	if status.Hash != hash {
		return errors.New("Broken file")
	}
	return nil
}
