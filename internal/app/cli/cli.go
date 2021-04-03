package cli

import (
	"context"
	"github.com/spf13/cobra"
	commonApp "github.com/yaruz/app/internal/app"
	"github.com/yaruz/app/internal/pkg/config"
)

// App is the application for CLI app
type App struct {
	*commonApp.App
	rootCmd *cobra.Command
}

// New func is a constructor for the App
func New(commonApp *commonApp.App, cfg config.Configuration) *App {
	app.App = commonApp
	return app
}

// Run is func to run the App
func (app *App) Run() error {
	return app.rootCmd.ExecuteContext(context.Background())
}

var app = &App{
	rootCmd: &cobra.Command{
		Use:   "cli",
		Short: "This is the short description.",
		Long:  `This is the long description.`,
	},
}
