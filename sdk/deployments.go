package sdk

import (
	"net/http"

	"github.com/tshak/riser-server/api/v1/model"
)

type DeploymentsClient interface {
	Save(deployment *model.DeploymentRequest, dryRun bool) (string, error)
}

type deploymentsClient struct {
	client *Client
}

func (c *deploymentsClient) Save(deployment *model.DeploymentRequest, dryRun bool) (string, error) {
	request, err := c.client.NewRequest(http.MethodPut, "/api/v1/deployments", deployment)
	if err != nil {
		return "", err
	}

	if dryRun {
		q := request.URL.Query()
		q.Add("dryRun", "true")
		request.URL.RawQuery = q.Encode()
	}

	responseModel := &model.APIResponse{}
	_, err = c.client.Do(request, responseModel)
	if err != nil {
		return "", err
	}

	// TODO: API should return a proper model here for dryrun support
	if dryRun {
		return "Dry run not yet implemented", nil
	} else {
		return responseModel.Message, nil
	}
}
