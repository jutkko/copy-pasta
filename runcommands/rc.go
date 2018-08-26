package runcommands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

// Target is the strut for a copy-pasta target
type Target struct {
	Name            string `yaml:"name"`
	Backend         string `yaml:"backend"`
	AccessKey       string `yaml:"accesskey"`
	SecretAccessKey string `yaml:"secretaccesskey"`
	BucketName      string `yaml:"bucketname"`
	Endpoint        string `yaml:"endpoint"`
	Location        string `yaml:"location"`
	GistToken       string `yaml:"gisttoken"`
	GistID          string `yaml:"gistID"`
}

// Config is the aggregation of currrent targets
type Config struct {
	CurrentTarget *Target            `yaml:"currenttarget"`
	Targets       map[string]*Target `yaml:"targets"`
}

// Update updates the config file
func Update(target, backend, accessKey, secretAccessKey, bucketName, endpoint, location, gistToken, gistID string) error {
	var config *Config
	var err error

	config, err = Load()
	if err != nil {
		config = &Config{
			CurrentTarget: &Target{},
			Targets:       make(map[string]*Target),
		}
	}

	currentTarget := &Target{
		Name:            target,
		Backend:         backend,
		AccessKey:       accessKey,
		SecretAccessKey: secretAccessKey,
		BucketName:      bucketName,
		Endpoint:        endpoint,
		Location:        location,
		GistToken:       gistToken,
		GistID:          gistID,
	}

	config.CurrentTarget = currentTarget
	config.Targets[target] = currentTarget

	configContents, err := yaml.Marshal(&config)
	if err != nil {
		// this err is not tested, but it should not happen either
		return err
	}

	// TODO error case
	ioutil.WriteFile(filepath.Join(os.Getenv("HOME"), ".copy-pastarc"), configContents, 0666)

	return nil
}

// Load loads the config from a runcommands file
func Load() (*Config, error) {
	var config *Config

	byteContent, err := ioutil.ReadFile(filepath.Join(os.Getenv("HOME"), ".copy-pastarc"))
	if err != nil {
		return nil, fmt.Errorf("Unable to load the targets, please check if ~/.copy-pastarc exists %s", err.Error())
	}
	err = yaml.Unmarshal(byteContent, &config)
	if err != nil {
		return nil, fmt.Errorf("Parsing failed %s", err.Error())
	}

	return config, nil
}
