package commands

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/statecrafthq/borg/commands/drivers"
	"github.com/statecrafthq/borg/commands/ops"
	"github.com/statecrafthq/borg/geometry"
	"github.com/statecrafthq/borg/utils"
	"github.com/urfave/cli"
	emoji "gopkg.in/kyokomi/emoji.v1"
)

func convertShapefile(c *cli.Context) error {
	src := c.String("src")
	dst := c.String("dst")
	simplify := c.Bool("simplify")
	if src == "" {
		return cli.NewExitError("Source file is not provided", 1)
	}
	if dst == "" {
		return cli.NewExitError("Destination file is not provided", 1)
	}

	//
	// Check if exists
	//

	exist := utils.FileExists(dst)
	if exist {
		if c.Bool("force") {
			e := os.Remove(dst)
			if e != nil {
				return e
			}
		} else {
			return cli.NewExitError("File already exists. Use --force for overwriting.", 1)
		}
	}

	//
	// Starting conversion
	//

	return utils.ShapefileToGeoJson(src, dst, simplify)
}

func converGeoJson(c *cli.Context) error {
	src := c.String("src")
	dst := c.String("dst")
	driverID := c.String("driver")
	strict := c.Bool("strict")
	noErrors := c.Bool("no-error-logging")
	fixAll := c.Bool("fix-all")
	if src == "" {
		return cli.NewExitError("Source file is not provided", 1)
	}
	if dst == "" {
		return cli.NewExitError("Destination file is not provided", 1)
	}
	if driverID == "" {
		return cli.NewExitError("driver is not provided", 1)
	}

	allDrivers := drivers.Drivers()
	if _, ok := allDrivers[strings.ToLower(driverID)]; !ok {
		return cli.NewExitError("Unable to find required driver", 1)
	}
	driver := allDrivers[strings.ToLower(driverID)]

	//
	// Existing file
	//

	exist := utils.FileExists(dst)
	if exist {
		if c.Bool("force") {
			e := os.Remove(dst)
			if e != nil {
				return e
			}
		} else {
			return cli.NewExitError("File already exists. Use --force for overwriting.", 1)
		}
	}

	//
	// Decoding Geometry
	//

	body, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	//
	// Generating of JSVC
	//

	file, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer file.Close()
	w := bufio.NewWriter(file)

	//
	// Iterating each feature
	//

	//
	// Detecting multiple features for same ID
	//

	emoji.Println(":hammer: Collecting stats about dataset")
	featureCounts := make(map[string]int32)
	err = utils.IterateFeatures(body, strict, !noErrors, func(feature *utils.Feature) error {

		// Record type
		recordType, err := driver.Record(feature)
		if err != nil {
			return err
		}

		// Ignored fields
		if recordType == drivers.Ignored {
			return nil
		}

		// ID
		idValue, err := driver.ID(feature)
		if err != nil {
			return err
		}

		// If multiple ID data source - handle just like separate ids
		if driver.MultipleID {
			for _, id := range idValue {
				if val, ok := featureCounts[id]; ok {
					featureCounts[id] = val + 1
				} else {
					featureCounts[id] = 1
				}
			}
		} else {
			if val, ok := featureCounts[idValue[0]]; ok {
				featureCounts[idValue[0]] = val + 1
			} else {
				featureCounts[idValue[0]] = 1
			}
		}

		return nil
	})

	//
	// Converting Features
	//
	emoji.Println(":hammer: Processing dataset")
	pendingFeatures := make(map[string][][][][]float64)
	pendingFeaturesCount := make(map[string]int32)
	err = utils.IterateFeatures(body, strict, !noErrors, func(feature *utils.Feature) error {

		// Record type
		recordType, err := driver.Record(feature)
		if err != nil {
			return err
		}

		// Ignored fields
		if recordType == drivers.Ignored {
			return nil
		}

		// Retired type
		retiredType, err := driver.Retired(feature)
		if err != nil {
			return err
		}

		// Parsing Coordinates
		// Ignore if geometry missing
		var coordinates [][][][]float64
		if feature.Geometry != nil {
			coordinates = *feature.Geometry
			// Fixing invalid polygons
			if fixAll {
				coordinates, err = utils.PolygonRepair(coordinates)
				if err != nil {
					return err
				}
			} else {
				err = utils.ValidateGeometry(coordinates)
				if err != nil {
					coordinates, err = utils.PolygonRepair(coordinates)
					if err != nil {
						return err
					}
					err = utils.ValidateGeometry(coordinates)
					if err != nil {
						return err
					}
				}
			}
		}

		// Loading ID
		idValue, err := driver.ID(feature)
		if err != nil {
			return err
		}

		// Iterate over all IDs
		primaryIDs := []string{idValue[0]}
		if driver.MultipleID {
			primaryIDs = idValue
		}
		for _, primaryID := range primaryIDs {
			// Check if present on counters
			var totlaCount int32
			if val, ok := featureCounts[primaryID]; ok {
				totlaCount = val
			} else {
				return errors.New("Internal inconsistency")
			}

			// Check how many features for this ID is already processed
			var currentCount int32
			if val, ok := pendingFeaturesCount[primaryID]; ok {
				currentCount = val + 1
			} else {
				currentCount = 1
			}
			pendingFeaturesCount[primaryID] = currentCount

			// Check if we are reached end for specific feature
			isLast := currentCount >= totlaCount

			//
			// Merging Geometry
			//

			currentCoordinates := make([][][][]float64, 0)
			if val, ok := pendingFeatures[primaryID]; ok {
				currentCoordinates = val
			}

			// Merge geometry only for primary records
			if recordType == drivers.Primary && feature.Geometry != nil {
				for _, poly := range coordinates {
					currentCoordinates = append(currentCoordinates, poly)
				}
			}

			// Update pending if is not last and delete from memory ASAP
			if !isLast {
				// Save pending geometry only for primary records
				if recordType == drivers.Primary {
					pendingFeatures[primaryID] = currentCoordinates
				}
				return nil
			}
			delete(pendingFeatures, primaryID)
			delete(pendingFeaturesCount, primaryID)

			// Loading Extras
			extras := ops.NewExtras()
			err = driver.Extras(feature, &extras)
			if err != nil {
				return err
			}

			// Preparing Bundle
			fields := make(map[string]interface{})
			fields["id"] = primaryID
			if len(idValue) > 1 {
				fields["displayId"] = idValue[1:]
			}
			if len(currentCoordinates) > 0 {
				fields["geometry"] = currentCoordinates
				extras.AppendFloat("area", geometry.NewGeoMultipolygon(currentCoordinates).Area())
			}
			if retiredType != drivers.Unkwnon {
				if retiredType == drivers.Retired {
					fields["retired"] = true
				} else {
					fields["retired"] = false
				}
			}
			fields["extras"] = extras

			// Writing
			marshaled, err := json.Marshal(fields)
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(w, string(marshaled))
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return w.Flush()
}

func CreateConvertingCommands() []cli.Command {
	return []cli.Command{
		{
			Name:  "convert",
			Usage: "Convert Datasets",
			Subcommands: []cli.Command{
				{
					Name:  "shapefile",
					Usage: "Converting Shapefile to GeoJSON",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "source, src",
							Usage: "path to source file",
						},
						cli.StringFlag{
							Name:  "dest, dst",
							Usage: "path to destination file",
						},
						cli.BoolFlag{
							Name:  "force, f",
							Usage: "Overwrite file if exists",
						},
						cli.BoolFlag{
							Name:  "simplify",
							Usage: "Simplify Geometry",
						},
					},
					Action: func(c *cli.Context) error {
						return convertShapefile(c)
					},
				},
				{
					Name:  "geojson",
					Usage: "Converting GeoJSON to ols file",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "source, src",
							Usage: "path to source file",
						},
						cli.StringFlag{
							Name:  "dest, dst",
							Usage: "path to destination file",
						},
						cli.StringFlag{
							Name:  "format,driver",
							Usage: "ny_blocks, ny_parcels, sf_blocks, sf_parcels",
						},
						cli.BoolFlag{
							Name:  "force, f",
							Usage: "Overwrite file if exists",
						},
						cli.BoolFlag{
							Name:  "strict",
							Usage: "Crash on invalid record",
						},
						cli.BoolFlag{
							Name:  "fix-all",
							Usage: "Fix all polygons",
						},
						cli.BoolFlag{
							Name:  "no-error-logging",
							Usage: "Disable error logging",
						},
					},
					Action: func(c *cli.Context) error {
						return converGeoJson(c)
					},
				},
			},
		},
	}
}
