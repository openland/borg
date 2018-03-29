package ops

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"sync"

	"github.com/statecrafthq/borg/utils"
	"golang.org/x/sync/semaphore"
	"gopkg.in/cheggaaa/pb.v1"
)

func DiffReader(a string, b string, handler func(a *map[string]interface{}, b *map[string]interface{}) error) error {

	//
	// Preflight
	//

	// Working folder
	e := utils.PrepareTemp()
	if e != nil {
		return e
	}
	defer utils.ClearTemp()

	// Sort
	aLines, e := SortFile(a, "./tmp/a.ols")
	if e != nil {
		return e
	}
	bLines, e := SortFile(b, "./tmp/b.ols")
	if e != nil {
		return e
	}

	// Init readers
	srcFile, e := os.Open("./tmp/a.ols")
	if e != nil {
		return e
	}
	defer srcFile.Close()
	updFile, e := os.Open("./tmp/b.ols")
	if e != nil {
		return e
	}
	defer updFile.Close()

	//
	// Differing
	//

	var srcLine map[string]interface{}
	var updLine map[string]interface{}
	srcLoaded := false
	srcEOF := false
	updLoaded := false
	updEOF := false
	srcReader := bufio.NewReader(srcFile)
	updReader := bufio.NewReader(updFile)
	read := 0
	bar := pb.StartNew(aLines + bLines)
	for {
		bar.Set(read)
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
					read++
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
					read++
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

			// Added
			e = handler(nil, &updLine)
			if e != nil {
				return e
			}

			// Move to next records
			updLoaded = false
		} else if !updLoaded && srcLoaded {

			// Removed
			e = handler(&srcLine, nil)
			if e != nil {
				return e
			}

			// Move to next records
			srcLoaded = false
		} else {
			sID := srcLine["id"].(string)
			uID := updLine["id"].(string)
			if sID != uID {
				// Updated or removed element
				if sID < uID {
					// sID was removed
					e = handler(&srcLine, nil)
					if e != nil {
						return e
					}
					srcLoaded = false
				} else {
					// uID was added
					e = handler(nil, &updLine)
					if e != nil {
						return e
					}
					updLoaded = false
				}
			} else {

				// Both are present
				e = handler(&srcLine, &updLine)
				if e != nil {
					return e
				}

				// Move to next records
				srcLoaded = false
				updLoaded = false
			}
		}
	}

	bar.Finish()

	return nil
}

func RecordReader(src string, handler func(row map[string]interface{}) error) error {
	// Opening file
	file, e := os.Open(src)
	if e != nil {
		return e
	}
	defer file.Close()

	// Line number
	lines, e := utils.CountLines(file)
	if e != nil {
		return e
	}

	//
	// Main loop
	//
	bar := pb.StartNew(lines)
	defer bar.Finish()
	rd := bufio.NewReader(file)

	//
	// Main Loop
	//
	linesRead := 0
	for {
		line, e := rd.ReadBytes('\n')
		if e != nil {
			if e == io.EOF {
				break
			}
			return e
		}

		var d map[string]interface{}
		e = json.Unmarshal(line, &d)
		if e != nil {
			return e
		}
		bar.Set(linesRead)
		linesRead = linesRead + 1
		e = handler(d)
		if e != nil {
			return e
		}
	}
	return nil
}

func RecordTransformer(src string, dst string, handler func(row map[string]interface{}) (map[string]interface{}, error)) error {
	dstFile, e := os.Create(dst)
	defer dstFile.Close()
	if e != nil {
		return e
	}
	writer := bufio.NewWriter(dstFile)
	var writerLock sync.Mutex
	ctx := context.Background()
	sem := semaphore.NewWeighted(10)

	var perror error
	e = RecordReader(src, func(row map[string]interface{}) error {
		if perror != nil {
			return perror
		}
		if e := sem.Acquire(ctx, 1); e != nil {
			return e
		}
		go func() {
			defer sem.Release(1)
			c, e := handler(row)

			// Result
			if e != nil {
				perror = e
				return
			}
			b, e := json.Marshal(c)
			if e != nil {
				perror = e
				return
			}

			// Writing
			writerLock.Lock()
			defer writerLock.Unlock()
			_, e = writer.Write(b)
			if e != nil {
				perror = e
				return
			}

			_, e = writer.WriteString("\n")
			if e != nil {
				perror = e
				return
			}
		}()
		return nil
	})
	if perror != nil {
		return e
	}
	if e != nil {
		return e
	}
	if e := sem.Acquire(ctx, 10); e != nil {
		return e
	}
	writerLock.Lock()

	e = writer.Flush()
	if e != nil {
		return e
	}

	return nil
}
