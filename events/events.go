package events

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"kudago/models"
	"log"
	"mime/multipart"
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

func (h *Handlers) GetOneEventByID(id int) (models.Event, error) {
	for _, event := range h.Events {
		if event.ID == uint64(id) {
			return event, nil
		}
	}
	return models.Event{}, errors.New("event not found")
}

func (h *Handlers) GetOneEvent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	event, err := h.GetOneEventByID(id)

	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, errors.New("Event with id "+fmt.Sprint(id)+" not found"))
	}
	if _, err = easyjson.MarshalToWriter(event, c.Response().Writer); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return nil

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

func (h *Handlers) SaveImage(src multipart.File, id int) (models.Event, error) {
	defer src.Close()

	fileName := fmt.Sprint(id) + ".jpg"

	dst, err := os.Create(fileName)
	if err != nil {
		return models.Event{}, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return models.Event{}, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	for i := range h.Events {
		if h.Events[i].ID == uint64(id) {
			h.Events[i].Image = fileName
			return h.Events[i], nil
		}
	}

	return models.Event{}, echo.NewHTTPError(http.StatusNotFound, errors.New("Event with id "+fmt.Sprint(id)+" not found"))
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

	event, errSave := h.SaveImage(src, id)
	if errSave != nil {
		return errSave
	}
	return c.JSON(http.StatusOK, event)
}

func (h *Handlers) GetImage(c echo.Context) error {
	defer c.Request().Body.Close()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	event, errFind := h.GetOneEventByID(id)

	if errFind != nil {
		log.Println("Cannot find file with id: " + fmt.Sprint(id))
	}

	file, err := ioutil.ReadFile(event.Image)
	if err != nil {
		log.Println("Cannot open file: " + event.Image)
	} else {
		c.Response().Write(file)
	}

	return nil
}
