package main

import (
	"bufio"
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
	if len(os.Args) > 1 {
		parseCommands()
	}

	config, err := runcommands.Load()
	if err != nil {
		fmt.Printf("Please log in")
		os.Exit(1)
	}

	stat, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal(err)
	}

	client, err := minioClient(config.CurrentTarget)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Failed initializing client: %s\n", err.Error()))
	}

	// stdin is pipe
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		bucketName, objectName, location := s3BucketInfo()
		err = store.S3Write(client, bucketName, objectName, location, os.Stdin)
		if err != nil {
			log.Fatalf(fmt.Sprintf("Failed writing to the bucket: %s\n", err.Error()))
		}
	} else {
		// stdin is tty
		bucketName, objectName, _ := s3BucketInfo()

		content, err := store.S3Read(client, bucketName, objectName)
		if err != nil {
			log.Fatalf(fmt.Sprintf("Failed reading the bucket: %s\n", err.Error()))
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

func s3BucketInfo() (string, string, string) {
	var bucketName, objectName, location string

	if bucketName = os.Getenv("S3BUCKETNAME"); bucketName == "" {
		bucketName = "default-bucket-name"
	}
	if objectName = os.Getenv("S3OBJECTNAME"); objectName == "" {
		objectName = "default-object-name"
	}
	if location = os.Getenv("S3LOCATION"); location == "" {
		location = "eu-west-2"
	}

	return bucketName, objectName, location
}

func parseCommands() {
	loginCommand := flag.NewFlagSet("login", flag.ExitOnError)
	loginTargetOption := loginCommand.String("target", "", "copy-pasta target name")

	switch os.Args[1] {
	case "login":
		loginCommand.Parse(os.Args[2:])
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

		fmt.Printf("Log in information saved")

		if *loginTargetOption == "" {
			runcommands.Update(accessKey, accessKey, secretAccessKey)
		} else {
			runcommands.Update(*loginTargetOption, accessKey, secretAccessKey)
		}
	}

	os.Exit(0)
}
