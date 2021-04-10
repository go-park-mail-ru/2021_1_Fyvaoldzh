package server

import (
	"context"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
	"github.com/tarantool/go-tarantool"
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
	"kudago/pkg/infrastructure"
	"log"
)

func NewServer() *echo.Echo {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
	}))

	pool, err := pgxpool.Connect(context.Background(), constants.DBConnect)
	if err != nil {
		log.Fatalln(err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	conn, err := tarantool.Connect(constants.TarantoolAddress, tarantool.Opts{
		User: constants.TarantoolUser,
		Pass: constants.TarantoolPassword,
	})

	if err != nil {
		log.Fatalln(err)
	}

	_, err = conn.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	userRep := repository.NewUserDatabase(pool)
	eventRep := erepository.NewEventDatabase(pool)
	subRep := srepository.NewSubscriptionDatabase(pool)

	userUC := usecase.NewUser(userRep, subRep)
	eventUC := eusecase.NewEvent(eventRep)
	subUC := susecase.NewSubscription(subRep)

	sm := infrastructure.SessionManager{}
	sm.Conn = conn

	http.CreateUserHandler(e, userUC, &sm)
	shttp.CreateSubscriptionsHandler(e, subUC, &sm)
	ehttp.CreateEventHandler(e, eventUC, &sm)

	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup:    constants.CSRFHeader,
		CookieHTTPOnly: true,
	}))

	return e
}

func ListenAndServe(e *echo.Echo) {
	e.Logger.Fatal(e.Start(":1323"))
}
