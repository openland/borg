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
	app.Version = "0.0.1"
	app.Usage = "Toolbelt to work with Statecraft API"

	//
	// Commands
	//

	app.Commands = append(append(append(commands.CreateImportingCommands(), commands.CreateConvertingCommands()...), commands.CreateSyncCommands()...), commands.CreateMergeCommands()...)

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
