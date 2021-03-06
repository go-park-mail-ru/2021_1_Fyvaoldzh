package events

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
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

func All(h *Handlers, c echo.Context) error {
	encoder := json.NewEncoder(c.Response().Writer)
	h.Mu.Lock()
	err := encoder.Encode(h.Events)
	h.Mu.Unlock()
	if err != nil {
		log.Println(err)
		c.Response().Write([]byte("{}"))
		return c.JSON(http.StatusOK, encoder)
	}
	return c.JSON(http.StatusOK, encoder)
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

func CreateEvent(c echo.Context) error {
	name := c.FormValue("name")
	img, err := c.FormFile("image")
	if err != nil {
		log.Println(err)
	}
	src, err := img.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(img.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	return c.HTML(http.StatusOK, "<b>Thank you! "+name+"</b>")
}
