package app

import (
	"github.com/womat/go-api-template/pkg/crypt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

const (
	ProdEnv = "prod"
	DevEnv  = "dev"

	DefaultEnv            = ProdEnv
	DefaultListenHost     = ""
	DefaultListenPort     = "443"
	DefaultMinTLS         = "1.2"
	DefaultLogLevel       = "info"
	DefaultLogDestination = "stdout"
	DefaultKeyFile        = "key.pem"
	DefaultCertFile       = "cert.pem"
	DefaultJwtSecret      = ""
	DefaultJwtID          = ""

	// DefaultAppKey is the default api key for the application (app_name --crypt secure_api_key)
	DefaultAppKey = "iIHj841o6kj1EaCfJqUVNss1oei6EIIiE+TJaF9TwNbD9rwJOmCHxYcx"
)

var (
	AppDir            = filepath.Join("/opt", MODULE)
	DefaultConfigFile = filepath.Join(AppDir, "etc/config.yaml")
)

// Config holds the application configuration
type Config struct {
	// Env is the app environment.
	// Env is read from APP_ENV environment variable.
	//  Allowed values: prod | dev
	//  It's used for:
	//  - jwt token expiration (1 day in dev, 5 minutes in prod)
	// Default is "prod".
	Env string

	// LogLevel is the log level, if set only message with at least this level is logged
	//  e.g.: debug -> means error, warning, info and debug messages are logged
	// Allowed values: debug | info | warning | error | trace
	// Default is info.
	LogLevel string `yaml:"logLevel"`

	// LogDestination defines the log destinations.
	//  supported values: stdout | stderr | /path/to/logfile
	LogDestination string `yaml:"logDestination"`

	// ApiKey is the global api key for the application.
	// ApiKey must be encrypted with "app --crypt <plaintext>"
	// Default is empty that means api key authentication is disabled.
	ApiKey crypt.EncryptedString `yaml:"apiKey"`

	// JwtSecret is a secret key used to sign jwt tokens.
	// JwtSecret must be encrypted with "app --crypt <plaintext>"
	JwtSecret crypt.EncryptedString `yaml:"jwtSecret"`

	// JwtID is a unique identifier for the jwt token used to prevent login with the same jwt token to another app.
	JwtID string `yaml:"jwtID"`

	// Webserver is the configuration of the webserver and webservice
	Webserver WebserverConfig `yaml:"webserver"`

	// add your application-specific configuration here
}

// WebserverConfig defines the struct of the webserver and webservice configuration and configuration file
type WebserverConfig struct {
	// ListenHost is the host address the https server listens for connections.
	// Default is empty, which means all available network interfaces.
	ListenHost string `yaml:"listenHost"`

	// ListenPort is the port the https server listens for connections.
	// Default is 443
	ListenPort string `yaml:"listenPort"`

	// MinTLS is the minimum TLS version the server accepts.
	// Default is "1.2"
	MinTLS string `yaml:"minTLS"`

	// KeyFile is the ssl certificate private key file
	// Default is key.pem
	KeyFile string `yaml:"keyFile"`

	// CertFile is the ssl certificate public key file
	// Default is cert.pem
	// Pfx files are supported as well, in which case KeyFile must be empty and CertFile must point to the pfx file, CertPassword must contain the password to decode the pfx file.
	CertFile string `yaml:"certFile"`

	// CertPassword is an optional certificate password that
	// CertPassword must be encrypted with "app --crypt <plaintext>"
	// Default is empty which means no password is used.
	CertPassword crypt.EncryptedString `yaml:"certPassword"`

	// BlockedIPs is a list of IP addresses or networks that are forbidden from accessing the application.
	// Default is empty, which means no IP addresses or networks are blocked.
	// Multiple IP addresses or networks can be defined separated by a comma
	// e.g.: 192.168.0.1,192.168.0.0/16,10.0.0.0/8,192.168.254.15
	BlockedIPs []string `yaml:"blockedIPs"`

	// AllowedIPs is a list of IP addresses that are allowed to access the application.
	// Default is empty, which means all IP addresses are allowed.
	// The value "ALL" allows access from all IP Addresses / IP Networks
	// multiple IP addresses or networks can be defined separated by a comma
	// e.g.: 127.0.0.1,::1,192.168.0.0/16,10.0.0.0/8
	// Note: '::1' is the IPv6 loopback address.
	AllowedIPs []string `yaml:"allowedIPs"`
}

func NewConfig() *Config {
	return &Config{
		Env:            DefaultEnv,
		LogLevel:       DefaultLogLevel,
		LogDestination: DefaultLogDestination,

		JwtSecret: crypt.NewEncryptedString(DefaultJwtSecret),
		JwtID:     DefaultJwtID,

		Webserver: WebserverConfig{
			ListenHost:   DefaultListenHost,
			ListenPort:   DefaultListenPort,
			MinTLS:       DefaultMinTLS,
			CertFile:     filepath.Join(AppDir, "etc", DefaultCertFile),
			KeyFile:      filepath.Join(AppDir, "etc", DefaultKeyFile),
			CertPassword: crypt.NewEncryptedString(""),
			BlockedIPs:   []string{},
			AllowedIPs:   []string{},
		},
	}
}

func (c *Config) LoadConfig(fileName string) (*Config, error) {

	content, err := os.ReadFile(fileName)
	if err != nil {
		return c, err
	}

	replaced := os.ExpandEnv(string(content))
	err = yaml.Unmarshal([]byte(replaced), c)
	return c, err
}

// IsDevEnv returns true if "dev" is configured as app environment.
func (c *Config) IsDevEnv() bool {
	return c.Env == DevEnv
}
