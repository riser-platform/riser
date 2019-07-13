package sdk

import (
	"encoding/json"
	"fmt"

	"github.com/tshak/riser-server/api/v1/model"
)

func (client *Client) ListSecretMetas(appName, stageName string) ([]model.SecretMetaStatus, error) {
	secretMetas := []model.SecretMetaStatus{}
	apiUri := client.uri(fmt.Sprintf("/api/v1/secrets/%s/%s", appName, stageName))
	response, err := doGet(apiUri.String())
	if err != nil {
		return nil, err
	}
	err = unmarshal(response, &secretMetas)
	if err != nil {
		return nil, err
	}
	return secretMetas, nil
}

func (client *Client) SaveSecret(appName, stageName, secretName, plainTextSecret string) error {
	secretJson, err := json.Marshal(model.UnsealedSecret{
		SecretMeta: model.SecretMeta{
			App:   appName,
			Stage: stageName,
			Name:  secretName,
		},
		PlainText: plainTextSecret,
	})
	if err != nil {
		return err
	}
	apiUri := client.uri("/api/v1/secrets")

	_, err = doBodyRequest(apiUri.String(), defaultContentType, "PUT", secretJson)
	return err
}
