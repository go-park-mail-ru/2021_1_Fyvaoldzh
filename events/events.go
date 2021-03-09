package events

import (
	"errors"
	"fmt"
	"io"
	"kudago/models"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/labstack/echo"
	"github.com/mailru/easyjson"
)

type Handlers struct {
	Events models.Events
	Mu     *sync.Mutex
	IdMax  uint64
}

func (h *Handlers) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	for i := range h.Events {
		if h.Events[i].ID == uint64(id) {
			h.Events = append(h.Events[:i], h.Events[i+1:]...)
			return c.String(http.StatusOK, "Event with id "+fmt.Sprint(id)+" succesfully deleted \n")
		}
	}

	return echo.NewHTTPError(http.StatusNotFound, errors.New("Event with id "+fmt.Sprint(id)+" not found"))
}

func (h *Handlers) GetAllEvents(c echo.Context) error {
	if _, err := easyjson.MarshalToWriter(h.Events, c.Response().Writer); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (h *Handlers) Create(c echo.Context) error {
	defer c.Request().Body.Close()

	newEvent := &models.Event{}

	if err := easyjson.UnmarshalFromReader(c.Request().Body, newEvent); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	h.Mu.Lock()

	h.IdMax++
	newEvent.ID = h.IdMax

	h.Events = append(h.Events, *newEvent)
	h.Mu.Unlock()
	return c.JSON(http.StatusOK, *newEvent)
}

func (h *Handlers) GetOneEvent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	for _, event := range h.Events {
		if event.ID == uint64(id) {
			if _, err = easyjson.MarshalToWriter(event, c.Response().Writer); err != nil {
				log.Println(err)
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			return nil
		}
	}

	return echo.NewHTTPError(http.StatusNotFound, errors.New("Event with id "+fmt.Sprint(id)+" not found"))
}

func (h *Handlers) GetEvents(c echo.Context) error {
	typeEvent := c.QueryParam("typeEvent")
	var showEvents models.Events
	for _, event := range h.Events {
		if event.TypeEvent == typeEvent {
			showEvents = append(showEvents, event)
		}
	}
	if _, err := easyjson.MarshalToWriter(showEvents, c.Response().Writer); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (h *Handlers) Save(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	img, err := c.FormFile("image")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	src, err := img.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer src.Close()

	fileName := fmt.Sprint(id) + img.Filename

	dst, err := os.Create(fileName)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	for i := range h.Events {
		if h.Events[i].ID == uint64(id) {
			h.Events[i].Image = "kudago/" + fileName
			return c.JSON(http.StatusOK, h.Events[i])
		}
	}

	return echo.NewHTTPError(http.StatusNotFound, errors.New("Event with id "+fmt.Sprint(id)+" not found"))
}
