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

			drawStatus(appName, !showAllRevisions, status)
		},
	}

	addAppFlag(cmd.Flags(), &appName)
	addNamespaceFlag(cmd.Flags(), &namespace)
	cmd.Flags().BoolVarP(&showAllRevisions, "all-revisions", "", false, "Shows all available revisions. Otherwise only shows the latest revision and older revisions with traffic")

	return cmd
}

func drawStatus(appName string, activeRevisionsOnly bool, appStatus *model.AppStatus) {
	if len(appStatus.Deployments) == 0 {
		fmt.Printf("There are no deployments for the app %q. Use \"riser deploy\" to make your first deployment.\n", appName)
		return
	}
	statusTable := table.Default().Header("Deployment", "Stage", "Traffic", "Rev", "Docker Tag", "Pods", "Status", "Reason")
	deploymentsPendingObservation := false
	for _, deploymentStatus := range appStatus.Deployments {
		if !deploymentObserved(deploymentStatus) {
			deploymentsPendingObservation = true
		}

		revisions := status.GetRevisionStatus(&deploymentStatus, activeRevisionsOnly)
		if len(revisions) > 0 {
			first := true
			for _, activeRevision := range revisions {
				if first {
					statusTable.AddRow(
						formatDeploymentName(deploymentStatus),
						deploymentStatus.StageName,
						formatTraffic(&activeRevision.Traffic),
						fmt.Sprintf("%d", activeRevision.RiserRevision),
						formatDockerTag(activeRevision.DockerImage),
						fmt.Sprintf("%d", activeRevision.AvailableReplicas),
						formatRevisionStatus(activeRevision.RevisionStatus),
						activeRevision.RevisionStatusReason,
					)
				} else {
					statusTable.AddRow(
						"", "",
						formatTraffic(&activeRevision.Traffic),
						fmt.Sprintf("%d", activeRevision.RiserRevision),
						formatDockerTag(activeRevision.DockerImage),
						fmt.Sprintf("%d", activeRevision.AvailableReplicas),
						formatRevisionStatus(activeRevision.RevisionStatus),
						activeRevision.RevisionStatusReason,
					)
				}
				first = false
			}
		}

	}

	fmt.Println(statusTable)
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

func formatTraffic(traffic *model.DeploymentTrafficStatus) string {
	// TODO: Determine if % is ever nil in practice and display as 100% if latest and only active revision
	if traffic.Percent != nil {
		return fmt.Sprintf("%d%%", *traffic.Percent)
	}

	return "0%"
}

func formatDeploymentName(deploymentStatus model.DeploymentStatus) string {
	name := deploymentStatus.DeploymentName
	if !deploymentObserved(deploymentStatus) {
		name = style.Emphasis("*") + name

	}
	return name
}

func deploymentObserved(deploymentStatus model.DeploymentStatus) bool {
	return deploymentStatus.RiserRevision == deploymentStatus.ObservedRiserRevision
}

func formatDockerTag(dockerImage string) string {
	idx := strings.Index(dockerImage, ":")
	if idx == -1 {
		return style.Warn("Unknown")
	}
	return dockerImage[idx+1:]
}

func formatRevisionStatus(rolloutStatus string) string {
	formatted := rolloutStatus
	switch rolloutStatus {
	case model.RevisionStatusReady:
		formatted = style.Good(rolloutStatus)
	case model.RevisionStatusWaiting:
		formatted = style.Emphasis(rolloutStatus)
	case model.RevisionStatusUnhealthy:
		formatted = style.Bad(rolloutStatus)
	case model.RevisionStatusUnknown:
		formatted = style.Warn(rolloutStatus)
	}

	return formatted
}
