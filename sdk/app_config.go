package sdk

import (
	"io"
	"text/template"

	"github.com/google/uuid"

	"github.com/pkg/errors"
)

const appConfigTemplate = `name: {{.AppName}}
namespace: {{.AppNamespace}}
id: {{.AppId}}
# TODO: Update to use your docker image registry/repo (without tag) here
image: your/image
expose:
  # TODO: Update the container port that gets exposed to the HTTPS gateway
  containerPort: 8000
`

type AppConfigTemplateData struct {
	AppName      string
	AppNamespace string
	AppId        string
}

// DefaultAppConfig generates a default app config yaml for an app
func DefaultAppConfig(writer io.Writer, appId uuid.UUID, appName, appNamespace string) error {
	parsedTemplate, err := template.New("appconfig").Parse(appConfigTemplate)
	if err != nil {
		return errors.Wrap(err, "Error parsing app config template")
	}

	err = parsedTemplate.Execute(writer, AppConfigTemplateData{
		AppName:      appName,
		AppNamespace: appNamespace,
		AppId:        appId.String(),
	})

	return err
}
