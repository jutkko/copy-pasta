package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jutkko/copy-pasta/runcommands"
	"github.com/jutkko/copy-pasta/store"
	minio "github.com/minio/minio-go"
)

func main() {
	var target *runcommands.Target
	var client *minio.Client
	var err error
	config := &runcommands.Config{}

	parseCommands(config)

	target = config.CurrentTarget
	client, err = minioClient(config.CurrentTarget)

	if err != nil {
		log.Fatalf(fmt.Sprintf("Failed initializing client: %s\n", err.Error()))
	}

	stat, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal(err)
	}
	// stdin is pipe
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		bucketName, objectName, location := s3BucketInfo(target)
		err = store.S3Write(client, bucketName, objectName, location, os.Stdin)
		if err != nil {
			log.Fatalf(fmt.Sprintf("Failed writing to the bucket: %s\n", err.Error()))
		}
	} else {
		// stdin is tty
		bucketName, objectName, _ := s3BucketInfo(target)

		content, err := store.S3Read(client, bucketName, objectName)
		if err != nil {
			log.Fatalf(fmt.Sprintf("Have you copied yet? Failed reading the bucket: %s\n", err.Error()))
		}
		fmt.Printf("%s", content)
	}
}

func minioClient(t *runcommands.Target) (*minio.Client, error) {
	var endpoint string
	if endpoint = os.Getenv("S3ENDPOINT"); endpoint == "" {
		endpoint = "s3.amazonaws.com"
	}
	accessKeyID := t.AccessKey
	secretAccessKey := t.SecretAccessKey
	useSSL := true

	// Initialize minio client object
	return minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
}

func s3BucketInfo(target *runcommands.Target) (string, string, string) {
	var bucketName, objectName, location string

	bucketName = target.BucketName
	if objectName = os.Getenv("S3OBJECTNAME"); objectName == "" {
		objectName = "default-object-name"
	}
	if location = os.Getenv("S3LOCATION"); location == "" {
		location = "eu-west-2"
	}

	return bucketName, objectName, location
}

func parseCommands(config *runcommands.Config) {
	if len(os.Args) == 1 {
		configTemp, err := runcommands.Load()
		if err != nil {
			fmt.Printf("Please log in\n")
			os.Exit(1)
		}
		*config = *configTemp

		return
	}

	loginCommand := flag.NewFlagSet("login", flag.ExitOnError)
	loginTargetOption := loginCommand.String("target", "", "the name for copy-pasta's target")

	switch os.Args[1] {
	case "login":
		loginCommand.Parse(os.Args[2:])
	case "target":
		if len(os.Args) > 2 {
			configTemp, err := runcommands.Load()
			if err != nil {
				fmt.Printf("Please log in\n")
				os.Exit(1)
			}
			*config = *configTemp

			if target, ok := config.Targets[os.Args[2]]; ok {
				err := runcommands.Update(target.Name, target.AccessKey, target.SecretAccessKey, target.BucketName)
				if err != nil {
					log.Fatalf(fmt.Sprintf("Failed to update the current target: %s\n", err.Error()))
				}
			} else {
				fmt.Printf("Target is invalid\n")
				os.Exit(3)
			}
			os.Exit(0)
		}

		fmt.Printf("No target provided\n")
		os.Exit(4)
	default:
		fmt.Printf("%s is not a valid command.\n", os.Args[1])
		os.Exit(2)
	}

	if loginCommand.Parsed() {
		var accessKey, secretAccessKey string

		reader := bufio.NewReader(os.Stdin)

		fmt.Printf("Please enter key:\n")
		accessKeyWithNewline, _ := reader.ReadString('\n')
		accessKey = strings.Trim(accessKeyWithNewline, "\n")

		fmt.Printf("Please enter secret key:\n")
		secretAccessKeyWithNewline, _ := reader.ReadString('\n')
		secretAccessKey = strings.Trim(secretAccessKeyWithNewline, "\n")

		err := runcommands.Update(*loginTargetOption, accessKey, secretAccessKey, getBucketName(accessKey+*loginTargetOption))
		if err != nil {
			log.Fatalf(fmt.Sprintf("Failed to update the current target: %s\n", err.Error()))
		}

		fmt.Printf("Log in information saved\n")
	}

	os.Exit(0)
}

func getBucketName(salt string) string {
	pastas := []string{
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

	suffix := md5.Sum([]byte(salt))
	pastaIndex := int(suffix[0]) % len(pastas)

	return fmt.Sprintf("%s-%s", pastas[pastaIndex], hex.EncodeToString(suffix[:]))
}
