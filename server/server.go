package server

import (
	"database/sql"
	_ "github.com/jackc/pgx/v4"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
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
	//e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
	//	TokenLookup: "header:X-XSRF-TOKEN",
	//}))





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

	userRep := repository.NewUserDatabase(getPostgres())
	userUC := usecase.NewUser(userRep)

	http.CreateUserHandler(e, userUC)
	e.Logger.Fatal(e.Start(":1323"))
}
