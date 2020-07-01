package main

import (
	"fmt"
	"os"
	"os/exec"
	"riser/pkg/assets"
	"riser/pkg/infra"
	"riser/pkg/rc"
	"riser/pkg/steps"
	"riser/pkg/ui"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	// DefaultKindNodeImage should roughly match the latest stable kubernetes version provided by GKE/AKS/EKS
	DefaultKindNodeImage = "kindest/node:v1.16.9"
	// DefaultKindName is the name of the kind cluster as well as the riser context by convention
	DefaultKindName = "riser-e2e"
)

func main() {
	var kindNodeImage string
	var kindName string
	var gitUrl string
	var gitSSHKeyPath string
	var destroy bool
	var riserServerImage string
	var riserControllerImage string
	cmd := &cobra.Command{}
	cmd.Flags().StringVar(&kindNodeImage, "image", DefaultKindNodeImage, "node docker image to use for booting the cluster")
	cmd.Flags().StringVar(&kindName, "name", DefaultKindName, "cluster context and riser context name")
	cmd.Flags().StringVar(&gitUrl, "git-url", "", "the git url for the state repo")
	cmd.Flags().StringVar(&gitSSHKeyPath, "git-ssh-key-path", "", "optional path to a git ssh key.")
	cmd.Flags().BoolVar(&destroy, "destroy", false, "destroy the cluster if it already exists")
	cmd.Flags().StringVar(&riserServerImage, "riser-server-image", infra.DefaultServerImage, "the riser server image")
	cmd.Flags().StringVar(&riserControllerImage, "riser-controller-image", infra.DefaultControllerImage, "the riser controller image")
	err := cobra.MarkFlagRequired(cmd.Flags(), "git-url")
	ui.ExitIfError(err)

	cmd.Run = func(_ *cobra.Command, args []string) {
		// TODO: Add support for alternate rc path
		config, err := rc.LoadRc()
		ui.ExitIfError(err)

		err = steps.Run(
			steps.NewFuncStep("Deploying Kind", func() error {
				kindDeployment := infra.NewKindDeployer(kindNodeImage, kindName)
				if destroy {
					err = kindDeployment.Destroy()
					if err != nil {
						return err
					}
				}
				err = kindDeployment.Deploy()
				if err != nil {
					return err
				}
				// TODO: Add support for loading a published container or a different local container name
				return kindDeployment.LoadLocalDockerImage("riser.dev/riser-e2e:local")
			}),
			steps.NewShellExecStep("Create riser-e2e namespace", "kubectl create namespace riser-e2e --dry-run=client -o yaml | kubectl apply -f -"),
			steps.NewShellExecStep("Add istio-injection label", "kubectl label namespace riser-e2e istio-injection=enabled --overwrite=true"),
			steps.NewFuncStep("Deploying Riser", func() error {
				riserDeployment := infra.NewRiserDeployment(
					assets.Assets,
					config,
					gitUrl)
				riserDeployment.EnvironmentName = kindName
				riserDeployment.ServerImage = riserServerImage
				riserDeployment.ControllerImage = riserControllerImage
				riserDeployment.GitSSHKeyPath = gitSSHKeyPath
				err = riserDeployment.Deploy()
				if err != nil {
					return err
				}
				riserCtx, err := config.CurrentContext()
				if err != nil {
					return errors.Wrap(err, "Error reading riser context")
				}
				return steps.NewShellExecStep("Create secret for e2e tests",
					"kubectl create secret generic riser-e2e --namespace=riser-e2e "+
						fmt.Sprintf("--from-literal=RISER_APIKEY=%s --dry-run=client -o yaml | kubectl apply -f -", riserCtx.Apikey)).Exec()
			}),
			steps.NewShellExecStep("Deploy e2e tests", "kubectl delete job riser-e2e --namespace=riser-e2e --ignore-not-found=true --wait=true && kubectl apply -f ./e2e/job.yaml"),
			steps.NewShellExecStep("Wait for test run to start", "kubectl wait --namespace=riser-e2e --for=condition=initialized --timeout=30s -l job-name=riser-e2e pod"),
			steps.NewRetryStep(func() steps.Step {
				return steps.NewFuncStep("Stream test results", func() error {
					jobCmd := exec.Command("kubectl", "logs", "-l=job-name=riser-e2e", "--namespace=riser-e2e", "-f", "-c=riser-e2e")
					// Stream logs to stdout
					jobCmd.Stdout = os.Stdout
					return jobCmd.Run()
				})
			}, 30, steps.AlwaysRetry()),
			// The job won't terminate because of the istio sidecar (https://github.com/kubernetes/kubernetes/issues/25908)
			// Grab the container exitCode to determine success or not.
			steps.NewFuncStep("Check test results", func() error {
				// The sleep hack is here due to a race condition between the container exiting and the containerStatus being updated.
				jobCmd := exec.Command("sh", "-c", `sleep 5 && kubectl get po --namespace=riser-e2e -l job-name=riser-e2e -o jsonpath='{.items[0].status.containerStatuses[?(@.name=="riser-e2e")].state.terminated.exitCode}'`)
				output, err := jobCmd.CombinedOutput()
				if err != nil {
					return fmt.Errorf("Error executing command: %s", string(output))
				}
				if string(output) == "0" {
					return nil
				}

				return fmt.Errorf("Received unexpected output: %s", output)
			}),
		)

		ui.ExitIfError(err)
	}

	err = cmd.Execute()
	ui.ExitIfError(err)
}
