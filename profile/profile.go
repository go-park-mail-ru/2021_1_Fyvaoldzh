package profile

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2021_1_Fyvaoldzh/auth"
	"github.com/go-park-mail-ru/2021_1_Fyvaoldzh/models"
	"github.com/labstack/echo"
	"github.com/mailru/easyjson"
	"log"
	"net/http"
	"sync"
)

type UserHandler struct {
	Mu *sync.Mutex
}

func (h *UserHandler) GetProfile(c echo.Context) (string, *echo.HTTPError) {
	defer c.Request().Body.Close()

	cookie, err := c.Cookie("SID")
	if err != nil {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	if auth.Store[cookie.Value] == 0 {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	user := auth.GetUser(auth.Store[cookie.Value])
	profile := auth.GetProfile(user.Id)

	var js []byte

	js, err = json.Marshal(user)

	var profileJs []byte
	profileJs, err = json.Marshal(profile)

	return string(js) + string(profileJs), nil
}

func (h *UserHandler) UpdateProfile(c echo.Context) (string, *echo.HTTPError) {
	defer c.Request().Body.Close()

	cookie, err := c.Cookie("SID")
	if err != nil {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	if auth.Store[cookie.Value] == 0 {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	user := auth.GetUser(auth.Store[cookie.Value])
	profile := auth.GetProfile(user.Id)

	ud := &models.UserData{}
	err = easyjson.UnmarshalFromReader(c.Request().Body, ud)
	if err != nil {
		return "", echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if len(ud.Name) != 0 {
		user.Name = ud.Name
	}

	if len(ud.Password) != 0 {
		user.Password = ud.Password
	}

	if len(ud.Email) != 0 {
		profile.Email = ud.Email
		//проверка на повторяемость почты
	}

	if len(ud.About) != 0 {
		profile.About = ud.About
	}

	if len(ud.Birthday) != 0 {
		profile.Birthday = ud.Birthday
		// код на изменение age
	}

	if len(ud.City) != 0 {
		profile.City = ud.City
	}

	log.Println(user)
	log.Println(profile)
	return "", nil
}
