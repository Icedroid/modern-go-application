package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/sagikazarmark/modern-go-application/internal/platform/database"
	"github.com/sagikazarmark/modern-go-application/internal/platform/jaeger"
	"github.com/sagikazarmark/modern-go-application/internal/platform/log"
	"github.com/sagikazarmark/modern-go-application/internal/platform/prometheus"
	"github.com/sagikazarmark/modern-go-application/internal/platform/redis"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config holds any kind of configuration that comes from the outside world and
// is necessary for running the application.
type Config struct {
	// Meaningful values are recommended (eg. production, development, staging, release/123, etc)
	Environment string

	// Turns on some debug functionality (eg. more verbose logs)
	Debug bool

	// Timeout for graceful shutdown
	ShutdownTimeout time.Duration

	// Log configuration
	Log log.Config

	// Instrumentation configuration
	Instrumentation InstrumentationConfig

	// App configuration
	App struct {
		// App server address
		Addr string
	}

	// Database connection information
	Database database.Config

	// Redis configuration
	Redis redis.Config
}

// Validate validates the configuration.
func (c Config) Validate() error {
	if c.Environment == "" {
		return errors.New("environment is required")
	}

	if err := c.Log.Validate(); err != nil {
		return err
	}

	if err := c.Instrumentation.Validate(); err != nil {
		return err
	}

	if c.App.Addr == "" {
		return errors.New("app server address is required")
	}

	if err := c.Database.Validate(); err != nil {
		return err
	}

	if err := c.Redis.Validate(); err != nil {
		return err
	}

	return nil
}

type InstrumentationConfig struct {
	// Instrumentation HTTP server address
	Addr string

	// Prometheus configuration
	Prometheus struct {
		Enabled           bool
		prometheus.Config `mapstructure:",squash"`
	}

	// Jaeger configuration
	Jaeger struct {
		Enabled       bool
		jaeger.Config `mapstructure:",squash"`
	}
}

// Validate validates the configuration.
func (c InstrumentationConfig) Validate() error {
	if c.Addr == "" {
		return errors.New("instrumentation http server address is required")
	}

	if c.Prometheus.Enabled {
		if err := c.Prometheus.Validate(); err != nil {
			return err
		}
	}

	if c.Jaeger.Enabled {
		if err := c.Jaeger.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// Configure configures some defaults in the Viper instance.
func Configure(v *viper.Viper, p *pflag.FlagSet) {
	v.AllowEmptyEnv(true)
	v.AddConfigPath(".")
	p.Init(FriendlyServiceName, pflag.ExitOnError)
	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", FriendlyServiceName)
		pflag.PrintDefaults()
	}
	v.BindPFlags(p) // nolint:errcheck

	// Global configuration
	v.SetDefault("environment", "production")
	v.SetDefault("debug", false)
	v.SetDefault("shutdownTimeout", 15*time.Second)

	// Log configuration
	v.SetDefault("log.format", "json")
	v.SetDefault("log.level", "info")

	// Instrumentation configuration
	p.String("instrumentation.addr", ":10000", "Instrumentation HTTP server address")
	v.SetDefault("instrumentation.addr", ":10000")

	v.SetDefault("instrumentation.prometheus.enabled", false)
	v.SetDefault("instrumentation.jaeger.enabled", false)
	v.SetDefault("instrumentation.jaeger.endpoint", "http://localhost:14268")
	v.SetDefault("instrumentation.jaeger.agentEndpoint", "localhost:6831")

	// App configuration
	p.String("app.addr", ":8000", "App HTTP server address")
	v.SetDefault("app.addr", ":8000")

	// Database configuration
	v.SetDefault("database.port", 3306)
	v.SetDefault("database.params", map[string]string{
		"charset": "utf8mb4",
	})

	// Redis configuration
	v.SetDefault("redis.port", 6379)
}
