package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/womat/go-api-template/app"
	"github.com/womat/go-api-template/pkg/crypt"
	"github.com/womat/go-api-template/pkg/xlog"
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
	cryptString := flags.String("crypt", "", "Encrypt the given string and exit")
	help := flags.Bool("help", false, "Print a help message and exit")
	version := flags.Bool("version", false, "Print the app version and exit")

	logLevel := flags.String("logLevel", "", "Set the log level (overrides the config file). Supported values: trace | debug | info | warning | error")
	logDestination := flags.String("logDestination", "", "Set the log destination (overrides the config file). Supported values: stdout | stderr | null | /path/to/logfile")
	configFile := flags.String("config", app.DefaultConfigFile, "Specify the path to the config file")

	if err := flags.Parse(os.Args[1:]); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}

	// Print about details and exit when about flag is used.
	if *about {
		printAbout()
		os.Exit(0)
	}

	switch {
	case *about:
		printAbout()
		os.Exit(0)
	case *cryptString != "":
		fmt.Println(crypt.NewEncryptedString(*cryptString).EncryptedValue())
		os.Exit(0)
	case *version:
		fmt.Println(app.VERSION)
		os.Exit(0)
	case *help:
		fmt.Println(Readme)
		os.Exit(0)
	}

	var logger *xlog.LoggerWrapper

	config, err := loadConfig(*configFile, *logLevel, *logDestination)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to load config file %s: %s\n", *configFile, err.Error())
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

			a, err := app.New(config, filepath.Join("/opt", app.MODULE)).Run()
			if err != nil {
				slog.Error("Critical error occurred, shutting down", "error", err)
				os.Exit(1)
			}

			select {
			case <-a.Restart():
				slog.Info("Reload configuration", "configFile", *configFile)
				if config, err = loadConfig(*configFile, *logLevel, *logDestination); err != nil {
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

func printAbout() {
	type ProgInfo struct {
		Author   string `yaml:"author"`
		Binary   string `yaml:"binary"`
		Comment  string `yaml:"comment"`
		Date     string `yaml:"date"`
		Desc     string `yaml:"desc"`
		Help     string `yaml:"help"`
		Libinfo  string `yaml:"libinfo"`
		Main     string `yaml:"main"`
		ProgLang string `yaml:"progLang"`
		Repo     string `yaml:"repo"`
		Version  string `yaml:"version"`
	}

	var p = ProgInfo{
		Author:   "Wolfgang Mathe",
		Binary:   "/opt/<MODUL_NAME>/bin/<MODUL_NAME>",
		Comment:  "config .env file see /opt/<MODUL_NAME>/.env  and config file /opt/<MODUL_NAME>/etc/config.yaml",
		Date:     "2024-10-04",
		Desc:     "Blueprint for Go applications",
		Help:     "/opt/<MODUL_NAME>/bin/<MODUL_NAME> --help",
		Libinfo:  "plain go with go modules from ITdesign golib",
		Main:     "/opt/src/<MODUL_NAME>/cmd/<MODUL_NAME>/main.go",
		ProgLang: runtime.Version(),
		Repo:     " https://github.com/womat/<MODUL_NAME>.git",
		Version:  app.VERSION,
	}
	b, _ := yaml.Marshal(p)
	fmt.Print(string(b))
}

// loadConfig loads the configuration from the given file.
func loadConfig(configFile, logLevel, logDestination string) (*app.Config, error) {

	config, err := app.NewConfig().LoadConfig(configFile)
	if err != nil {
		return nil, err
	}

	switch logLevel {
	case "": // if no log level is provided, use the one from the config
	case "debug", "trace", "info", "warning", "error", "err", "warn":
		config.LogLevel = logLevel
	default:
		return nil, fmt.Errorf("invalid log level: %s", logLevel)
	}

	// Set log destination if provided
	if logDestination != "" {
		config.LogDestination = logDestination
	}

	return config, nil
}
