package sdk

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/stretchr/testify/assert"
)

func Test_Namespaces_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/namespaces", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		response := `
		[
			{"name": "myns01"},
			{"name": "myns02"}
		]`

		fmt.Fprint(w, response)
	})

	namespaces, err := client.Namespaces.List()

	assert.NoError(t, err)
	assert.Len(t, namespaces, 2)
	assert.EqualValues(t, "myns01", namespaces[0].Name)
	assert.EqualValues(t, "myns02", namespaces[1].Name)
}

func Test_Namespaces_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/namespaces", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		actualModel := &model.Namespace{}
		mustUnmarshalR(r.Body, actualModel)
		assert.EqualValues(t, "myns", actualModel.Name)
		fmt.Fprint(w, "")
	})

	err := client.Namespaces.Create("myns")

	assert.NoError(t, err)
}
