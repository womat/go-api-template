package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/womat/go-api-template/app"
	"github.com/womat/golib/xlog"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
)

//go:embed README.md
var Readme string

func main() {
	// Parse command line flags.
	flags := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flags.SetOutput(os.Stdout)

	about := flags.Bool("about", false, "Print app details and exit")
	help := flags.Bool("help", false, "Print a help message and exit")
	version := flags.Bool("version", false, "Print the app version and exit")
	debug := flags.Bool("debug", false, "Enable debug logging to stdout (overrides log settings from the config file)")
	configFile := flags.String("config", filepath.Join("/opt", app.MODULE, "etc", "config.yaml"), "Specify the path to the config file")

	if err := flags.Parse(os.Args[1:]); err != nil {
		fmt.Printf("error: %s\n", err.Error())
		os.Exit(1)
	}

	switch {
	case *about:
		fmt.Println(About())
		os.Exit(0)
	case *version:
		fmt.Println(app.VERSION)
		os.Exit(0)
	case *help:
		fmt.Println(Readme)
		os.Exit(0)
	}

	var logger *xlog.LoggerWrapper

	config, err := loadConfig(*configFile, *debug)
	if err != nil {
		fmt.Printf("Failed to load config file %s: %s\n", *configFile, err.Error())
		os.Exit(1)
	}

	for {
		// run the app in a function to be able to restart it and reload the config
		// possible open log files are always closed before the function exits
		func() {
			if logger, err = xlog.Init(config.LogDestination, config.LogLevel); err != nil {
				fmt.Printf("Failed to initialize logger: %s\n", err.Error())
				os.Exit(1)
			}
			defer logger.Close()

			// set slog logger as default logger
			slog.SetDefault(logger.Logger)
			slog.Info("Logging initialized", "logLevel", config.LogLevel)
			slog.Debug("Starting with configuration", "config", config)

			a, err := app.New(config).Run()
			if err != nil {
				slog.Error("Critical error occurred, shutting down", "error", err)
				os.Exit(1)
			}

			select {
			case <-a.Restart():
				slog.Info("Reload configuration", "configFile", *configFile)
				if config, err = loadConfig(*configFile, *debug); err != nil {
					slog.Error("Failed to reload config file, shutting down",
						"configFile", *configFile,
						"error", err)
					os.Exit(1)
				}

			case <-a.Shutdown():
				os.Exit(0)
			}
		}()
	}
}

func About() string {
	p := map[string]string{
		"Author":   "Wolfgang Mathe",
		"Binary":   filepath.Join("/opt", app.MODULE, "bin", app.MODULE),
		"Date":     "2025-02-24",
		"Desc":     "Blueprint for Go applications",
		"Help":     filepath.Join("/opt", app.MODULE, "bin", app.MODULE) + " --help",
		"Libinfo":  "plain go with go modules from ITdesign golib",
		"Main":     filepath.Join("/opt/src", app.MODULE, "cmd", app.MODULE, "main.go"),
		"ProgLang": runtime.Version(),
		"Repo":     "https://github.com/womat/" + app.MODULE + ".git",
		"Version":  app.VERSION,
	}
	b, _ := yaml.Marshal(p)
	return string(b)
}

// loadConfig loads the configuration from the given file.
func loadConfig(configFile string, debug bool) (*app.Config, error) {

	config, err := app.NewConfig().LoadConfig(configFile)
	if err != nil {
		return nil, err
	}

	if debug {
		config.LogLevel = "debug"
		config.LogDestination = "stdout"
	}

	return config, nil
}
