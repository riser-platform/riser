package sdk

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/stretchr/testify/assert"
)

func Test_Secrets_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/secrets/dev/myns/myapp", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `[{"name":"mysecret"}]`)
	})

	secretMetas, err := client.Secrets.List("myapp", "myns", "dev")

	assert.NoError(t, err)
	assert.Len(t, secretMetas, 1)
	assert.Equal(t, "mysecret", secretMetas[0].Name)
}

func Test_Secrets_Save(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/secrets", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		unsealed := &model.UnsealedSecret{}
		mustUnmarshalR(r.Body, unsealed)
		assert.EqualValues(t, "myapp", unsealed.AppName)
		assert.EqualValues(t, "myns", unsealed.Namespace)
		assert.Equal(t, "dev", unsealed.Stage)
		assert.Equal(t, "mysecret", unsealed.Name)
		assert.Equal(t, "myval", unsealed.PlainText)
	})

	err := client.Secrets.Save("myapp", "myns", "dev", "mysecret", "myval")

	assert.NoError(t, err)
}
