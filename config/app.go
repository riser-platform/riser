// Package config provides config types and related utilities
package config

// TODO: Either import model from riser-server and/or convert to using json tags and use https://github.com/bronze1man/yaml2json

// App describes the root of an app config
type App struct {
	AppCore `yaml:",inline"`
	Stages  map[string]AppStage `yaml:"stages,omitempty"`
}

type AppCore struct {
	Name        string            `yaml:"name"`
	Deploy      *AppDeploy        `yaml:"deploy,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty"`
	Expose      *AppExpose        `yaml:"expose,omitempty"`
}

// AppDeploy describes the deploy section of an app config
type AppDeploy struct {
	Replicas int `yaml:"replicas"`
}

// AppExpose describes the expose section of an app config
type AppExpose struct {
	ContainerPort int `yaml:"containerPort"`
}

// AppStage contains overrides for fields in AppCore
type AppStage struct {
	AppCore `yaml:",inline"`
}
