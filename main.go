package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/jutkko/copy-pasta/store"
	minio "github.com/minio/minio-go"
)

func main() {
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
		var input []string
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			input = append(input, scanner.Text())
		}

		bucketName, location := s3BucketInfo()
		err = store.S3Write(client, bucketName, location, input)
		if err != nil {
			log.Fatalf(fmt.Sprintf("Failed writing to the bucket: %s\n", err.Error()))
		}
	} else {
		// stdin is tty
		println("Getting the last copied item...")
		bucketName, _ := s3BucketInfo()

		content, err := store.S3Read(client, bucketName, "zhou-test-object-real-shit")
		if err != nil {
			log.Fatalf(fmt.Sprintf("Failed writing to read bucket: %s\n", err.Error()))
		}
		fmt.Printf("%#+v\n", content)
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

func s3BucketInfo() (string, string) {
	bucketName := os.Getenv("S3BUCKETNAME")
	location := os.Getenv("S3LOCATION")
	return bucketName, location
}
