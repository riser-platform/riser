package sdk

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tshak/riser-server/api/v1/model"
)

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

	message, err := client.Deployments.Save(requestModel, false)

	assert.NoError(t, err)
	assert.Equal(t, "saved", message)
}

// TODO: Test_Deployments_Save_Dryrun
