package utils

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/urfave/cli"
)

func ShapefileToGeoJson(src string, dst string, simplify bool) error {
	exist := FileExists(dst)
	if exist {
		return cli.NewExitError("File already exists", 1)
	}
	_, e := exec.LookPath("ogr2ogr")
	if e != nil {
		return e
	}
	var command *exec.Cmd
	if simplify {
		command = exec.Command("ogr2ogr", "-f", "GeoJSON", "-t_srs", "crs:84", "-simplify", "0.00001", dst, src)
	} else {
		command = exec.Command("ogr2ogr", "-f", "GeoJSON", "-t_srs", "crs:84", dst, src)
	}
	var out bytes.Buffer
	command.Stdout = &out
	err := command.Run()
	if err != nil {
		fmt.Printf("Execution failed: \n%s\n", out.String())
	}
	return err
}

func GeoJsonToShapefile(src string, dst string, simplify bool) error {
	exist := FileExists(dst)
	if exist {
		return cli.NewExitError("File already exists", 1)
	}
	_, e := exec.LookPath("ogr2ogr")
	if e != nil {
		return e
	}
	var command *exec.Cmd
	if simplify {
		command = exec.Command("ogr2ogr", "-f", "ESRI Shapefile", "-simplify", "0.00001", dst, src)
	} else {
		command = exec.Command("ogr2ogr", "-f", "ESRI Shapefile", dst, src)
	}

	var out bytes.Buffer
	command.Stdout = &out
	err := command.Run()
	if err != nil {
		fmt.Printf("Execution failed: \n%s\n", out.String())
	}
	return err
}
