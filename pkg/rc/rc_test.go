package rc

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_GetRcPath(t *testing.T) {
	os.Setenv("HOME", "/foo")

	result, err := getRcPath()

	assert.Nil(t, err)
	assert.Equal(t, "/foo/.riserrc", result)
}

func Test_GetRcPath_ErrorIfNoHome(t *testing.T) {
	os.Setenv("HOME", "")

	result, err := getRcPath()

	assert.Empty(t, result)
	assert.Equal(t, "the $HOME environment variable must be set to a writeable directory", err.Error())
}

func Test_CurrentContext_ReturnsCurrentContext(t *testing.T) {
	rc := RuntimeConfiguration{
		contextMap: toContextMap([]Context{
			{Name: "a"},
			{Name: "b"},
		}),
		CurrentContextName: "b",
	}

	result, err := rc.CurrentContext()

	assert.Nil(t, err)
	assert.Equal(t, result.Name, "b")
}

func Test_CurrentContext_ReturnsErr_WhenNoCurrentContext(t *testing.T) {
	rc := RuntimeConfiguration{
		contextMap: toContextMap([]Context{
			{Name: "a"},
		}),
	}

	result, err := rc.CurrentContext()

	assert.Nil(t, result)
	assert.Equal(t, "Unable to load current context: no context set. Use \"riser context current <contextName>\" to set the context", err.Error())
}

func Test_CurrentContext_ReturnsError_WhenCurrentContextDoesNotExist(t *testing.T) {
	rc := RuntimeConfiguration{
		CurrentContextName: "invalid",
	}

	result, err := rc.CurrentContext()

	assert.Nil(t, result)
	assert.Equal(t, err.Error(), "Unable to load current context: context \"invalid\" does not exist")
}

func Test_SetCurrentContext(t *testing.T) {
	rc := RuntimeConfiguration{
		contextMap: toContextMap([]Context{
			{Name: "a"},
			{Name: "b"},
		}),

		CurrentContextName: "a",
	}

	err := rc.SetCurrentContext("b")

	assert.Nil(t, err)
	assert.Equal(t, "b", rc.CurrentContextName)
}

func Test_SetCurrentContext_ReturnsError_WhenInvalidContext(t *testing.T) {
	rc := RuntimeConfiguration{
		contextMap: toContextMap([]Context{
			{Name: "a"},
		}),

		CurrentContextName: "a",
	}

	err := rc.SetCurrentContext("invalid")

	assert.Equal(t, "a", rc.CurrentContextName)
	assert.Equal(t, "Context \"invalid\" does not exist", err.Error())
}

func Test_loadAndParse(t *testing.T) {
	rc := `currentContext: a
contexts:
  - name: a
    serverUrl: https://riser.up
`
	tmpfile, err := ioutil.TempFile("", ".testrc")
	assert.Nil(t, err)

	defer os.Remove(tmpfile.Name())
	_, err = tmpfile.WriteString(rc)
	assert.Nil(t, err)

	result, err := loadAndParseRc(tmpfile.Name())

	assert.NoError(t, err)
	assert.Equal(t, "a", result.CurrentContextName)
	require.Equal(t, 1, len(result.GetContexts()))
	assert.Equal(t, "a", result.contextMap["a"].Name)
	assert.Equal(t, "https://riser.up", result.contextMap["a"].ServerURL)
}

func Test_loadAndParse_ReturnsEmptyIfRcFileIsMissing(t *testing.T) {
	result, err := loadAndParseRc("/missing")

	assert.NotNil(t, result)
	assert.NoError(t, err)
}

func Test_SaveContext_AddsNewContext(t *testing.T) {
	rc := &RuntimeConfiguration{}

	rc.SetContext(&Context{Name: "a"})

	assert.Equal(t, "a", rc.contextMap["a"].Name)
	assert.Equal(t, "a", rc.Contexts[0].Name)
}

func Test_SaveContext_SavesIfContextExists(t *testing.T) {
	rc := &RuntimeConfiguration{
		contextMap: toContextMap([]Context{
			{Name: "a", ServerURL: "URLa"},
		}),

		CurrentContextName: "a",
	}

	rc.SetContext(&Context{Name: "a", ServerURL: "modified"})

	assert.Equal(t, "a", rc.contextMap["a"].Name)
	assert.Equal(t, "a", rc.Contexts[0].Name)
	assert.Equal(t, "modified", rc.contextMap["a"].ServerURL)
	assert.Equal(t, "modified", rc.Contexts[0].ServerURL)
}

func Test_SaveContext_SetsCurrentContext(t *testing.T) {
	rc := &RuntimeConfiguration{CurrentContextName: "b"}

	rc.SetContext(&Context{Name: "a"})

	assert.Equal(t, "a", rc.CurrentContextName)
}

func Test_SaveContext_SetsCurrentContext_WhenNotSet(t *testing.T) {
	rc := &RuntimeConfiguration{}

	rc.SetContext(&Context{Name: "a"})

	assert.Equal(t, "a", rc.CurrentContextName)
}

func Test_RemoveContext(t *testing.T) {
	rc := RuntimeConfiguration{
		contextMap: toContextMap([]Context{
			{Name: "a"},
			{Name: "b"},
		}),

		CurrentContextName: "a",
	}

	result := rc.RemoveContext("a")

	assert.Nil(t, result)
	_, found := rc.contextMap["a"]
	assert.False(t, found, "context should be removed")
	assert.Equal(t, "", rc.CurrentContextName)
	assert.Equal(t, "b", rc.Contexts[0].Name)
}

func Test_RemoveContext_ReturnsError_WhenContextDoesNotExist(t *testing.T) {
	rc := RuntimeConfiguration{}

	result := rc.RemoveContext("a")

	assert.Equal(t, "a context with the name \"a\" does not exist", result.Error())
}
