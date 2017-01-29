package store

import (
	"bytes"
	"io"
)

//go:generate counterfeiter . MinioClient

type MinioClient interface {
	MakeBucket(string, string) error
	BucketExists(string) (bool, error)
	PutObject(string, string, io.Reader, string) (int64, error)
}

func S3Write(client MinioClient, bucketName, location string, content []string) error {
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

	_, err = client.PutObject(bucketName, "zhou-test-object-real-shit", bytes.NewReader([]byte(content[0])), "text/html")
	if err != nil {
		return err
	}
	return nil
}

func S3Read() ([]string, error) {
	return nil, nil
}
