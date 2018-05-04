package commands

import (
	"github.com/statecrafthq/borg/commands/ops"
	"github.com/statecrafthq/borg/utils"
	"github.com/urfave/cli"
)

func normalizeOls(c *cli.Context) error {
	src := c.String("src")
	dst := c.String("dst")
	if src == "" {
		return cli.NewExitError("You should provide latest file", 1)
	}
	if dst == "" {
		return cli.NewExitError("You should provide previous file", 1)
	}
	e := utils.AssumeNotExists(dst, c.Bool("force"))
	if e != nil {
		return e
	}
	_, e = ops.SortFile(src, dst)
	if e != nil {
		return e
	}
	return nil
}

func CreateNormalizeCommands() []cli.Command {
	return []cli.Command{
		{
			Name:  "normalize",
			Usage: "Normalize files",
			Subcommands: []cli.Command{
				{
					Name: "ols",
					Flags: []cli.Flag{
						cli.StringSliceFlag{
							Name:  "source, src",
							Usage: "Path to dataset",
						},
						cli.StringFlag{
							Name:  "dest, dst",
							Usage: "Path to destination file",
						},
						cli.BoolFlag{
							Name:  "force, f",
							Usage: "Overwrite file if exists",
						},
					},
					Action: func(c *cli.Context) error {
						return normalizeOls(c)
					},
				},
			},
		},
	}
}
