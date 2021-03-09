package server

import (
	"github.com/go-park-mail-ru/2021_1_Fyvaoldzh/events"
	"github.com/go-park-mail-ru/2021_1_Fyvaoldzh/models"
	"github.com/go-park-mail-ru/2021_1_Fyvaoldzh/profile"
	"log"
	"net/http"
	"sync"

	"github.com/go-park-mail-ru/2021_1_Fyvaoldzh/auth"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)


func NewServer() *echo.Echo {
	e := echo.New()
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
	}))

	regHandler := auth.RegisterHandler{Mu: &sync.Mutex{}}
	loginHandler := auth.LoginHandler{Mu: &sync.Mutex{}}
	profileHandler := profile.UserHandler{Mu: &sync.Mutex{}}

	handlers := events.Handlers{
		Events: events.BaseEvents,
		Mu:     &sync.Mutex{},
	}


	e.POST("/api/v1/register", func(c echo.Context) error {
		err := regHandler.CreateUser(c)
		if err != nil {
			return c.JSON(err.Code, err.Error())
		}
		log.Println(models.UserBase)
		return c.JSON(http.StatusOK, "ok")
	})

	e.POST("/api/v1/login", func(c echo.Context) error {
		err := loginHandler.Login(c)
		if err != nil {
			return c.JSON(err.Code, err.Error())
		}

		return c.JSON(http.StatusOK, "login successful")
	})

	e.GET("/api/v1/profile/:id", func(c echo.Context) error {
		err := profileHandler.GetUserProfile(c)
		if err != nil {
			return c.JSON(err.Code, err.Error())
		}

		return c.JSON(http.StatusOK, "ok")
	})

	e.GET("/api/v1/profile", func(c echo.Context) error {
		err := profileHandler.GetProfile(c)
		if err != nil {
			return c.JSON(err.Code, err.Error())
		}

		return c.JSON(http.StatusOK, "ok")
	})

	e.GET("/api/v1/logout", func(c echo.Context) error {
		err := loginHandler.Logout(c)
		if err != nil {
			return c.JSON(err.Code, err.Error())
		}

		return c.JSON(http.StatusOK, "ok")
	})

	e.PUT("/api/v1/profile", func(c echo.Context) error {
		err := profileHandler.UpdateProfile(c)
		if err != nil {
			return c.JSON(err.Code, err.Error())
		}

		return c.JSON(http.StatusOK, "ok")
	})

	e.PUT("/api/v1/upload_avatar", func(c echo.Context) error {
		err := profileHandler.UploadAvatar(c)
		if err != nil {
			return c.JSON(err.Code, err.Error())
		}

		return c.JSON(http.StatusOK, "ok")
	})

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

