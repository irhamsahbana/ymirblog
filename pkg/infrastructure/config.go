// Package infrastructure is implements an adapter to talk to low-level modules.
package infrastructure

import (
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/shared/config"
)

// Development is debugging the constant.
const Development string = "development"

// Config is data structure for dynamic configuration and env variables.
type Config struct {
	Ports struct {
		Grpc   int `yaml:"grpc" env:"GRPC_PORT" env-description:"port for grpc"`
		HTTP   int `yaml:"http" env:"HTTP_PORT" env-description:"port for http"`
		HTTPS  int `yaml:"https" env:"HTTPS_PORT" env-description:"port for https"`
		Metric int `yaml:"metric" env:"METRIC_PORT" env-description:"port for metric"`
	} `yaml:"Ports"`
	App struct {
		Environment string `yaml:"environment" env:"ENV_STAGE"`
		ServiceName string `yaml:"serviceName"`
		Logger      string `yaml:"logger" env:"LOG_LEVEL" env-description:"log level debug, info, warn, error, fatal or panic"`
	} `yaml:"App"`
	Server struct {
		Timeout time.Duration `yaml:"timeout" env:"SERVER_TIMEOUT" env-description:"server timeout"`
	} `yaml:"Server"`
	Telemetry struct {
		CollectorEnable   bool   `yaml:"collector_enable" env:"COLLECTOR_ENABLE" env-description:"exporter tracing monitoring"`
		CollectorDebug    bool   `yaml:"collector_debug" env:"COLLECTOR_DEBUG" env-description:"exporter debug collector"`
		CollectorGrpcAddr string `yaml:"collector_grpc_addr" env:"COLLECTOR_GRPC_ADDR" env-description:"exporter addr tracing monitoring"`
	} `yaml:"Telemetry"`
	DB struct {
		ConnectionTimeout int `yaml:"connection_timeout" env:"CONN_TIMEOUT" env-description:"database timeout"`
		MaxOpenCons       int `yaml:"max_open_cons" env:"MAX_OPEN_CONS" env-description:"database max open conn"`
		MaxIdleCons       int `yaml:"max_idle_cons" env:"MAX_IDLE_CONS" env-description:"database max idle conn"`
		ConnMaxLifetime   int `yaml:"conn_max_lifetime" env:"CONN_MAX_LIFETIME" env-description:"database max lifetime"`
	} `yaml:"DB"`
	YmirBlogMySQL struct {
		Database string `yaml:"database" env:"YMIRBLOG_MYSQL_DATABASE" env-description:"database name"`
		User     string `yaml:"user" env:"YMIRBLOG_MYSQL_USER" env-description:"database user"`
		Password string `yaml:"password" env:"YMIRBLOG_MYSQL_PASSWORD" env-description:"database password"`
		Host     string `yaml:"host" env:"YMIRBLOG_MYSQL_HOST" env-description:"database host"`
		Port     uint16 `yaml:"port" env:"YMIRBLOG_MYSQL_PORT" env-description:"database port"`
	} `yaml:"YmirBlogMySQL"`
}

var (
	Envs *Config // Envs is global vars Config.
	once sync.Once
)

// Option is Configure type return func.
type Option = func(c *Configure) error

// Configure is the data struct.
type Configure struct {
	path     string
	filename string
}

// Configuration create instance.
func Configuration(opts ...Option) *Configure {
	c := &Configure{}

	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			panic(err)
		}
	}
	return c
}

// Initialize will create instance of Configure.
func (c *Configure) Initialize() {
	once.Do(func() {
		Envs = &Config{}
		if err := config.Load(config.Opts{
			Config:    Envs,
			Paths:     []string{c.path},
			Filenames: []string{c.filename},
		}); err != nil {
			log.Error().Err(err).Msg("get config error")
		}
	})
}

// WithPath will assign to field path Configure.
func WithPath(path string) Option {
	return func(c *Configure) error {
		c.path = path
		return nil
	}
}

// WithFilename will assign to field name Configure.
func WithFilename(name string) Option {
	return func(c *Configure) error {
		c.filename = name
		return nil
	}
}
