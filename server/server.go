package server

import (
	"myapp/events"
	mytype "myapp/events"
	"net/http"
	"sync"

	"github.com/labstack/echo"
)

//Закомменчено, так как это часть пока что висит на Насте
/*type User struct {
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
}*/

func NewServer() *echo.Echo {
	e := echo.New()
	handlers := mytype.Handlers{
		Events: make([]mytype.Event, 0),
		Mu:     &sync.Mutex{},
	}
	e.GET("/", func(c echo.Context) error {
		handlers.All(c)
		return c.JSON(http.StatusOK, "ok")
	})
	e.GET("/events/:id", events.GetEvent)
	e.GET("/show", events.Show)
	//e.POST("/users", UserJson)
	e.POST("/create", func(c echo.Context) error {
		handlers.Create(c)
		return c.JSON(http.StatusOK, "ok")
	})

	return e
}

func ListenAndServe(e *echo.Echo) {
	e.Logger.Fatal(e.Start(":1323"))
}

//curl -v -X POST -H "Content-Type: application/json" -d '{"Title": "dada", "Description": "yaya"}' http://localhost:1323/create
