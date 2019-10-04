package sdk

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/tshak/riser-server/api/v1/model"

	"github.com/stretchr/testify/assert"
)

func Test_Status_Save(t *testing.T) {
	setup()
	defer teardown()

	requestModel := &model.DeploymentStatus{AppName: "myapp"}

	mux.HandleFunc("/api/v1/status", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		actualModel := &model.DeploymentStatus{}
		mustUnmarshalR(r.Body, actualModel)
		assert.Equal(t, requestModel, actualModel)
		response := `{"message": "saved"}`

		fmt.Fprint(w, response)
	})

	err := client.Status.Save(requestModel)

	assert.NoError(t, err)
}

func Test_Status_Save_Error(t *testing.T) {
	setup()
	defer teardown()

	requestModel := &model.DeploymentStatus{AppName: "myapp"}

	mux.HandleFunc("/api/v1/status", func(w http.ResponseWriter, r *http.Request) {
		response := `{"message": "err"}`
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, response)
	})

	err := client.Status.Save(requestModel)

	assert.IsType(t, &ClientError{}, err)
	ce := err.(*ClientError)
	assert.Equal(t, "err", ce.Message)
}

func Test_Status_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/status", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "myapp", r.URL.Query().Get("app"))

		response := `{"deployments":[{"app":"myapp"}]}`
		fmt.Fprint(w, response)
	})

	status, err := client.Status.Get("myapp")

	assert.NoError(t, err)
	assert.Equal(t, "myapp", status.Deployments[0].AppName)
}
