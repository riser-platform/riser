// +build e2e

package e2e

import (
	"encoding/json"
	"riser/pkg/rc"
	"testing"

	"github.com/riser-platform/riser-server/pkg/sdk"
	"github.com/stretchr/testify/require"
)

var currentTestContext *singleStageTestContext

type singleStageTestContext struct {
	KubeContext   string
	RiserContext  string
	RiserStage    string
	IngressIP     string
	IngressDomain string
	Riser         *sdk.Client
	Http          *ingressClient
}

func setupSingleStageTestContext(t *testing.T) *singleStageTestContext {
	// Cache since this must not change between tests
	if currentTestContext != nil {
		return currentTestContext
	}
	riserClient, err := getRiserClient()
	require.NoError(t, err)
	ctx := &singleStageTestContext{
		KubeContext:   shellOrFail(t, "kubectl config current-context"),
		RiserContext:  shellOrFail(t, "riser context current"),
		RiserStage:    shellOrFail(t, `kubectl get cm riser-controller -n riser-system -o jsonpath="{.data['RISER_STAGE']}"`),
		IngressIP:     shellOrFail(t, "kubectl get service istio-ingressgateway -n istio-system -o jsonpath='{.status.loadBalancer.ingress[0].ip}'"),
		IngressDomain: getRiserDomain(t),
		Riser:         riserClient,
	}
	ctx.Http = NewIngressClient(ctx.IngressIP)
	currentTestContext = ctx
	return currentTestContext
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
