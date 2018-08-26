package commands

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
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

func getBucketName(salt string) string {
	suffix := md5.Sum([]byte(salt))
	pastaIndex := int(suffix[0]) % len(pastas)

	return fmt.Sprintf("%s-%s", pastas[pastaIndex], hex.EncodeToString(suffix[:]))
}
