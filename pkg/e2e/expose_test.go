// +build e2e

package e2e

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func Test_ExposeClusterOnly(t *testing.T) {
	t.Parallel()
	var testContext *SingleEnvTestContext

	Step(t, "setup test context", func() {
		testContext = SetupSingleEnvTestContext(t)
	})

	namespace := "apps"
	appContext := NewRandomAppContext(t, namespace, testContext.IngressDomain)
	defer appContext.Cleanup()

	Step(t, fmt.Sprintf("create app %q", appContext.Name), func() {
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
				Scope:         model.AppExposeScope_Cluster,
			},
			HealthCheck: &model.AppConfigHealthCheck{
				Path: "/health",
			},
			OverrideableAppConfig: model.OverrideableAppConfig{
				Environment: map[string]intstr.IntOrString{
					"env1": intstr.FromString("val1"),
				},
			},
		}

		appCfgBytes, err := yaml.Marshal(appCfg)
		require.NoError(t, err)
		appCfgPath := path.Join(appContext.AppDir, "app.yaml")
		err = ioutil.WriteFile(appCfgPath, appCfgBytes, 0644)
		require.NoError(t, err)
	})

	versionA := "0.0.15"
	Step(t, fmt.Sprintf("deploy version %q", versionA), func() {
		DeployOrFail(t, appContext.AppDir, versionA, testContext.RiserEnvironment)
		response, err := testContext.Http.Get(appContext.Url("/"))
		require.Error(t, err)
		assert.Equal(t, response.StatusCode, http.NotFound)
		// TODO: Test cluster.local URL.
		// We must first setup a second service to prove that cluster local routing is working as we can't guarantee that
		// the test runner is running on the mesh
	})

	Step(t, fmt.Sprintf("delete deployment %q", appContext.Name), func() {
		DeleteDeploymentOrFail(t, appContext.AppDir, appContext.Name, testContext.RiserEnvironment)
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