// +build e2e

package e2e

import (
	"fmt"
	"io/ioutil"
	"path"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// See smoke_test for common paths. These tests are for less common paths
func Test_DeploymentName(t *testing.T) {
	var testContext *singleEnvTestContext

	step("setup test context", func() {
		testContext = setupSingleEnvTestContext(t)
	})

	namespace := "apps"
	appContext := newRandomAppContext(t, namespace, testContext.IngressDomain)
	defer appContext.Cleanup()

	step(fmt.Sprintf("create app %q", appContext.Name), func() {
		var err error

		shellOrFail(t, "riser apps new %s", appContext.Name)

		app, err := testContext.Riser.Apps.Get(appContext.Name, namespace)
		require.NoError(t, err)

		appCfg := model.AppConfig{
			Id:        app.Id,
			Name:      model.AppName(appContext.Name),
			Namespace: model.NamespaceName(namespace),
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
	step(fmt.Sprintf("deploy %q version %q", deploymentName, versionA), func() {
		deployArgsOrFail(t, appContext.AppDir, versionA, testContext.RiserEnvironment, fmt.Sprintf("--name %s", deploymentName))

		err := testContext.Http.RetryGet(appContext.UrlByName("/version", deploymentName), func(r *httpResult) bool {
			return string(r.body) == versionA
		})
		require.NoError(t, err)
	})

	step(fmt.Sprintf("delete deployment %q", deploymentName), func() {
		deleteDeploymentOrFail(t, appContext.AppDir, deploymentName, testContext.RiserEnvironment)

		// Wait until no deployments in status
		err := Retry(func() (bool, error) {
			appStatus, err := testContext.Riser.Apps.GetStatus(appContext.Name, namespace)
			if err != nil {
				return true, err
			}

			return len(appStatus.Deployments) == 0, err
		})
		require.NoError(t, err)

		// Check kube resources
		err = Retry(func() (bool, error) {
			configResult := shellOrFail(t, fmt.Sprintf("kubectl get config %s -n %s --ignore-not-found", appContext.Name, namespace))
			return configResult == "", nil
		})
		assert.NoError(t, err)

		err = Retry(func() (bool, error) {
			routeResult := shellOrFail(t, fmt.Sprintf("kubectl get route %s -n %s --ignore-not-found", appContext.Name, namespace))
			return routeResult == "", nil
		})
		assert.NoError(t, err)
	})
}
