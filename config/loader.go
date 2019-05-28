package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// LoadApp loads an app yaml from a file unmarshalled into an API model.
func LoadApp(pathToAppConfig string) (*App, error) {
	rawFile, err := ioutil.ReadFile(pathToAppConfig)
	if err != nil {
		return nil, err
	}
	app := &App{}
	err = yaml.Unmarshal(rawFile, app)
	if err != nil {
		return nil, err
	}

	return app, nil
}
