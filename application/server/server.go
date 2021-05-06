package server

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"

	chhttp "kudago/application/chat/delivery/http"
	ehttp "kudago/application/event/delivery/http"
	erepository "kudago/application/event/repository"
	eusecase "kudago/application/event/usecase"
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

	"github.com/microcosm-cc/bluemonday"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	middleware1 "kudago/application/server/middleware"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type Server struct {
	rpcAuth clientAuth.IAuthClient
	rpcSub  clientSub.ISubscriptionClient
	rpcChat clientChat.IChatClient
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
	logger := logger.NewLogger(l)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
	}))

	pool, err := pgxpool.Connect(context.Background(), constants.DBConnect)
	if err != nil {
		logger.Fatal(err)
	}
	err = pool.Ping(context.Background())
	if err != nil {
		logger.Fatal(err)
	}


	rpcAuth, err := clientAuth.NewAuthClient(constants.AuthServicePort, logger, tracer)
	if err != nil {
		logger.Fatal(err)
	}

	rpcSub, err := clientSub.NewSubscriptionClient(constants.SubscriptionServicePort, logger, tracer)
	if err != nil {
		logger.Fatal(err)
	}

	rpcChat, err := clientChat.NewChatClient(constants.ChatServicePort, logger, tracer)
	if err != nil {
		logger.Fatal(err)
	}

	userRep := repository.NewUserDatabase(pool, logger)
	eventRep := erepository.NewEventDatabase(pool, logger)
	subRep := srepository.NewSubscriptionDatabase(pool, logger)

	userUC := usecase.NewUser(userRep, subRep, logger)
	eventUC := eusecase.NewEvent(eventRep, subRep, logger)
	subscriptionUC := subusecase.NewSubscription(subRep, logger)

	sanitizer := custom_sanitizer.NewCustomSanitizer(bluemonday.UGCPolicy())

	auth := middleware1.NewAuth(rpcAuth)

	http.CreateUserHandler(e, userUC, rpcAuth, sanitizer, logger, auth)
	shttp.CreateSubscriptionsHandler(e, rpcAuth, rpcSub, subscriptionUC, sanitizer, logger, auth)
	ehttp.CreateEventHandler(e, eventUC, rpcAuth, sanitizer, logger, auth)
	chhttp.CreateChatHandler(e, rpcAuth, sanitizer, logger, auth, rpcChat)

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
	return &server
}

func (s Server) ListenAndServe() {
	s.e.Logger.Fatal(s.e.Start(":1323"))
	defer s.rpcAuth.Close()
	defer s.rpcSub.Close()
	defer s.rpcChat.Close()
}
