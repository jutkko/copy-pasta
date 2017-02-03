package runcommands

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

type Target struct {
	Name            string `yaml:"name"`
	AccessKey       string `yaml:"accesskey"`
	SecretAccessKey string `yaml:"secretaccesskey"`
}

type Config struct {
	CurrentTarget *Target            `yaml:"currenttarget"`
	Targets       map[string]*Target `yaml:"targets"`
}

func Update(target, accessKey, secretAccessKey string) error {
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
		AccessKey:       accessKey,
		SecretAccessKey: secretAccessKey,
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

func Load() (*Config, error) {
	var config *Config

	byteContent, err := ioutil.ReadFile(filepath.Join(os.Getenv("HOME"), ".copy-pastarc"))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to load the targets, please check if ~/.copy-pastarc exists %s", err.Error()))
	}
	err = yaml.Unmarshal(byteContent, &config)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Parsing failed %s", err.Error()))
	}

	return config, nil
}
