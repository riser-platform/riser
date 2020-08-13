package cmd

import (
	"os"
	"path"
	"riser/pkg/rc"
	"riser/pkg/ui"
	"strings"

	"github.com/riser-platform/riser-server/pkg/sdk"
)

// safeCurrentContext loads the CurrentContext and exits if there is any error
func safeCurrentContext(cfg *rc.RuntimeConfiguration) *rc.Context {
	context, err := cfg.CurrentContext()
	ui.ExitIfError(err)
	return context
}

func getRiserClient(c *rc.Context) *sdk.Client {
	client, err := sdk.NewClient(c.ServerURL, c.Apikey)
	ui.ExitIfErrorMsg(err, "Error instantiating riser SDK")

	if c.Secure != nil && !*c.Secure {
		client.MakeInsecure()
	}
	return client
}

// expandTildeInPath expands the tilde to the user's home dir if specified. Whereever possible, use the
// underlying OS's shell to do this. This has not been tested against Windows.
func expandTildeInPath(pathToExpand string) string {
	if strings.HasPrefix(pathToExpand, "~") {
		pathToExpand = path.Join(os.Getenv("HOME"), pathToExpand[1:])
	}
	return pathToExpand
}
