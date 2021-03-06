package server

import (
	"net/http"
	"sync"

	model "github.com/go-park-mail-ru/2021_1_Fyvaoldzh/models"

	"github.com/labstack/echo"
)

func NewServer() *echo.Echo {
	e := echo.New()
	regHandler :=


	e.GET("/", func(c echo.Context) error {
		handlers.All(c)
		return c.JSON(http.StatusOK, "ok")
	})

	//e.POST("/users", UserJson)


	return e
}

func ListenAndServe(e *echo.Echo) {
	e.Logger.Fatal(e.Start(":1323"))
}
