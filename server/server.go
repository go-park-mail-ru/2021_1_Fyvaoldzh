package server

import (
	"context"
	chhttp "kudago/application/chat/delivery/http"
	chrepository "kudago/application/chat/repository"
	chusecase "kudago/application/chat/usecase"
	ehttp "kudago/application/event/delivery/http"
	erepository "kudago/application/event/repository"
	eusecase "kudago/application/event/usecase"
	shttp "kudago/application/subscription/delivery/http"
	srepository "kudago/application/subscription/repository"
	susecase "kudago/application/subscription/usecase"
	"kudago/application/user/delivery/http"
	"kudago/application/user/repository"
	"kudago/application/user/usecase"
	"kudago/pkg/constants"
	"kudago/pkg/custom_sanitizer"
	tarantool2 "kudago/pkg/infrastructure/session_manager"
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

func NewServer(l *zap.SugaredLogger) *echo.Echo {
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

	userRep := repository.NewUserDatabase(pool, logger)
	eventRep := erepository.NewEventDatabase(pool, logger)
	subRep := srepository.NewSubscriptionDatabase(pool, logger)
	chatRep := chrepository.NewChatDatabase(pool, logger)

	userUC := usecase.NewUser(userRep, subRep, logger)
	eventUC := eusecase.NewEvent(eventRep, subRep, logger)
	subUC := susecase.NewSubscription(subRep, logger)
	chatUC := chusecase.NewChat(chatRep, subRep, userRep, logger)

	sm := tarantool2.NewSessionManager(conn)

	sanitizer := custom_sanitizer.NewCustomSanitizer(bluemonday.UGCPolicy())

	http.CreateUserHandler(e, userUC, sm, sanitizer, logger)
	shttp.CreateSubscriptionsHandler(e, subUC, sm, logger)
	ehttp.CreateEventHandler(e, eventUC, sm, sanitizer, logger)
	chhttp.CreateChatHandler(e, chatUC, sm, sanitizer, logger)

	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: constants.CSRFHeader,
		CookiePath:  "/",
	}))

	return e
}

func ListenAndServe(e *echo.Echo) {
	e.Logger.Fatal(e.Start(":1323"))
}
