package config

import (
	"io/ioutil"

	"github.com/tshak/riser-server/api/v1/model"

	"github.com/ghodss/yaml"
)

// LoadApp loads an app yaml from a file unmarshalled into an API model.
func LoadApp(pathToAppConfig string) (*model.AppConfigWithOverrides, error) {
	rawFile, err := ioutil.ReadFile(pathToAppConfig)
	if err != nil {
		return nil, err
	}
	app := &model.AppConfigWithOverrides{}
	err = yaml.Unmarshal(rawFile, app)
	if err != nil {
		return nil, err
	}

	return app, nil
}
