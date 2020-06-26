package cmd

import (
	"riser/pkg/rc"
	"riser/pkg/ui"

	"github.com/spf13/cobra"
)

func newStatusCommand(runtimeConfig *rc.RuntimeConfiguration) *cobra.Command {
	var appName string
	var namespace string
	showAllRevisions := false
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Gets the status for a deployment.",
		Run: func(cmd *cobra.Command, args []string) {
			currentContext := safeCurrentContext(runtimeConfig)
			riserClient := getRiserClient(currentContext)

			status, err := riserClient.Apps.GetStatus(appName, namespace)
			ui.ExitIfErrorMsg(err, "Error getting status")

			view := &statusView{
				appName:             appName,
				activeRevisionsOnly: !showAllRevisions,
				status:              status,
			}
			ui.RenderView(view)
		},
	}

	addAppFlag(cmd.Flags(), &appName)
	addNamespaceFlag(cmd.Flags(), &namespace)
	addOutputFlag(cmd.Flags())
	cmd.Flags().BoolVarP(&showAllRevisions, "all-revisions", "", false, "Shows all available revisions. Otherwise only shows the latest revision and older revisions with traffic")

	return cmd
}
