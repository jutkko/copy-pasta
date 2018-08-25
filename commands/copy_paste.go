package commands

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jutkko/cli"
	"github.com/jutkko/copy-pasta/runcommands"
	"github.com/jutkko/copy-pasta/store"
)

// InvalidConfig is the custom error struct for invalid configuration files
type InvalidConfig struct {
	error  string
	status int
}

func (ic *InvalidConfig) Error() string {
	return ic.error
}

// CopyPasteCommand is the command that is responsible for the actual copying
// and pasting
type CopyPasteCommand struct {
	Ui cli.Ui
}

// Help string
func (c *CopyPasteCommand) Help() string {
	return `Usage to paste: copy-pasta [--paste]
Usage to copy: <some command with output> | copy-pasta

    Copy or paste using copy-pasta. Use --paste to force copy-pasta to
		ignore its stdin and output from the current target.
`
}

// Run function for the command
func (c *CopyPasteCommand) Run(args []string) int {
	config, invalidConfig := loadRunCommands()
	if invalidConfig != nil {
		c.Ui.Error(fmt.Sprintf("Failed to load the runcommands file: %s", invalidConfig.Error()))
		os.Exit(invalidConfig.status)
	}

	copyPasteCommand := flag.NewFlagSet("", flag.ExitOnError)
	copyPastePasteOption := copyPasteCommand.Bool("paste", false, "")

	// not tested, may be too hard
	err := copyPasteCommand.Parse(args)
	if err != nil {
		return 10
	}

	if config != nil {
		content, err := copyPaste(config.CurrentTarget, *copyPastePasteOption)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Failed to load the runcommands file: %s", err.Error()))
			os.Exit(-15)
		}

		// cannot use c.Ui since it prints a newline
		fmt.Print(content)
	}

	return 0
}

// Synopsis is the short help string
func (c *CopyPasteCommand) Synopsis() string {
	return "Copy or paste using copy-pasta"
}

// copyPaste function deals with both copying and pasting
func copyPaste(target *runcommands.Target, paste bool) (string, error) {
	store, _ := store.NewS3Store(target)

	if isFromAPipe() && !paste {
		if err := store.Write(target.BucketName, "default-object-name", target.Location, os.Stdin); err != nil {
			return "", fmt.Errorf("failed writing to the bucket: %s", err.Error())
		}

		return "", nil
	} else {
		content, err := store.Read(target.BucketName, "default-object-name")
		if err != nil {
			return "", fmt.Errorf("Have you copied yet? Failed reading the bucket: %s", err.Error())
		}

		return content, nil
	}
}

func isFromAPipe() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal(err)
	}

	return (stat.Mode() & os.ModeCharDevice) == 0
}

func loadRunCommands() (*runcommands.Config, *InvalidConfig) {
	loadedConfig, err := runcommands.Load()
	if err != nil {
		return nil, &InvalidConfig{
			error:  "Please log in",
			status: 1,
		}
	}

	return loadedConfig, nil
}
