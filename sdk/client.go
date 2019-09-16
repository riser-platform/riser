package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/pkg/errors"

	"github.com/tshak/riser-server/api/v1/model"
)

const defaultContentType = "application/json"

// Client is a API v1 client
type Client struct {
	baseURI *url.URL
	apikey  string
}

// ClientError provides the error message, status code.
type ClientError struct {
	StatusCode int
	Message    string
}

func (e *ClientError) Error() string {
	return fmt.Sprintf("Received HTTP %d (%s) %s", e.StatusCode, http.StatusText(e.StatusCode), e.Message)
}

type ClientValidationError struct {
	ClientError
	ValidationErrors map[string]interface{}
}

func (e *ClientValidationError) Error() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%s:", e.Message))
	builder.WriteString("\n")
	for fieldName, errorMessage := range e.ValidationErrors {
		builder.WriteString(fmt.Sprintf(" • %s: %s\n", fieldName, errorMessage))
	}
	return builder.String()
}

func NewClient(baseURI string, apikey string) (*Client, error) {
	baseURIParsed, err := url.Parse(baseURI)
	if err != nil {
		return nil, err
	}

	return &Client{baseURI: baseURIParsed, apikey: apikey}, nil
}

func (client *Client) ListApps() ([]model.App, error) {
	apps := []model.App{}
	apiUri := client.uri("/api/v1/apps")
	response, err := doGet(apiUri.String(), client.apikey)
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
	response, err := doBodyRequest(apiUri.String(), client.apikey, defaultContentType, "POST", appJson)
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

func (client *Client) PutDeployment(deployment *model.DeploymentRequest, dryRun bool) (string, error) {
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

	responseBody, err := doBodyRequest(apiUri.String(), client.apikey, defaultContentType, "PUT", deploymentJson)
	if err != nil {
		return "", err
	}

	if dryRun {
		return fmt.Sprintf("%s", responseBody), nil
	} else {
		return safeUnmarshalApiResponse(responseBody).Message, nil
	}
}

func (client *Client) PutStatus(status *model.DeploymentStatus) error {
	statusJson, err := json.Marshal(status)
	if err != nil {
		return err
	}
	apiUri := client.uri("/api/v1/status")

	_, err = doBodyRequest(apiUri.String(), client.apikey, defaultContentType, "PUT", statusJson)
	return err
}

func (client *Client) GetStatus(appName string) (*model.Status, error) {
	apiUri := client.uri(fmt.Sprintf("/api/v1/status/%s", appName))

	responseBody, err := doGet(apiUri.String(), client.apikey)
	if err != nil {
		return nil, err
	}

	status := &model.Status{}
	err = unmarshal(responseBody, status)
	if err != nil {
		return nil, err
	}
	return status, nil
}

func (client *Client) PostStagePing(stageName string) error {
	apiUri := client.uri(fmt.Sprintf("/api/v1/stage/%s/ping", stageName))

	request, err := http.NewRequest("POST", apiUri.String(), nil)
	if err != nil {
		return err
	}
	request.Header.Add("Accept", defaultContentType)
	request.Header.Add(apikeyHeader(client.apikey))
	_, err = doRequest(request)
	return err
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

func doGet(uri string, apikey string) ([]byte, error) {
	request, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Accept", defaultContentType)
	request.Header.Add(apikeyHeader(apikey))
	return doRequest(request)
}

func apikeyHeader(apikey string) (string, string) {
	return "Authorization", fmt.Sprintf("Apikey: %s", apikey)
}

func doBodyRequest(uri, apikey, contentType, httpMethod string, body []byte) ([]byte, error) {
	request, err := http.NewRequest(httpMethod, uri, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", contentType)
	request.Header.Add(apikeyHeader(apikey))
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

		return nil, getErrorFromResponse(response)
	}
}

func getErrorFromResponse(response *http.Response) error {
	errorMessage := fmt.Sprintf("%s", response.Body)
	responseMap := map[string]interface{}{}

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err == nil {
		err := unmarshal(responseBytes, &responseMap)
		if err == nil {
			errorMessage = fmt.Sprintf("%s", responseMap["message"])
		}
	}
	clientErr := &ClientError{StatusCode: response.StatusCode, Message: errorMessage}
	if validationErrors, ok := responseMap["validationErrors"]; ok {
		if validationErrorMap, ok := validationErrors.(map[string]interface{}); ok {
			return &ClientValidationError{ClientError: *clientErr, ValidationErrors: validationErrorMap}
		}
	}
	return clientErr
}

func safeUnmarshalApiResponse(responseBytes []byte) *model.APIResponse {
	message := string(responseBytes)
	responseMap := map[string]interface{}{}
	err := unmarshal(responseBytes, &responseMap)
	if err == nil {
		message = fmt.Sprintf("%s", responseMap["message"])
	}
	return &model.APIResponse{
		Message: message,
	}
}

func isSuccess(statusCode int) bool {
	return statusCode >= 200 && statusCode <= 300
}
