package sdk

import (
	"net/http"

	"github.com/riser-platform/riser-server/api/v1/model"
)

type ValidateClient interface {
	AppConfig(appConfig *model.AppConfigWithOverrides) error
}

type validateClient struct {
	client *Client
}

func (c *validateClient) AppConfig(appConfig *model.AppConfigWithOverrides) error {
	request, err := c.client.NewRequest(http.MethodPost, "/api/v1/validate/appconfig", appConfig)
	if err != nil {
		return err
	}

	_, err = c.client.Do(request, nil)
	if err != nil {
		return err
	}

	return nil
}
