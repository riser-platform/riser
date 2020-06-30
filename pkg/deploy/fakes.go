package deploy

import (
	"github.com/riser-platform/riser-server/api/v1/model"
)

type fakeAppsClient struct {
	GetStatusFn        func(name, namespace string) (*model.AppStatus, error)
	GetStatusCallCount int
}

func (fake *fakeAppsClient) List() ([]model.App, error) {
	panic("NI")
}

func (fake *fakeAppsClient) Create(newApp *model.NewApp) (*model.App, error) {
	panic("NI")
}

func (fake *fakeAppsClient) Get(name, namespace string) (*model.App, error) {
	panic("NI")
}

func (fake *fakeAppsClient) GetStatus(name, namespace string) (*model.AppStatus, error) {
	fake.GetStatusCallCount++
	return fake.GetStatusFn(name, namespace)
}
