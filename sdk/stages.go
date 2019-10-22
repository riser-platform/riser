package sdk

import (
	"fmt"
	"net/http"

	"github.com/riser-platform/riser-server/api/v1/model"
)

type StagesClient interface {
	Ping(stageName string) error
	List() ([]model.StageMeta, error)
	SetConfig(stageName string, config *model.StageConfig) error
}

type stagesClient struct {
	client *Client
}

func (c *stagesClient) Ping(stageName string) error {
	request, err := c.client.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/stages/%s/ping", stageName), nil)
	if err != nil {
		return err
	}

	_, err = c.client.Do(request, nil)
	return err
}

func (c *stagesClient) List() ([]model.StageMeta, error) {
	request, err := c.client.NewGetRequest("/api/v1/stages")
	if err != nil {
		return nil, err
	}

	stages := []model.StageMeta{}
	_, err = c.client.Do(request, &stages)
	if err != nil {
		return nil, err
	}

	return stages, nil
}

// SetConfig sets configuration for a stage. Empty values are ignored and merged with existing config values.
func (c *stagesClient) SetConfig(stageName string, config *model.StageConfig) error {
	request, err := c.client.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/stages/%s/config", stageName), config)
	if err != nil {
		return err
	}

	_, err = c.client.Do(request, nil)
	return err
}
