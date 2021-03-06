package server

import (
	"log"
	"myapp/events"
	mytype "myapp/events"
	"net/http"
	"sync"

	"github.com/labstack/echo"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func UserJson(c echo.Context) error {
	u := new(User)
	err := c.Bind(u)
	if err != nil {
		log.Println(err)
	}
	return c.JSON(http.StatusCreated, u.Name+u.Email)
}

func NewServer() *echo.Echo {
	e := echo.New()
	handlers := mytype.Handlers{
		Events: make([]mytype.Event, 0),
		Mu:     &sync.Mutex{},
	}
	e.GET("/", handlers.All(c echo.Context))
	e.GET("/events/:id", events.GetEvent)
	e.GET("/show", events.Show)
	e.POST("/users", UserJson)
	e.POST("/create", events.CreateEvent)

	return e
}

func ListenAndServe(e *echo.Echo) {
	e.Logger.Fatal(e.Start(":1323"))
}
