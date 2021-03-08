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
		if value.Login == (*user).Login {
			return true
		}
	}
	return false
}

func (h *RegisterHandler) CreateUser(c echo.Context) *echo.HTTPError {
	defer c.Request().Body.Close()
	newData := &models.RegData{}

	log.Println(c.Request().Body)
	err := easyjson.UnmarshalFromReader(c.Request().Body, newData)
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	newUser := &models.User{}
	newUser.Login = newData.Login
	newUser.Password = newData.Password

	if isExistingUser(newUser) {
		return echo.NewHTTPError(http.StatusBadRequest, "user with this login does exist")
	}

	newUser.Id = id
	id++

	newProfile := &models.UserOwnProfile{}
	newProfile.Uid = newUser.Id
	newProfile.Name = newData.Name


	h.Mu.Lock()
	UserBase = append(UserBase, newUser)
	ProfileBase = append(ProfileBase, newProfile)
	h.Mu.Unlock()
	return nil
}