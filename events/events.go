package events

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/labstack/echo"
)

type EventInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Event struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Handlers struct {
	Events []Event
	Mu     *sync.Mutex
}

func (h *Handlers) All(c echo.Context) {
	encoder := json.NewEncoder(c.Response().Writer)
	h.Mu.Lock()
	err := encoder.Encode(h.Events)
	h.Mu.Unlock()
	if err != nil {
		log.Println(err)
		c.Response().Write([]byte("{}"))
		return
	}
}

func (h *Handlers) Create(c echo.Context) {
	defer c.Request().Body.Close()

	decoder := json.NewDecoder(c.Request().Body)

	newEventInput := new(EventInput)
	err := decoder.Decode(newEventInput)
	if err != nil {
		log.Println(err)
		c.Response().Write([]byte("{}"))
		return
	}

	fmt.Println(newEventInput)
	h.Mu.Lock()

	var id uint64 = 0
	if len(h.Events) > 0 {
		id = h.Events[len(h.Events)-1].ID + 1
	}

	h.Events = append(h.Events, Event{
		ID:          id,
		Title:       newEventInput.Title,
		Description: newEventInput.Description,
	})
	h.Mu.Unlock()
}

func GetEvent(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusOK, "event "+id)
}

func Show(c echo.Context) error {
	city := c.QueryParam("city")
	typeEvent := c.QueryParam("typeEvent")
	return c.JSON(http.StatusOK, "type "+typeEvent+" in city "+city)
}
