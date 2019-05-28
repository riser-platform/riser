package client

import (
	"riser/config"
)

type Deployment struct {
	Name   string           `json:"name"`
	Stage  string           `json:"stage"`
	Docker DeploymentDocker `json:docker`
	App    config.App       `json:app`
}

type DeploymentDocker struct {
	Tag string `json:tag`
}
