package commands

import (
	"github.com/statecrafthq/borg/geometry"

	"github.com/statecrafthq/borg/commands/ops"
	"github.com/statecrafthq/borg/utils"
	"github.com/urfave/cli"
	emoji "gopkg.in/kyokomi/emoji.v1"
)

func overlay(c *cli.Context) error {
	src := c.String("src")
	dst := c.String("dst")
	zoning := c.String("zoning")
	if src == "" {
		return cli.NewExitError("You should provide source file", 1)
	}
	if dst == "" {
		return cli.NewExitError("You should provide dest file", 1)
	}
	if zoning == "" {
		return cli.NewExitError("You should provide zoning file", 1)
	}
	e := utils.AssumeNotExists(dst, c.Bool("force"))
	if e != nil {
		return e
	}

	//
	// Loading zoning data
	//

	emoji.Println(":file_cabinet: Loading zoning data")
	minLat := 10000.0
	minLon := 10000.0
	maxLat := -10000.0
	maxLon := -10000.0
	zoningDataGeo := make(map[string]geometry.MultipolygonGeo)
	e = ops.RecordReader(zoning, func(row map[string]interface{}) error {
		if geom, ok := row["geometry"]; ok {
			coords := utils.ParseFloat4(geom.([]interface{}))
			g := geometry.NewGeoMultipolygon(coords)
			b := g.Bounds()
			if b.MaxLatitude > maxLat {
				maxLat = b.MaxLatitude
			}
			if b.MaxLongitude > maxLon {
				maxLon = b.MaxLongitude
			}
			if b.MinLongitude < minLon {
				minLon = b.MinLongitude
			}
			if b.MinLatitude < minLat {
				minLat = b.MinLatitude
			}

			mainID := row["id"].(string)
			if displayIds, ok := row["displayId"]; ok {
				d2 := displayIds.([]interface{})
				if len(d2) > 0 {
					for _, d := range d2 {
						if ex, ok := zoningDataGeo[d.(string)]; ok {
							zoningDataGeo[d.(string)] = ex.Merge(g)
						} else {
							zoningDataGeo[d.(string)] = g
						}
					}
				} else {
					if ex, ok := zoningDataGeo[mainID]; ok {
						zoningDataGeo[mainID] = ex.Merge(g)
					} else {
						zoningDataGeo[mainID] = g
					}
				}
			} else {
				if ex, ok := zoningDataGeo[mainID]; ok {
					zoningDataGeo[mainID] = ex.Merge(g)
				} else {
					zoningDataGeo[mainID] = g
				}
			}
		}
		return nil
	})
	if e != nil {
		return e
	}

	//
	// Prepare projection from center of zoning data
	//

	centerLatitude := (maxLat + minLat) / 2.0
	centerLongitude := (maxLon + minLon) / 2.0
	proj := geometry.NewProjection(geometry.PointGeo{Latitude: centerLatitude, Longitude: centerLongitude})

	//
	// Project zoning
	//

	zoningData := make(map[string]geometry.Multipolygon2D)
	for k, v := range zoningDataGeo {
		zoningData[k] = v.Project(proj)
	}

	//
	// Mapping zoning map
	//

	e = ops.RecordTransformer(src, dst, func(row map[string]interface{}) (map[string]interface{}, error) {
		// Reading extras
		extras, e := ops.LoadExtras(row["extras"])
		if e != nil {
			return row, e
		}

		// Searching for zoning codes
		zkeys := make([]string, 0)
		if g, ok := row["geometry"]; ok {
			multipoly := geometry.NewGeoMultipolygon(utils.ParseFloat4(g.([]interface{})))
			projected := multipoly.Project(proj)
			for k, v := range zoningData {
				if projected.Intersects(v) {
					zkeys = append(zkeys, k)
				}
			}
		}
		extras.AppendEnum("zoning", zkeys)
		row["extras"] = extras

		return row, nil
	})
	if e != nil {
		return e
	}

	return nil
}

func CreateZoningCommands() []cli.Command {
	return []cli.Command{
		{
			Name:  "zoning",
			Usage: "Apply zoning maps",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "source, src",
					Usage: "Path to dataset",
				},
				cli.StringFlag{
					Name:  "zoning",
					Usage: "Path to zoning file",
				},
				cli.StringFlag{
					Name:  "dest,dst",
					Usage: "Path to destination file",
				},
				cli.BoolFlag{
					Name:  "force, f",
					Usage: "Overwrite file if exists",
				},
			},
			Action: func(c *cli.Context) error {
				return overlay(c)
			},
		},
	}
}
