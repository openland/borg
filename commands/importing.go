package commands

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/statecrafthq/borg/utils"
	"github.com/urfave/cli"
	"gopkg.in/cheggaaa/pb.v1"
)

func builtInQueries() map[string]string {
	res := make(map[string]string)
	res["blocks"] = "mutation($data: [BlockInput!]!, $state: String!, $county: String!, $city: String!) { importBlocks(tate: $state, county: $county, city: $city, blocks: $data) }"
	res["parcels"] = "mutation($data: [ParcelInput!]!, $state: String!, $county: String!, $city: String!) { importParcels(state: $state, county: $county, city: $city, parcels: $data) }"
	return res
}

func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func doQuery(c *cli.Context, streaming bool) error {

	//
	// Checking arguments
	//

	var body string
	if c.String("body") != "" {
		body = c.String("body")
	} else if c.String("file") != "" {
		data, err := ioutil.ReadFile(c.String("file"))
		if err != nil {
			return err
		}
		body = string(data)
	} else if c.String("query") != "" {
		queries := builtInQueries()
		query := c.String("query")
		body = queries[query]
	} else {
		return cli.NewExitError("You should provide query or file argument", 1)
	}

	if streaming && c.String("source") == "" {
		return cli.NewExitError("Dataset for importing is not provided", 1)
	}

	//
	// Parsing Variables
	//

	variables := c.StringSlice("variable")
	queryVariables := make(map[string]interface{})
	for v := range variables {
		args := strings.SplitN(variables[v], "=", 2)
		if len(args) < 2 {
			return cli.NewExitError(fmt.Sprintf("Query variable mailformed: %s", variables[v]), 1)
		}
		queryVariables[args[0]] = args[1]
	}

	//
	// Create Client
	//

	serverOpt := c.String("server")
	var serverURL string
	if serverOpt == "production" || serverOpt == "prod" {
		serverURL = "https://api.statecrafthq.com/api"
	} else if serverOpt == "local" {
		serverURL = "http://localhost:9000/api"
	} else {
		serverURL = serverOpt
	}

	//
	// Performing Reques
	//

	if streaming {
		srcFileName := c.String("source")
		batchSize := c.Int("batch")
		faultTolerant := c.Bool("fault-tolerant")

		// Opening file
		file, e := os.Open(srcFileName)
		if e != nil {
			return e
		}
		defer file.Close()

		// Line number
		lines, e := lineCounter(file)
		if e != nil {
			return e
		}
		file.Seek(0, 0)

		//
		// Reading And Importing
		//

		bar := pb.StartNew(lines)
		pending := make([]interface{}, 0)

		rd := bufio.NewReader(file)
		linesRead := 0
		for {
			line, e := rd.ReadBytes('\n')
			if e != nil {
				if e == io.EOF {
					break
				}
				return e
			}

			var d interface{}
			e = json.Unmarshal(line, &d)
			if e != nil {
				return e
			}
			pending = append(pending, d)
			linesRead = linesRead + 1
			bar.Set(linesRead)
			if len(pending) >= batchSize {
				queryVariables["data"] = pending
				if faultTolerant {
					for {
						_, e := utils.GraqhQLRequest(serverURL, body, queryVariables)
						if e != nil {
							fmt.Println(e)
							time.Sleep(1000)
						} else {
							break
						}
					}
				} else {
					_, e := utils.GraqhQLRequest(serverURL, body, queryVariables)
					if e != nil {
						return e
					}
				}
				pending = make([]interface{}, 0)
			}
		}
		if len(pending) >= batchSize {
			queryVariables["data"] = pending
			if faultTolerant {
				for {
					_, e := utils.GraqhQLRequest(serverURL, body, queryVariables)
					if e != nil {
						fmt.Println(e)
						time.Sleep(1000)
					} else {
						break
					}
				}
			} else {
				_, e := utils.GraqhQLRequest(serverURL, body, queryVariables)
				if e != nil {
					return e
				}
			}
		}
		bar.FinishPrint("Importing completed")
	} else {
		// Non-streaming request
		r, e := utils.GraqhQLRequest(serverURL, body, queryVariables)
		if e != nil {
			return e
		}
		fmt.Printf("%s\n", r)
	}

	return nil
}

func CreateImportingCommands() []cli.Command {
	return []cli.Command{
		{
			Name:    "query",
			Aliases: []string{"q"},
			Usage:   "Query server",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "server, s",
					Value:  "prod",
					Usage:  "prod, local or direct URL to server",
					EnvVar: "STATECRAFT_SERVER",
				},
				cli.StringFlag{
					Name:  "body, b",
					Usage: "Body of query",
				},
				cli.StringFlag{
					Name:  "file, f",
					Usage: "Body of query from file",
				},
				cli.StringSliceFlag{
					Name:  "variable, v",
					Usage: "Variables to query key=value",
				},
			},
			Action: func(c *cli.Context) error {
				return doQuery(c, false)
			},
		},
		{
			Name:    "import",
			Aliases: []string{"i"},
			Usage:   "Import dataset to server",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "server, s",
					Value:  "prod",
					Usage:  "prod, local or direct URL to server",
					EnvVar: "STATECRAFT_SERVER",
				},
				cli.StringFlag{
					Name:  "source, src",
					Usage: "[REQUIRED] File for importing",
				},
				cli.StringFlag{
					Name:  "body, b",
					Usage: "Body of query",
				},
				cli.StringFlag{
					Name:  "file, f",
					Usage: "Body of query from file",
				},
				cli.StringFlag{
					Name:  "query",
					Usage: "Built-in query: blocks, parcels",
				},
				cli.StringSliceFlag{
					Name:  "variable, v",
					Usage: "Variables to query key=value",
				},
				cli.IntFlag{
					Name:  "batch",
					Value: 50,
					Usage: "Batch size",
				},
				cli.BoolFlag{
					Name:  "fault-tolerant",
					Usage: "Set this flag to repeat on errors",
				},
			},
			Action: func(c *cli.Context) error {
				return doQuery(c, true)
			},
		},
	}
}
