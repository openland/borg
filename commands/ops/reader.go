package ops

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/statecrafthq/borg/utils"
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
