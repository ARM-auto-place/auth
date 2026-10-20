package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	goRedis "github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	"github.com/ARM-auto-place/auth/internal/client"
	"github.com/ARM-auto-place/auth/internal/config"
	"github.com/ARM-auto-place/auth/internal/service"
	"github.com/ARM-auto-place/auth/internal/storage/postgres"
	storageRedis "github.com/ARM-auto-place/auth/internal/storage/redis"
	"github.com/ARM-auto-place/auth/internal/transport/jsonRPC/externalapi"
)

func main() {

	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(".env"); err != nil {
			panic("Error loading .env file: " + err.Error())
		}
	}

	cfg := config.Values()
	logger := cfg.Logger()
	log.Logger = logger

	// postgres
	connectManager, err := postgres.New(logger, cfg.Postgres)
	if err != nil {
		logger.Err(err).Msg("Failed to connect to postgres")
		return
	}
	defer func() {
		connectManager.GetMasterConn().Close()
		connectManager.GetReplicaConn().Close()
	}()
	authStorage := postgres.NewAuthStorage(logger, connectManager)

	// redis
	redisClient := goRedis.NewUniversalClient(&goRedis.UniversalOptions{
		Addrs:      cfg.Redis.Addrs,
		Username:   cfg.Redis.Username,
		Password:   cfg.Redis.Password,
		DB:         cfg.Redis.DB,
		MaxRetries: cfg.Redis.MaxRetries,
		MasterName: cfg.Redis.MasterName,
	})
	defer redisClient.Close()

	cacheAuth, err := storageRedis.New(logger, redisClient)
	if err != nil {
		logger.Err(err).Msg("Failed to connect to redis")
		return
	}

	// уведомления на почту
	mailer := client.NewSMTPMailer(
		logger,
		cfg.SMTP.Host,
		cfg.SMTP.Port,
		cfg.SMTP.Username,
		cfg.SMTP.Password,
		cfg.SMTP.From,
		cfg.SMTP.TLSEnabled,
		cfg.SMTP.Timeout,
	)

	authService := service.NewAuthService(logger, authStorage, mailer, cacheAuth)

	app := externalapi.New(
		logger,
		externalapi.SetReadBufferSize(cfg.MaxRequestHeaderSize),
		externalapi.MaxBodySize(cfg.MaxRequestBodySize),
		externalapi.ReadTimeout(cfg.ReadTimeout),
		externalapi.WriteTimeout(cfg.WriteTimeout),
		externalapi.WithRequestID("X-Request-Id"),
		externalapi.Use(cors.New(cors.Config{
			AllowOriginsFunc: isAllowedOrigin(cfg.AllowedOrigins),
			AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
			AllowHeaders:     "Origin,Content-Type,Accept,Authorization,Cookie,X-Requested-With,X-Log-Level",
			AllowCredentials: true,
		})),
		externalapi.ExternalAPI(externalapi.NewExternalAPI(authService)),
	).WithLog().WithMetrics()

	app.ServeHealth(cfg.HealthBind, map[string]string{"status": "ok"})
	app.ServeMetrics(logger, cfg.MetricsPath, cfg.MetricsBind)

	shutdownCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Info().Msgf("HTTP server started address:%s", cfg.ServiceBind)
		if err := app.Fiber().Listen(cfg.ServiceBind); err != nil {
			logger.Err(err).Msg("Error serving HTTP server")
			stop()
		}
	}()

	<-shutdownCtx.Done()
	logger.Info().Msg("Server stopping")

	done := make(chan struct{})
	go func() {
		app.Shutdown()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(cfg.ShutdownTimeout):
		logger.Warn().Msg("Server shutdown timed out")
	}
	logger.Info().Msg("Server stopped")
}

func isAllowedOrigin(allowedOrigins []string) func(origin string) bool {
	return func(origin string) bool {
		for _, allowed := range allowedOrigins {
			if allowed == origin {
				return true
			}
		}
		return false
	}
}
