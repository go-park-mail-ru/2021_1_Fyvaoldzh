package server

import (
	"context"
	ehttp "kudago/application/event/delivery/http"
	erepository "kudago/application/event/repository"
	eusecase "kudago/application/event/usecase"
	clientAuth "kudago/application/microservices/auth/client"
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

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
	"github.com/tarantool/go-tarantool"
	"go.uber.org/zap"
)

type Server struct {
	rpcAuth *clientAuth.AuthClient
	rpcSub *clientSub.SubscriptionClient
	e *echo.Echo
}

func NewServer(l *zap.SugaredLogger) *Server {
	var server Server
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

	rpcAuth, err := clientAuth.NewAuthClient(constants.AuthServicePort, logger)
	if err != nil {
		logger.Fatal(err)
	}

	rpcSub, err := clientSub.NewSubscriptionClient(constants.SubscriptionServicePort, logger)
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

	http.CreateUserHandler(e, userUC, *rpcAuth, sanitizer, logger)
	shttp.CreateSubscriptionsHandler(e, *rpcAuth, *rpcSub, subscriptionUC, sanitizer, logger)
	ehttp.CreateEventHandler(e, eventUC, *rpcAuth, sanitizer, logger)

	//e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
	//	TokenLookup: constants.CSRFHeader,
	//	CookiePath:  "/",
	//}))

	server.e = e
	server.rpcAuth = rpcAuth
	return &server
}

func (s Server) ListenAndServe() {
	s.e.Logger.Fatal(s.e.Start(":1323"))
	defer s.rpcAuth.Close()
	defer s.rpcSub.Close()
}
