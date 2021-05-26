package kudago_server

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
	kudago_repository "kudago/application/microservices/api_kudago/kudago/repository"
	kudago_usecase "kudago/application/microservices/api_kudago/kudago/usecase"
	"kudago/application/microservices/api_kudago/kudago_proto"
	"kudago/pkg/constants"
	"kudago/pkg/logger"
	"log"
	"net"
	"os"
)

type Server struct {
	port string
	kudago *KudagoServer
}

func NewServer(port string, logger *logger.Logger) *Server {
	pool, err := pgxpool.Connect(context.Background(),
		"user=" + os.Getenv("POSTGRE_USER") +
		" password=" + os.Getenv("DB_PASSWORD") + constants.DBConnect)
	if err != nil {
		logger.Fatal(err)
	}
	err = pool.Ping(context.Background())
	if err != nil {
		logger.Fatal(err)
	}

	kRepo := kudago_repository.NewKudagoDatabase(pool, *logger)
	kUsecase := kudago_usecase.NewKudagoUsecase(kRepo, *logger)

	return &Server{
		port: port,
		kudago: NewKudagoServer(kUsecase, logger),
	}
}

func (s *Server) ListenAndServe() error {
	jaegerCfgInstance := jaegercfg.Configuration{
		ServiceName: "kudago_server",
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
	kudago_proto.RegisterKudagoServer(gServer, s.kudago)
	log.Println("starting server at " + s.port)
	err = gServer.Serve(listener)

	if err != nil {
		return nil
	}

	return nil
}
