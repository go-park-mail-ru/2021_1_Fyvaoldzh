package server

import (
	"myapp/events"
	"sync"

	"github.com/labstack/echo"
)

func NewServer() *echo.Echo {
	e := echo.New()
	handlers := events.Handlers{
		Events: events.BaseEvents,
		Mu:     &sync.Mutex{},
	}
	e.GET("/api/v1/", handlers.GetAllEvents)
	e.GET("/api/v1/event/:id", handlers.GetOneEvent)
	e.GET("/api/v1/event", handlers.GetEvents)
	e.POST("/api/v1/create", handlers.Create)
	e.DELETE("/api/v1/event/:id", handlers.Delete)
	e.POST("/api/v1/save/:id", handlers.Save)

	return e
}

func ListenAndServe(e *echo.Echo) {
	e.Logger.Fatal(e.Start(":1323"))
}

//curl -v -X POST -H "Content-Type: application/json" -d '{"title": "dada", "description": "yaya", "typeEvent": "rave"}' http://localhost:1323/api/v1/create
//curl -X DELETE localhost:1323/api/v1/event/3
