package commands

import (
	"flag"
	"fmt"

	"github.com/jutkko/cli"
	"github.com/jutkko/copy-pasta/runcommands"
)

// GistLoginCommand is the command that is responsible for logging in, the
// size effect is that it saves the config file locally
type GistLoginCommand struct {
	Ui cli.Ui
}

// Help string
func (l *GistLoginCommand) Help() string {
	return `Usage: copy-pasta gist-login [--target] [<target>]

		Prompts to login interactively. The command expects github gist token. If no
		target is provided, the "default" target name is provided.

Options:
    --target       Specify the new target name.
`
}

// Run function for the command
func (l *GistLoginCommand) Run(args []string) int {
	loginCommand := flag.NewFlagSet("login", flag.ExitOnError)
	loginTargetOption := loginCommand.String("target", "default", "the name for copy-pasta's target")

	// not tested, may be too hard
	err := loginCommand.Parse(args)
	if err != nil {
		l.Ui.Error(err.Error())
		return 10
	}

	token, err := l.Ui.Ask("Please enter github auth token:")
	if err != nil {
		l.Ui.Error(err.Error())
		return 10
	}

	if err := runcommands.Update(*loginTargetOption, "gist", "", "", "", "", "", token, ""); err != nil {
		l.Ui.Error(fmt.Sprintf("Failed to update the current target: %s\n", err.Error()))
		return 9
	}

	fmt.Println("Log in information saved")

	return 0
}

// Synopsis is the short help string
func (l *GistLoginCommand) Synopsis() string {
	return fmt.Sprintf("Login to copy-pasta with github gist backend")
}
