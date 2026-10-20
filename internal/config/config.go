package config

import (
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ARM-auto-place/auth/internal/metrics"
	"github.com/caarlos0/env/v10"
	"github.com/rs/zerolog"
)

type PostgresConfig struct {
	Addr            string `env:"PG_ADDR"`
	DB              string `env:"PG_DB"`
	User            string `env:"PG_USER"`
	Password        string `env:"PG_PASSWORD"`
	MaxConn         int    `env:"PG_MAX_CONN" envDefault:"10"`
	MaxIdleLifetime string `env:"PG_MAX_IDLE_LIFETIME" envDefault:"30s"`
	MaxLifetime     string `env:"PG_MAX_LIFETIME" envDefault:"3m"`
	PrepareCacheCap int    `env:"PG_PREPARE_CACHE_CAP" envDefault:"128"`
	ReplicaAddr     string `env:"PG_REPLICA_ADDR"`
	UserRO          string `env:"PG_REPLICA_USER"`
	PasswordRO      string `env:"PG_REPLICA_PASSWORD"`
}

type RedisConfig struct {
	Addrs      []string `env:"REDIS_ADDRS" envDefault:"localhost:6379" envSeparator:","`
	Username   string   `env:"REDIS_USERNAME" envDefault:""`
	Password   string   `env:"REDIS_PASSWORD" envDefault:""`
	DB         int      `env:"REDIS_DB" envDefault:"0"`
	MaxRetries int      `env:"REDIS_MAX_RETRIES" envDefault:"5"`
	MasterName string   `env:"REDIS_MASTER_NAME"` // заполняется при использовании sentinel
}

type SMTPConfig struct {
	Host       string        `env:"SMTP_HOST"`
	Port       string        `env:"SMTP_PORT" envDefault:"587"`
	Username   string        `env:"SMTP_USERNAME"`
	Password   string        `env:"SMTP_PASSWORD"`
	From       string        `env:"SMTP_FROM"`
	TLSEnabled bool          `env:"SMTP_TLS_ENABLED" envDefault:"true"`
	Timeout    time.Duration `env:"SMTP_TIMEOUT" envDefault:"10s"`
}

type Config struct {
	ServiceName string `env:"SERVICE_NAME" envDefault:"auth"`

	LogLevel  string `env:"LOGGER_LEVEL" envDefault:"debug"`
	LogFormat string `env:"LOGGER_FORMAT" envDefault:""`

	// common env vars
	ServiceBind    string   `env:"BIND_ADDR" envDefault:":9000"`
	HealthBind     string   `env:"BIND_HEALTH" envDefault:":9091"`
	MetricsBind    string   `env:"BIND_METRICS" envDefault:":9090"`
	MetricsPath    string   `env:"METRICS_PATH" envDefault:"/metrics"`
	AllowedOrigins []string `env:"ALLOWED_ORIGINS" envSeparator:","`

	MaxRequestBodySize   int           `env:"MAX_REQUEST_BODY_SIZE" envDefault:"104857600"` // 100 MB
	MaxRequestHeaderSize int           `env:"MAX_REQUEST_HEADER_SIZE" envDefault:"16384"`   // 16 KB
	ReadTimeout          time.Duration `env:"READ_TIMEOUT" envDefault:"120s"`
	WriteTimeout         time.Duration `env:"WRITE_TIMEOUT" envDefault:"30s"`
	ShutdownTimeout      time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"10s"`

	Postgres PostgresConfig
	Redis    RedisConfig
	SMTP     SMTPConfig
}

var (
	configuration *Config
	configMu      sync.Once
)

func internalConfig() *Config {
	configMu.Do(func() {
		configuration = &Config{}
		if err := env.Parse(configuration); err != nil {
			panic(err)
		}

		// реплика необязательна: если не задана, читаем из мастера
		if configuration.Postgres.ReplicaAddr == "" {
			configuration.Postgres.ReplicaAddr = configuration.Postgres.Addr
		}
		if configuration.Postgres.UserRO == "" {
			configuration.Postgres.UserRO = configuration.Postgres.User
			configuration.Postgres.PasswordRO = configuration.Postgres.Password
		}

		for i, addr := range configuration.Redis.Addrs {
			configuration.Redis.Addrs[i] = strings.TrimSpace(addr)
		}
		for i, origin := range configuration.AllowedOrigins {
			configuration.AllowedOrigins[i] = strings.TrimSpace(origin)
		}
	})
	return configuration
}

func Values() Config {
	return *internalConfig()
}

func (cfg Config) Logger() zerolog.Logger {
	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil || cfg.LogLevel == "" {
		level = zerolog.DebugLevel
	}

	var logger zerolog.Logger
	if cfg.LogFormat == "console" {
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	} else {
		logger = zerolog.New(os.Stdout)
	}

	return logger.Level(level).With().Timestamp().Str("service", cfg.ServiceName).Logger()
}

var (
	metricsValue *metrics.Metrics
	metricsMu    sync.Once
)

func Metrics() *metrics.Metrics {
	metricsMu.Do(func() {
		metricsValue = metrics.NewMetrics(internalConfig().ServiceName)
	})

	return metricsValue
}
