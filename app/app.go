package app

import (
	"context"
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
		baseDir: baseDir,
		config:  config,
		web:     &http.Server{},

		restart:  make(chan struct{}),
		shutdown: make(chan struct{}),
	}
}

// Run starts the application.
//   - Initialize the application.
//   - start the web server.
func (app *App) Run() (*App, error) {
	slog.Info("Initializing application")

	if err := app.Init(); err != nil {
		return app, err
	}

	// handle the OS signals
	app.HandleOSSignals()

	webServerAddress := net.JoinHostPort(app.config.HttpsServer.ListenHost, app.config.HttpsServer.ListenPort)
	slog.Info("Starting web server", "url", webServerAddress)
	err := app.StartWebServer()
	if err != nil {
		slog.Error("Web server failed to start", "url", webServerAddress, "error", err)
		return app, err
	}

	slog.Info(fmt.Sprintf("%s started successfully", MODULE), "version", VERSION, "pid", os.Getpid())
	return app, nil
}

// Init is called by Run() and should be used to initialize the application.
func (app *App) Init() (err error) {

	// initRoutes should always be called at the end
	slog.Info("Initializing API routes")
	app.InitRoutes()

	return nil
}

// Restart returns the read-only restart channel.
// Restart is used to be able to react on application restart.
func (app *App) Restart() <-chan struct{} {
	return app.restart
}

// Shutdown returns the read-only shutdown channel.
// Shutdown is used to be able to react to application shutdown.
func (app *App) Shutdown() <-chan struct{} {
	return app.shutdown
}

// HandleOSSignals runs the os signal handler to react on os signals (SIGHUP, SIGTERM, SIGINT).
func (app *App) HandleOSSignals() {

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT)

		slog.Info("Starting signal handler")

		receivedSignal := <-sig
		slog.Warn("Received OS signal", "signal", receivedSignal)

		switch receivedSignal {
		case syscall.SIGHUP:
			slog.Info("SIGHUP received, initiating restart")
			app.shutdownProcedure("restart")
			// reset the signal registration before the program restarts.
			// with program restarts, the HandleOSSignals function is called again and re-registers the signals.
			signal.Reset()

		case syscall.SIGTERM:
			slog.Info("SIGTERM received, gracefully shutting down")
			app.shutdownProcedure("shutdown")

		case syscall.SIGINT:
			slog.Info("SIGINT received, exiting")
			app.shutdownProcedure("terminate")
		}
	}()
}

// shutdownProcedure Handles SIGTERM, SIGINT and SIGHUP (restart) for a graceful shutdown.
//   - terminate: Cleanup app resources and terminates the application.
//   - shutdown: graceful shutdown the web server, Cleanup app resources and exit the application.
//   - restart: graceful shutdown the web server and Cleanup app resources and restart the application.
func (app *App) shutdownProcedure(mode string) {
	slog.Info("Initiating shutdown", "mode", mode)

	if mode == "shutdown" || mode == "restart" {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := app.web.Shutdown(ctx); err != nil {
			slog.Error("Web server shutdown failed", "error", err)
		}
	}

	if err := app.Cleanup(); err != nil {
		slog.Error("Cleanup failed", "error", err)
	}

	if mode == "restart" {
		slog.Info("Shutdown completed, preparing to restart")
		app.restart <- struct{}{}
		return
	}

	slog.Info(fmt.Sprintf("%s stopped", MODULE), "version", VERSION, "pid", os.Getpid())
	app.shutdown <- struct{}{}
	close(app.shutdown)
	close(app.restart)
}

// Cleanup free's application resources.
// It's called when application is shutdown or restarted.
// Should be used to free up resources.
func (app *App) Cleanup() error {
	var err error
	return err
}
