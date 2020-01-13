package sdk

import (
	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func Test_Rollouts_Save(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/rollout/myapp/dev", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		rollout := &model.RolloutRequest{}
		mustUnmarshalR(r.Body, rollout)
		assert.Len(t, rollout.Traffic, 2)
		assert.EqualValues(t, 1, rollout.Traffic[0].RiserGeneration)
		assert.EqualValues(t, 10, rollout.Traffic[0].Percent)
		assert.EqualValues(t, 2, rollout.Traffic[1].RiserGeneration)
		assert.EqualValues(t, 90, rollout.Traffic[1].Percent)

	})

	err := client.Rollouts.Save("myapp", "dev", "1:10", "2:90")

	assert.NoError(t, err)
}

func Test_Rollouts_Save_ReturnsError_WhenBadRule(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/rollout/myapp/dev", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		rollout := &model.RolloutRequest{}
		mustUnmarshalR(r.Body, rollout)
		assert.Len(t, rollout.Traffic, 2)
		assert.EqualValues(t, 1, rollout.Traffic[0].RiserGeneration)
		assert.EqualValues(t, 10, rollout.Traffic[0].Percent)
		assert.EqualValues(t, 2, rollout.Traffic[1].RiserGeneration)
		assert.EqualValues(t, 90, rollout.Traffic[1].Percent)

	})

	err := client.Rollouts.Save("myapp", "dev", "1:10", "bad:90")

	assert.Equal(t, `Rules must be in the format of "(rev):(percentage)" e.g. "1:100" routes 100% of traffic to rev 1`, err.Error())
}
