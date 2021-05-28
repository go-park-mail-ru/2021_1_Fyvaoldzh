package server

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tarantool/go-tarantool"
	kudago_http "kudago/application/api_kudago/delivery/http"
	chhttp "kudago/application/chat/delivery/http"
	ehttp "kudago/application/event/delivery/http"
	erepository "kudago/application/event/repository"
	eusecase "kudago/application/event/usecase"
	kudago_client "kudago/application/microservices/api_kudago/client"
	clientAuth "kudago/application/microservices/auth/client"
	clientChat "kudago/application/microservices/chat/client"
	clientSub "kudago/application/microservices/subscription/client"
	shttp "kudago/application/subscription/delivery/http"
	srepository "kudago/application/subscription/repository"
	subusecase "kudago/application/subscription/usecase"
	"kudago/application/user/delivery/http"
	"kudago/application/user/repository"
	"kudago/application/user/usecase"
	"kudago/pkg/constants"
	"kudago/pkg/custom_sanitizer"
	"kudago/pkg/logger"
	"log"
	"os"

	middleware1 "kudago/application/server/middleware"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
	"github.com/microcosm-cc/bluemonday"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
)

type Server struct {
	rpcAuth clientAuth.IAuthClient
	rpcSub  clientSub.ISubscriptionClient
	rpcChat clientChat.IChatClient
	rpcKudago *kudago_client.KudagoClient
	e       *echo.Echo
}

func NewServer(l *zap.SugaredLogger) *Server {
	var server Server

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

	e := echo.New()
	lg := logger.NewLogger(l)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
	}))

	pool, err := pgxpool.Connect(context.Background(),
		"user=" + os.Getenv("POSTGRE_USER") +
			" password=" + os.Getenv("DB_PASSWORD") + constants.DBConnect)
	if err != nil {
		lg.Fatal(err)
	}
	err = pool.Ping(context.Background())
	if err != nil {
		lg.Fatal(err)
	}

	rpcAuth, err := clientAuth.NewAuthClient(constants.AuthServicePort, lg, tracer)
	if err != nil {
		lg.Fatal(err)
	}

	rpcSub, err := clientSub.NewSubscriptionClient(constants.SubscriptionServicePort, lg, tracer)
	if err != nil {
		lg.Fatal(err)
	}

	rpcChat, err := clientChat.NewChatClient(constants.ChatServicePort, lg, tracer)
	if err != nil {
		lg.Fatal(err)
	}

	rpcKudago, err := kudago_client.NewKudagoClient(constants.KudagoServicePort, lg, tracer)
	if err != nil {
		lg.Fatal(err)
	}
	conn, err := tarantool.Connect(constants.TarantoolAddress, tarantool.Opts{
		User: os.Getenv("TARANTOOL_USER"),
		Pass: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		lg.Fatal(err)
	}

	_, err = conn.Ping()
	if err != nil {
		lg.Fatal(err)
	}

	userRep := repository.NewUserDatabase(pool, conn, lg)
	eventRep := erepository.NewEventDatabase(pool, lg)
	subRep := srepository.NewSubscriptionDatabase(pool, lg)

	userUC := usecase.NewUser(userRep, subRep, lg)
	eventUC := eusecase.NewEvent(eventRep, subRep, lg)
	subscriptionUC := subusecase.NewSubscription(subRep, lg)

	sanitizer := custom_sanitizer.NewCustomSanitizer(bluemonday.UGCPolicy())

	auth := middleware1.NewAuth(rpcAuth)

	http.CreateUserHandler(e, userUC, rpcAuth, sanitizer, lg, auth)
	shttp.CreateSubscriptionsHandler(e, rpcAuth, rpcSub, subscriptionUC, sanitizer, lg, auth)
	ehttp.CreateEventHandler(e, eventUC, rpcAuth, sanitizer, lg, auth)
	chhttp.CreateChatHandler(e, rpcAuth, sanitizer, lg, auth, rpcChat)
	kudago_http.CreateKudagoHandler(e, rpcKudago, lg)

	//e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
	//	TokenLookup: constants.CSRFHeader,
	//	CookiePath:  "/",
	//}))

	prometheus.MustRegister(middleware1.FooCount, middleware1.Hits)
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	server.e = e
	server.rpcAuth = rpcAuth
	server.rpcSub = rpcSub
	server.rpcChat = rpcChat
	server.rpcKudago = rpcKudago
	return &server
}

func (s Server) ListenAndServe() {
	s.e.Logger.Fatal(s.e.Start(":1323"))
	defer s.rpcAuth.Close()
	defer s.rpcSub.Close()
	defer s.rpcChat.Close()
	defer s.rpcKudago.Close()
}
