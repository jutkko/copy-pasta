package commands

import (
	"fmt"

	"github.com/jutkko/cli"
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

	t.Ui.Output("copy-pasta current target:")
	t.Ui.Output("  " + config.CurrentTarget.Name)

	t.Ui.Output("copy-pasta saved targets:")
	for _, target := range config.Targets {
		t.Ui.Output(fmt.Sprintf("  %s", target.Name))
	}

	return 0
}

func (t *TargetsCommand) Synopsis() string {
	return "Lists the targets on this machine"
}
