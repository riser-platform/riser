package cmd

import (
	"testing"

	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_newDeploymentsDescribeView(t *testing.T) {
	app := &model.App{}
	appStatus := &model.AppStatus{
		Environments: []model.EnvironmentStatus{
			{EnvironmentName: "env1"},
			{EnvironmentName: "env2"},
		},
		Deployments: []model.DeploymentStatus{
			{DeploymentName: "mydep", EnvironmentName: "env1", RiserRevision: 1},
			{DeploymentName: "mydep", EnvironmentName: "env2", RiserRevision: 1},
			{DeploymentName: "mydep-pr-1", EnvironmentName: "env2", RiserRevision: 1},
			{DeploymentName: "mydep", EnvironmentName: "env1", RiserRevision: 2},
			{DeploymentName: "mydep", EnvironmentName: "env2", RiserRevision: 2},
		},
	}

	result, err := newDeploymentsDescribeView(app, appStatus, "mydep", "env2", "demo.riser")

	assert.NoError(t, err)
	assert.Equal(t, app, result.app)
	assert.Equal(t, "mydep", result.deploymentName)
	assert.Equal(t, "demo.riser", result.externalHost)
	assert.Equal(t, "env2", result.environmentStatus.EnvironmentName)
	// Filter all deployments by name and environment
	require.Len(t, result.deploymentStatus, 2)
	assert.Equal(t, "mydep", result.deploymentStatus[0].DeploymentName)
	assert.Equal(t, "env2", result.deploymentStatus[0].EnvironmentName)
	assert.Equal(t, int64(1), result.deploymentStatus[0].RiserRevision)
	assert.Equal(t, "mydep", result.deploymentStatus[1].DeploymentName)
	assert.Equal(t, "env2", result.deploymentStatus[1].EnvironmentName)
	assert.Equal(t, int64(2), result.deploymentStatus[1].RiserRevision)
}

func Test_newDeploymentsDescribeView_InvalidDeployment(t *testing.T) {
	app := &model.App{
		Namespace: model.NamespaceName("apps"),
	}
	appStatus := &model.AppStatus{
		Environments: []model.EnvironmentStatus{
			{EnvironmentName: "env1"},
		},
		Deployments: []model.DeploymentStatus{
			{DeploymentName: "mydep", EnvironmentName: "env1", RiserRevision: 1},
		},
	}

	result, err := newDeploymentsDescribeView(app, appStatus, "mydep", "env2", "demo.riser")

	assert.Nil(t, result)
	assert.Equal(t, `The environment "env2" does not contain the deployment "mydep" in the "apps" namespace`, err.Error())
}

func Test_formatExternalUrl(t *testing.T) {
	result := formatExternalUrl("mydep", "apps", "demo.riser")

	assert.Equal(t, "https://mydep.apps.demo.riser", result)
}

func Test_formatClusterLocalUrl(t *testing.T) {
	result := formatClusterLocalUrl("mydep", "apps")

	assert.Equal(t, "http://mydep.apps.svc.cluster.local", result)
}
