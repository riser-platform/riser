package sdk

import (
	"fmt"
	"net/http"

	"github.com/riser-platform/riser-server/api/v1/model"
)

type AppsClient interface {
	List() ([]model.App, error)
	Create(newApp *model.NewApp) (*model.App, error)
	Get(name, namespace string) (*model.App, error)
	GetStatus(name, namespace string) (*model.AppStatus, error)
}

type appsClient struct {
	client *Client
}

func (c *appsClient) Get(name, namespace string) (*model.App, error) {
	app := &model.App{}
	request, err := c.client.NewGetRequest(fmt.Sprintf("/api/v1/apps/%s/%s", namespace, name))
	if err != nil {
		return nil, err
	}
	_, err = c.client.Do(request, &app)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (c *appsClient) List() ([]model.App, error) {
	apps := []model.App{}
	request, err := c.client.NewGetRequest("/api/v1/apps")
	if err != nil {
		return nil, err
	}
	_, err = c.client.Do(request, &apps)
	if err != nil {
		return nil, err
	}
	return apps, nil
}

func (c *appsClient) Create(newApp *model.NewApp) (*model.App, error) {
	request, err := c.client.NewRequest(http.MethodPost, "/api/v1/apps", newApp)
	if err != nil {
		return nil, err
	}

	app := &model.App{}
	_, err = c.client.Do(request, &app)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (c *appsClient) GetStatus(name, namespace string) (*model.AppStatus, error) {
	request, err := c.client.NewGetRequest(fmt.Sprintf("/api/v1/apps/%s/%s/status", namespace, name))
	if err != nil {
		return nil, err
	}

	status := &model.AppStatus{}
	_, err = c.client.Do(request, status)
	if err != nil {
		return nil, err
	}
	return status, nil
}
