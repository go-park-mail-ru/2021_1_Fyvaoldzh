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

// временно айдишники
var id = 4

func isExistingUser(user *models.User) bool {
	for _, value := range UserBase {
		if value.Email == (*user).Email {
			return true
		}
	}
	return false
}

func (h *RegisterHandler) CreateUser(c echo.Context) *echo.HTTPError {
	defer c.Request().Body.Close()
	newUser := &models.User{}

	log.Println(c.Request().Body)
	err := easyjson.UnmarshalFromReader(c.Request().Body, newUser)
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if isExistingUser(newUser) {
		return echo.NewHTTPError(http.StatusBadRequest, "user with this email does exist")
	}

	newUser.Id = id
	id++

	h.Mu.Lock()
	UserBase = append(UserBase, *newUser)
	h.Mu.Unlock()
	return nil
}