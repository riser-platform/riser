package cmd

import (
	"github.com/riser-platform/riser/sdk"
	"riser/pkg/rc"
	"riser/pkg/ui"
)

func getRiserClient(c *rc.Context) *sdk.Client {
	client, err := sdk.NewClient(c.ServerURL, c.Apikey)
	ui.ExitIfErrorMsg(err, "Error instantiating riser SDK")

	if c.Secure != nil && !*c.Secure {
		client.MakeInsecure()
	}
	return client
}
