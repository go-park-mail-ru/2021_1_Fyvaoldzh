package server

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	traceutils "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"google.golang.org/grpc"
	srepo "kudago/application/microservices/subscription/subscription/repository"
	subscriptions "kudago/application/microservices/subscription/subscription/usecase"

	"kudago/application/microservices/subscription/proto"
	"kudago/pkg/constants"
	"kudago/pkg/logger"
	"log"
	"net"
)

type Server struct {
	port string
	ss   *SubscriptionServer
}

func NewServer(port string, logger *logger.Logger) *Server {
	pool, err := pgxpool.Connect(context.Background(), constants.DBConnect)
	if err != nil {
		logger.Fatal(err)
	}
	err = pool.Ping(context.Background())
	if err != nil {
		logger.Fatal(err)
	}

	sRepo := srepo.NewSubscriptionDatabase(pool, *logger)
	sUseCase := subscriptions.NewSubscription(sRepo, *logger)

	return &Server{
		port: port,
		ss:   NewSubscriptionServer(sUseCase),
	}
}

func (s *Server) ListenAndServe() error {
	jaegerCfgInstance := jaegercfg.Configuration{
		ServiceName: "sub_server",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "localhost:6831",
		},
	}

	tracer, _, err := jaegerCfgInstance.NewTracer(
		jaegercfg.Logger(jaegerlog.StdLogger),
		jaegercfg.Metrics(metrics.NullFactory),
	)
	if err != nil {
		log.Fatal("cannot create tracer", err)
	}
	opentracing.SetGlobalTracer(tracer)
	gServer := grpc.NewServer(grpc.
		UnaryInterceptor(traceutils.OpenTracingServerInterceptor(tracer)))

	listener, err := net.Listen("tcp", s.port)
	defer listener.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	proto.RegisterSubscriptionServer(gServer, s.ss)
	log.Println("starting server at " + s.port)
	err = gServer.Serve(listener)

	if err != nil {
		return nil
	}

	return nil
}
