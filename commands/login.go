package commands

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/jutkko/copy-pasta/runcommands"
	"github.com/mitchellh/cli"
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

type LoginCommand struct {
	Ui cli.Ui
}

func (l *LoginCommand) Help() string {
	return "Lists the targets on this machine"
}

func (l *LoginCommand) Run(args []string) int {
	loginCommand := flag.NewFlagSet("login", flag.ExitOnError)
	loginTargetOption := loginCommand.String("target", "", "the name for copy-pasta's target")

	// not tested, may be too hard
	err := loginCommand.Parse(args)
	if err != nil {
		l.Ui.Error(err.Error())
		return 10
	}

	reader := bufio.NewReader(os.Stdin)
	accessKey, _ := prompt("Please enter key: ", reader)
	secretAccessKey, _ := prompt("Please enter secret key: ", reader)

	if err := runcommands.Update(*loginTargetOption, accessKey, secretAccessKey, getBucketName(accessKey+*loginTargetOption)); err != nil {
		l.Ui.Error(fmt.Sprintf("Failed to update the current target: %s\n", err.Error()))
		return 9
	}

	fmt.Println("Log in information saved")
	return 0
}

func (l *LoginCommand) Synopsis() string {
	return "Lists the targets on this machine"
}

func getBucketName(salt string) string {
	suffix := md5.Sum([]byte(salt))
	pastaIndex := int(suffix[0]) % len(pastas)

	return fmt.Sprintf("%s-%s", pastas[pastaIndex], hex.EncodeToString(suffix[:]))
}

func prompt(message string, reader *bufio.Reader) (string, error) {
	fmt.Print(message)
	resultWithNewLine, err := reader.ReadString('\n')
	// TODO test this?
	if err != nil {
		return "", err
	}
	return strings.Trim(resultWithNewLine, "\n"), nil
}
