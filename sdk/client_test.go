package sdk

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/stretchr/testify/assert"
)

const testApiKey = "ad4a79e70859f3166d7232c67e6bb6b853fc4ff4"

var (
	mux    *http.ServeMux
	client *Client
	server *httptest.Server
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	client, _ = NewClient(server.URL, testApiKey)
}

func teardown() {
	server.Close()
}

func Test_NewRequest(t *testing.T) {
	setup()
	defer teardown()

	newApp := &model.NewApp{Name: model.AppName("myapp"), Namespace: model.NamespaceName("myns")}
	requestJson := `{	"name": "myapp", "namespace": "myns"}`

	request, err := client.NewRequest(http.MethodGet, "/v1/test", newApp)
	body, _ := ioutil.ReadAll(request.Body)

	assert.NoError(t, err)
	assert.Equal(t, server.URL+"/v1/test", request.URL.String())
	assert.JSONEq(t, requestJson, string(body))
	assert.Equal(t, defaultAccept, request.Header.Get("Accept"))
	assert.Equal(t, defaultContentType, request.Header.Get("Content-Type"))
	assert.Equal(t, userAgent, request.Header.Get("User-Agent"))
}

func Test_NewGetRequest(t *testing.T) {
	setup()
	defer teardown()

	request, err := client.NewGetRequest("/v1/test")

	assert.NoError(t, err)
	assert.Empty(t, request.Body)
	assert.Equal(t, server.URL+"/v1/test", request.URL.String())
	assert.Equal(t, defaultAccept, request.Header.Get("Accept"))
	assert.Equal(t, defaultContentType, request.Header.Get("Content-Type"))
	assert.Equal(t, userAgent, request.Header.Get("User-Agent"))
}

type testResponse struct {
	Field string
}

func Test_Do(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, `{"field": "val"}`)
	})

	request, _ := client.NewGetRequest("/")
	responseBody := testResponse{}
	response, err := client.Do(request, &responseBody)

	assert.NoError(t, err)
	assert.Equal(t, "val", responseBody.Field)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func Test_Do_ErrorMessage(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"message": "val"}`)
	})

	request, _ := client.NewGetRequest("/")
	responseBody := testResponse{}
	_, err := client.Do(request, &responseBody)

	assert.Error(t, err)
	assert.Empty(t, responseBody)
	assert.IsType(t, &ClientError{}, err)
	clientError := err.(*ClientError)
	assert.Equal(t, "val", clientError.Message)
}

func Test_Do_ClientError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `
		{
			"message": "val",
			"validationErrors": {
				"field": "fieldVal"
			}
		}`)
	})

	request, _ := client.NewGetRequest("/")
	responseBody := testResponse{}
	_, err := client.Do(request, &responseBody)

	assert.Error(t, err)
	assert.Empty(t, responseBody)
	assert.IsType(t, &ClientError{}, err)
	clientError := err.(*ClientError)
	assert.Equal(t, "val", clientError.Message)
	assert.Equal(t, "fieldVal", clientError.ValidationErrors["field"])
}

func mustReadAll(r io.Reader) []byte {
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}

	return bytes
}

func mustUnmarshalR(r io.ReadCloser, v interface{}) {
	err := json.Unmarshal(mustReadAll(r), v)
	if err != nil {
		panic(err)
	}
}
