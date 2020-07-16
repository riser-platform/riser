// +build e2e

package e2e

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"testing"

	"github.com/ghodss/yaml"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/riser-platform/riser-server/api/v1/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Initial attempt at e2e testing. Just run through a smoke test of a simple happy path. Lots of refactoring to do as we add more tests.
// Kube and Riser context must be pointing to the correct location
func Test_Smoke(t *testing.T) {
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

		err := testContext.Http.RetryGet(appContext.Url("/version"), func(r *httpResult) bool {
			return string(r.body) == versionA
		})
		require.NoError(t, err)

		envResponse, err := testContext.Http.Get(appContext.Url("/env"))
		require.NoError(t, err)
		assert.Equal(t, envResponse.StatusCode, http.StatusOK)

		envBody, err := ioutil.ReadAll(envResponse.Body)
		require.NoError(t, err)

		envMap := ParseTestDummyEnv(envBody)
		require.Equal(t, "val1", envMap["ENV1"])
	})

	secretName := "secret1"
	secretValue := "secretVal1"
	Step(t, "create secret", func() {
		shellOrFail(t, "cd %s && riser secrets save %s %s %s", appContext.AppDir, secretName, secretValue, testContext.RiserEnvironment)
		// We do not wait for the secret to be available in k8s. The next deployment should have the secret ref and
		// not become available until the secret is present.
	})

	versionB := "0.0.16"
	Step(t, fmt.Sprintf("deploy version %q", versionB), func() {
		DeployOrFail(t, appContext.AppDir, versionB, testContext.RiserEnvironment)

		err := testContext.Http.RetryGet(appContext.Url("/version"), func(r *httpResult) bool {
			return string(r.body) == versionB
		})
		require.NoError(t, err)

		envResponse, err := testContext.Http.Get(appContext.Url("/env"))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, envResponse.StatusCode)

		envBody, err := ioutil.ReadAll(envResponse.Body)
		require.NoError(t, err)

		envMap := ParseTestDummyEnv(envBody)
		require.Equal(t, "val1", envMap["ENV1"])
		require.Equal(t, secretValue, envMap[strings.ToUpper(secretName)])
	})

	Step(t, "rollout 50/50 with previous deployment", func() {
		riserOrFail(t, appContext.AppDir, fmt.Sprintf("rollout %s r1:50 r2:50", testContext.RiserEnvironment))
		// Wait until we get one hit from versionA to ensure that the rollout is working before we start taking samples
		err := testContext.Http.RetryGet(appContext.Url("/version"), func(r *httpResult) bool {
			return string(r.body) == versionA
		})
		require.NoError(t, err)

		const samples = 100
		sum := 0.0
		for i := 0; i < samples; i++ {
			err := testContext.Http.RetryGet(appContext.Url("/version"), func(r *httpResult) bool {
				if string(r.body) == versionB {
					sum += 1.0
				}
				return true
			})
			require.NoError(t, err)
		}

		mean := sum / float64(samples)
		// Approximate that the traffic splitting is correct. We are just validating e2e configuration, not the precision of istio's traffic splitting.
		assert.InDelta(t, 0.5, mean, 0.2, "%v")
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

func Test_Namespace(t *testing.T) {
	t.Parallel()
	var testContext *SingleEnvTestContext

	Step(t, "setup test context", func() {
		testContext = SetupSingleEnvTestContext(t)
	})

	namespace := fmt.Sprintf("e2e-ns-%s", RandomString(6))
	appContext := NewRandomAppContext(t, namespace, testContext.IngressDomain)
	defer appContext.Cleanup()

	Step(t, fmt.Sprintf("create namespace %q", namespace), func() {
		riserOrFail(t, appContext.AppDir, fmt.Sprintf("namespaces create %s", namespace))
	})

	Step(t, fmt.Sprintf("create app %q in namespace %q", appContext.Name, namespace), func() {
		var err error

		shellOrFail(t, "riser apps new %s -n %s", appContext.Name, namespace)

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
		err := testContext.Http.RetryGet(appContext.Url("/version"), func(r *httpResult) bool {
			return string(r.body) == versionA
		})
		require.NoError(t, err)

		envResponse, err := testContext.Http.Get(appContext.Url("/env"))
		require.NoError(t, err)
		assert.Equal(t, envResponse.StatusCode, http.StatusOK)

		envBody, err := ioutil.ReadAll(envResponse.Body)
		require.NoError(t, err)

		envMap := ParseTestDummyEnv(envBody)
		require.Equal(t, "val1", envMap["ENV1"])
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
