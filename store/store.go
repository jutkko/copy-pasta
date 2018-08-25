package store

import (
	"fmt"
	"io"

	"github.com/jutkko/copy-pasta/runcommands"
	"github.com/jutkko/copy-pasta/store/s3"
	minio "github.com/minio/minio-go"
)

type Store interface {
	Write(content io.Reader) error
	Read() (string, error)
}

// Only do s3 for now
func NewStore(target *runcommands.Target) (Store, error) {
	client, err := minioClient(target)
	if err != nil {
		return nil, fmt.Errorf("Failed initializing client: %s", err.Error())
	}

	return s3.NewS3Store(client, target), nil
}

func minioClient(t *runcommands.Target) (*minio.Client, error) {
	endpoint := t.Endpoint
	accessKeyID := t.AccessKey
	secretAccessKey := t.SecretAccessKey
	useSSL := true

	// Initialize minio client object
	return minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
}
