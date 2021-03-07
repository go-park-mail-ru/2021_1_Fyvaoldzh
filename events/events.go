package events

import (
	"fmt"
	"log"
	"myapp/models"
	"net/http"
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

//Как тут будет выгоднее? Такой способ или же честно найти, копировать последний элемент на место удаляемого, а затем усечь слайс?
func (h *Handlers) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	newEvents := h.Events[:0]

	for _, event := range h.Events {
		if event.ID != uint64(id) {
			newEvents = append(newEvents, event)
		}
	}

	h.Events = newEvents
	//Нужно ли здесь почистить старый слайс? Я че-то уже под вечер не соображаю

	return nil
}

func (h *Handlers) All(c echo.Context) error {
	if _, err := easyjson.MarshalToWriter(h.Events, c.Response().Writer); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (h *Handlers) Create(c echo.Context) error {
	defer c.Request().Body.Close()

	newEventInput := &models.EventInput{}

	if err := easyjson.UnmarshalFromReader(c.Request().Body, newEventInput); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	fmt.Println(newEventInput)
	h.Mu.Lock()

	//id := h.IdMax++   Почему так не работает?
	h.IdMax++
	id := h.IdMax

	h.Events = append(h.Events, models.Event{
		ID:          id,
		Title:       newEventInput.Title,
		Description: newEventInput.Description,
		TypeEvent:   newEventInput.TypeEvent,
	})
	h.Mu.Unlock()
	return nil
}

func (h *Handlers) GetEvent(c echo.Context) error {
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

//Не думаю, что это правильный подход, но пока что не могу придумать ничего лучше, по крайней мере это работает
func (h *Handlers) Show(c echo.Context) error {
	typeEvent := c.QueryParam("typeEvent")
	//showEvents := new(models.Events)   Почему на такую запись го ругается, что это не слайс?
	var showEvents []models.Event
	for _, event := range h.Events {
		if event.TypeEvent == typeEvent {
			showEvents = append(showEvents, event)
		}
	}
	if _, err := easyjson.MarshalToWriter(models.Events(showEvents), c.Response().Writer); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return nil
}
