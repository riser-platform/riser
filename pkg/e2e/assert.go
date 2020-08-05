package e2e

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func AssertDeploymentDeleted(t *testing.T, testContext *SingleEnvTestContext, appContext *AppContext, deploymentName string) {
	// Wait until no deployments in status
	err := Retry(func() (bool, error) {
		appStatus, err := testContext.Riser.Apps.GetStatus(appContext.Name, appContext.Namespace)
		if err != nil {
			return true, err
		}

		return len(appStatus.Deployments) == 0, err
	})
	require.NoError(t, err)

	// Check kube resources
	err = Retry(func() (bool, error) {
		configResult := shellOrFail(t, fmt.Sprintf("kubectl get config %s -n %s --ignore-not-found", deploymentName, appContext.Namespace))
		return configResult == "", nil
	})
	assert.NoError(t, err)

	err = Retry(func() (bool, error) {
		routeResult := shellOrFail(t, fmt.Sprintf("kubectl get route %s -n %s --ignore-not-found", deploymentName, appContext.Namespace))
		return routeResult == "", nil
	})
	assert.NoError(t, err)
}
