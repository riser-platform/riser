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

func newStatusCommand(currentContext *rc.Context) *cobra.Command {
	var appName string
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Gets the status for a deployment.",
		Run: func(cmd *cobra.Command, args []string) {
			apiClient, err := sdk.NewClient(currentContext.ServerURL, currentContext.Apikey)
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

func drawStatusSummary(statuses []model.DeploymentStatus) {
	table := table.Default().Header("Deployment", "Stage", "Rev", "Docker Tag", "Rollout", "Rollout Details", "Problems")

	for _, status := range statuses {
		table.AddRow(
			status.DeploymentName,
			status.StageName,
			fmt.Sprintf("%d", status.RolloutRevision),
			getDockerTag(status.DockerImage),
			formatRolloutStatus(status.RolloutStatus),
			status.RolloutStatusReason,
			formatProblems(status.Problems))
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

func formatProblems(problems []model.DeploymentStatusProblem) string {
	if len(problems) == 0 {
		return fmt.Sprint(ctc.ForegroundBrightGreen, "None Found", ctc.Reset)
	}

	message := ""
	first := true
	for _, problem := range problems {
		newline := "\n"
		if first {
			newline = ""
			first = false
		}
		message = fmt.Sprintf("%s%s%s", message, newline, formatProblem(problem))
	}

	return fmt.Sprint(ctc.ForegroundBrightRed, message, ctc.Reset)
}

func formatProblem(problem model.DeploymentStatusProblem) string {
	if problem.Count == 1 {
		return problem.Message
	}
	return fmt.Sprintf("(x%d) %s", problem.Count, problem.Message)
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
