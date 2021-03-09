package auth

import (
	"github.com/labstack/echo"
	"github.com/mailru/easyjson"
	"kudago/models"
	"log"
	"net/http"
	"sync"
	"time"
)

type RegisterHandler struct {
	Mu *sync.Mutex
}

// временно айдишники
var id = 4

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

	if models.IsExistingUser(newUser) {
		return echo.NewHTTPError(http.StatusBadRequest, "user with this login does exist")
	}

	newUser.Id = id
	id++

	newProfile := &models.UserOwnProfile{}
	newProfile.Uid = newUser.Id
	newProfile.Name = newData.Name

	h.Mu.Lock()
	models.UserBase = append(models.UserBase, newUser)
	models.ProfileBase = append(models.ProfileBase, newProfile)
	h.Mu.Unlock()

	key := RandStringRunes(32)

	newCookie := &http.Cookie{
		Name:     "SID",
		Value:    key,
		Expires:  time.Now().Add(10 * time.Hour),
		SameSite: http.SameSiteNoneMode,
	}

	Store[key] = newUser.Id
	c.SetCookie(newCookie)

	return nil
}
