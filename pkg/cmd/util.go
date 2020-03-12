package cmd

import (
	"riser/pkg/rc"
	"riser/pkg/ui"

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
