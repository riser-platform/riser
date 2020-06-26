package cmd

import (
	"fmt"
	"os"
	"riser/pkg/logger"
	"riser/pkg/rc"
	"riser/pkg/ui"
	"riser/pkg/ui/style"

	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/riser-platform/riser-server/pkg/sdk"
	"github.com/spf13/cobra"
)

const AppConfigPath = "./app.yaml"

func newAppsCommand(config *rc.RuntimeConfiguration) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apps",
		Short: "Commands for managing apps",
	}

	cmd.AddCommand(newAppsListCommand(config))
	cmd.AddCommand(newAppsNewCommand(config))
	cmd.AddCommand(newAppsInitCommand(config))

	return cmd
}

func newAppsListCommand(config *rc.RuntimeConfiguration) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all apps",
		Run: func(cmd *cobra.Command, args []string) {
			currentContext := safeCurrentContext(config)
			riserClient := getRiserClient(currentContext)
			apps, err := riserClient.Apps.List()
			ui.ExitIfError(err)
			view := &ui.BasicTableView{}
			view.Header("Name", "Namespace", "Id")

			for _, app := range apps {
				view.AddRow(app.Name, app.Namespace, app.Id)
			}

			ui.RenderView(view)
		},
	}

	addOutputFlag(cmd.Flags())
	return cmd
}

func newAppsInitCommand(config *rc.RuntimeConfiguration) *cobra.Command {
	var namespace string
	cmd := &cobra.Command{
		Use:   "init (app name)",
		Short: "Creates a new app with a default app.yaml file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			appName := args[0]
			app := createNewApp(config, appName, namespace)

			file, err := os.OpenFile(AppConfigPath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644)
			ui.ExitIfErrorMsg(err, "Error creating default app config")
			defer file.Close()

			err = sdk.DefaultAppConfig(file, app.Id, appName, namespace)
			ui.ExitIfErrorMsg(err, "Error creating default app config")
			logger.Log().Info(fmt.Sprintf("App %s created with a default app config file %q. Please review the TODO's before deploying your app.", style.Emphasis(appName), AppConfigPath))
		},
	}

	addNamespaceFlag(cmd.Flags(), &namespace)

	return cmd
}

func newAppsNewCommand(config *rc.RuntimeConfiguration) *cobra.Command {
	var namespace string
	cmd := &cobra.Command{
		Use:   "new (app name)",
		Short: "Creates a new app",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			appName := args[0]
			app := createNewApp(config, appName, namespace)

			fmt.Printf("App %s created. Please add the following id to your manifest: %s", app.Name, app.Id)
		},
	}

	addNamespaceFlag(cmd.Flags(), &namespace)

	return cmd
}

func createNewApp(config *rc.RuntimeConfiguration, appName, namespace string) *model.App {
	currentContext := safeCurrentContext(config)
	riserClient := getRiserClient(currentContext)
	app, err := riserClient.Apps.Create(&model.NewApp{Name: model.AppName(appName), Namespace: model.NamespaceName(namespace)})
	ui.ExitIfError(err)
	return app
}
