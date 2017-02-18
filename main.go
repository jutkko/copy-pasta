package main

import (
	"log"
	"os"

	"github.com/jutkko/copy-pasta/commands"
	"github.com/mitchellh/cli"
)

func main() {
	ui := &cli.BasicUi{Writer: os.Stdout}

	c := cli.NewCLI("copy-pasta", "0.0.1")
	// no copy-pasta is passed down
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"": func() (cli.Command, error) {
			return &commands.CopyPasteCommand{}, nil
		},

		"login": func() (cli.Command, error) {
			return &commands.LoginCommand{Ui: ui}, nil
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
