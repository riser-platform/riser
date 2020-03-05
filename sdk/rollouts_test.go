package sdk

import (
	"net/http"
	"testing"

	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/stretchr/testify/assert"
)

func Test_Rollouts_Save(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/rollout/dev/myns/myapp", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		rollout := &model.RolloutRequest{}
		mustUnmarshalR(r.Body, rollout)
		assert.Len(t, rollout.Traffic, 2)
		assert.EqualValues(t, 1, rollout.Traffic[0].RiserRevision)
		assert.EqualValues(t, 10, rollout.Traffic[0].Percent)
		assert.EqualValues(t, 2, rollout.Traffic[1].RiserRevision)
		assert.EqualValues(t, 90, rollout.Traffic[1].Percent)

	})

	err := client.Rollouts.Save("myapp", "myns", "dev", "rev-1:10", "rev-2:90")

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
		assert.EqualValues(t, 1, rollout.Traffic[0].RiserRevision)
		assert.EqualValues(t, 10, rollout.Traffic[0].Percent)
		assert.EqualValues(t, 2, rollout.Traffic[1].RiserRevision)
		assert.EqualValues(t, 90, rollout.Traffic[1].Percent)

	})

	err := client.Rollouts.Save("myapp", "dev", "rev-1:10", "bad:90")

	assert.Equal(t, `Rules must be in the format of "rev-(rev):(percentage)" e.g. "rev-1:100" routes 100% of traffic to rev 1`, err.Error())
}
