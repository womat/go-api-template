package app

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// VERSION holds the version information with the following logic in mind
//
//	1 ... fixed
//	0 ... year 2020, 1->year 2021, etc.
//	7 ... month of year (7=July)
//	the date format after the + is always the first of the month
//
// VERSION differs from semantic versioning as described in https://semver.org/
// but we keep the correct syntax.
// TODO: increase version number to 1.0.1+2020xxyy
const (
	VERSION = "0.0.0+yyyymmdd"
	MODULE  = "<MODUL_NAME>"
)

//go:embed README.md
var Readme string

// App is the main application struct.
// App is where the application is wired up.
type App struct {

	// baseDir is the application working directory
	baseDir string

	// config is the application configuration
	config *Config

	// web is the web server.
	web *http.Server

	// restart signals application restart
	restart chan struct{}

	// shutdown signals application shutdown
	shutdown chan struct{}
}

// New checks the Web server URL and initializes the main app structure
func New(config *Config, baseDir string) *App {

	return &App{
		baseDir:  baseDir,
		config:   config,
		web:      &http.Server{},
		restart:  make(chan struct{}),
		shutdown: make(chan struct{}),
	}
}

// Run starts the application.
func (app *App) Run() error {
	if err := app.Init(); err != nil {
		return err
	}

	// start your application here

	app.runSignalHandler()

	slog.Info("Starting web server", "url", net.JoinHostPort(app.config.Webserver.ListenHost, app.config.Webserver.ListenPort))
	err := app.runWebServer()
	if err != nil {
		return err
	}

	slog.Info(fmt.Sprintf("%s started", MODULE), "version", VERSION, "pid", os.Getpid())

	return nil
}

// Init initializes the application.
func (app *App) Init() (err error) {

	// init your application here

	// initRoutes should always be called at the end
	app.InitRoutes()

	return nil
}

// Restart returns the read-only restart channel.
// Restart is used to be able to react on application restart. (see cmd/main.go)
func (app *App) Restart() <-chan struct{} {
	return app.restart
}

// Shutdown returns the read only shutdown channel.
// Shutdown is used to be able to react to application shutdown (see cmd/main.go)
func (app *App) Shutdown() <-chan struct{} {
	return app.shutdown
}

// runSignalHandler runs the os signal handler to react on os signals (SIGHUP, SIGTERM, SIGINT).
func (app *App) runSignalHandler() {

	go func() {

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT)

		slog.Info("Starting signal handler")

		switch <-sig {

		case syscall.SIGHUP:
			slog.Info("SIGHUP received, restarting")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := app.web.Shutdown(ctx); err != nil {
				slog.Error("Shutdown failed", "error", err)
			}
			if err := app.cleanup(); err != nil {
				slog.Error("Cleanup failed", "error", err)
			}
			signal.Reset()
			slog.Info("Shutdown completed, restarting")
			app.restart <- struct{}{}
			close(sig)

		case syscall.SIGTERM:
			slog.Info("SIGTERM received, gracefully shutting down")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := app.web.Shutdown(ctx); err != nil {
				slog.Error("Shutdown failed", "error", err)
			}
			if err := app.cleanup(); err != nil {
				slog.Error("Cleanup failed", "error", err)
			}
			slog.Info("Shutdown completed, exiting")
			slog.Log(context.Background(), slog.LevelInfo, fmt.Sprintf("%s stopped", MODULE), "version", VERSION, "pid", os.Getpid())
			app.shutdown <- struct{}{}
			close(app.shutdown)
			close(app.restart)
			close(sig)

		case syscall.SIGINT:
			slog.Info("SIGINT received, exiting")
			_ = app.cleanup()
			slog.Log(context.Background(), slog.LevelInfo, fmt.Sprintf("%s stopped", MODULE), "version", VERSION, "pid", os.Getpid())
			app.shutdown <- struct{}{}
			close(app.shutdown)
			close(app.restart)
			close(sig)
		}

	}()
}

// cleanup free's application resources.
// It's called when application is shutdown or restarted.
// Should be used to free up resources.
func (app *App) cleanup() error {

	//err := app.meters.Close()
	return nil
}
