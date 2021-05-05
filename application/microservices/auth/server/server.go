package server

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	traceutils "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"github.com/tarantool/go-tarantool"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"google.golang.org/grpc"
	"kudago/application/microservices/auth/proto"
	srepo "kudago/application/microservices/auth/session/repository"
	sessions "kudago/application/microservices/auth/session/usecase"
	urepo "kudago/application/microservices/auth/user/repository"
	users "kudago/application/microservices/auth/user/usecase"
	"kudago/pkg/constants"
	"kudago/pkg/logger"
	"log"
	"net"
)

type Server struct {
	port string
	auth *AuthServer
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

	conn, err := tarantool.Connect(constants.TarantoolAddress, tarantool.Opts{
		User: constants.TarantoolUser,
		Pass: constants.TarantoolPassword,
	})
	if err != nil {
		logger.Fatal(err)
	}

	_, err = conn.Ping()
	if err != nil {
		logger.Fatal(err)
	}

	sRepo := srepo.NewSessionRepository(conn, *logger)
	uRepo := urepo.NewUserDatabase(pool, *logger)
	uUseCase := users.NewUserUseCase(uRepo, *logger)
	s := sessions.NewSessionUseCase(sRepo, uUseCase, *logger)

	return &Server{
		port: port,
		auth: NewAuthServer(s),
	}
}

func (s *Server) ListenAndServe() error {
	jaegerCfgInstance := jaegercfg.Configuration{
		ServiceName: "main_server",
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
	proto.RegisterAuthServer(gServer, s.auth)
	log.Println("starting server at " + s.port)
	err = gServer.Serve(listener)

	if err != nil {
		return nil
	}

	return nil
}
