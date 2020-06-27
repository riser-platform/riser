// +build e2e

package e2e

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type appContext struct {
	Name          string
	Namespace     string
	AppDir        string
	IngressDomain string
}

func newRandomAppContext(t *testing.T, namespace string, ingressDomain string) *appContext {
	appName := fmt.Sprintf("e2e-app-%s", randomString(6))

	tmpDir, err := ioutil.TempDir(os.TempDir(), "riser-e2e-")
	require.NoError(t, err)

	return &appContext{
		Name:          appName,
		Namespace:     namespace,
		IngressDomain: ingressDomain,
		AppDir:        tmpDir,
	}
}

func (ctx *appContext) Url(pathAndQuery string) string {
	return ctx.UrlByName(pathAndQuery, ctx.Name)
}

func (ctx *appContext) UrlByName(pathAndQuery, deploymentName string) string {
	pqParsed, _ := url.Parse(pathAndQuery)
	baseUrl, _ := url.Parse(fmt.Sprintf("https://%s.%s.%s", deploymentName, ctx.Namespace, ctx.IngressDomain))
	return baseUrl.ResolveReference(pqParsed).String()
}

// Cleanup currently only cleans up temporary file system resources used. It does not clean up any riser resources (e.g. deployments)
func (ctx *appContext) Cleanup() {
	os.RemoveAll(ctx.AppDir)
}
