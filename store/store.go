package store

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/google/go-github/github"
	"github.com/jutkko/copy-pasta/runcommands"
	"github.com/jutkko/copy-pasta/store/gist"
	"github.com/jutkko/copy-pasta/store/s3"
	minio "github.com/minio/minio-go"
	"golang.org/x/oauth2"
)

type Store interface {
	Write(content io.Reader) error
	Read() (string, error)
}

// Only do s3 for now
func NewStore(target *runcommands.Target) (Store, error) {
	if target.Backend == "s3" {
		client, err := minioClient(target)
		if err != nil {
			return nil, fmt.Errorf("Failed initializing client: %s", err.Error())
		}
		return s3.NewS3Store(client, target), nil
	}

	if target.Backend == "gist" {
		return gist.NewGistStore(gistClient(target), target), nil
	}

	return nil, errors.New(fmt.Sprintf("Invalid backend: %s", target.Backend))
}

func minioClient(t *runcommands.Target) (*minio.Client, error) {
	endpoint := t.Endpoint
	accessKeyID := t.AccessKey
	secretAccessKey := t.SecretAccessKey
	useSSL := true

	// Initialize minio client object
	return minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
}

func gistClient(t *runcommands.Target) gist.GistClient {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: t.GistToken},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	client := github.NewClient(tc)
	return client.Gists
}
