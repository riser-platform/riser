package cmd

import (
	"fmt"

	"github.com/alexeyco/simpletable"
	"github.com/spf13/cobra"
	"github.com/tshak/riser-server/sdk"
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
					defaultCell("Healthy?"),
					defaultCell("Available Replicas"),
				},
			}

			for _, status := range statuses {
				row := []*simpletable.Cell{
					defaultCell(status.DeploymentName),
					defaultCell(status.StageName),
					defaultCell(formatHealthStatus(status.Healthy)),
					defaultCell(fmt.Sprintf("%d/%d", status.AvailableReplicas, status.DesiredReplicas)),
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

func formatHealthStatus(healthy string) string {
	if healthy == "true" {
		return fmt.Sprint(ctc.ForegroundBrightGreen, healthy, ctc.Reset)
	}
	if healthy == "false" {
		return fmt.Sprint(ctc.ForegroundBrightRed, healthy, ctc.Reset)
	}

	return healthy
}
