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
			return c.JSON(err.Code, err.Error())
		}
		log.Println(auth.UserBase)
		return c.JSON(http.StatusOK, "ok")
	})

	e.POST("/login", func(c echo.Context) error {
		err := loginHandler.Login(c)
		if err != nil {
			return c.JSON(err.Code, err.Error())
		}

		return c.JSON(http.StatusOK, "login successful")
	})

	e.GET("/profile", func(c echo.Context) error {
		json, err := profileHandler.GetProfile(c)
		if err != nil {
			return c.JSON(err.Code, err.Error())
		}

		return c.JSON(http.StatusOK, json)
	})

	e.GET("/logout", func(c echo.Context) error {
		err := loginHandler.Logout(c)
		if err != nil {
			return c.JSON(err.Code, err.Error())
		}

		return c.JSON(http.StatusOK, "ok")
	})

	return e
}

func ListenAndServe(e *echo.Echo) {
	e.Logger.Fatal(e.Start(":1323"))
}
