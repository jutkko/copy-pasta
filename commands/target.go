package commands

import (
	"fmt"

	"github.com/jutkko/cli"
	"github.com/jutkko/copy-pasta/runcommands"
)

type TargetCommand struct {
	Ui cli.Ui
}

func (t *TargetCommand) Help() string {
	return "Changes the current target to the provided target"
}

func (t *TargetCommand) Run(args []string) int {
	config, err := loadRunCommands()
	if err != nil {
		return 1
	}

	if len(args) > 0 {
		if target, ok := config.Targets[args[0]]; ok {
			if err := runcommands.Update(target.Name, target.AccessKey, target.SecretAccessKey, target.BucketName); err != nil {
				t.Ui.Error(fmt.Sprintf("Failed to update the current target: %s", err.Error()))
				return 2
			} else {
				return 0
			}
		} else {
			t.Ui.Error("Target is invalid")
			return 3
		}
	} else {
		t.Ui.Output("copy-pasta current target:")
		t.Ui.Output("  " + config.CurrentTarget.Name)
		return 0
	}
}

func (t *TargetCommand) Synopsis() string {
	return "Changes the current target to the provided target"
}
