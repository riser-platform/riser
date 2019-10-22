package sdk

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/riser-platform/riser-server/api/v1/model"
)

func Test_Validate_AppConfig(t *testing.T) {
	setup()
	defer teardown()

	requestModel := &model.AppConfigWithOverrides{}

	mux.HandleFunc("/api/v1/validate/appconfig", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		actualModel := &model.AppConfigWithOverrides{}
		mustUnmarshalR(r.Body, actualModel)
		assert.Equal(t, requestModel, actualModel)
		w.WriteHeader(http.StatusOK)
	})

	err := client.Validate.AppConfig(requestModel)

	assert.NoError(t, err)
}
