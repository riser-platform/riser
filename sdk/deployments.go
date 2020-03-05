package sdk

import (
	"fmt"
	"net/http"

	"github.com/riser-platform/riser-server/api/v1/model"
)

type DeploymentsClient interface {
	Delete(deploymentName, namespace, stageName string) (*model.DeploymentResponse, error)
	Save(deployment *model.DeploymentRequest, dryRun bool) (*model.DeploymentResponse, error)
	SaveStatus(deploymentName, namespace, stageName string, status *model.DeploymentStatusMutable) error
}

type deploymentsClient struct {
	client *Client
}

func (c *deploymentsClient) Delete(deploymentName, namespace, stageName string) (*model.DeploymentResponse, error) {
	request, err := c.client.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/deployments/%s/%s/%s", stageName, namespace, deploymentName), nil)
	if err != nil {
		return nil, err
	}

	responseModel := &model.DeploymentResponse{}
	_, err = c.client.Do(request, responseModel)
	if err != nil {
		return nil, err
	}

	return responseModel, nil
}

func (c *deploymentsClient) Save(deployment *model.DeploymentRequest, dryRun bool) (*model.DeploymentResponse, error) {
	request, err := c.client.NewRequest(http.MethodPut, "/api/v1/deployments", deployment)
	if err != nil {
		return nil, err
	}

	if dryRun {
		q := request.URL.Query()
		q.Add("dryRun", "true")
		request.URL.RawQuery = q.Encode()
	}

	responseModel := &model.DeploymentResponse{}
	_, err = c.client.Do(request, responseModel)
	if err != nil {
		return nil, err
	}

	return responseModel, nil
}

func (c *deploymentsClient) SaveStatus(deploymentName, namespace, stageName string, status *model.DeploymentStatusMutable) error {
	request, err := c.client.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/deployments/%s/%s/%s/status", stageName, namespace, deploymentName), status)
	if err != nil {
		return err
	}
	_, err = c.client.Do(request, nil)
	return err
}
