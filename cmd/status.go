package cmd

import (
	"fmt"
	"riser/rc"
	"riser/ui/table"
	"strings"

	"github.com/tshak/riser-server/api/v1/model"

	"github.com/spf13/cobra"
	"github.com/tshak/riser/sdk"
	"github.com/wzshiming/ctc"
)

func newStatusCommand(currentContext *rc.RuntimeContext) *cobra.Command {
	var appName string
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Gets the status for a deployment.",
		Run: func(cmd *cobra.Command, args []string) {
			apiClient, err := sdk.NewClient(currentContext.ServerURL)
			if err != nil {
				panic(err)
			}

			statuses, err := apiClient.GetStatus(appName)
			if err != nil {
				panic(err)
			}

			drawStatusSummary(statuses)
		},
	}

	addAppFlag(cmd.Flags(), &appName)

	return cmd
}

func drawStatusSummary(statuses []model.StatusSummary) {
	table := table.Default().Header("Deployment", "Stage", "Rev", "Docker Tag", "Rollout", "Rollout Details")

	for _, status := range statuses {
		table.AddRow(
			status.DeploymentName,
			status.StageName,
			fmt.Sprintf("%d", status.RolloutRevision),
			getDockerTag(status.DockerImage),
			formatRolloutStatus(status.RolloutStatus),
			status.RolloutStatusReason)
	}

	fmt.Println(table)
}

func getDockerTag(dockerImage string) string {
	idx := strings.Index(dockerImage, ":")
	if idx == -1 {
		// This should never happen since we don't allow images without tags or with digests
		return "Unknown"
	}
	return dockerImage[idx+1:]
}

// TODO: add health status back
// func formatHealthStatus(healthStatus string) string {
// 	if healthStatus == model.HealthStatusTrue {
// 		return fmt.Sprint(ctc.ForegroundBrightGreen, healthStatus, ctc.Reset)
// 	}
// 	if healthStatus == model.HealthStatusFalse {
// 		return fmt.Sprint(ctc.ForegroundBrightRed, healthStatus, ctc.Reset)
// 	}
// 	if healthStatus == model.HealthStatusUnknown {
// 		return fmt.Sprint(ctc.ForegroundBrightYellow, healthStatus, ctc.Reset)
// 	}

// 	return healthStatus
// }

func formatRolloutStatus(rolloutStatus string) string {
	if rolloutStatus == model.RolloutStatusInProgress {
		return fmt.Sprint(ctc.ForegroundBrightCyan, rolloutStatus, ctc.Reset)
	}
	if rolloutStatus == model.RolloutStatusFailed {
		return fmt.Sprint(ctc.ForegroundBrightRed, rolloutStatus, ctc.Reset)
	}

	return rolloutStatus
}
