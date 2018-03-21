package ops

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"os"
	"sort"

	"github.com/cheggaaa/pb"
	"github.com/statecrafthq/borg/utils"
)

func SortFile(src string, dst string) error {

	// Open files
	srcFile, e := os.Open(src)
	if e != nil {
		return e
	}
	defer srcFile.Close()
	dstFile, e := os.Create(dst)
	if e != nil {
		return e
	}
	defer dstFile.Close()

	// Preflight configuration
	totalLines, e := utils.CountLines(srcFile)
	if e != nil {
		return e
	}
	reader := bufio.NewReader(srcFile)
	writer := bufio.NewWriter(dstFile)
	records := make(map[string][]byte)
	keys := make([]string, 0)

	// Reading all records and ids
	bar := pb.StartNew(totalLines)
	linesRead := 0
	for {
		line, e := reader.ReadBytes('\n')
		if e != nil {
			if e == io.EOF {
				break
			}
			return e
		}
		linesRead = linesRead + 1
		bar.Set(linesRead)

		// Parsing id
		var dst map[string]interface{}
		e = json.Unmarshal(line, &dst)
		if e != nil {
			return e
		}
		id, ok := dst["id"]
		if !ok {
			return errors.New("Unable to find ID field in the dataset")
		}
		ids := id.(string)
		_, ok = records[ids]
		if ok {
			return errors.New("Duplicate records with id=" + ids)
		}
		records[ids] = line
		keys = append(keys, ids)
	}
	bar.Finish()

	// Sorting
	sort.Strings(keys)

	// Writing results
	for _, key := range keys {
		record := records[key]
		_, e := writer.Write(record)
		if e != nil {
			return e
		}
	}
	e = writer.Flush()
	if e != nil {
		return e
	}

	return nil
}
