package events

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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

func (h *Handlers) DeleteByID(id int) bool {
	for i := range h.Events {
		if h.Events[i].ID == uint64(id) {
			h.Events = append(h.Events[:i], h.Events[i+1:]...)
			return true
		}
	}
	return false
}

func (h *Handlers) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if h.DeleteByID(id) {
		return c.String(http.StatusOK, "Event with id "+fmt.Sprint(id)+" succesfully deleted \n")
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

func (h *Handlers) CreateNewEvent(newEvent *models.Event) {
	h.Mu.Lock()

	h.IdMax++
	newEvent.ID = h.IdMax

	h.Events = append(h.Events, *newEvent)
	h.Mu.Unlock()
}

func (h *Handlers) Create(c echo.Context) error {
	defer c.Request().Body.Close()

	newEvent := &models.Event{}

	if err := easyjson.UnmarshalFromReader(c.Request().Body, newEvent); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	h.CreateNewEvent(newEvent)

	return c.JSON(http.StatusOK, *newEvent)
}

//Как здесь сделать нормальный возврат ошибки? Приходится делать костыли на костыли, чтобы затестить функции. Так со всеми)
func (h *Handlers) GetOneEventByID(id int) models.Event {
	for _, event := range h.Events {
		if event.ID == uint64(id) {
			return event
		}
	}
	return models.Event{}
}

func (h *Handlers) GetOneEvent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	//Здесь из-за такого подхода вообще хаос какой-то. Дважды вызываем функцию, в общем бред
	if h.GetOneEventByID(id).ID != 0 {
		if _, err = easyjson.MarshalToWriter(h.GetOneEventByID(id), c.Response().Writer); err != nil {
			log.Println(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return nil
	}

	return echo.NewHTTPError(http.StatusNotFound, errors.New("Event with id "+fmt.Sprint(id)+" not found"))
}

func (h *Handlers) GetEventsByType(typeEvent string) models.Events {
	var showEvents models.Events
	for _, event := range h.Events {
		if event.TypeEvent == typeEvent {
			showEvents = append(showEvents, event)
		}
	}
	return showEvents
}

func (h *Handlers) GetEvents(c echo.Context) error {
	typeEvent := c.QueryParam("typeEvent")
	if _, err := easyjson.MarshalToWriter(h.GetEventsByType(typeEvent), c.Response().Writer); err != nil {
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
			h.Events[i].Image = fileName
			return c.JSON(http.StatusOK, h.Events[i])
		}
	}

	return echo.NewHTTPError(http.StatusNotFound, errors.New("Event with id "+fmt.Sprint(id)+" not found"))
}

func (h *Handlers) GetImage(c echo.Context) error {
	defer c.Request().Body.Close()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	event := h.GetOneEventByID(id)

	file, err := ioutil.ReadFile(event.Image)
	if err != nil {
		log.Println("Cannot open file: " + event.Image)
	} else {
		c.Response().Write(file)
	}

	return nil
}
