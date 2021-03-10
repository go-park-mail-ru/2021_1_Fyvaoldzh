package server

import (
	"kudago/events"
	"kudago/models"
	"kudago/profile"
	"kudago/user"
	"sync"

	"kudago/auth"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func NewServer() *echo.Echo {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
	}))

	handlers := events.Handlers{
		Events: models.BaseEvents,
		Mu:     &sync.Mutex{},
	}
	userHandler := user.HandlerUser{UserBase: user.UserBase, ProfileBase: user.ProfileBase, PlanningEvent: user.PlanningEvent, Store: make(map[string]int), Mu: &sync.Mutex{}}


	e.POST("/api/v1/login", userHandler.Login)
	e.GET("/api/v1/logout", userHandler.Logout)
	e.POST("/api/v1/register", userHandler.CreateUser)

	e.GET("/api/v1/profile", userHandler.GetProfile)
	e.POST("/api/v1/profile/:id", userHandler.GetUserProfile)
	e.GET("/api/v1/avatar/:id", userHandler.GetAvatar)
	e.PUT("/api/v1/profile", userHandler.UpdateProfile)
	e.PUT("/api/v1/upload_avatar", userHandler.UploadAvatar)

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
