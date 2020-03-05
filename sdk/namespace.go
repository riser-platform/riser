package sdk

import (
	"net/http"

	"github.com/riser-platform/riser-server/api/v1/model"
)

type NamespacesClient interface {
	List() ([]model.Namespace, error)
	Create(namespaceName string) error
}

type namespacesClient struct {
	client *Client
}

func (c *namespacesClient) List() ([]model.Namespace, error) {
	namespaces := []model.Namespace{}
	request, err := c.client.NewGetRequest("/api/v1/namespaces")
	if err != nil {
		return nil, err
	}
	_, err = c.client.Do(request, &namespaces)
	if err != nil {
		return nil, err
	}
	return namespaces, nil
}

func (c *namespacesClient) Create(namespaceName string) error {
	namespace := &model.Namespace{Name: model.NamespaceName(namespaceName)}
	request, err := c.client.NewRequest(http.MethodPost, "/api/v1/namespaces", namespace)
	if err != nil {
		return err
	}
	_, err = c.client.Do(request, nil)
	if err != nil {
		return err
	}

	return nil
}
