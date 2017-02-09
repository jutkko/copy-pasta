package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jutkko/copy-pasta/runcommands"
	"github.com/jutkko/copy-pasta/store"
	minio "github.com/minio/minio-go"
)

type InvalidConfig struct {
	error  string
	status int
}

func (ic *InvalidConfig) Error() string {
	return ic.error
}

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

func main() {
	config, invalidConfig := parseCommands()
	if invalidConfig != nil {
		fmt.Println(invalidConfig)
		os.Exit(invalidConfig.status)
	}

	if config != nil {
		if err := copyPaste(config.CurrentTarget); err != nil {
			log.Fatal(err)
		}
	}
}

func copyPaste(target *runcommands.Target) error {
	client, err := minioClient(target)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed initializing client: %s\n", err.Error()))
	}

	bucketName, objectName, location := s3BucketInfo(target)
	if isFromAPipe() {
		if err = store.S3Write(client, bucketName, objectName, location, os.Stdin); err != nil {
			return errors.New(fmt.Sprintf("Failed writing to the bucket: %s\n", err.Error()))
		}
	} else {
		content, err := store.S3Read(client, bucketName, objectName)
		if err != nil {
			return errors.New(fmt.Sprintf("Have you copied yet? Failed reading the bucket: %s\n", err.Error()))
		}
		fmt.Print(content)
	}
	return nil
}

func isFromAPipe() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal(err)
	}

	return (stat.Mode() & os.ModeCharDevice) == 0
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

func getOrElse(key, defaultValue string) string {
	result := os.Getenv(key)
	if result == "" {
		return defaultValue
	}
	return result
}

func s3BucketInfo(target *runcommands.Target) (string, string, string) {
	return target.BucketName,
		getOrElse("S3OBJECTNAME", "default-object-name"),
		getOrElse("S3LOCATION", "eu-west-2")
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

func prompt(message string, reader *bufio.Reader) (string, error) {
	fmt.Print(message)
	resultWithNewLine, err := reader.ReadString('\n')
	// TODO test this?
	if err != nil {
		return "", err
	}
	return strings.Trim(resultWithNewLine, "\n"), nil

}
func parseCommands() (*runcommands.Config, *InvalidConfig) {
	if len(os.Args) == 1 {
		return loadRunCommands()
	}

	loginCommand := flag.NewFlagSet("login", flag.ExitOnError)
	loginTargetOption := loginCommand.String("target", "", "the name for copy-pasta's target")

	switch os.Args[1] {
	case "login":
		loginCommand.Parse(os.Args[2:])

		reader := bufio.NewReader(os.Stdin)
		accessKey, _ := prompt("Please enter key: ", reader)
		secretAccessKey, _ := prompt("Please enter secret key: ", reader)

		if err := runcommands.Update(*loginTargetOption, accessKey, secretAccessKey, getBucketName(accessKey+*loginTargetOption)); err != nil {
			return nil, &InvalidConfig{
				error:  fmt.Sprintf("Failed to update the current target: %s\n", err.Error()),
				status: 1,
			}
		}

		fmt.Println("Log in information saved")
	case "target":
		if len(os.Args) > 2 {
			config, err := loadRunCommands()
			if err != nil {
				return nil, err
			}

			if target, ok := config.Targets[os.Args[2]]; ok {
				if err := runcommands.Update(target.Name, target.AccessKey, target.SecretAccessKey, target.BucketName); err != nil {
					return nil, &InvalidConfig{
						error:  fmt.Sprintf("Failed to update the current target: %s", err.Error()),
						status: 1,
					}
				}
			} else {
				return nil, &InvalidConfig{
					error:  "Target is invalid",
					status: 3,
				}
			}
		} else {
			return nil, &InvalidConfig{
				error:  "No target provided",
				status: 4,
			}
		}
	case "targets":
		config, err := loadRunCommands()
		if err != nil {
			return nil, err
		}

		fmt.Println("copy-pasta targets:")
		for _, target := range config.Targets {
			fmt.Printf("  %s\n", target.Name)
		}
	default:
		return nil, &InvalidConfig{
			error:  fmt.Sprintf("%s is not a valid command.\n", os.Args[1]),
			status: 2,
		}
	}

	return nil, nil
}

func getBucketName(salt string) string {
	suffix := md5.Sum([]byte(salt))
	pastaIndex := int(suffix[0]) % len(pastas)

	return fmt.Sprintf("%s-%s", pastas[pastaIndex], hex.EncodeToString(suffix[:]))
}
