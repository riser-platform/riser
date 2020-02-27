package sdk

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/stretchr/testify/assert"
)

func Test_Secrets_List(t *testing.T) {
	setup()
	defer teardown()

	appId := uuid.New()

	mux.HandleFunc(fmt.Sprintf("/api/v1/secrets/%s/dev", appId), func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fmt.Sprintf(`[{"appId":"%s"}]`, appId))
	})

	secretMetas, err := client.Secrets.List(appId, "dev")

	assert.NoError(t, err)
	assert.Len(t, secretMetas, 1)
	assert.Equal(t, appId, secretMetas[0].AppId)
}

func Test_Secrets_Save(t *testing.T) {
	setup()
	defer teardown()

	appId := uuid.New()

	mux.HandleFunc("/api/v1/secrets", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		unsealed := &model.UnsealedSecret{}
		mustUnmarshalR(r.Body, unsealed)
		assert.Equal(t, appId, unsealed.AppId)
		assert.Equal(t, "dev", unsealed.Stage)
		assert.Equal(t, "mysecret", unsealed.Name)
		assert.Equal(t, "myval", unsealed.PlainText)
	})

	err := client.Secrets.Save(appId, "dev", "mysecret", "myval")

	assert.NoError(t, err)
}
