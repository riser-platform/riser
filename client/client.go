package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

const defaultContentType = "application/json"

type Client struct {
	baseURI *url.URL
}

func NewClient(baseURI string) (*Client, error) {
	baseURIParsed, err := url.Parse(baseURI)
	if err != nil {
		return nil, err
	}

	return &Client{baseURI: baseURIParsed}, nil
}

func (client *Client) PutDeployment(deployment Deployment, dryRun bool) error {
	deploymentJson, err := json.Marshal(deployment)
	if err != nil {
		return err
	}
	apiUri := client.uri("/api/v1/deployments")
	if dryRun {
		q := apiUri.Query()
		q.Add("dryRun", "true")
		apiUri.RawQuery = q.Encode()
	}

	fmt.Println(apiUri.String())
	responseBody, err := doPut(apiUri.String(), defaultContentType, deploymentJson)
	if err != nil {
		return err
	}
	fmt.Printf("%s", responseBody)
	return nil
}

func (client *Client) uri(apiPath string) url.URL {
	apiURI := *client.baseURI
	apiURI.Path = path.Join(apiURI.Path, apiPath)
	return apiURI
}

func doPut(uri, contentType string, body []byte) ([]byte, error) {
	client := &http.Client{}
	request, err := http.NewRequest("PUT", uri, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", contentType)

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	} else {
		defer response.Body.Close()
		return ioutil.ReadAll(response.Body)
	}
}
