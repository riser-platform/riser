package sdk

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tshak/riser-server/api/v1/model"
)

func Test_Secrets_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/secrets/myapp/dev", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `[{"app":"myapp"}]`)
	})

	secretMetas, err := client.Secrets.List("myapp", "dev")

	assert.NoError(t, err)
	assert.Len(t, secretMetas, 1)
	assert.Equal(t, "myapp", secretMetas[0].App)
}

func Test_Secrets_Save(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/secrets", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		unsealed := &model.UnsealedSecret{}
		mustUnmarshalR(r.Body, unsealed)
		assert.Equal(t, "myapp", unsealed.App)
		assert.Equal(t, "dev", unsealed.Stage)
		assert.Equal(t, "mysecret", unsealed.Name)
		assert.Equal(t, "myval", unsealed.PlainText)
	})

	err := client.Secrets.Save("myapp", "dev", "mysecret", "myval")

	assert.NoError(t, err)
}
