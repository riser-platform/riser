package e2e

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

type appContext struct {
	Name      string
	Namespace string
	BaseUrl   url.URL
}

func newRandomAppContext(t *testing.T, namespace string, ingressDomain string) *appContext {
	appName := fmt.Sprintf("e2e-app-%s", randomString(6))
	baseUrl, err := url.Parse(fmt.Sprintf("https://%s.%s.%s", appName, namespace, ingressDomain))
	require.NoError(t, err)

	return &appContext{
		Name:      appName,
		Namespace: namespace,
		BaseUrl:   *baseUrl,
	}
}

func (ctx *appContext) Url(pathAndQuery string) string {
	pqParsed, _ := url.Parse(pathAndQuery)
	return ctx.BaseUrl.ResolveReference(pqParsed).String()
}
