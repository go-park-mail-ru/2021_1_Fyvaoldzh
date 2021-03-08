package profile

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2021_1_Fyvaoldzh/auth"
	"github.com/labstack/echo"
	"net/http"
	"sync"
)

type UserHandler struct {
	Mu     *sync.Mutex
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


	return "", nil
}