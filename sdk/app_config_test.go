package sdk

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

const expectedAppConfig = `name: myapp
id: "myid"
# TODO: Update to use your docker image registry/repo (without tag) here
image: your/image
expose:
	# TODO: Update the container port that gets exposed to the HTTPS gateway
  containerPort: 8000
`

func Test_DefaultAppConfig(t *testing.T) {
	buffer := &bytes.Buffer{}

	err := DefaultAppConfig(buffer, "myapp", "myid")

	assert.NoError(t, err)
	assert.Equal(t, expectedAppConfig, buffer.String())
}
