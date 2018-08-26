package gist

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/google/go-github/github"
	"github.com/jutkko/copy-pasta/runcommands"
)

const filename = "copy-pasta"

//go:generate counterfeiter . GistClient
type GistClient interface {
	Get(context.Context, string) (*github.Gist, *github.Response, error)
	Create(context.Context, *github.Gist) (*github.Gist, *github.Response, error)
	Edit(context.Context, string, *github.Gist) (*github.Gist, *github.Response, error)
}

type GistStore struct {
	gistClient GistClient
	target     *runcommands.Target
}

func NewGistStore(client GistClient, target *runcommands.Target) *GistStore {
	return &GistStore{
		gistClient: client,
		target:     target,
	}
}

func (g *GistStore) Write(content io.Reader) error {
	b, err := ioutil.ReadAll(content)
	if err != nil {
		return errors.New("Failed to read from the content")
	}
	if len(b) > 10*1024 {
		return errors.New("Size too large")
	}
	if len(b) == 0 {
		return nil
	}

	contentText := string(b)
	gistFilename := github.GistFilename(filename)
	localFilename := filename
	public := false
	gist := &github.Gist{
		Files: map[github.GistFilename]github.GistFile{
			gistFilename: github.GistFile{
				Content:  &contentText,
				Filename: &localFilename,
			},
		},
		Public: &public,
	}

	if g.target.GistID != "" {
		_, _, err = g.gistClient.Edit(context.Background(), g.target.GistID, gist)
		if err == nil {
			return nil
		}

		// Otherwise continue to create a new one
	}

	newGist, _, err := g.gistClient.Create(context.Background(), gist)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to create the gist: %s", err.Error()))
	}

	// Update the config for command
	err = runcommands.Update(g.target.Name, g.target.Backend, g.target.AccessKey, g.target.SecretAccessKey, g.target.BucketName, g.target.Endpoint, g.target.Location, g.target.GistToken, *newGist.ID)
	if err != nil {
		return err
	}

	return nil
}

func (g *GistStore) Read() (string, error) {
	ctx := context.Background()
	gist, _, err := g.gistClient.Get(ctx, g.target.GistID)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to get gist %s", err.Error()))
	}

	rawURL := *gist.Files[filename].RawURL
	resp, err := http.Get(rawURL)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to get content: %s", err.Error()))
	}

	content, _ := ioutil.ReadAll(resp.Body)
	return string(content), nil
}
