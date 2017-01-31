package main

import (
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
	bucketName := os.Getenv("S3BUCKETNAME")
	objectName := os.Getenv("S3OBJECTNAME")
	location := os.Getenv("S3LOCATION")
	return bucketName, objectName, location
}
