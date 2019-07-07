package cmd

import (
	"riser/ui/table"
	"fmt"
	"riser/rc"

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
	table := table.Default().Header("Deployment", "Stage", "Rollout", "Healthy")

	for _, status := range statuses {
		table.AddRow(
			status.DeploymentName,
			status.StageName,
			formatRolloutStatus(status.RolloutStatus),
			formatHealthStatus(status.HealthStatus))
		}

	fmt.Println(table)
}

func formatHealthStatus(healthStatus string) string {
	if healthStatus == model.HealthStatusTrue {
		return fmt.Sprint(ctc.ForegroundBrightGreen, healthStatus, ctc.Reset)
	}
	if healthStatus == model.HealthStatusFalse {
		return fmt.Sprint(ctc.ForegroundBrightRed, healthStatus, ctc.Reset)
	}
	if healthStatus == model.HealthStatusUnknown {
		return fmt.Sprint(ctc.ForegroundBrightYellow, healthStatus, ctc.Reset)
	}

	return healthStatus
}

func formatRolloutStatus(rolloutStatus string) string {
	if rolloutStatus == model.RolloutStatusInProgress {
		return fmt.Sprint(ctc.ForegroundBrightCyan, rolloutStatus, ctc.Reset)
	}
	if rolloutStatus == model.RolloutStatusFailed {
		return fmt.Sprint(ctc.ForegroundBrightRed, rolloutStatus, ctc.Reset)
	}

	return rolloutStatus
}
