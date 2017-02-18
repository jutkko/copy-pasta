package commands

import (
	"fmt"

	"github.com/mitchellh/cli"
)

type TargetsCommand struct {
	Ui cli.Ui
}

func (t *TargetsCommand) Help() string {
	return "Lists the targets on this machine"
}

func (t *TargetsCommand) Run(args []string) int {
	config, err := loadRunCommands()
	if err != nil {
		t.Ui.Error(err.Error())
		return 6
	}

	fmt.Println("copy-pasta targets:")

	for _, target := range config.Targets {
		fmt.Printf("  %s\n", target.Name)
	}

	return 0
}

func (t *TargetsCommand) Synopsis() string {
	return "Lists the targets on this machine"
}
