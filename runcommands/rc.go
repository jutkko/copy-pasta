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
	AccessKey       string `yaml:"accesskey"`
	SecretAccessKey string `yaml:"secretaccesskey"`
}

func Update(target, accessKey, secretAccessKey string) error {
	var targets map[string]*Target
	var err error

	targets, err = Load()
	if err != nil {
		targets = make(map[string]*Target)
	}

	targets[target] = &Target{
		AccessKey:       accessKey,
		SecretAccessKey: secretAccessKey,
	}

	targetsContents, err := yaml.Marshal(&targets)
	if err != nil {
		// this err is not tested, but it should not happen either
		return err
	}

	// TODO error case
	ioutil.WriteFile(filepath.Join(os.Getenv("HOME"), ".copy-pastarc"), targetsContents, 0666)
	return nil
}

func Load() (map[string]*Target, error) {
	var targets map[string]*Target

	byteContent, err := ioutil.ReadFile(filepath.Join(os.Getenv("HOME"), ".copy-pastarc"))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to load the targets, please check if ~/.copy-pastarc exists %s", err.Error()))
	}
	err = yaml.Unmarshal(byteContent, &targets)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Parsing failed %s", err.Error()))
	}

	return targets, nil
}
