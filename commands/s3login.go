package commands

import (
	"flag"
	"fmt"

	"github.com/jutkko/cli"
	"github.com/jutkko/copy-pasta/runcommands"
)

// S3LoginCommand is the command that is responsible for logging in, the
// size effect is that it saves the config file locally
type S3LoginCommand struct {
	Ui cli.Ui
}

// Help string
func (l *S3LoginCommand) Help() string {
	return `Usage: copy-pasta s3-login [--target] [<target>] [--endpoint] [<endpoint>] [--location] [<location>]

		Prompts to login interactively. The command expects S3 credentials. If no
		target is provided, the "default" target name is provided.

Options:
    --target       Specify the new target name.
    --endpoint     Specify the new target's endpoint, defaults to s3.amazonaws.com.
    --location     Specify the new target's location, defaults to eu-west-2.
`
}

// Run function for the command
func (l *S3LoginCommand) Run(args []string) int {
	loginCommand := flag.NewFlagSet("login", flag.ExitOnError)
	loginTargetOption := loginCommand.String("target", "default", "the name for copy-pasta's target")
	loginEndpointOption := loginCommand.String("endpoint", "s3.amazonaws.com", "the endpoint for copy-pasta's backend")
	loginLocationOption := loginCommand.String("location", "eu-west-2", "the location for the backend bucket")

	// not tested, may be too hard
	err := loginCommand.Parse(args)
	if err != nil {
		l.Ui.Error(err.Error())
		return 10
	}

	accessKey, err := l.Ui.Ask("Please enter key:")
	if err != nil {
		l.Ui.Error(err.Error())
		return 10
	}

	secretAccessKey, err := l.Ui.AskSecret("Please enter secret key:")
	if err != nil {
		l.Ui.Error(err.Error())
		return 10
	}

	if err := runcommands.Update(*loginTargetOption, "s3", accessKey, secretAccessKey, getBucketName(accessKey+*loginTargetOption), *loginEndpointOption, *loginLocationOption, "", ""); err != nil {
		l.Ui.Error(fmt.Sprintf("Failed to update the current target: %s\n", err.Error()))
		return 9
	}

	fmt.Println("Log in information saved")

	return 0
}

// Synopsis is the short help string
func (l *S3LoginCommand) Synopsis() string {
	return fmt.Sprintf("Login to copy-pasta with S3 backend")
}
