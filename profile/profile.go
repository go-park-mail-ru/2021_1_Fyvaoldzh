package profile

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"kudago/auth"
	"kudago/models"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/labstack/echo"
	"github.com/mailru/easyjson"
)

type UserHandler struct {
	Mu *sync.Mutex
}

func (h *UserHandler) GetProfile(c echo.Context) error {
	defer c.Request().Body.Close()

	cookie, err := c.Cookie("SID")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	if auth.Store[cookie.Value] == 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	profile := models.GetProfile(auth.Store[cookie.Value])

	// некрасиво, но пока
	for _, value := range models.PlanningEvent {
		if value.Uid == profile.Uid {
			profile.Event = append(profile.Event, value.Eid)
		}
	}

	log.Println(profile)

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

	user := models.GetUser(auth.Store[cookie.Value])
	profile := models.GetProfile(user.Id)

	ud := &models.UserData{}
	err = easyjson.UnmarshalFromReader(c.Request().Body, ud)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if len(ud.Name) != 0 {
		profile.Name = ud.Name
	}

	if len(ud.Password) != 0 {
		user.Password = ud.Password
	}

	if len(ud.Email) != 0 {
		if models.IsExistingEMail(ud.Email) {
			return echo.NewHTTPError(http.StatusBadRequest, "this email does exist")
		}
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

	return nil
}

func (h *UserHandler) UploadAvatar(c echo.Context) *echo.HTTPError {
	defer c.Request().Body.Close()

	cookie, err := c.Cookie("SID")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	if auth.Store[cookie.Value] == 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	user := models.GetUser(auth.Store[cookie.Value])
	profile := models.GetProfile(user.Id)

	img, err := c.FormFile("avatar")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if img == nil {
		return nil
	}

	src, err := img.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer src.Close()

	fileName := fmt.Sprint(user.Id) + img.Filename
	dst, err := os.Create(fileName)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	profile.Avatar = fileName
	return nil
}

func (h *UserHandler) GetUserProfile(c echo.Context) error {
	defer c.Request().Body.Close()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	user := models.GetOtherUserProfile(id)

	if user.Uid == 0 {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New("user does not exist"))
	}

	for _, value := range models.PlanningEvent {
		if value.Uid == user.Uid {
			user.Event = append(user.Event, value.Eid)
		}
	}

	log.Println(user)

	if _, err = easyjson.MarshalToWriter(user, c.Response().Writer); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (h *UserHandler) GetAvatar(c echo.Context) error {
	defer c.Request().Body.Close()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	profile := models.GetOtherUserProfile(id)

	if profile.Uid == 0 {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New("user does not exist"))
	}

	file, err := ioutil.ReadFile(profile.Avatar)
	if err != nil {
		log.Println("Cann't open file: " + profile.Avatar)
	} else {
			c.Response().Write(file)
	}

	return nil
}
