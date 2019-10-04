package sdk

import (
	"net/http"

	"github.com/tshak/riser-server/api/v1/model"
)

type StatusClient interface {
	Save(status *model.DeploymentStatus) error
	Get(appName string) (*model.Status, error)
}

type statusClient struct {
	client *Client
}

func (c *statusClient) Save(status *model.DeploymentStatus) error {
	request, err := c.client.NewRequest(http.MethodPut, "/api/v1/status", status)
	if err != nil {
		return err
	}
	_, err = c.client.Do(request, nil)
	return err
}

func (c *statusClient) Get(appName string) (*model.Status, error) {
	request, err := c.client.NewGetRequest("/api/v1/status")
	if err != nil {
		return nil, err
	}

	q := request.URL.Query()
	q.Add("app", appName)
	request.URL.RawQuery = q.Encode()

	status := &model.Status{}
	_, err = c.client.Do(request, status)
	if err != nil {
		return nil, err
	}
	return status, nil
}
