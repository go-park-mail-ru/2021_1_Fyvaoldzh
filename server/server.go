package server

import (
	"net/http"
	"sync"

	"github.com/go-park-mail-ru/2021_1_Fyvaoldzh/auth"

	"github.com/labstack/echo"
)


func NewServer() *echo.Echo {
	e := echo.New()
	regHandler := auth.RegisterHandler{Mu: &sync.Mutex{}}


	e.POST("/register", func(c echo.Context) error {
		regHandler.CreateUser()
		return c.JSON(http.StatusOK, "ok")
	})

	return e
}

func ListenAndServe(e *echo.Echo) {
	e.Logger.Fatal(e.Start(":1323"))
}
