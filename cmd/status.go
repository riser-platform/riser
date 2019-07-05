package cmd

import (
	"fmt"

	"github.com/tshak/riser-server/api/v1/model"

	"github.com/alexeyco/simpletable"
	"github.com/spf13/cobra"
	"github.com/tshak/riser/sdk"
	"github.com/wzshiming/ctc"
)

func newStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status (app name)",
		Short: "Gets the status for a deployment.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			appName := args[0]
			apiClient, err := sdk.NewClient("http://localhost:8000")
			if err != nil {
				panic(err)
			}

			statuses, err := apiClient.GetStatus(appName)

			if err != nil {
				panic(err)
			}

			table := simpletable.New()
			table.SetStyle(simpletable.StyleCompactLite)
			table.Header = &simpletable.Header{
				Cells: []*simpletable.Cell{
					defaultCell("Deployment"),
					defaultCell("Stage"),
					defaultCell("Rollout"),
					defaultCell("Healthy"),
				},
			}

			for _, status := range statuses {
				row := []*simpletable.Cell{
					defaultCell(status.DeploymentName),
					defaultCell(status.StageName),
					defaultCell(formatRolloutStatus(status.RolloutStatus)),
					defaultCell(formatHealthStatus(status.HealthStatus)),
				}
				table.Body.Cells = append(table.Body.Cells, row)
			}

			fmt.Println(table.String())
		},
	}
}

func defaultCell(text string) *simpletable.Cell {
	return &simpletable.Cell{Align: simpletable.AlignLeft, Text: text}
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
