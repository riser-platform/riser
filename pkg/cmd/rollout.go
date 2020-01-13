package cmd

import (
	"fmt"
	"riser/pkg/config"
	"riser/pkg/rc"
	"riser/pkg/ui"

	"github.com/spf13/cobra"
)

func newRolloutCommand(currentContext *rc.Context) *cobra.Command {
	var deploymentName string
	cmd := &cobra.Command{
		Use:     "rollout (stage) (trafficRule0) [trafficRuleN...]",
		Short:   "Manually controls traffic for a deployment's rollout",
		Long:    "Manually controls traffic for a deployment's rollout. Typically only used when a deployment is deployed with the \"--manual-rollout\" flag. Traffic rules are in the format \"(rev):(%traffic)\" where \"rev\" is the riser revision as shown in \"riser status\"",
		Args:    cobra.MinimumNArgs(2),
		Example: "riser rollout prod 1:90 2:10 // Canary \nriser rollout prod 2:100 // Route all traffic to rev 2",
		Run: func(cmd *cobra.Command, args []string) {
			stage := args[0]
			riserClient := getRiserClient(currentContext)
			err := riserClient.Rollouts.Save(deploymentName, stage, args[1:]...)
			ui.ExitIfError(err)
			fmt.Println("Rollout requested")
		},
	}

	cmd.Flags().StringVarP(&deploymentName, "name", "n", config.SafeLoadDefaultAppName(), "The name of the deployment")
	return cmd
}
