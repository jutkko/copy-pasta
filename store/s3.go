package store

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

//go:generate counterfeiter . MinioClient

type MinioClient interface {
	MakeBucket(string, string) error
	BucketExists(string) (bool, error)
	PutObject(string, string, io.Reader, string) (int64, error)
	FGetObject(string, string, string) error
}

func S3Write(client MinioClient, bucketName, objectName, location string, content io.Reader) error {
	exists, err := client.BucketExists(bucketName)
	if err != nil {
		return err
	}

	if exists == false {
		err := client.MakeBucket(bucketName, location)
		if err != nil {
			return err
		}
	}

	_, err = client.PutObject(bucketName, objectName, content, "text/html")
	if err != nil {
		return err
	}

	return nil
}

func S3Read(client MinioClient, bucketName, objectName string) (string, error) {
	tempFile, err := ioutil.TempFile("/tmp", "tempS3ObjectFile")
	if err != nil {
		return "", err
	}

	defer tempFile.Close()
	defer func() {
		err = os.Remove(tempFile.Name())
		if err != nil {
			log.Fatal(err)
		}
	}()

	err = client.FGetObject(bucketName, objectName, tempFile.Name())
	if err != nil {
		return "", err
	}

	byteContent, err := ioutil.ReadFile(tempFile.Name())
	if err != nil {
		return "", err
	}

	return string(byteContent), nil
}
