package infra

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
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
	clusterExists, err := checkClusterExists(deployer.Name)
	if err != nil {
		return errors.Wrap(err, "Error checking cluster existence")
	}
	if clusterExists {
		return nil
	}

	args := []string{"create", "cluster",
		fmt.Sprintf("--image=%s", deployer.NodeImage),
		fmt.Sprintf("--name=%s", deployer.Name),
		fmt.Sprintf("--wait=%s", WaitForClusterDuration)}
	err = execStreamOutput("kind", args...)
	return err
}

func (deployer *KindDeployer) Destroy() error {
	args := []string{"delete", "cluster", fmt.Sprintf("--name=%s", deployer.Name)}
	return execStreamOutput("kind", args...)
}

func (deployer *KindDeployer) LoadLocalDockerImage(imageName string) error {
	args := []string{"load", "docker-image", "-q", fmt.Sprintf("--name=%s", deployer.Name), imageName}
	return execStreamOutput("kind", args...)
}

func execStreamOutput(cmdName string, arg ...string) error {
	cmd := exec.Command(cmdName, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func checkClusterExists(name string) (bool, error) {
	cmd := exec.Command("kind", "get", "clusters")
	stderr := &bytes.Buffer{}
	cmd.Stderr = stderr
	outBytes, err := cmd.Output()
	if err != nil {
		return false, errors.Wrap(err, stderr.String())
	}

	for _, clusterName := range strings.Split(string(outBytes), "\n") {
		if clusterName == name {
			return true, nil
		}
	}

	return false, nil
}
