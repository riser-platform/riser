// +build e2e

package e2e

import (
	"fmt"
	"testing"
	"time"
)

const deployTimeoutSeconds = 60

// deployOrFail calls riser deploy with --wait
func deployOrFail(t *testing.T, appDir, dockerTag, environment string) {
	riserOrFail(t, appDir, fmt.Sprintf("deploy %s %s --wait --wait-seconds=%d", dockerTag, environment, deployTimeoutSeconds))
}

func deleteDeploymentOrFail(t *testing.T, appDir, deploymentName, environment string) {
	riserOrFail(t, appDir, fmt.Sprintf("deployments delete %s %s --no-prompt", deploymentName, environment))
}

func riserOrFail(t *testing.T, appDir, command string) {
	shellOrFailTimeout(t, time.Duration(deployTimeoutSeconds)*time.Second, "cd %s && riser %s", appDir, command)
}
