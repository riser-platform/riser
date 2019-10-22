package cmd

import (
	"fmt"

	version "github.com/hashicorp/go-version"
	"github.com/spf13/cobra"
)

func newVersionCmd(currentVersion *version.Version) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(currentVersion.String())
		},
	}
}
