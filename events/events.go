package events

import (
	"fmt"
	"io"
	"log"
	"myapp/models"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/labstack/echo"
	"github.com/mailru/easyjson"
)

var BaseEvents = models.Events{
	{125, "Идущий к реке", "Я в своем познаии настолько преисполнился", "cognition", "myapp/125"},
	{126, "Димон заминированный тапок", "Мне абсолютно все равно", "cognition", "myapp/126"},
	{127, "exampleTitle", "exampleText", "1", "myapp/127"},
	{128, "Пример", "Пример без картинки", "noImg", ""},
}

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
			//Или здесь запоминать удаляемый элемент, а затем ретернить его json'ом?
			return c.String(http.StatusOK, "Event with id "+fmt.Sprint(id)+" succesfully deleted \n")
		}
	}

	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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

	//fmt.Println(newEvent)
	h.Mu.Lock()

	h.IdMax++
	id := h.IdMax
	newEvent.ID = id

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
		}
	}

	return nil
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

//Не знаю, где какие ошибки правильно выкидывать
func (h *Handlers) Save(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	img, err := c.FormFile("image")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	src, err := img.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer src.Close()

	dst, err := os.Create(fmt.Sprint(id))
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	for i := range h.Events {
		if h.Events[i].ID == uint64(id) {
			h.Events[i].Image = "myapp/" + fmt.Sprint(id)
		}
	}

	return nil
}
