package auth

import (
	"github.com/labstack/echo"
	"github.com/mailru/easyjson"
	"log"
	"net/http"
	"sync"

	"github.com/go-park-mail-ru/2021_1_Fyvaoldzh/models"

)

type RegisterHandler struct {
	Mu     *sync.Mutex
}

func (h *RegisterHandler) CreateUser(c echo.Context) error {
	defer c.Request().Body.Close()
	newUser := &models.User{}

	err := easyjson.UnmarshalFromReader(c.Request().Body, newUser)
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	userBase = append(userBase, *newUser)

	return nil
}