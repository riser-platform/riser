package rc

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"

	"github.com/pkg/errors"
)

// RuntimeConfiguration provides configuration for the client
type RuntimeConfiguration struct {
	CurrentContextName string    `yaml:"currentContext,omitempty"`
	Contexts           []Context `yaml:"contexts,omitempty"`
	contextMap         map[string]Context
}

// Context represents all configuration related to a particular environment
type Context struct {
	Name      string `yaml:"name"`
	ServerURL string `yaml:"serverUrl"`
	// TODO: Consider introducing an "auth" section to better support other auth methods in the future (e.g. OIDC, LDAP)
	Apikey string `yaml:"apikey,omitempty"`
}

// SaveRc saves a runtime configuration
func SaveRc(rc *RuntimeConfiguration) error {
	rcPath, err := getRcPath()
	if err != nil {
		return err
	}

	rcBytes, err := yaml.Marshal(rc)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(rcPath, rcBytes, 0644)
}

// LoadRc loads runtime configuration from the HOME directory
func LoadRc() (*RuntimeConfiguration, error) {
	rcPath, err := getRcPath()
	if err != nil {
		return nil, err
	}

	return loadAndParseRc(rcPath)
}

func loadAndParseRc(rcPath string) (*RuntimeConfiguration, error) {
	var rc *RuntimeConfiguration

	if _, err := os.Stat(rcPath); err == nil {
		rcBytes, err := ioutil.ReadFile(rcPath)
		if err != nil {
			return nil, err
		}

		err = yaml.UnmarshalStrict(rcBytes, &rc)
		if err != nil {
			return nil, err
		}

		rc.contextMap = toContextMap(rc.Contexts)
		return rc, nil
	}

	return &RuntimeConfiguration{}, nil
}

// CurrentContext returns the current context.
func (rc *RuntimeConfiguration) CurrentContext() (*Context, error) {
	var context *Context
	if rc.CurrentContextName == "" {
		return nil, contextError("no context set. Use \"riser context current <contextName>\" to set the context")
	} else {
		context = rc.getContextByName(rc.CurrentContextName)
	}
	if context != nil {
		return context, nil
	}

	return nil, contextError(fmt.Sprintf("context \"%s\" does not exist", rc.CurrentContextName))
}

func (rc *RuntimeConfiguration) getContextByName(name string) *Context {
	context, found := rc.contextMap[name]
	if found {
		return &context
	}
	return nil
}

func contextError(errorMessage string) error {
	return fmt.Errorf("Unable to load current context: %s", errorMessage)
}

// SetCurrentContext sets the current context
func (rc *RuntimeConfiguration) SetCurrentContext(contextName string) error {
	_, found := rc.contextMap[contextName]
	if found {
		rc.CurrentContextName = contextName
		return nil
	}

	return fmt.Errorf("Context \"%s\" does not exist", contextName)
}

// GetContexts returns all configured contexts
func (rc *RuntimeConfiguration) GetContexts() []Context {
	values := []Context{}
	for _, value := range rc.contextMap {
		values = append(values, value)
	}

	return values
}

// AddContext adds or updates a context. Sets the passed in context as the currentContext if not set
func (rc *RuntimeConfiguration) AddContext(context *Context) error {
	if rc.contextMap == nil {
		rc.contextMap = map[string]Context{}
	}

	_, found := rc.contextMap[context.Name]
	if found {
		return fmt.Errorf("a context with the name \"%s\" already exists", context.Name)
	}

	rc.contextMap[context.Name] = *context
	if rc.CurrentContextName == "" {
		rc.CurrentContextName = context.Name
	}

	rc.Contexts = rc.GetContexts()

	return nil
}

// RemoveContext removes a context
func (rc *RuntimeConfiguration) RemoveContext(contextName string) error {
	_, found := rc.contextMap[contextName]
	if !found {
		return fmt.Errorf("a context with the name \"%s\" does not exist", contextName)
	}

	delete(rc.contextMap, contextName)
	if rc.CurrentContextName == contextName {
		rc.CurrentContextName = ""
	}

	rc.Contexts = rc.GetContexts()
	return nil
}

func toContextMap(contexts []Context) map[string]Context {
	contextMap := map[string]Context{}

	for i := range contexts {
		context := contexts[i]
		contextMap[context.Name] = context
	}

	return contextMap
}

func getRcPath() (string, error) {
	// TODO: Provide a flag for al alternate path
	home := os.Getenv("HOME")
	if home != "" {
		return path.Join(home, ".riserrc"), nil
	}

	return "", errors.New("the $HOME environment variable must be set to a writeable directory")
}
