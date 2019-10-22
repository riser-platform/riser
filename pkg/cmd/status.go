package cmd

import (
	"fmt"
	"riser/pkg/rc"
	"riser/pkg/ui"
	"riser/pkg/ui/table"
	"strings"

	"github.com/riser-platform/riser-server/api/v1/model"

	"github.com/spf13/cobra"
	"github.com/wzshiming/ctc"
)

func newStatusCommand(currentContext *rc.Context) *cobra.Command {
	var appName string
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Gets the status for a deployment.",
		Run: func(cmd *cobra.Command, args []string) {
			riserClient := getRiserClient(currentContext)

			status, err := riserClient.Status.Get(appName)
			ui.ExitIfErrorMsg(err, "Error getting status")

			drawStatus(status)
		},
	}

	addAppFlag(cmd.Flags(), &appName)

	return cmd
}

func drawStatus(status *model.Status) {
	table := table.Default().Header("Deployment", "Stage", "Rev", "Docker Tag", "Rollout", "Rollout Details", "Problems")
	for _, deploymentStatus := range status.Deployments {
		table.AddRow(
			deploymentStatus.DeploymentName,
			deploymentStatus.StageName,
			fmt.Sprintf("%d", deploymentStatus.RolloutRevision),
			getDockerTag(deploymentStatus.DockerImage),
			formatRolloutStatus(deploymentStatus.RolloutStatus),
			deploymentStatus.RolloutStatusReason,
			formatProblems(deploymentStatus.Problems))
	}

	fmt.Println(table)
	fmt.Print("\n")

	for _, stageStatus := range status.Stages {
		if !stageStatus.Healthy {
			fmt.Print(ctc.ForegroundBrightYellow)
			fmt.Printf("Warning: stage %q is not healthy. %s\n", stageStatus.StageName, stageStatus.Reason)
			fmt.Print(ctc.Reset)
		}
	}

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
