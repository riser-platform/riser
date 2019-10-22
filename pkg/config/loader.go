package config

import (
	"io/ioutil"
	"os"

	"github.com/riser-platform/riser-server/api/v1/model"

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
	err = yaml.UnmarshalStrict(rawFile, app, yaml.DisallowUnknownFields)
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
// Returns an empty string if the file does not exist, cannot be be parsed, or if any other error ocurrs.
func SafeLoadDefaultAppName() string {
	for _, pathToAppConfig := range DefaultAppConfigPaths {
		appName := SafeLoadAppName(pathToAppConfig)
		if len(appName) > 0 {
			return appName
		}
	}

	return ""
}

// GetAppConfigPathFromDefaults searches for an app config from the default locations and returns the first found
// Returns an empty string if no file is found.
func GetAppConfigPathFromDefaults() string {
	for _, pathToAppConfig := range DefaultAppConfigPaths {
		if _, err := os.Stat(pathToAppConfig); err == nil {
			return pathToAppConfig
		}
	}
	return ""
}
