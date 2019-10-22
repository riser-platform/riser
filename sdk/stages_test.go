package sdk

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/riser-platform/riser-server/api/v1/model"

	"github.com/stretchr/testify/assert"
)

func Test_Stages_Ping(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/stages/dev/ping", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.WriteHeader(http.StatusNoContent)
	})

	err := client.Stages.Ping("dev")

	assert.NoError(t, err)
}

func Test_Stages_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/stages", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, `[{"name":"dev"},{"name":"prod"}]`)
	})

	stages, err := client.Stages.List()

	assert.NoError(t, err)
	assert.Len(t, stages, 2)
	assert.Equal(t, "dev", stages[0].Name)
	assert.Equal(t, "prod", stages[1].Name)
}

func Test_Stages_SetConfig(t *testing.T) {
	setup()
	defer teardown()

	config := &model.StageConfig{PublicGatewayHost: "tempuri.org"}

	mux.HandleFunc("/api/v1/stages/dev/config", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		actualConfig := &model.StageConfig{}
		mustUnmarshalR(r.Body, actualConfig)
		assert.Equal(t, config, actualConfig)
		w.WriteHeader(http.StatusAccepted)
	})

	err := client.Stages.SetConfig("dev", config)

	assert.NoError(t, err)
}
