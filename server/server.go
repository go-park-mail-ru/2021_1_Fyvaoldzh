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
	"kudago/application/user/delivery/http"
	"kudago/application/user/repository"
	"kudago/application/user/usecase"
	"kudago/pkg/constants"
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

	userUC := usecase.NewUser(userRep)

	conn, err := tarantool.Connect("127.0.0.1:3301", tarantool.Opts{
		User: "admin",
		Pass: "fyvaoldzh",
	})

	if err != nil {
		log.Fatalf("Connection refused")
	}

	defer conn.Close()

	http.CreateUserHandler(e, userUC)




	//subRep :=

	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "header:X-XSRF-TOKEN",
		CookieHTTPOnly: true,
	}))





	/*handlers := events.Handlers{
		Events: models.BaseEvents,
		Mu:     &sync.Mutex{},
	}*/



	/*
	e.GET("/api/v1/", handlers.GetAllEvents)
	e.GET("/api/v1/event/:id", handlers.GetOneEvent)
	e.GET("/api/v1/event", handlers.GetEvents)
	e.POST("/api/v1/create", handlers.Create)
	e.DELETE("/api/v1/event/:id", handlers.Delete)
	e.POST("/api/v1/save/:id", handlers.Save)
	e.GET("api/v1/event/:id/image", handlers.GetImage)
	 */
	return e
}

func ListenAndServe(e *echo.Echo) {
	e.Logger.Fatal(e.Start(":1323"))
}
