package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/kyokomi/emoji.v1"

	"github.com/statecrafthq/borg/commands"
	"github.com/urfave/cli"
)

func main() {

	//
	// Basic Info
	//

	app := cli.NewApp()
	app.Name = "borg toolbelt"
	app.Version = "0.0.2"
	app.Usage = "Toolbelt to work with Statecraft API"

	//
	// Commands
	//

	app.Commands = []cli.Command{}
	app.Commands = append(app.Commands, commands.CreateImportingCommands()...)
	app.Commands = append(app.Commands, commands.CreateConvertingCommands()...)
	app.Commands = append(app.Commands, commands.CreateSyncCommands()...)
	app.Commands = append(app.Commands, commands.CreateMergeCommands()...)
	app.Commands = append(app.Commands, commands.CreateDiffCommands()...)
	app.Commands = append(app.Commands, commands.CreateFinalizeCommands()...)
	app.Commands = append(app.Commands, commands.CreateAnalyzeCommands()...)
	app.Commands = append(app.Commands, commands.CreateCursorCommands()...)
	app.Commands = append(app.Commands, commands.CreateNormalizeCommands()...)
	app.Commands = append(app.Commands, commands.CreateZoningCommands()...)
	app.Commands = append(app.Commands, commands.CreateExportCommands()...)

	//
	// Starting
	//

	startTime := time.Now()
	err := app.Run(os.Args)
	endTime := time.Now()
	if err != nil {
		fmt.Println(emoji.Sprintf(":warning: Failed in %d s", endTime.Sub(startTime)/time.Second))
		log.Fatal(err)
	} else {
		fmt.Println(emoji.Sprintf(":beer: Completed in %d s", endTime.Sub(startTime)/time.Second))
	}
}
