package sdk

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
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
			{"name": "myapp01", "id":"e29bf621-4da7-4df1-8c04-6609b9eb2447"},
			{"name": "myapp02", "id":"e29bf621-4da7-4df1-8c04-6609b9eb2448"}
		]`

		fmt.Fprint(w, response)
	})

	apps, err := client.Apps.List()

	assert.NoError(t, err)
	assert.Len(t, apps, 2)
	assert.EqualValues(t, "myapp01", apps[0].Name)
	assert.Equal(t, uuid.MustParse("e29bf621-4da7-4df1-8c04-6609b9eb2447"), apps[0].Id)
	assert.EqualValues(t, "myapp02", apps[1].Name)
	assert.Equal(t, uuid.MustParse("e29bf621-4da7-4df1-8c04-6609b9eb2448"), apps[1].Id)
}

func Test_Apps_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/apps", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		response := `{"name": "myapp01", "id":"e29bf621-4da7-4df1-8c04-6609b9eb2447"}`

		fmt.Fprint(w, response)
	})

	app, err := client.Apps.Create(&model.NewApp{Name: "myapp01"})

	assert.NoError(t, err)
	assert.EqualValues(t, "myapp01", app.Name)
	assert.Equal(t, uuid.MustParse("e29bf621-4da7-4df1-8c04-6609b9eb2447"), app.Id)
}

func Test_Apps_GetStatus(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/apps/myns/myapp/status", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)

		response := `{"deployments":[{"deployment":"mydeployment"}]}`
		fmt.Fprint(w, response)
	})

	status, err := client.Apps.GetStatus("myapp", "myns")

	assert.NoError(t, err)
	assert.Equal(t, "mydeployment", status.Deployments[0].DeploymentName)
}

func Test_Apps_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/apps/myns/myapp", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)

		response := `{"name": "myapp"}`
		fmt.Fprint(w, response)
	})

	app, err := client.Apps.Get("myapp", "myns")

	assert.NoError(t, err)
	assert.EqualValues(t, "myapp", app.Name)
}
