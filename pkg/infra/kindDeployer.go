package infra

import (
	"fmt"
	"os"
	"os/exec"
)

const WaitForClusterDuration = "5m"

type KindDeployer struct {
	NodeImage string
	Name      string
}

func NewKindDeployer(nodeImage, name string) *KindDeployer {
	return &KindDeployer{nodeImage, name}
}

func (deployer *KindDeployer) Deploy() error {
	args := []string{"create", "cluster",
		fmt.Sprintf("--image=%s", deployer.NodeImage),
		fmt.Sprintf("--name=%s", deployer.Name),
		fmt.Sprintf("--wait=%s", WaitForClusterDuration)}
	cmd := exec.Command("kind", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (deployer *KindDeployer) Destroy() error {
	args := []string{"delete", "cluster", fmt.Sprintf("--name=%s", deployer.Name)}
	cmd := exec.Command("kind", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
