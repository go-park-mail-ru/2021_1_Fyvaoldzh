package server

import (
	"myapp/events"
	"myapp/models"
	"sync"

	"github.com/labstack/echo"
)

func NewServer() *echo.Echo {
	e := echo.New()
	handlers := events.Handlers{
		Events: make([]models.Event, 0),
		Mu:     &sync.Mutex{},
	}
	e.GET("/", handlers.All)
	e.GET("/event/:id", handlers.GetEvent)
	e.GET("/show", handlers.Show)
	e.POST("/create", handlers.Create)
	e.DELETE("/event/:id", handlers.Delete)

	return e
}

func ListenAndServe(e *echo.Echo) {
	e.Logger.Fatal(e.Start(":1323"))
}

//curl -v -X POST -H "Content-Type: application/json" -d '{"title": "dada", "description": "yaya", "typeEvent": "rave"}' http://localhost:1323/create
//curl -X DELETE localhost:1323/event/3
