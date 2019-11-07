package sdk

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/riser-platform/riser-server/api/v1/model"

	"github.com/stretchr/testify/assert"
)

func Test_Apps_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/apps", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		response := `
		[
			{"name": "myapp01", "id":"1"},
			{"name": "myapp02", "id":"2"}
		]`

		fmt.Fprint(w, response)
	})

	apps, err := client.Apps.List()

	assert.NoError(t, err)
	assert.Len(t, apps, 2)
	assert.Equal(t, "myapp01", apps[0].Name)
	assert.Equal(t, "1", apps[0].Id)
	assert.Equal(t, "myapp02", apps[1].Name)
	assert.Equal(t, "2", apps[1].Id)
}

func Test_Apps_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/apps", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		response := `{"name": "myapp01", "id":"1"}`

		fmt.Fprint(w, response)
	})

	app, err := client.Apps.Create(&model.NewApp{Name: "myapp01"})

	assert.NoError(t, err)
	assert.Equal(t, "myapp01", app.Name)
	assert.Equal(t, "1", app.Id)
}

func Test_Apps_GetStatus(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/apps/myapp/status", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)

		response := `{"deployments":[{"deployment":"mydeployment"}]}`
		fmt.Fprint(w, response)
	})

	status, err := client.Apps.GetStatus("myapp")

	assert.NoError(t, err)
	assert.Equal(t, "mydeployment", status.Deployments[0].DeploymentName)
}
