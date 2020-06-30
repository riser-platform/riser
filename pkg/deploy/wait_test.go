package deploy

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/riser-platform/riser-server/api/v1/model"
)

// revision is a compact version of model.DeploymentRevisionStatus for testing
type revision struct {
	riserRevision        int64
	revisionStatus       string
	revisionStatusReason string
}

func Test_waitForReady_polls(t *testing.T) {
	app := model.App{
		Name:      model.AppName("myapp"),
		Namespace: model.NamespaceName("apps"),
	}

	// We return an empty status every time and stub out isReady as it's easier to test separately
	returnedStatuses := []model.DeploymentStatus{}

	apps := &fakeAppsClient{
		GetStatusFn: func(name, namespace string) (*model.AppStatus, error) {
			assert.Equal(t, "myapp", name)
			assert.Equal(t, "apps", namespace)
			return &model.AppStatus{Deployments: returnedStatuses}, nil
		},
	}

	fakeIsReadyResults := []bool{false, false, true}
	fakeIsReadyCallIdx := 0
	fakeIsReady := func(statuses []model.DeploymentStatus, deploymentName string, environmentName string, riserRevision int64) (bool, string) {
		assert.Equal(t, statuses, returnedStatuses)
		assert.Equal(t, "mydep", deploymentName)
		assert.Equal(t, "myenv", environmentName)
		assert.EqualValues(t, 1, riserRevision)

		isReady := fakeIsReadyResults[fakeIsReadyCallIdx]
		fakeIsReadyCallIdx++
		return isReady, ""
	}

	start := time.Now()
	err := waitForReady(fakeIsReady, apps, app, "mydep", "myenv", 1, 3*time.Second)
	elapsed := time.Since(start)

	assert.NoError(t, err)
	// Ensure that we're actually pausing between retries.
	// TODO Add a fake clock to speed up this test and make it more deterministic
	assert.Equal(t, 2*time.Second, elapsed.Round(time.Second))
	assert.Equal(t, apps.GetStatusCallCount, 3)
}

func Test_waitForReady_ReturnsErrorAfterTimeout(t *testing.T) {
	app := model.App{
		Name:      model.AppName("myapp"),
		Namespace: model.NamespaceName("apps"),
	}

	// We return an empty status every time and stub out isReady as it's easier to test separately
	returnedStatuses := []model.DeploymentStatus{}

	apps := &fakeAppsClient{
		GetStatusFn: func(name, namespace string) (*model.AppStatus, error) {
			return &model.AppStatus{Deployments: returnedStatuses}, nil
		},
	}

	fakeIsReady := func(statuses []model.DeploymentStatus, deploymentName string, environmentName string, riserRevision int64) (bool, string) {
		return false, "Unhealthy"
	}

	err := waitForReady(fakeIsReady, apps, app, "mydep", "myenv", 1, 1*time.Second)

	assert.Equal(t, `Timeout of 1s exceeded waiting for the new revision to become ready: Revision status is "Unhealthy"`, err.Error())
}

func Test_waitForReady_RetriesAppsClientError(t *testing.T) {
	app := model.App{
		Name:      model.AppName("myapp"),
		Namespace: model.NamespaceName("apps"),
	}

	// We return an empty status every time and stub out isReady as it's easier to test separately
	returnedStatuses := []model.DeploymentStatus{}

	getStatusCallCount := 0
	apps := &fakeAppsClient{
		GetStatusFn: func(name, namespace string) (*model.AppStatus, error) {
			if getStatusCallCount > 0 {
				return &model.AppStatus{Deployments: returnedStatuses}, nil
			}
			getStatusCallCount++
			return nil, errors.New("busted")
		},
	}

	fakeIsReady := func(statuses []model.DeploymentStatus, deploymentName string, environmentName string, riserRevision int64) (bool, string) {
		return true, ""
	}

	err := waitForReady(fakeIsReady, apps, app, "mydep", "myenv", 1, 1*time.Second)

	assert.NoError(t, err)
}

func Test_waitForReady_ReturnsAppsClientErrorAfterTimeout(t *testing.T) {
	app := model.App{
		Name:      model.AppName("myapp"),
		Namespace: model.NamespaceName("apps"),
	}

	apps := &fakeAppsClient{
		GetStatusFn: func(name, namespace string) (*model.AppStatus, error) {
			return nil, errors.New("busted")
		},
	}

	fakeIsReady := func(statuses []model.DeploymentStatus, deploymentName string, environmentName string, riserRevision int64) (bool, string) {
		return false, ""
	}

	err := waitForReady(fakeIsReady, apps, app, "mydep", "myenv", 1, 100*time.Millisecond)

	assert.Equal(t, "Timeout of 100ms exceeded waiting for the new revision to become ready: busted", err.Error())
}

func Test_isReady(t *testing.T) {
	tests := []struct {
		test            string
		status          []model.DeploymentStatus
		deploymentName  string
		environmentName string
		riserRevision   int64
		expectedReady   bool
		expectedReason  string
	}{
		{
			test:   "No status",
			status: []model.DeploymentStatus{},
		},
		{
			test: "Not yet observed",
			status: []model.DeploymentStatus{
				makeTestDeploymentStatus("mydep", "dev", 1, 0),
			},
			deploymentName:  "mydep",
			environmentName: "dev",
			riserRevision:   1,
			expectedReason:  "The revision has not yet been observed",
		},
		{
			test: "Observed, revision is waiting",
			status: []model.DeploymentStatus{
				makeTestDeploymentStatus("mydep", "dev", 1, 1, revision{1, model.RevisionStatusWaiting, ""}),
			},
			deploymentName:  "mydep",
			environmentName: "dev",
			riserRevision:   1,
			expectedReason:  model.RevisionStatusWaiting,
		},
		{
			test: "Observed, revision is unhealthy with reason",
			status: []model.DeploymentStatus{
				makeTestDeploymentStatus("mydep", "dev", 1, 1, revision{1, model.RevisionStatusUnhealthy, "ImagePullError"}),
			},
			deploymentName:  "mydep",
			environmentName: "dev",
			riserRevision:   1,
			expectedReason:  "Unhealthy (ImagePullError)",
		},
		{
			test: "Observed, revision is ready",
			status: []model.DeploymentStatus{
				makeTestDeploymentStatus("mydep", "dev", 1, 1, revision{1, model.RevisionStatusReady, ""}),
			},
			deploymentName:  "mydep",
			environmentName: "dev",
			riserRevision:   1,
			expectedReady:   true,
			expectedReason:  model.RevisionStatusReady,
		},
		{
			test: "Observed, revision not ready, different env ready",
			status: []model.DeploymentStatus{
				makeTestDeploymentStatus("mydep", "dev", 1, 1, revision{1, model.RevisionStatusWaiting, ""}),
				makeTestDeploymentStatus("mydep", "prod", 1, 1, revision{1, model.RevisionStatusReady, ""}),
			},
			deploymentName:  "mydep",
			environmentName: "dev",
			riserRevision:   1,
			expectedReason:  model.RevisionStatusWaiting,
		},
		{
			test: "Observed, revision not ready, different deployment ready",
			status: []model.DeploymentStatus{
				makeTestDeploymentStatus("mydep", "dev", 1, 1, revision{1, model.RevisionStatusWaiting, ""}),
				makeTestDeploymentStatus("mydep-2", "dev", 1, 1, revision{1, model.RevisionStatusReady, ""}),
			},
			deploymentName:  "mydep",
			environmentName: "dev",
			riserRevision:   1,
			expectedReason:  model.RevisionStatusWaiting,
		},
	}

	for _, tt := range tests {
		result, reason := isReady(tt.status, tt.deploymentName, tt.environmentName, tt.riserRevision)
		assert.Equal(t, tt.expectedReady, result, fmt.Sprintf("%s (Ready)", tt.test))
		assert.Equal(t, tt.expectedReason, reason, fmt.Sprintf("%s (Reason)", tt.test))
	}
}

func makeTestDeploymentStatus(deploymentName, environmentName string, riserRevision int64, observedRiserRevision int64, revisionStatus ...revision) model.DeploymentStatus {
	status := model.DeploymentStatus{
		DeploymentName:  deploymentName,
		EnvironmentName: environmentName,
		RiserRevision:   riserRevision,
		DeploymentStatusMutable: model.DeploymentStatusMutable{
			ObservedRiserRevision: observedRiserRevision,
		},
	}

	for _, rev := range revisionStatus {
		status.Revisions = append(status.Revisions,
			model.DeploymentRevisionStatus{
				RiserRevision:        rev.riserRevision,
				RevisionStatus:       rev.revisionStatus,
				RevisionStatusReason: rev.revisionStatusReason,
			},
		)
	}

	return status
}
