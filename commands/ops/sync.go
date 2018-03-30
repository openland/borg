package ops

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"

	"github.com/statecrafthq/borg/utils"
	"gopkg.in/cheggaaa/pb.v1"
)

type PassThru struct {
	io.Reader
	progress *pb.ProgressBar
	read     int64
}

func (pt *PassThru) Read(p []byte) (int, error) {
	n, err := pt.Reader.Read(p)
	pt.read += int64(n)

	if err == nil {
		pt.progress.Set(int(pt.read))
	}

	return n, err
}

func Copy(writer io.Writer, reader io.Reader, total int64) error {
	progress := pb.New(int(total))
	progress.Start()
	_, e := io.Copy(writer, &PassThru{progress: progress, Reader: reader})
	progress.Finish()
	return e
}

type CurrentSyncStatus struct {
	Hash   string `json:"hash"`
	Latest string `json:"latest"`
}

func ReadStatus(fullPath string) (*CurrentSyncStatus, error) {
	bucket, err := CreateBucket()
	if err != nil {
		return nil, err
	}

	reader, err := bucket.Object(fullPath).NewReader(context.Background())
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

func WriteStatus(fullPath string, hash string, fileName string) error {
	bucket, err := CreateBucket()
	if err != nil {
		return err
	}
	writer := bucket.Object(fullPath).NewWriter(context.Background())
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

func UploadFile(name string, fileName string, src string) error {
	bucket, err := CreateBucket()
	if err != nil {
		return err
	}
	writer := bucket.Object("imports/" + name + "/" + fileName).NewWriter(context.Background())
	stat, err := os.Stat(src)
	if err != nil {
		return err
	}
	size := stat.Size()
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()
	err = Copy(writer, file, size)
	if err != nil {
		return err
	}
	return writer.Close()
}

func DownloadFile(name string, status CurrentSyncStatus, dst string) error {
	bucket, err := CreateBucket()
	if err != nil {
		return err
	}
	object := bucket.Object("imports/" + name + "/" + status.Latest)
	attr, err := object.Attrs(context.Background())
	if err != nil {
		return err
	}
	reader, err := object.NewReader(context.Background())
	if err != nil {
		return err
	}
	defer reader.Close()
	file, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer file.Close()

	err = Copy(file, reader, attr.Size)
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
