package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/pkg/errors"

	"github.com/tshak/riser-server/api/v1/model"
)

const defaultContentType = "application/json"

// Client is a API v1 client
type Client struct {
	baseURI *url.URL
}

// HttpError provides the status code.
type HttpError struct {
	StatusCode int
	Message    string
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("Received HTTP %d (%s) %s", e.StatusCode, http.StatusText(e.StatusCode), e.Message)
}

func NewClient(baseURI string) (*Client, error) {
	baseURIParsed, err := url.Parse(baseURI)
	if err != nil {
		return nil, err
	}

	return &Client{baseURI: baseURIParsed}, nil
}

func (client *Client) ListApps() ([]model.App, error) {
	apps := []model.App{}
	apiUri := client.uri("/api/v1/apps")
	response, err := doGet(apiUri.String())
	if err != nil {
		return nil, err
	}
	err = unmarshal(response, &apps)
	if err != nil {
		return nil, err
	}
	return apps, nil
}

func (client *Client) PostApp(newApp *model.NewApp) (*model.App, error) {
	appJson, err := json.Marshal(newApp)
	if err != nil {
		return nil, err
	}

	apiUri := client.uri("/api/v1/apps")
	response, err := doBodyRequest(apiUri.String(), defaultContentType, "POST", appJson)
	if err != nil {
		return nil, errors.Wrap(err, string(response))
	}

	app := &model.App{}
	err = unmarshal(response, app)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (client *Client) PutDeployment(deployment *model.RawDeployment, dryRun bool) (string, error) {
	deploymentJson, err := json.Marshal(deployment)
	if err != nil {
		return "", err
	}
	apiUri := client.uri("/api/v1/deployments")
	if dryRun {
		q := apiUri.Query()
		q.Add("dryRun", "true")
		apiUri.RawQuery = q.Encode()
	}

	response, err := doBodyRequest(apiUri.String(), defaultContentType, "PUT", deploymentJson)
	if err != nil {
		return "", errors.Wrap(err, string(response))
	}

	return safeDeserializeMessage(response), nil
}

func (client *Client) PutStatus(status *model.RawStatus) error {
	statusJson, err := json.Marshal(status)
	if err != nil {
		return err
	}
	apiUri := client.uri("/api/v1/status")

	_, err = doBodyRequest(apiUri.String(), defaultContentType, "PUT", statusJson)
	return err
}

func (client *Client) GetStatus(appName string) ([]model.StatusSummary, error) {
	apiUri := client.uri(fmt.Sprintf("/api/v1/status/%s", appName))

	responseBody, err := doGet(apiUri.String())
	if err != nil {
		return nil, err
	}

	statuses := []model.StatusSummary{}
	err = unmarshal(responseBody, &statuses)
	if err != nil {
		return nil, err
	}
	return statuses, nil
}

// unmarshal ignores nil response data and wraps the error with response content for easier debugging
func unmarshal(responseBody []byte, v interface{}) error {
	if responseBody != nil {
		err := json.Unmarshal(responseBody, v)
		if err != nil {
			return errors.Wrap(err, string(responseBody))
		}
	}

	return nil
}

func (client *Client) uri(apiPath string) url.URL {
	apiURI := *client.baseURI
	apiURI.Path = path.Join(apiURI.Path, apiPath)
	return apiURI
}

func doGet(uri string) ([]byte, error) {
	request, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Accept", defaultContentType)
	return doRequest(request)
}

func doBodyRequest(uri, contentType, httpMethod string, body []byte) ([]byte, error) {
	request, err := http.NewRequest(httpMethod, uri, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", contentType)
	return doRequest(request)
}

func doRequest(request *http.Request) ([]byte, error) {
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	} else {
		defer response.Body.Close()
		if isSuccess(response.StatusCode) {
			return ioutil.ReadAll(response.Body)
		}
		if response.StatusCode == 404 {
			return nil, nil
		}

		errorMessage := ""
		responseBody, err := ioutil.ReadAll(response.Body)
		if err == nil {
			errorMessage = safeDeserializeMessage(responseBody)
		}
		return nil, &HttpError{StatusCode: response.StatusCode, Message: errorMessage}
	}
}

func safeDeserializeMessage(body []byte) string {
	response := map[string]interface{}{}
	err := unmarshal(body, &response)
	if err == nil {
		return fmt.Sprintf("%s", response["message"])
	}
	return ""
}

func isSuccess(statusCode int) bool {
	return statusCode >= 200 && statusCode <= 300
}
