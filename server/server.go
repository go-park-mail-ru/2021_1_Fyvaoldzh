package server

import (
	"github.com/go-park-mail-ru/2021_1_Fyvaoldzh/profile"
	"log"
	"net/http"
	"sync"

	"github.com/go-park-mail-ru/2021_1_Fyvaoldzh/auth"

	"github.com/labstack/echo"
)


func NewServer() *echo.Echo {
	e := echo.New()

	regHandler := auth.RegisterHandler{Mu: &sync.Mutex{}}
	loginHandler := auth.LoginHandler{Mu: &sync.Mutex{}}
	profileHandler := profile.UserHandler{Mu: &sync.Mutex{}}


	e.POST("/register", func(c echo.Context) error {
		err := regHandler.CreateUser(c)
		if err != nil {
			// do some code
			// return err
		}
		log.Println(auth.UserBase)
		return c.JSON(http.StatusOK, "ok")
	})

	e.GET("/login", func(c echo.Context) error {
		err := loginHandler.Login(c)
		if err != nil {
			// do some code
			// return err
		}

		return c.JSON(http.StatusOK, "ok")
	})

	e.GET("/profile", func(c echo.Context) error {
		json, err := profileHandler.GetProfile(c)
		if err != nil {
			// do some code
			// return err
		}

		return c.JSON(http.StatusOK, json)
	})

	return e
}

func ListenAndServe(e *echo.Echo) {
	e.Logger.Fatal(e.Start(":1323"))
}
