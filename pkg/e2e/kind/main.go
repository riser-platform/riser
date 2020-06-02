package main

import (
	"flag"
)

const (
	// DefaultKindNodeImage should roughly match the latest stable kubernetes version provided by GKE/AKS/EKS
	DefaultKindNodeImage = "kindest/node:v1.16.9"
	// DefaultKindName      = "riser-e2e"
)

func main() {
	var kindNodeImage string
	flag.StringVar(&kindNodeImage, "name", DefaultKindNodeImage, "help message for flagname")

	// TODO: After Assets are refactored
	// deployment := infra.NewDeployment(
	// 	assets http.FileSystem,
	// 	riserConfig *rc.RuntimeConfiguration,
	// 	gitUrl *url.URL)

}
