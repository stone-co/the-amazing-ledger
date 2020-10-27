package configuration

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ShutdownTimeout time.Duration `envconfig:"APP_SHUTDOWN_TIMEOUT" default:"5s"`
	API             HTTPConfig
	GRPC            GRPCConfig
	Postgres        PostgresConfig
}

func LoadConfig() (*Config, error) {
	var config Config
	noPrefix := ""
	err := envconfig.Process(noPrefix, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

type HTTPConfig struct {
	Port int `envconfig:"HTTP_PORT" default:"3000"`
}

type GRPCConfig struct {
	Port            int           `envconfig:"GRPC_PORT" default:"50051"`
	ShutdownTimeout time.Duration `envconfig:"GRPC_SHUTDOWN_TIMEOUT" default:"5s"`
}

type PostgresConfig struct {
	DatabaseName string `envconfig:"DATABASE_NAME" default:"dev"`
	User         string `envconfig:"DATABASE_USER" default:"postgres"`
	Password     string `envconfig:"DATABASE_PASSWORD" default:"postgres"`
	Host         string `envconfig:"DATABASE_HOST" default:"localhost"`
	Port         string `envconfig:"DATABASE_PORT" default:"5432"`
	PoolSize     string `envconfig:"DATABASE_POOL_SIZE" default:"10"`
	SSLMode      string `envconfig:"DATABASE_SSLMODE" default:"disable"`
	SSLRootCert  string `envconfig:"DATABASE_SSL_ROOTCERT"`
	SSLCert      string `envconfig:"DATABASE_SSL_CERT"`
	SSLKey       string `envconfig:"DATABASE_SSL_KEY"`
}

func (c PostgresConfig) DSN() string {
	connectString := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s pool_max_conns=%s",
		c.User, c.Password, c.Host, c.Port, c.DatabaseName, c.PoolSize)

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
