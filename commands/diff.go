package commands

import "github.com/urfave/cli"

func diff(c *cli.Context) error {
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
