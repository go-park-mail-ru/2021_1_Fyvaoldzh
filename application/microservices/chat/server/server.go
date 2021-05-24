package server

import (
	"context"
	"github.com/tarantool/go-tarantool"
	erepo "kudago/application/event/repository"
	"kudago/application/microservices/chat/chat/repository"
	"kudago/application/microservices/chat/chat/usecase"
	"kudago/application/microservices/chat/proto"
	srepo "kudago/application/subscription/repository"
	urepo "kudago/application/user/repository"
	"kudago/pkg/constants"
	"net"

	"github.com/jackc/pgx/v4/pgxpool"
	traceutils "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"google.golang.org/grpc"

	"kudago/pkg/logger"
	"log"
)

type Server struct {
	port string
	ss   *ChatServer
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

	cRepo := repository.NewChatDatabase(pool, conn, *logger)
	sRepo := srepo.NewSubscriptionDatabase(pool, *logger)
	eRepo := erepo.NewEventDatabase(pool, *logger)
	uRepo := urepo.NewUserDatabase(pool, *logger)
	cUseCase := usecase.NewChat(cRepo, sRepo, uRepo, eRepo, *logger)

	return &Server{
		port: port,
		ss:   NewChatServer(cUseCase),
	}
}

func (s *Server) ListenAndServe() error {
	jaegerCfgInstance := jaegercfg.Configuration{
		ServiceName: "chat_server",
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
	if err != nil {
		log.Println(err)
		return err
	}
	defer listener.Close()
	proto.RegisterChatServer(gServer, s.ss)
	log.Println("starting server at " + s.port)
	err = gServer.Serve(listener)

	if err != nil {
		return nil
	}

	return nil
}
