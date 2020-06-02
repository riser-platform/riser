package infra

import (
	"fmt"
	"os/exec"

	"github.com/pkg/errors"
)

type KindDeployer struct {
	NodeImage string
	Name      string
}

func NewKindDeployer(nodeImage, name string) *KindDeployer {
	return &KindDeployer{nodeImage, name}
}

func (deployer *KindDeployer) Deploy() error {
	args := []string{"create", "cluster", fmt.Sprintf("--image=%s", deployer.NodeImage), fmt.Sprintf("--name=%s", deployer.Name)}
	cmd := exec.Command("kind", args...)
	output, _ := cmd.CombinedOutput()

	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error executing kind: %v", string(output)))
	}

	return nil
}
