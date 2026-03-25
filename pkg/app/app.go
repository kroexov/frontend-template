package app

import (
	"context"
	"frontend/pkg/frontend"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/vmkteam/appkit"
	"github.com/vmkteam/embedlog"
)

type Config struct {
	Server struct {
		Host    string
		Port    int
		IsDevel bool
	}
	Sentry struct {
		Environment string
		DSN         string
	}
}

type App struct {
	embedlog.Logger
	appName string
	cfg     Config
	echo    *echo.Echo

	wm *frontend.WidgetManager
}

func New(appName string, sl embedlog.Logger, cfg Config) *App {
	a := &App{
		appName: appName,
		cfg:     cfg,
		echo:    appkit.NewEcho(),
		Logger:  sl,
		wm:      frontend.NewWidgetManager(sl),
	}

	return a
}

// Run is a function that runs application.
func (a *App) Run(ctx context.Context) error {
	a.registerHandlers()
	a.registerDebugHandlers()
	a.registerMetadata()

	err := a.wm.Init()
	if err != nil {
		a.Error(ctx, "init widgetManager failed", "err", err)
		return err
	}

	return a.runHTTPServer(ctx, a.cfg.Server.Host, a.cfg.Server.Port)
}

// Shutdown is a function that gracefully stops HTTP server.
func (a *App) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return a.echo.Shutdown(ctx)
}

// registerMetadata is a function that registers meta info from service. Must be updated.
func (a *App) registerMetadata() {
	opts := appkit.MetadataOpts{
		HasPublicAPI:  true,
		HasPrivateAPI: true,
		Services:      []appkit.ServiceMetadata{
			// NewServiceMetadata("srv", MetadataServiceTypeAsync),
		},
	}

	md := appkit.NewMetadataManager(opts)
	md.RegisterMetrics()

	a.echo.GET("/debug/metadata", md.Handler)
}
