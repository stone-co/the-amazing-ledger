package app

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	RPCServer  RPCServerConfig
	HttpServer HttpServerConfig
	Postgres   PostgresConfig
	NewRelic   NewRelicConfig
}

func LoadConfig() (*Config, error) {
	var config Config
	noPrefix := ""

	err := envconfig.Process(noPrefix, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}

type RPCServerConfig struct {
	Host            string        `envconfig:"GRPC_HOST" default:"0.0.0.0"`
	Port            int           `envconfig:"GRPC_PORT" default:"3000"`
	ShutdownTimeout time.Duration `envconfig:"APP_SHUTDOWN_TIMEOUT" default:"5s"`
	ReadTimeout     time.Duration `envconfig:"GRPC_READ_TIMEOUT" default:"30s"`
	WriteTimeout    time.Duration `envconfig:"GRPC_WRITE_TIMEOUT" default:"10s"`
}

type HttpServerConfig struct {
	Host            string        `envconfig:"HTTP_HOST" default:"0.0.0.0"`
	Port            int           `envconfig:"HTTP_PORT" default:"3001"`
	ShutdownTimeout time.Duration `envconfig:"HTTP_SHUTDOWN_TIMEOUT" default:"1s"`
	ReadTimeout     time.Duration `envconfig:"HTTP_READ_TIMEOUT" default:"30s"`
	WriteTimeout    time.Duration `envconfig:"HTTP_WRITE_TIMEOUT" default:"10s"`
}

type PostgresConfig struct {
	DatabaseName string `envconfig:"DATABASE_NAME" default:"dev"`
	User         string `envconfig:"DATABASE_USER" default:"postgres"`
	Password     string `envconfig:"DATABASE_PASSWORD" default:"postgres"`
	Host         string `envconfig:"DATABASE_HOST" default:"localhost"`
	Port         string `envconfig:"DATABASE_PORT" default:"5432"`
	PoolMinSize  string `envconfig:"DATABASE_POOL_MIN_SIZE" default:"2"`
	PoolMaxSize  string `envconfig:"DATABASE_POOL_MAX_SIZE" default:"10"`
	SSLMode      string `envconfig:"DATABASE_SSLMODE" default:"disable"`
	SSLRootCert  string `envconfig:"DATABASE_SSL_ROOTCERT"`
	SSLCert      string `envconfig:"DATABASE_SSL_CERT"`
	SSLKey       string `envconfig:"DATABASE_SSL_KEY"`
}

type NewRelicConfig struct {
	AppName    string `envconfig:"NEW_RELIC_APP_NAME"`
	LicenseKey string `envconfig:"NEW_RELIC_LICENSE_KEY"`
}

func (c PostgresConfig) DSN() string {
	connectString := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s pool_min_conns=%s pool_max_conns=%s",
		c.User, c.Password, c.Host, c.Port, c.DatabaseName, c.PoolMinSize, c.PoolMaxSize)

	if c.SSLMode != "" {
		connectString = fmt.Sprintf("%s sslmode=%s",
			connectString, c.SSLMode)
	}

	if c.SSLRootCert != "" {
		connectString = fmt.Sprintf("%s sslrootcert=%s sslcert=%s sslkey=%s",
			connectString, c.SSLRootCert, c.SSLCert, c.SSLKey)
	}

	return connectString
}

func (c PostgresConfig) URL() string {
	if c.SSLMode == "" {
		c.SSLMode = "disable"
	}

	connectString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DatabaseName, c.SSLMode)

	if c.SSLRootCert != "" {
		connectString = fmt.Sprintf("%s&sslrootcert=%s&sslcert=%s&sslkey=%s",
			connectString, c.SSLRootCert, c.SSLCert, c.SSLKey)
	}

	return connectString
}
