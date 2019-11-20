package cmd

import (
	"fmt"
	"riser/pkg/rc"
	"riser/pkg/status"
	"riser/pkg/ui"
	"riser/pkg/ui/style"
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

			status, err := riserClient.Apps.GetStatus(appName)
			ui.ExitIfErrorMsg(err, "Error getting status")

			drawStatus(appName, status)
		},
	}

	addAppFlag(cmd.Flags(), &appName)

	return cmd
}

func drawStatus(appName string, appStatus *model.AppStatus) {
	if len(appStatus.Deployments) == 0 {
		fmt.Printf("There are no deployments for the app %q. Use \"riser deploy\" to make your first deployment.\n", appName)
		return
	}
	table := table.Default().Header("Deployment", "Stage", "Rev", "Docker Tag", "Replicas", "Problems")
	deploymentsPendingObservation := false
	for _, deploymentStatus := range appStatus.Deployments {
		if !deploymentObserved(deploymentStatus) {
			deploymentsPendingObservation = true
		}

		revision := status.GetLatestReadyRevision(&deploymentStatus)
		if len(deploymentStatus.Revisions) > 0 {
			table.AddRow(
				formatDeploymentName(deploymentStatus),
				deploymentStatus.StageName,
				fmt.Sprintf("%d", revision.RiserGeneration),
				formatDockerTag(revision.DockerImage),
				fmt.Sprintf("%d", revision.AvailableReplicas),
				formatProblems(deploymentStatus.Problems))
		}
	}

	fmt.Println(table)
	fmt.Print("\n")

	if deploymentsPendingObservation {
		fmt.Println(style.Emphasis("* This deployment has changes that have not yet been observed."))
	}

	for _, stageStatus := range appStatus.Stages {
		if !stageStatus.Healthy {
			fmt.Print(ctc.ForegroundBrightYellow)
			fmt.Printf("Warning: stage %q is not healthy. %s\n", stageStatus.StageName, stageStatus.Reason)
			fmt.Print(ctc.Reset)
		}
	}
}

func formatDeploymentName(deploymentStatus model.DeploymentStatus) string {
	name := deploymentStatus.DeploymentName
	if !deploymentObserved(deploymentStatus) {
		name = style.Emphasis("*") + name

	}
	return name
}

func deploymentObserved(deploymentStatus model.DeploymentStatus) bool {
	return deploymentStatus.RiserGeneration == deploymentStatus.ObservedRiserGeneration
}

func formatDockerTag(dockerImage string) string {
	idx := strings.Index(dockerImage, ":")
	if idx == -1 {
		return style.Warn("Unknown")
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
	formatted := rolloutStatus
	switch rolloutStatus {
	case model.RolloutStatusInProgress:
		formatted = style.Emphasis(rolloutStatus)
	case model.RolloutStatusFailed:
		formatted = style.Bad(rolloutStatus)
	case model.RolloutStatusUnknown:
		formatted = style.Warn(rolloutStatus)
	}

	return formatted
}
