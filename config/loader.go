package config

import (
	"io/ioutil"

	"github.com/tshak/riser-server/api/v1/model"

	"github.com/ghodss/yaml"
)

var DefaultAppConfigPaths = []string{"./app.yml", "./app.yaml"}

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

// SafeLoadAppName attempts to retrieve the name of the app in the specified path.
// An empty string is returned if the file does not exist, cannot be be parsed, or if any other error ocurrs.
func SafeLoadAppName(pathToAppConfig string) string {
	appConfig, err := LoadApp(pathToAppConfig)
	if err == nil {
		return appConfig.Name
	}
	return ""
}

// SafeLoadDefaultAppName attempts to retrieve the name of the app in the default app config locations
// An empty string is returned if the file does not exist, cannot be be parsed, or if any other error ocurrs.
func SafeLoadDefaultAppName() string {
	for _, pathToAppConfig := range DefaultAppConfigPaths {
		appName := SafeLoadAppName(pathToAppConfig)
		if len(appName) > 0 {
			return appName
		}
	}

	return ""
}
