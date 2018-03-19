package commands

import (
	"bufio"
	"errors"
	"io"
	"os"
	"sort"

	"github.com/statecrafthq/borg/utils"
	"github.com/urfave/cli"
	"gopkg.in/cheggaaa/pb.v1"

	"encoding/json"
)

func sortFile(src string, dst string) error {

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

	// Prepare TEMP directory
	err := utils.PrepareTemp()
	if err != nil {
		return err
	}

	// Sorting
	e := sortFile(src, "./tmp/src.ols")
	if e != nil {
		return e
	}
	e = sortFile(updated, "./tmp/dst.ols")
	if e != nil {
		return e
	}

	// Preflight operations
	srcFile, e := os.Open("./tmp/src.ols")
	if e != nil {
		return e
	}
	defer srcFile.Close()
	updFile, e := os.Open("./tmp/dst.ols")
	if e != nil {
		return e
	}
	defer updFile.Close()
	dstFile, e := os.Create(out)
	if e != nil {
		return e
	}
	defer dstFile.Close()

	// Diffing
	var srcLine map[string]interface{}
	var updLine map[string]interface{}
	srcLoaded := false
	srcEOF := false
	updLoaded := false
	updEOF := false
	srcReader := bufio.NewReader(srcFile)
	updReader := bufio.NewReader(updFile)
	writer := bufio.NewWriter(dstFile)

	for {
		if srcLoaded && updLoaded {
			return errors.New("Invariant broken")
		}
		// Loading next chunk
		if !srcLoaded {
			if !srcEOF {
				line, e := srcReader.ReadBytes('\n')
				if e != nil {
					if e == io.EOF {
						srcEOF = true
					} else {
						return e
					}
				} else {
					srcLoaded = true
					e = json.Unmarshal(line, &srcLine)
					if e != nil {
						return e
					}
				}
			}
		}
		if !updLoaded {
			if !updEOF {
				line, e := updReader.ReadBytes('\n')
				if e != nil {
					if e == io.EOF {
						updEOF = true
					} else {
						return e
					}
				} else {
					updLoaded = true
					e = json.Unmarshal(line, &updLine)
					if e != nil {
						return e
					}
				}
			}
		}

		// Handling cases
		if !updLoaded && !srcLoaded {
			// All records are read
			break
		} else if updLoaded && !srcLoaded {

			// New Record
			e = writeDiffAdded(writer, updLine)
			if e != nil {
				return e
			}

			// Move to next records
			updLoaded = false
		} else if !updLoaded && srcLoaded {

			// Removed record
			e = writeDiffRemoved(writer, srcLine)
			if e != nil {
				return e
			}

			// Move to next records
			srcLoaded = false
		} else {
			sId := srcLine["id"].(string)
			uId := updLine["id"].(string)
			if sId != uId {
				// Updated or removed element
				if sId < uId {
					// sID was removed
					e = writeDiffRemoved(writer, srcLine)
					if e != nil {
						return e
					}
					srcLoaded = false
				} else {
					// uID was added
					e = writeDiffAdded(writer, updLine)
					if e != nil {
						return e
					}
					updLoaded = false
				}
			} else {
				changed, e := utils.IsChanged(srcLine, updLine)
				if e != nil {
					return e
				}
				if changed {
					// Record changed
					e = writeDiffUpdated(writer, srcLine, updLine)
					if e != nil {
						return e
					}
				} else {
					// Record is same
				}

				// Move to next records
				srcLoaded = false
				updLoaded = false
			}
		}
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
