//go:build e2e
// +build e2e

package e2e

import (
	"fmt"
	"io/ioutil"
	"path"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/stretchr/testify/require"
)

// See smoke_test for common paths. These tests are for less common paths
func Test_DeploymentName(t *testing.T) {
	t.Parallel()
	var testContext *SingleEnvTestContext

	Step(t, "setup test context", func() {
		testContext = SetupSingleEnvTestContext(t)
	})

	appContext := NewRandomAppContext(t, "apps", testContext.IngressDomain)
	defer appContext.Cleanup()

	Step(t, fmt.Sprintf("create app %q", appContext.Name), func() {
		var err error

		shellOrFail(t, "riser apps new %s", appContext.Name)

		app, err := testContext.Riser.Apps.Get(appContext.Name, appContext.Namespace)
		require.NoError(t, err)

		appCfg := model.AppConfig{
			Id:        app.Id,
			Name:      model.AppName(appContext.Name),
			Namespace: model.NamespaceName(appContext.Namespace),
			Image:     "tshak/testdummy",
			Expose: &model.AppConfigExpose{
				ContainerPort: 8000,
			},
		}

		appCfgBytes, err := yaml.Marshal(appCfg)
		require.NoError(t, err)
		appCfgPath := path.Join(appContext.AppDir, "app.yaml")
		err = ioutil.WriteFile(appCfgPath, appCfgBytes, 0644)
		require.NoError(t, err)
	})

	versionA := "0.0.15"
	deploymentName := fmt.Sprintf("%s-alt1", appContext.Name)

	Step(t, fmt.Sprintf("deploy %q version %q", deploymentName, versionA), func() {
		DeployArgsOrFail(t, appContext.AppDir, versionA, testContext.RiserEnvironment, fmt.Sprintf("--name %s", deploymentName))

		err := testContext.Http.RetryGet(appContext.UrlByName("/version", deploymentName), func(r *httpResult) bool {
			return string(r.body) == versionA
		})
		require.NoError(t, err)
	})

	Step(t, fmt.Sprintf("delete deployment %q", deploymentName), func() {
		DeleteDeploymentOrFail(t, appContext.AppDir, deploymentName, testContext.RiserEnvironment)
		AssertDeploymentDeleted(t, testContext, appContext, deploymentName)
	})

	// Ensure that we can deploy a deleted deployment again
	Step(t, fmt.Sprintf("redeploy %q version %q", deploymentName, versionA), func() {
		DeployArgsOrFail(t, appContext.AppDir, versionA, testContext.RiserEnvironment, fmt.Sprintf("--name %s", deploymentName))

		err := testContext.Http.RetryGet(appContext.UrlByName("/version", deploymentName), func(r *httpResult) bool {
			return string(r.body) == versionA
		})
		require.NoError(t, err)
	})

	Step(t, fmt.Sprintf("delete deployment %q", deploymentName), func() {
		DeleteDeploymentOrFail(t, appContext.AppDir, deploymentName, testContext.RiserEnvironment)
		AssertDeploymentDeleted(t, testContext, appContext, deploymentName)
	})
}
