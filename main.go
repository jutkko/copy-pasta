package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/jutkko/copy-pasta/store"
	minio "github.com/minio/minio-go"
)

func main() {
	loginCommand := flag.NewFlagSet("login", flag.ExitOnError)

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "login":
			loginCommand.Parse(os.Args[2:])
		default:
			fmt.Printf("%s is not a valid command.\n", os.Args[1])
			os.Exit(2)
		}

		if loginCommand.Parsed() {
			fmt.Printf("Please enter key:\n")
			fmt.Printf("Please enter secret key:\n")
		}
	}

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	// var profile *runcommands.Rc
	if _, err := os.Stat(filepath.Join(usr.HomeDir, ".copy-pastarc")); os.IsNotExist(err) {
		// err := runcommands.Initialize()
		// if err != nil {
		// 	log.Fatal(err)
		// }
		fmt.Printf("Please log in")
		os.Exit(1)

		// reader := bufio.NewReader(os.Stdin)
		// fmt.Print("Enter text: ")
		// text, _ := reader.ReadString('\n')
		// fmt.Println(text)
	}

	stat, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal(err)
	}

	client, err := minioClient()
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
		println("Getting the last copied item...")
		bucketName, objectName, _ := s3BucketInfo()

		content, err := store.S3Read(client, bucketName, objectName)
		if err != nil {
			log.Fatalf(fmt.Sprintf("Failed writing to read bucket: %s\n", err.Error()))
		}
		fmt.Printf("%s", content)
	}
}

func minioClient() (*minio.Client, error) {
	endpoint := os.Getenv("S3ENDPOINT")
	accessKeyID := os.Getenv("S3ACCESSKEYID")
	secretAccessKey := os.Getenv("S3SECRETACCESSKEY")
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
