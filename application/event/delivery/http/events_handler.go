package http

import (
	"errors"
	"fmt"
	"github.com/labstack/echo"
	"github.com/mailru/easyjson"
	"kudago/application/event"
	"kudago/pkg/infrastructure"
	"log"
	"net/http"
	"strconv"
)

type EventHandler struct {
	UseCase event.UseCase
	Sm      *infrastructure.SessionManager
}

func CreateEventHandler(e *echo.Echo, uc event.UseCase, sm *infrastructure.SessionManager) {

	eventHandler := EventHandler{UseCase: uc, Sm: sm}

	e.GET("/api/v1/", eventHandler.GetAllEvents)
	e.GET("/api/v1/event/:id", eventHandler.GetOneEvent)
	e.GET("/api/v1/event", eventHandler.GetEvents)
	e.POST("/api/v1/create", eventHandler.Create)
	e.DELETE("/api/v1/event/:id", eventHandler.Delete)
	e.POST("/api/v1/save/:id", eventHandler.Save)
	e.GET("api/v1/event/:id/image", eventHandler.GetImage)
}

func (eh EventHandler) GetAllEvents(c echo.Context) error {
	defer c.Request().Body.Close()

	events, err := eh.UseCase.GetAllEvents()
	if err != nil {
		log.Println(err)
		return err
	}

	if _, err = easyjson.MarshalToWriter(events, c.Response().Writer); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (eh EventHandler) GetOneEvent(c echo.Context) error {
	defer c.Request().Body.Close()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	ev, err := eh.UseCase.GetOneEvent(uint64(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, errors.New("Event with id "+fmt.Sprint(id)+" not found"))
	}

	if _, err = easyjson.MarshalToWriter(ev, c.Response().Writer); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (eh EventHandler) GetEvents(c echo.Context) error {

	return nil

}

func (eh EventHandler) Create(c echo.Context) error {

	return nil

}

func (eh EventHandler) Delete(c echo.Context) error {

	return nil

}

func (eh EventHandler) Save(c echo.Context) error {

	return nil

}

func (eh EventHandler) GetImage(c echo.Context) error {

	return nil

}

