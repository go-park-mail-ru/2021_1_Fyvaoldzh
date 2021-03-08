package profile

import (
	"errors"
	"github.com/go-park-mail-ru/2021_1_Fyvaoldzh/auth"
	"github.com/go-park-mail-ru/2021_1_Fyvaoldzh/models"
	"github.com/labstack/echo"
	"github.com/mailru/easyjson"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type UserHandler struct {
	Mu *sync.Mutex
}

func (h *UserHandler) GetProfile(c echo.Context) *echo.HTTPError {
	defer c.Request().Body.Close()

	cookie, err := c.Cookie("SID")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	if auth.Store[cookie.Value] == 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	profile := auth.GetProfile(auth.Store[cookie.Value])

	if _, err = easyjson.MarshalToWriter(profile, c.Response().Writer); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (h *UserHandler) UpdateProfile(c echo.Context) *echo.HTTPError {
	defer c.Request().Body.Close()

	cookie, err := c.Cookie("SID")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	if auth.Store[cookie.Value] == 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	user := auth.GetUser(auth.Store[cookie.Value])
	profile := auth.GetProfile(user.Id)

	ud := &models.UserData{}
	err = easyjson.UnmarshalFromReader(c.Request().Body, ud)
	if err != nil {
		return  echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if len(ud.Name) != 0 {
		profile.Name = ud.Name
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
	return nil
}


func (h *UserHandler) GetUserProfile(c echo.Context) *echo.HTTPError {
	defer c.Request().Body.Close()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	user := auth.GetOtherUserProfile(id)

	if user.Uid == 0 {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New("user does not exist"))
	}

	if _, err = easyjson.MarshalToWriter(user, c.Response().Writer); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil

}