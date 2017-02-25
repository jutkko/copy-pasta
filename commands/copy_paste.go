package commands

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/jutkko/copy-pasta/runcommands"
	"github.com/jutkko/copy-pasta/store"
	minio "github.com/minio/minio-go"
)

type InvalidConfig struct {
	error  string
	status int
}

func (ic *InvalidConfig) Error() string {
	return ic.error
}

type CopyPasteCommand struct{}

func (c *CopyPasteCommand) Help() string {
	return `Usage to paste: copy-pasta
Usage to copy: <some command> | copy-pasta

    Copy or paste using copy-pasta.
`
}

func (c *CopyPasteCommand) Run(args []string) int {
	config, invalidConfig := loadRunCommands()

	if invalidConfig != nil {
		fmt.Println(invalidConfig)
		os.Exit(invalidConfig.status)
	}

	if config != nil {
		if err := copyPaste(config.CurrentTarget); err != nil {
			log.Fatal(err)
		}
	}

	return 0
}

func (c *CopyPasteCommand) Synopsis() string {
	return "Copy or paste using copy-pasta"
}

func copyPaste(target *runcommands.Target) error {
	client, err := minioClient(target)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed initializing client: %s\n", err.Error()))
	}

	if isFromAPipe() {
		if err = store.S3Write(client, target.BucketName, "default-object-name", target.Location, os.Stdin); err != nil {
			return errors.New(fmt.Sprintf("Failed writing to the bucket: %s\n", err.Error()))
		}
	} else {
		content, err := store.S3Read(client, target.BucketName, "default-object-name")
		if err != nil {
			return errors.New(fmt.Sprintf("Have you copied yet? Failed reading the bucket: %s\n", err.Error()))
		}
		fmt.Print(content)
	}

	return nil
}

func isFromAPipe() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal(err)
	}

	return (stat.Mode() & os.ModeCharDevice) == 0
}

func minioClient(t *runcommands.Target) (*minio.Client, error) {
	endpoint := t.Endpoint
	accessKeyID := t.AccessKey
	secretAccessKey := t.SecretAccessKey
	useSSL := true

	// Initialize minio client object
	return minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
}

func getOrElse(key, defaultValue string) string {
	result := os.Getenv(key)
	if result == "" {
		return defaultValue
	}

	return result
}

func loadRunCommands() (*runcommands.Config, *InvalidConfig) {
	loadedConfig, err := runcommands.Load()
	if err != nil {
		return nil, &InvalidConfig{
			error:  "Please log in",
			status: 1,
		}
	}

	return loadedConfig, nil
}
