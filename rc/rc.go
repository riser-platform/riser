package rc

// RuntimeConfiguration is a stub for the future context support to switch between riser instances
type RuntimeConfiguration struct {
	/*
		TODO:
		CurrentContext string           `yaml:"currentContext,omitempty"`
		Contexts       []RuntimeContext `yaml:"contexts,omitempty"`
	*/
}

// CurrentContext is currently hardcoded for local development only
func (rc *RuntimeConfiguration) CurrentContext() (*RuntimeContext, error) {
	return &RuntimeContext{
		Name:      "LocalDev",
		ServerURL: "http://localhost:8000",
	}, nil
}

// RuntimeContext represents all configuration related to a particular environment
type RuntimeContext struct {
	Name      string `yaml:"name"`
	ServerURL string `yaml:"url,omitempty"`
}
