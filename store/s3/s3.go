package s3

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/jutkko/copy-pasta/runcommands"
	minio "github.com/minio/minio-go"
)

//go:generate counterfeiter . MinioClient

// MinioClient is the s3 client
type MinioClient interface {
	MakeBucket(string, string) error
	BucketExists(string) (bool, error)
	PutObject(string, string, io.Reader, string) (int64, error)
	FGetObject(string, string, string) error
}

type Store interface {
	Write(*runcommands.Target) error
	Read(*runcommands.Target) error
}

type S3Store struct {
	MinioClient MinioClient
}

func NewS3Store(target *runcommands.Target) (*S3Store, error) {
	client, err := minioClient(target)
	if err != nil {
		return nil, fmt.Errorf("failed initializing client: %s", err.Error())
	}

	return &S3Store{MinioClient: client}, nil
}

// S3Write is the function that writes to s3
func (s *S3Store) Write(bucketName, objectName, location string, content io.Reader) error {
	exists, err := s.MinioClient.BucketExists(bucketName)
	if err != nil {
		return err
	}

	if exists == false {
		err := s.MinioClient.MakeBucket(bucketName, location)
		if err != nil {
			return err
		}
	}

	_, err = s.MinioClient.PutObject(bucketName, objectName, content, "text/html")
	if err != nil {
		return err
	}

	return nil
}

// S3Read is the function that reads from s3
func (s *S3Store) Read(bucketName, objectName string) (string, error) {
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

	err = s.MinioClient.FGetObject(bucketName, objectName, tempFile.Name())
	if err != nil {
		return "", err
	}

	byteContent, err := ioutil.ReadFile(tempFile.Name())
	if err != nil {
		return "", err
	}

	return string(byteContent), nil
}

func minioClient(t *runcommands.Target) (*minio.Client, error) {
	endpoint := t.Endpoint
	accessKeyID := t.AccessKey
	secretAccessKey := t.SecretAccessKey
	useSSL := true

	// Initialize minio client object
	return minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
}
