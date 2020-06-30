package deploy

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/riser-platform/riser-server/pkg/sdk"
)

type isReadyFunc func(statuses []model.DeploymentStatus, deploymentName string, environmentName string, riserRevision int64) (bool, string)

// WaitForReady waits for a deployment to become ready for a specified riserRevision. It returns an error
// if the deployment's revision is not ready within the specified timeout.
func WaitForReady(apps sdk.AppsClient, app model.App, deploymentName string, environmentName string, riserRevision int64, timeout time.Duration) error {
	return waitForReady(isReady, apps, app, deploymentName, environmentName, riserRevision, timeout)
}

func waitForReady(isReady isReadyFunc, apps sdk.AppsClient, app model.App, deploymentName string, environmentName string, riserRevision int64, timeout time.Duration) error {
	var resultErr error

	done := make(chan bool)
	start := time.Now()

	go func() {
		for {
			if time.Since(start) >= timeout {
				err := fmt.Errorf("Timeout of %s exceeded waiting for the new revision to become ready", timeout)
				if resultErr == nil {
					resultErr = err
				} else {
					resultErr = errors.Wrap(resultErr, err.Error())
				}
				close(done)
				break
			}
			// Note: Since the SDK does not currently support context, a slow request can cause the wait to take longer than the timeout
			appStatus, err := apps.GetStatus(string(app.Name), string(app.Namespace))
			// If we timeout we want to capture the last error.
			resultErr = err
			if err != nil {
				continue
			}
			isReadyResult, reason := isReady(appStatus.Deployments, deploymentName, environmentName, riserRevision)
			if isReadyResult {
				// We don't care about the last error since we've successfully received the status
				resultErr = nil
				close(done)
				break
			}
			resultErr = errors.New(fmt.Sprintf("Revision status is %q", reason))

			time.Sleep(1 * time.Second)
		}
	}()

	<-done

	return resultErr
}

// isReady determines if a deployment at a specific revision is ready from a set of deployment statuses.
// Assumes that all status's are from the same app.
func isReady(statuses []model.DeploymentStatus, deploymentName string, environmentName string, riserRevision int64) (ready bool, reason string) {
	lastRevisionStatus := ""
	lastRevisionReason := ""
	for _, status := range statuses {
		if status.DeploymentName != deploymentName || status.EnvironmentName != environmentName {
			continue
		}

		if riserRevision > status.ObservedRiserRevision {
			return false, "The revision has not yet been observed"
		}

		for _, revStatus := range status.Revisions {
			if revStatus.RiserRevision == riserRevision && revStatus.RevisionStatus == model.RevisionStatusReady {
				return true, revStatus.RevisionStatus
			}
			lastRevisionStatus = revStatus.RevisionStatus
			lastRevisionReason = revStatus.RevisionStatusReason
		}
	}
	reason = lastRevisionStatus
	if lastRevisionReason != "" {
		reason = fmt.Sprintf("%s (%s)", lastRevisionStatus, lastRevisionReason)
	}
	return false, reason
}
