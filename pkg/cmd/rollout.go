package cmd

import (
	"fmt"
	"riser/pkg/rc"
	"riser/pkg/ui"

	"github.com/spf13/cobra"
)

func newRolloutCommand(runtimeConfig *rc.RuntimeConfiguration) *cobra.Command {
	var deploymentName string
	var namespace string
	cmd := &cobra.Command{
		Use:     "rollout (targetEnvironment) (trafficRule0) [trafficRuleN...]",
		Short:   "Manually controls traffic for a deployment's rollout",
		Long:    "Manually controls traffic for a deployment's rollout. Typically only used when a deployment is deployed with the \"--manual-rollout\" flag. Traffic rules are in the format \"r(rev#):(traffic%)\" where \"rev\" is the riser revision as shown in \"riser status\"",
		Args:    cobra.MinimumNArgs(2),
		Example: "  riser rollout prod r1:90 r2:10 // Canary routing 10% of traffic to a new revision \n  riser rollout prod r2:100 // Route all traffic to rev 2",
		Run: func(cmd *cobra.Command, args []string) {
			currentContext := safeCurrentContext(runtimeConfig)
			environmentName := args[0]
			riserClient := getRiserClient(currentContext)
			err := riserClient.Rollouts.Save(deploymentName, namespace, environmentName, args[1:]...)
			ui.ExitIfError(err)
			fmt.Println("Rollout requested")
		},
	}

	addDeploymentNameFlag(cmd.Flags(), &deploymentName)
	addNamespaceFlag(cmd.Flags(), &namespace)
	return cmd
}
