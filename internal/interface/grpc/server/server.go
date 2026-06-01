package server

import (
	"net"
	taskgrpc "task-service/gen/go/task-service/task"
	"task-service/internal/interface/grpc/handler"
	"task-service/internal/interface/grpc/interceptor"
	"time"

	"buf.build/go/protovalidate"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
)

type Opts struct {
	Logger                    *zerolog.Logger
	Handler                   *handler.TaskHandler
	MeterProvider             metric.Meter
	Addr                      string
	MaxConnectionLimit        uint32
	MaxConnectionIdle         time.Duration
	PingHeartbeat             time.Duration
	PingResponseTime          time.Duration
	MinConsecutivePingTimeout time.Duration
}

type GRPCServer struct {
	Server *grpc.Server
	Addr   string
}

func NewMustGRPCServer(opts Opts) *GRPCServer {
	validator, err := protovalidate.New(protovalidate.WithFailFast())
	if err != nil {
		panic(err)
	}

	keepaliveOptions := keepalive.ServerParameters{
		MaxConnectionIdle: opts.MaxConnectionIdle,
		Time:              opts.PingHeartbeat,
		Timeout:           opts.PingResponseTime,
	}

	pingPolicy := keepalive.EnforcementPolicy{
		MinTime:             opts.MinConsecutivePingTimeout,
		PermitWithoutStream: true,
	}

	metricsInterceptor, err := interceptor.NewMetricsInterceptor(opts.MeterProvider)
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.NewRecoverInterceptor(opts.Logger),
			metricsInterceptor.UnaryServerInterceptor(),
			interceptor.NewRateLimitInterceptor(),
			interceptor.NewAuthInterceptor(),
			interceptor.NewValidateInterceptor(validator),
			interceptor.NewErrTranslateInterceptor(),
			interceptor.NewLoggingInterceptor(opts.Logger),
		),
		grpc.MaxConcurrentStreams(opts.MaxConnectionLimit),
		grpc.KeepaliveParams(keepaliveOptions),
		grpc.KeepaliveEnforcementPolicy(pingPolicy),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	reflection.Register(server)
	taskgrpc.RegisterTaskServiceServer(
		server,
		opts.Handler,
	)
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(server, healthServer)
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus("tasks.v1.DatabasePool", healthpb.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus("tasks.v1.EventBus", healthpb.HealthCheckResponse_SERVING)

	return &GRPCServer{
		Server: server,
		Addr:   opts.Addr,
	}
}

func (s *GRPCServer) Serve() error {
	lis, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	return s.Server.Serve(lis)
}
