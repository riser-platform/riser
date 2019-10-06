package sdk

import (
	"net/http"

	"github.com/tshak/riser-server/api/v1/model"
)

type DeploymentsClient interface {
	Save(deployment *model.DeploymentRequest, dryRun bool) (*model.DeploymentResponse, error)
}

type deploymentsClient struct {
	client *Client
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
