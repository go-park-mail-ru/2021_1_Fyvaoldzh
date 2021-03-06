package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/labstack/echo"
)

type Handlers struct {
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

type LoginHandler struct{

}