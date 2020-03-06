package sdk

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/stretchr/testify/assert"
)

func Test_Deployments_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/deployments/mystage/myns/mydep", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		fmt.Fprint(w, `{"message": "deleted"}`)
	})

	result, err := client.Deployments.Delete("mydep", "myns", "mystage")

	assert.NoError(t, err)
	assert.Equal(t, "deleted", result.Message)
}

func Test_Deployments_Save(t *testing.T) {
	setup()
	defer teardown()

	requestModel := &model.DeploymentRequest{
		DeploymentMeta: model.DeploymentMeta{
			Name: "mydeployment",
		},
	}

	mux.HandleFunc("/api/v1/deployments", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		actualModel := &model.DeploymentRequest{}
		mustUnmarshalR(r.Body, actualModel)
		assert.Equal(t, requestModel, actualModel)
		fmt.Fprint(w, `{"message": "saved"}`)
	})

	result, err := client.Deployments.Save(requestModel, false)

	assert.NoError(t, err)
	assert.Equal(t, "saved", result.Message)
}

func Test_Deployments_Save_DryRun(t *testing.T) {
	setup()
	defer teardown()

	requestModel := &model.DeploymentRequest{
		DeploymentMeta: model.DeploymentMeta{
			Name: "mydeployment",
		},
	}

	mux.HandleFunc("/api/v1/deployments", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "true", r.URL.Query().Get("dryRun"))
		actualModel := &model.DeploymentRequest{}
		mustUnmarshalR(r.Body, actualModel)
		assert.Equal(t, requestModel, actualModel)
		fmt.Fprint(w, `{"message": "dryRun", "dryRunCommits": [{ "message": "test"}]}`)
	})

	result, err := client.Deployments.Save(requestModel, true)

	assert.NoError(t, err)
	assert.Equal(t, "dryRun", result.Message)
	assert.EqualValues(t, "test", result.DryRunCommits[0].Message)
}

func Test_Deployments_SaveStatus(t *testing.T) {
	setup()
	defer teardown()

	requestModel := &model.DeploymentStatusMutable{
		ObservedRiserRevision: 1,
	}

	mux.HandleFunc("/api/v1/deployments/mystage/myns/myapp/status", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		actualModel := &model.DeploymentStatusMutable{}
		mustUnmarshalR(r.Body, actualModel)
		assert.Equal(t, requestModel, actualModel)
		response := `{"message": "saved"}`

		fmt.Fprint(w, response)
	})

	err := client.Deployments.SaveStatus("myapp", "myns", "mystage", requestModel)

	assert.NoError(t, err)
}

func Test_Deployments_SaveStatus_Error(t *testing.T) {
	setup()
	defer teardown()

	requestModel := &model.DeploymentStatusMutable{}

	mux.HandleFunc("/api/v1/deployments/mystage/myns/myapp/status", func(w http.ResponseWriter, r *http.Request) {
		response := `{"message": "err"}`
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, response)
	})

	err := client.Deployments.SaveStatus("myapp", "myns", "mystage", requestModel)

	assert.IsType(t, &ClientError{}, err)
	ce := err.(*ClientError)
	assert.Equal(t, "err", ce.Message)
}
