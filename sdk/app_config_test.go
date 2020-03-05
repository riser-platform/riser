package sdk

import (
	"bytes"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const expectedAppConfig = `name: myapp
namespace: myns
id: e29bf621-4da7-4df1-8c04-6609b9eb2447
# TODO: Update to use your docker image registry/repo (without tag) here
image: your/image
expose:
  # TODO: Update the container port that gets exposed to the HTTPS gateway
  containerPort: 8000
`

func Test_DefaultAppConfig(t *testing.T) {
	buffer := &bytes.Buffer{}

	err := DefaultAppConfig(buffer, uuid.MustParse("E29BF621-4DA7-4DF1-8C04-6609B9EB2447"), "myapp", "myns")

	assert.NoError(t, err)
	assert.Equal(t, expectedAppConfig, buffer.String())
}
