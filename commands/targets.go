package commands

import (
	"fmt"

	"github.com/jutkko/cli"
)

// TargetsCommand is the command that lists the current targets
type TargetsCommand struct {
	Ui cli.Ui
}

// Help string
func (t *TargetsCommand) Help() string {
	return `Usage: copy-pasta targets

    List the current as well as all the saved targets.
`
}

// Run function for the command
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

// Synopsis is the short help string
func (t *TargetsCommand) Synopsis() string {
	return "List the current as well as the saved targets"
}
