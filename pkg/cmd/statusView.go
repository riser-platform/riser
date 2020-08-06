package cmd

import (
	"fmt"
	"io"
	"riser/pkg/status"
	"riser/pkg/ui"
	"riser/pkg/ui/style"
	"riser/pkg/ui/table"
	"strings"

	"github.com/riser-platform/riser-server/api/v1/model"
)

type statusView struct {
	appName             string
	status              *model.AppStatus
	activeRevisionsOnly bool
}

func (view *statusView) RenderHuman(writer io.Writer) error {
	outStr := ""
	if len(view.status.Deployments) == 0 {
		outStr += fmt.Sprintf("There are no deployments for the app %q. Use \"riser deploy\" to make your first deployment.\n", view.appName)
	} else {

		statusTable := table.Default().Header("Deployment", "Env", "Traffic", "Rev", "Docker Tag", "Status", "Reason")

		for _, deploymentStatus := range view.status.Deployments {
			revisions := status.GetRevisionStatus(&deploymentStatus, view.activeRevisionsOnly)

			if len(revisions) > 0 {
				first := true
				for _, activeRevision := range revisions {
					if first {
						statusTable.AddRow(
							deploymentStatus.DeploymentName,
							deploymentStatus.EnvironmentName,
							formatTraffic(&activeRevision.Traffic),
							fmt.Sprintf("%d", activeRevision.RiserRevision),
							formatDockerTag(activeRevision.DockerImage),
							formatRevisionStatus(activeRevision.RevisionStatus),
							activeRevision.RevisionStatusReason,
						)
					} else {
						statusTable.AddRow(
							"", "",
							formatTraffic(&activeRevision.Traffic),
							fmt.Sprintf("%d", activeRevision.RiserRevision),
							formatDockerTag(activeRevision.DockerImage),
							formatRevisionStatus(activeRevision.RevisionStatus),
							activeRevision.RevisionStatusReason,
						)
					}
					first = false
				}
			}
		}

		outStr += statusTable.String()
		outStr += "\n\n"

		for _, environmentStatus := range view.status.Environments {
			if !environmentStatus.Healthy {
				outStr += style.Warn(fmt.Sprintf("Warning: environment %q is not healthy. %s\n", environmentStatus.EnvironmentName, environmentStatus.Reason))
			}
		}
	}
	_, err := writer.Write([]byte(outStr))
	return err
}

func (view *statusView) RenderJson(writer io.Writer) error {
	return ui.RenderJson(view.status, writer)
}

func formatTraffic(traffic *model.DeploymentTrafficStatus) string {
	// TODO: Determine if % is ever nil in practice and display as 100% if latest and only active revision
	if traffic.Percent != nil {
		return fmt.Sprintf("%d%%", *traffic.Percent)
	}

	return "0%"
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
		formatted = style.Emphasis(rolloutStatus)
	}

	return formatted
}
