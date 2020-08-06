package cmd

import (
	"fmt"
	"io"
	"riser/pkg/status"
	"riser/pkg/ui/style"
	"riser/pkg/ui/table"

	"github.com/riser-platform/riser-server/api/v1/model"
)

const deploymentsDescribeMaxRevisions = 5

type deploymentDescribeView struct {
	deploymentName    string
	externalHost      string
	app               *model.App
	environmentStatus *model.EnvironmentStatus
	deploymentStatus  []model.DeploymentStatus
}

func newDeploymentsDescribeView(app *model.App, appStatus *model.AppStatus, deploymentName, envName, externalHost string) (*deploymentDescribeView, error) {
	view := &deploymentDescribeView{
		app:            app,
		deploymentName: deploymentName,
		externalHost:   externalHost,
	}
	for _, envStatus := range appStatus.Environments {
		if envStatus.EnvironmentName == envName {
			view.environmentStatus = &envStatus
			break
		}
	}

	view.deploymentStatus = []model.DeploymentStatus{}
	for _, deploymentStatus := range appStatus.Deployments {
		if deploymentStatus.DeploymentName == deploymentName && deploymentStatus.EnvironmentName == envName {
			view.deploymentStatus = append(view.deploymentStatus, deploymentStatus)
		}
	}
	if len(view.deploymentStatus) == 0 {
		return nil, fmt.Errorf(`The environment "%s" does not contain the deployment "%s" in the "%s" namespace`, envName, deploymentName, app.Namespace)
	}
	return view, nil
}

func (view *deploymentDescribeView) RenderHuman(writer io.Writer) error {
	outStr := ""
	outStr += fmt.Sprintf("Name: %s\n", view.deploymentName)
	outStr += fmt.Sprintf("Namespace: %s\n", view.app.Namespace)
	outStr += fmt.Sprintf("Environment: %s\n", view.environmentStatus.EnvironmentName)
	outStr += fmt.Sprintf("App: %s (%s)\n", view.app.Name, view.app.Id)

	outStr += "\nIngress URLs:\n"
	// TODO: This can be confusing as it's possible that an app is not exposed externally
	outStr += fmt.Sprintf("  External: %s\n", formatExternalUrl(view.deploymentName, string(view.app.Namespace), view.externalHost))
	outStr += fmt.Sprintf("  Cluster: %s\n", formatClusterLocalUrl(view.deploymentName, string(view.app.Namespace)))
	outStr += fmt.Sprintf("\nRecent %d Revisions:\n", deploymentsDescribeMaxRevisions)
	statusTable := table.Default().Header("Traffic", "Rev", "Docker Tag", "Status", "Reason")

	for _, deploymentStatus := range view.deploymentStatus {
		revisions := status.GetRevisionStatus(&deploymentStatus, false)

		if len(revisions) > 0 {
			for revIdx, revision := range revisions {
				statusTable.AddRow(
					formatTraffic(&revision.Traffic),
					fmt.Sprintf("%d", revision.RiserRevision),
					formatDockerTag(revision.DockerImage),
					formatRevisionStatus(revision.RevisionStatus),
					revision.RevisionStatusReason,
				)
				if revIdx+1 >= deploymentsDescribeMaxRevisions {
					break
				}
			}
		}
	}
	outStr += statusTable.String()
	outStr += "\n\n"
	if !view.environmentStatus.Healthy {
		outStr += style.Warn(fmt.Sprintf("Warning: environment %q is not healthy. %s\n", view.environmentStatus.EnvironmentName, view.environmentStatus.Reason))
	}
	_, err := writer.Write([]byte(outStr))
	return err
}

func (view *deploymentDescribeView) RenderJson(writer io.Writer) error {
	return nil
}

func formatExternalUrl(deploymentName, namespace, externalHost string) string {
	return fmt.Sprintf("https://%s.%s.%s", deploymentName, namespace, externalHost)
}

func formatClusterLocalUrl(deploymentName, namespace string) string {
	return fmt.Sprintf("http://%s.%s.svc.cluster.local", deploymentName, namespace)
}
