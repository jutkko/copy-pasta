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
	return "Use echo $something | copy-pasta to copy and copy-pasta to paste"
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
	return "Use echo $something | copy-pasta to copy and copy-pasta to paste"
}

func copyPaste(target *runcommands.Target) error {
	client, err := minioClient(target)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed initializing client: %s\n", err.Error()))
	}

	bucketName, objectName, location := s3BucketInfo(target)
	if isFromAPipe() {
		if err = store.S3Write(client, bucketName, objectName, location, os.Stdin); err != nil {
			return errors.New(fmt.Sprintf("Failed writing to the bucket: %s\n", err.Error()))
		}
	} else {
		content, err := store.S3Read(client, bucketName, objectName)
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
	var endpoint string
	if endpoint = os.Getenv("S3ENDPOINT"); endpoint == "" {
		endpoint = "s3.amazonaws.com"
	}
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

func s3BucketInfo(target *runcommands.Target) (string, string, string) {
	return target.BucketName,
		getOrElse("S3OBJECTNAME", "default-object-name"),
		getOrElse("S3LOCATION", "eu-west-2")
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
