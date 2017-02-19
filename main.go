package main

import (
	"log"
	"os"

	"github.com/jutkko/cli"
	"github.com/jutkko/copy-pasta/commands"
)

func main() {
	ui := &cli.BasicUi{
		Writer:      os.Stdout,
		Reader:      os.Stdin,
		ErrorWriter: os.Stdout,
	}

	uiColored := &cli.ColoredUi{
		OutputColor: cli.UiColorNone,
		InfoColor:   cli.UiColorNone,
		ErrorColor:  cli.UiColorRed,
		Ui:          ui,
	}

	c := cli.NewCLI("copy-pasta", "0.0.1")

	// no copy-pasta is passed down
	c.Args = os.Args[1:]

	c.Commands = map[string]cli.CommandFactory{
		"": func() (cli.Command, error) {
			return &commands.CopyPasteCommand{}, nil
		},

		"login": func() (cli.Command, error) {
			return &commands.LoginCommand{Ui: uiColored}, nil
		},

		"target": func() (cli.Command, error) {
			return &commands.TargetCommand{Ui: ui}, nil
		},

		"targets": func() (cli.Command, error) {
			return &commands.TargetsCommand{Ui: ui}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
