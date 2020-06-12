// +build e2e

package e2e

import (
	"encoding/json"
	"fmt"
	"os"
	"riser/pkg/rc"
	"strings"
	"testing"

	"github.com/riser-platform/riser-server/pkg/sdk"
	"github.com/stretchr/testify/require"
)

const (
	RiserApiKeyEnv          = "RISER_APIKEY"
	DefaultRiserServerUrl   = "http://riser-server.riser-system.svc.cluster.local"
	DefaultRiserContextName = "e2e"
)

var currentTestContext *singleEnvTestContext

type singleEnvTestContext struct {
	KubeContext      string
	RiserContext     string
	RiserEnvironment string
	IngressIP        string
	IngressDomain    string
	Riser            *sdk.Client
	Http             *ingressClient
}

func setupSingleEnvTestContext(t *testing.T) *singleEnvTestContext {
	// Cache since this must not change between tests
	if currentTestContext != nil {
		return currentTestContext
	}
	// Use riser context current by default. If it doesn't exist, construct one
	// TODO: Allow alternate riserrc file
	riserContext := shellOrFail(t, "riser context current")
	if strings.TrimSpace(riserContext) == "" {
		riserContext = setupE2ERiserContext(t)
	}

	riserClient, err := getRiserClient()
	require.NoError(t, err)

	ingressIp := shellOrFail(t, "kubectl get service istio-ingressgateway -n istio-system -o jsonpath='{.status.loadBalancer.ingress[0].ip}'")
	if ingressIp == "" {
		ingressIp = shellOrFail(t, "kubectl get service istio-ingressgateway -n istio-system -o jsonpath='{.spec.clusterIP}'")
		t.Log("Warning: istio gateway loadbalancer IP not found. Proceeding with clusterIP.")
	}

	ctx := &singleEnvTestContext{
		RiserContext:     riserContext,
		RiserEnvironment: shellOrFail(t, `kubectl get cm riser-controller -n riser-system -o jsonpath="{.data['RISER_ENVIRONMENT']}"`),
		IngressIP:        ingressIp,
		IngressDomain:    getRiserDomain(t),
		Riser:            riserClient,
	}

	ctx.Http = NewIngressClient(ctx.IngressIP)
	currentTestContext = ctx

	// Validate riser environment is setup
	err = Retry(func() (bool, error) {
		// ... grep until we support json output
		_, err = shell("riser environments list | grep %s", ctx.RiserEnvironment)
		return err == nil, err
	})
	if err != nil {
		t.Fatalf("Environment %q does not exist in riser: %v", ctx.RiserEnvironment, err)
	}
	return currentTestContext
}

func setupE2ERiserContext(t *testing.T) string {
	apiKey := os.Getenv(RiserApiKeyEnv)
	if apiKey == "" {
		t.Fatalf("No riser context found. Either create a riser context or specify the env var %s to use the default riser e2e context", RiserApiKeyEnv)
	}
	shellOrFail(t, fmt.Sprintf("riser context save %s %s %s",
		DefaultRiserContextName,
		DefaultRiserServerUrl,
		apiKey))

	return DefaultRiserContextName
}

func getRiserClient() (*sdk.Client, error) {
	cfg, err := rc.LoadRc()
	if err != nil {
		return nil, err
	}

	ctx, err := cfg.CurrentContext()
	if err != nil {
		return nil, err
	}

	client, err := sdk.NewClient(ctx.ServerURL, ctx.Apikey)
	if err != nil {
		return nil, err
	}

	if ctx.Secure != nil && !*ctx.Secure {
		client.MakeInsecure()
	}

	return client, err
}

type configMap struct {
	Data map[string]string `json:"data"`
}

func getRiserDomain(t *testing.T) string {
	// We can't use jsonpath because of how knative stores domain config
	domainConfigJson := shellOrFail(t, "kubectl get cm config-domain -n knative-serving -o json")
	domainConfigMap := configMap{}
	err := json.Unmarshal([]byte(domainConfigJson), &domainConfigMap)
	require.NoError(t, err)
	var domain string
	for key := range domainConfigMap.Data {
		domain = key
		break
	}
	require.NotEmpty(t, domain)
	return domain
}
