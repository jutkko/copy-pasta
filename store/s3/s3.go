package s3

import (
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/jutkko/copy-pasta/runcommands"
)

const DefaultObjectName = "default-object-name"

//go:generate counterfeiter . MinioClient

// MinioClient is the s3 client
type MinioClient interface {
	MakeBucket(string, string) error
	BucketExists(string) (bool, error)
	PutObject(string, string, io.Reader, string) (int64, error)
	FGetObject(string, string, string) error
}

type S3Store struct {
	MinioClient MinioClient
	target      *runcommands.Target
}

func NewS3Store(client MinioClient, target *runcommands.Target) *S3Store {
	return &S3Store{
		MinioClient: client,
		target:      target,
	}
}

// S3Write is the function that writes to s3
func (s *S3Store) Write(content io.Reader) error {
	exists, err := s.MinioClient.BucketExists(s.target.BucketName)
	if err != nil {
		return err
	}

	if exists == false {
		err := s.MinioClient.MakeBucket(s.target.BucketName, s.target.Location)
		if err != nil {
			return err
		}
	}

	_, err = s.MinioClient.PutObject(s.target.BucketName, DefaultObjectName, content, "text/html")
	if err != nil {
		return err
	}

	return nil
}

// S3Read is the function that reads from s3
func (s *S3Store) Read() (string, error) {
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

	err = s.MinioClient.FGetObject(s.target.BucketName, DefaultObjectName, tempFile.Name())
	if err != nil {
		return "", err
	}

	byteContent, err := ioutil.ReadFile(tempFile.Name())
	if err != nil {
		return "", err
	}

	return string(byteContent), nil
}
