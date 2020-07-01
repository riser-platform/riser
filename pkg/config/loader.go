package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"riser/pkg/logger"

	"github.com/ghodss/yaml"

	"github.com/riser-platform/riser-server/api/v1/model"
)

// DefaultNamespace is the default namespace to use when a namespace is not specified.
const DefaultNamespace = "apps"

var DefaultAppConfigPaths = []string{"./app.yml", "./app.yaml"}

var appConfigCache = map[string]*model.AppConfigWithOverrides{}
var alreadyWarned bool

type unmarshalAppError struct {
	path           string
	unmarshalError error
}

func (err *unmarshalAppError) Error() string {
	return err.unmarshalError.Error()
}

func warnIfUnmartialAppError(err error) {
	if alreadyWarned {
		return
	}
	if v, ok := err.(*unmarshalAppError); ok {
		alreadyWarned = true
		logger.Log().Warn(fmt.Sprintf("Error parsing app config %q: %s", v.path, err))
	}
}

// LoadAppFromConfig loads an app yaml from a file unmarshalled into an API model.
func LoadAppFromConfig(pathToAppConfig string) (*model.AppConfigWithOverrides, error) {
	// This func is called from multiple cobra arguments so prevent constantly hitting the disk and unmarshalling
	if fromCache, ok := appConfigCache[pathToAppConfig]; ok {
		return fromCache, nil
	}
	rawFile, err := ioutil.ReadFile(pathToAppConfig)
	if err != nil {
		return nil, err
	}
	app := &model.AppConfigWithOverrides{}
	err = yaml.UnmarshalStrict(rawFile, app, yaml.DisallowUnknownFields)
	if err != nil {
		return nil, &unmarshalAppError{path: pathToAppConfig, unmarshalError: err}
	}
	appConfigCache[pathToAppConfig] = app

	return app, nil
}

// SafeLoadAppName attempts to retrieve the name of the app in the specified path.
// An empty string is returned if the file does not exist, cannot be be parsed, or if any other error occurs.
func SafeLoadAppName(pathToAppConfig string) string {
	appConfig, err := LoadAppFromConfig(pathToAppConfig)
	warnIfUnmartialAppError(err)
	if err == nil {
		return string(appConfig.Name)
	}

	return ""
}

// SafeLoadDefaultAppName attempts to retrieve the name of the app in the default app config locations
// Returns an empty string if the file does not exist, cannot be be parsed, or if any other error occurs.
func SafeLoadDefaultAppName() string {
	for _, pathToAppConfig := range DefaultAppConfigPaths {
		result := SafeLoadAppName(pathToAppConfig)
		if result != "" {
			return result
		}
	}

	return ""
}

// SafeLoadAppNamespace attempts to retrieve the namespace of the app in the specified path.
// Returns an empty string if the namespace is not specified, the file does not exist, cannot be be parsed,
// or if any other error occurs.
func SafeLoadAppNamespace(pathToAppConfig string) string {
	appConfig, err := LoadAppFromConfig(pathToAppConfig)
	warnIfUnmartialAppError(err)
	if err == nil {
		return string(appConfig.Namespace)
	}
	return ""
}

// SafeLoadDefaultAppNamespace attempts to retrieve the namespace of the app in the default app config locations
// Returns the default namespace "apps" if the namespace is not specified, the file does not exist, cannot be be parsed,
// or if any other error occurs.
func SafeLoadDefaultAppNamespace() string {
	for _, pathToAppConfig := range DefaultAppConfigPaths {
		result := SafeLoadAppNamespace(pathToAppConfig)
		if result != "" {
			return result
		}
	}

	return DefaultNamespace
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
