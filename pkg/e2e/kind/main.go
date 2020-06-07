package main

import (
	"net/url"
	"riser/pkg/assets"
	"riser/pkg/infra"
	"riser/pkg/rc"
	"riser/pkg/steps"
	"riser/pkg/ui"

	"github.com/spf13/cobra"
)

const (
	// DefaultKindNodeImage should roughly match the latest stable kubernetes version provided by GKE/AKS/EKS
	DefaultKindNodeImage = "kindest/node:v1.16.9"
	// DefaultKindName is the name of the kind cluster as well as the riser context by convention
	DefaultKindName = "riser-e2e"
)

func deploy() {

}

func main() {
	var kindNodeImage string
	var kindName string
	var gitUrlRaw string
	cmd := &cobra.Command{}
	cmd.Flags().StringVar(&kindNodeImage, "image", DefaultKindNodeImage, "node docker image to use for booting the cluster")
	cmd.Flags().StringVar(&kindName, "name", DefaultKindName, "cluster context and riser context name")
	cmd.Flags().StringVar(&gitUrlRaw, "git-url", "", "the git url for the state repo")
	err := cobra.MarkFlagRequired(cmd.Flags(), "git-url")
	ui.ExitIfError(err)

	cmd.Run = func(_ *cobra.Command, args []string) {
		gitUrl, err := url.Parse(gitUrlRaw)
		ui.ExitIfErrorMsg(err, "Error parsing git url")

		// TODO: Add support for alternate rc path
		config, err := rc.LoadRc()
		ui.ExitIfError(err)

		err = steps.Run(
			steps.NewFuncStep("Deploying Kind", func() error {
				kindDeployment := infra.NewKindDeployer(kindNodeImage, kindName)
				err = kindDeployment.Destroy()
				if err != nil {
					return err
				}
				return kindDeployment.Deploy()
			}),
			steps.NewFuncStep("Deploying Riser", func() error {
				riserDeployment := infra.NewRiserDeployment(
					assets.Assets,
					config,
					gitUrl)
				riserDeployment.EnvironmentName = kindName
				return riserDeployment.Deploy()
			}),
		)
		ui.ExitIfError(err)
	}

	err = cmd.Execute()
	ui.ExitIfError(err)
}
