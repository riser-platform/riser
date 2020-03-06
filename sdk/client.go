package sdk

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

const (
	defaultContentType = "application/json"
	defaultAccept      = "application/json"
	// TODO: Add version here
	userAgent = "risercli"
)

// Client is a API v1 client
type Client struct {
	BaseURL *url.URL
	apikey  string
	client  *http.Client

	// Model clients
	Apps        AppsClient
	Deployments DeploymentsClient
	Namespaces  NamespacesClient
	Rollouts    RolloutsClient
	Secrets     SecretsClient
	Stages      StagesClient
	Validate    ValidateClient
}

func NewClient(baseURI string, apikey string) (*Client, error) {
	baseURIParsed, err := url.Parse(baseURI)
	if err != nil {
		return nil, err
	}

	client := &Client{BaseURL: baseURIParsed, apikey: apikey}
	client.client = &http.Client{}

	client.Apps = &appsClient{client}
	client.Deployments = &deploymentsClient{client}
	client.Namespaces = &namespacesClient{client}
	client.Rollouts = &rolloutsClient{client}
	client.Secrets = &secretsClient{client}
	client.Stages = &stagesClient{client}
	client.Validate = &validateClient{client}

	return client, nil
}

func (c *Client) MakeInsecure() {
	c.client = &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
}

func (c *Client) NewGetRequest(relativeUrl string) (*http.Request, error) {
	return c.NewRequest(http.MethodGet, relativeUrl, nil)
}

func (c *Client) NewRequest(method, relativeUrl string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(relativeUrl)
	if err != nil {
		return nil, err
	}

	bodyBuffer := &bytes.Buffer{}
	if body != nil {
		err = json.NewEncoder(bodyBuffer).Encode(body)
		if err != nil {
			return nil, errors.Wrap(err, "Error marshaling json")
		}
	}

	req, err := http.NewRequest(method, c.BaseURL.ResolveReference(rel).String(), bodyBuffer)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", defaultAccept)
	req.Header.Add(authorizationHeader(c.apikey))
	req.Header.Add("Content-Type", defaultContentType)
	req.Header.Add("User-Agent", userAgent)
	return req, nil
}

func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	response, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	err = validateResponse(response)
	if err != nil {
		return response, err
	}

	if v != nil {
		// TODO: Try to use same error handling logic from validateResponse and return a ClientError instead
		responseBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return response, err
		}
		err = json.Unmarshal(responseBytes, v)
		if err != nil {
			return response, errors.Wrap(err, string(responseBytes))
		}
	}

	return response, err
}

func authorizationHeader(apikey string) (string, string) {
	return "Authorization", fmt.Sprintf("Apikey: %s", apikey)
}

func validateResponse(response *http.Response) error {
	if isSuccess(response.StatusCode) {
		return nil
	}

	clientErr := &ClientError{StatusCode: response.StatusCode}
	// TODO: Handle 404 in a better way (it's not really an error!)
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		clientErr.Message = fmt.Sprintf("Unable to read response: %s", err)
	}
	err = json.Unmarshal(responseBytes, clientErr)
	if err != nil {
		clientErr.Message = fmt.Sprintf("Unable to parse response: %s", responseBytes)
	}

	return clientErr
}

func isSuccess(statusCode int) bool {
	return statusCode >= 200 && statusCode <= 299
}
