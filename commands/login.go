package commands

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"

	"github.com/jutkko/cli"
	"github.com/jutkko/copy-pasta/runcommands"
)

var pastas = []string{
	"acinidipepe",
	"agnolotti",
	"alphabetpasta",
	"anelli",
	"anellini",
	"bigoli",
	"bucatini",
	"calamarata",
	"campanelle",
	"cannelloni",
	"capellini",
	"casarecce",
	"casoncelli",
	"casunziei",
	"cavatappi",
	"cavatelli",
	"cencioni",
	"conchiglie",
	"corzetti",
	"croxetti",
	"ditalini",
	"fagottini",
	"farfalle",
	"fettuccine",
	"fiori",
	"fogliedulivo",
	"fregula",
	"fusi",
	"fusilli",
	"garganelli",
	"gemelli",
	"lanterne",
	"lasagne",
	"lasagnette",
	"linguettine",
	"linguine",
	"macaroni",
	"mafalde",
	"mafaldine",
	"mezzelune",
	"occhidilupo",
	"orecchiette",
	"orzo",
	"pappardelle",
	"passatelli",
	"pastina",
	"penne",
	"pici",
	"pillus",
	"pizzoccheri",
	"radiatori",
	"ravioli",
	"rigatoni",
	"rotelle",
	"rotini",
	"sacchettoni",
	"sagnarelli",
	"scialatelli",
	"spaghetti",
	"stringozzi",
	"strozzapreti",
	"tagliatelle",
	"taglierini",
	"testaroli",
	"tortellini",
	"tortelli",
	"tortelloni",
	"trenette",
	"tripoline",
	"troccoli",
	"trofie",
	"vermicelli",
}

// LoginCommand is the command that is responsible for logging in, the
// size effect is that it saves the config file locally
type LoginCommand struct {
	Ui cli.Ui
}

// Help string
func (l *LoginCommand) Help() string {
	return `Usage: copy-pasta login [--target] [<target>] [--endpoint] [<endpoint>] [--location] [<location>]

    Prompts to login interactively. If no target is provided,
    the  "default" target naem is provided.

Options:
    --target       Specify the new target name.
    --endpoint     Specify the new target's endpoint, defaults to s3.amazonaws.com.
    --location     Specify the new target's location, defaults to eu-west-2.
`
}

// Run function for the command
func (l *LoginCommand) Run(args []string) int {
	loginCommand := flag.NewFlagSet("login", flag.ExitOnError)
	loginTargetOption := loginCommand.String("target", "", "the name for copy-pasta's target")
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

	if err := runcommands.Update(*loginTargetOption, accessKey, secretAccessKey, getBucketName(accessKey+*loginTargetOption), *loginEndpointOption, *loginLocationOption); err != nil {
		l.Ui.Error(fmt.Sprintf("Failed to update the current target: %s\n", err.Error()))
		return 9
	}

	fmt.Println("Log in information saved")

	return 0
}

// Synopsis is the short help string
func (l *LoginCommand) Synopsis() string {
	return fmt.Sprintf("Login to copy-pasta")
}

func getBucketName(salt string) string {
	suffix := md5.Sum([]byte(salt))
	pastaIndex := int(suffix[0]) % len(pastas)

	return fmt.Sprintf("%s-%s", pastas[pastaIndex], hex.EncodeToString(suffix[:]))
}
