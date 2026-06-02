package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"task-service/internal/application/event"
	apptask "task-service/internal/application/task"
	"task-service/internal/infrastructure/kafka"
	"task-service/internal/infrastructure/postgres"
	"task-service/internal/interface/grpc/handler"
	taskserver "task-service/internal/interface/grpc/server"
	"task-service/internal/interface/observability"
	"task-service/internal/utility"
	"time"

	_ "net/http/pprof"

	"go.opentelemetry.io/contrib/bridges/otelzerolog"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/sdk/metric"
)

func main() {
	logger := zerolog.New(os.Stderr).
		With().
		Timestamp().
		Logger()

	defer func() {
		if err := recover(); err != nil {
			logger.Error().
				Stack().
				Msg("recover")
		}
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	tracerProvider, err := observability.InitTracerProvider(
		ctx,
		observability.TracerOptions{
			ServiceName:      utility.GetEnv("TASK_SERVICE_NAME", "service.task"),
			ServiceVersion:   utility.GetEnv("TASK_SERVICE_VERSION", ""),
			CollectorAddress: utility.GetEnv("TRACE_COLLECTOR_ADDRESS", "localhost:4317"),
			TracerTimeout:    utility.GetEnvDuration("TRACE_TIMEOUT", 5*time.Second),
			Env:              utility.GetEnv("ENV", "dev"),
		},
	)
	if err != nil {
		logger.Fatal().
			Err(err).
			Msg("observability.InitTracer")
	}

	defer func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			logger.Error().
				Err(err).
				Msg("tracerProvider.Shutdown")
		}
	}()

	meterProvider, err := observability.InitMetricsProvider(
		ctx,
		observability.MetricsProviderOptions{
			MetricsExporterAddress:      utility.GetEnv("METRICS_EXPORTER_ADDRESS", "localhost:4317"),
			MetricsExporterTimeout:      utility.GetEnvDuration("METRICS_EXPORTER_TIMEOUT", 5*time.Second),
			MeterPeriodicReaderInterval: utility.GetEnvDuration("METRICS_EXPORTER_READER_INTERVAL", 5*time.Second),
			MeterPeriodicReaderTimeout:  utility.GetEnvDuration("METRICS_EXPORTER_READER_TIMEOUT", 5*time.Second),
		},
	)
	if err != nil {
		logger.Fatal().
			Err(err).
			Msg("observability.NewMetricsProvider")
	}

	defer func() {
		if err := meterProvider.Shutdown(ctx); err != nil {
			logger.Error().
				Err(err).
				Msg("meterProvider.Shutdown")
		}
	}()

	loggerProvider, err := observability.InitLoggerProvider(
		ctx,
		observability.LoggerProviderOptions{
			ServiceName:           utility.GetEnv("TASK_SERVICE_NAME", "service.task"),
			ServiceVersion:        utility.GetEnv("TASK_SERVICE_VERSION", ""),
			Env:                   utility.GetEnv("ENV", "dev"),
			LoggerExporterAddress: utility.GetEnv("LOGGER_EXPORTER_ADDRESS", "localhost:4317"),
			LoggerExporterTimeout: utility.GetEnvDuration("LOGGER_EXPORTER_TIMEOUT", 5*time.Second),
		},
	)
	if err != nil {
		logger.Fatal().
			Err(err).
			Msg("observability.NewLoggerProvider")
	}

	defer func() {
		if err := loggerProvider.Shutdown(ctx); err != nil {
			logger.Error().
				Err(err).
				Msg("loggerProvider.Shutdown")
		}
	}()

	otelLoggerHook := otelzerolog.NewHook("task-service", otelzerolog.WithLoggerProvider(loggerProvider))
	otelLogger := new(logger.Hook(otelLoggerHook))

	pool, err := postgres.NewPool(
		ctx,
		postgres.PoolOptions{
			ConnectionString:      utility.GetEnv("PG_CONNECTION_STRING", ""),
			MinConnections:        utility.GetEnvInt("PG_MIN_CONNECTIONS", 10),
			MaxConnections:        utility.GetEnvInt("PG_MAX_CONNECTIONS", 100),
			MaxIdleConnectionTime: utility.GetEnvDuration("PG_MAX_IDLE_CONNECTION_TIME", 30*time.Minute),
		},
	)
	if err != nil {
		logger.Fatal().
			Err(err).
			Msg("postgres.NewPool")
	}

	defer pool.Close()

	langDetector := postgres.NewLanguageDetector(utility.GetEnvFloat("LANG_DETECTION_ACCURACY", 0.9))

	repository, err := postgres.NewTaskRepository(pool, *langDetector)
	if err != nil {
		logger.Fatal().
			Err(err).
			Msg("postgres.NewTaskRepository")
	}

	searcher := postgres.NewTaskSearcher(pool, *langDetector)
	accessGuard := apptask.NewAccessGuard()

	getHandler, err := apptask.NewGetHandler(repository, accessGuard)
	if err != nil {
		logger.Fatal().
			Err(err).
			Msg("apptask.NewGetHandler")
	}

	searchHandler, err := apptask.NewSearchHandler(searcher, accessGuard)
	if err != nil {
		logger.Fatal().
			Err(err).
			Msg("apptask.NewSearchHandler")
	}

	createHandler, err := apptask.NewCreateHandler(repository, accessGuard)
	if err != nil {
		logger.Fatal().
			Err(err).
			Msg("apptask.NewCreateHandler")
	}

	updateHandler, err := apptask.NewUpdateHandler(repository, accessGuard)
	if err != nil {
		logger.Fatal().
			Err(err).
			Msg("apptask.NewUpdateHandler")
	}

	completeHandler, err := apptask.NewCompleteHandler(repository, accessGuard)
	if err != nil {
		logger.Fatal().
			Err(err).
			Msg("apptask.NewCompleteHandler")
	}

	restoreHandler, err := apptask.NewRestoreHandler(repository, accessGuard)
	if err != nil {
		logger.Fatal().
			Err(err).
			Msg("apptask.NewRestoreHandler")
	}

	deleteHandler, err := apptask.NewDeleteHandler(repository, accessGuard)
	if err != nil {
		logger.Fatal().
			Err(err).
			Msg("apptask.NewDeleteHandler")
	}

	listDeletedHandler, err := apptask.NewListDeletedHandler(repository, accessGuard)
	if err != nil {
		logger.Fatal().
			Err(err).
			Msg("apptask.NewListDeletedHandler")
	}

	grpcHandler, err := handler.NewTaskHandlerBuilder().
		WithGetHandler(getHandler).
		WithSearchHandler(searchHandler).
		WithCreateHandler(createHandler).
		WithUpdateHandler(updateHandler).
		WithCompleteHandler(completeHandler).
		WithRestoreHandler(restoreHandler).
		WithDeleteHandler(deleteHandler).
		WithListDeletedHandler(listDeletedHandler).
		Build()

	if err != nil {
		logger.Fatal().
			Err(err).
			Msg("handler.NewTaskHandlerBuilder.Build")
	}

	server := taskserver.NewMustGRPCServer(taskserver.Opts{
		Logger:                    otelLogger,
		Handler:                   grpcHandler,
		MeterProvider:             metric.NewMeterProvider().Meter(utility.GetEnv("METER_PROVIDER_NAME", "")),
		Addr:                      utility.GetEnv("SERVER_ADDRESS", ":8080"),
		MaxConnectionLimit:        utility.GetEnvUInt("SERVER_MAX_CONNECTION_LIMIT", 2000),
		MaxConnectionIdle:         utility.GetEnvDuration("SERVER_MAX_CONNECTION_IDLE", 5*time.Minute),
		PingHeartbeat:             utility.GetEnvDuration("SERVER_PING_HEARDBEAT", 2*time.Minute),
		PingResponseTime:          utility.GetEnvDuration("SERVER_PING_RESPONSE_TIME", 20*time.Second),
		MinConsecutivePingTimeout: utility.GetEnvDuration("SERVER_MIN_CONSECUTIVE_PING_TIME", 5*time.Second),
	})

	deleteUserTasksHandler, err := event.NewDeleteUserTasksHandler(repository)
	if err != nil {
		logger.Fatal().
			Err(err).
			Msg("event.NewDeleteUserTasksHandler")
	}

	deletedUserConsumer, err := kafka.NewUserDeletedConsumer(deleteUserTasksHandler, kafka.UserDeletedConsumerOptions{
		Logger:            otelLogger,
		Brokers:           strings.Split(utility.GetEnv("EVENT_BUS_BROKERS", "localhost:9092"), ","),
		GroupID:           utility.GetEnv("EVENT_BUS_GROUP_ID", ""),
		Topic:             utility.GetEnv("EVENT_BUS_TOPIC", ""),
		MinBytes:          int(utility.GetEnvInt("EVENT_BUS_MIN_BYTES", 1)),
		MaxBytes:          int(utility.GetEnvInt("EVENT_BUS_MAX_BYTES", 100)),
		MaxWait:           utility.GetEnvDuration("EVENT_BUS_MAX_WAIT", 1*time.Second),
		HeartbeatInterval: utility.GetEnvDuration("EVENT_BUS_HEARTBEAT_INTERVAL", 30*time.Second),
		SessionTimeout:    utility.GetEnvDuration("EVENT_BUS_SESSION_TIME", 30*time.Minute),
	})
	if err != nil {
		logger.Fatal().
			Err(err).
			Msg("event.NewUserDeletedConsumer")
	}

	wg := &sync.WaitGroup{}
	wg.Add(func() int {
		if utility.GetEnv("ENV", "dev") == "dev" {
			return 3
		}
		return 2
	}())

	go func() {
		if err := server.Serve(); err != nil {
			logger.Error().
				Err(err).
				Msg("grpc.server.Serve")
		}
		wg.Done()
	}()

	if utility.GetEnv("ENV", "dev") == "dev" {
		go func() {
			if err := http.ListenAndServe(utility.GetEnv("PPROF_SERVER_ADDRESS", ":8081"), nil); err != nil {
				logger.Error().
					Err(err).
					Msg("pprof http.Server")
			}
			wg.Done()
		}()
	}

	go func() {
		if err := deletedUserConsumer.Handle(ctx); err != nil {
			logger.Error().
				Err(err).
				Msg("deletedUserConsumer.Handle")
		}
		wg.Done()
	}()

	logger.Info().Msg("task service started")

	<-ctx.Done()
	server.Server.GracefulStop()
	_ = deletedUserConsumer.Stop()
	wg.Wait()

	logger.Info().Msg("task service stopped")
}
