package server

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
	"github.com/tarantool/go-tarantool"
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

func getPostgres() *sql.DB {
	dsn := constants.DBConnect
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalln("cant parse config", err)
	}
	err = db.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		log.Fatalln(err)
	}
	db.SetMaxOpenConns(10)
	return db
}

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

	err = pool.Ping(context.Background()) // вот тут будет первое подключение к базе
	if err != nil {
		log.Fatalln(err)
	}

	userRep := repository.NewUserDatabase(pool)
	subRep := srepository.NewSubscriptionDatabase(pool)

	userUC := usecase.NewUser(userRep)
	subUC := susecase.NewSubscription(subRep)

	conn, err := tarantool.Connect("127.0.0.1:3301", tarantool.Opts{
		User: "admin",
		Pass: "fyvaoldzh",
	})

	if err != nil {
		log.Fatalln(err)
	}

	_, err = conn.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	sm := infrastructure.SessionManager{}
	sm.Conn = conn

	http.CreateUserHandler(e, userUC, &sm)
	shttp.CreateSubscriptionsHandler(e, subUC, &sm)

	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "header:X-XSRF-TOKEN",
		CookieHTTPOnly: true,
	}))

	return e
}

func ListenAndServe(e *echo.Echo) {
	e.Logger.Fatal(e.Start(":1323"))
}
