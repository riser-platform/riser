package config

import (
	"github.com/tshak/riser-server/api/v1/model"
	"io/ioutil"

	"github.com/ghodss/yaml"
)

// LoadApp loads an app yaml from a file unmarshalled into an API model.
func LoadApp(pathToAppConfig string) (*model.AppWithOverrides, error) {
	rawFile, err := ioutil.ReadFile(pathToAppConfig)
	if err != nil {
		return nil, err
	}
	app := &model.AppWithOverrides{}
	err = yaml.Unmarshal(rawFile, app)
	if err != nil {
		return nil, err
	}

	return app, nil
}
