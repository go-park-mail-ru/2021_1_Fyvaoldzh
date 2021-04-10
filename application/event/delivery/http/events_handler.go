package http

import (
	"errors"
	"fmt"
	"kudago/application/event"
	"kudago/application/models"
	"kudago/pkg/constants"
	"kudago/pkg/infrastructure"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/mailru/easyjson"
)

type EventHandler struct {
	UseCase event.UseCase
	Sm      *infrastructure.SessionManager
}

func CreateEventHandler(e *echo.Echo, uc event.UseCase, sm *infrastructure.SessionManager) {

	eventHandler := EventHandler{UseCase: uc, Sm: sm}

	e.GET("/api/v1/", eventHandler.GetAllEvents)
	e.GET("/api/v1/event/:id", eventHandler.GetOneEvent)
	e.GET("/api/v1/event", eventHandler.GetEvents)
	e.GET("/api/v1/search", eventHandler.FindEvents)
	//create & delete & save вообще не должно быть, пользователь НИКАК не может создавать и удалять что-либо, только админ работает с БД
	e.POST("/api/v1/create", eventHandler.Create)
	e.DELETE("/api/v1/event/:id", eventHandler.Delete)
	e.POST("/api/v1/save/:id", eventHandler.Save)
	e.GET("api/v1/event/:id/image", eventHandler.GetImage)
}

func (eh EventHandler) GetAllEvents(c echo.Context) error {
	defer c.Request().Body.Close()

	page, err := strconv.Atoi(c.QueryParam("page"))
	if c.QueryParam("page") == "" {
		page = 1
	} else {
		if err != nil {
			log.Println(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	if page == 0 {
		page = 1
	}

	events, err := eh.UseCase.GetAllEvents(page)
	if err != nil {
		log.Println(err)
		return err
	}

	if _, err = easyjson.MarshalToWriter(events, c.Response().Writer); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (eh EventHandler) GetUserID(c echo.Context) (uint64, error) {
	cookie, err := c.Cookie(constants.SessionCookieName)
	if err != nil && cookie != nil {
		log.Println(err)
		return 0, errors.New("user is not authorized")
	}

	var uid uint64
	var exists bool

	if cookie != nil {
		exists, uid, err = eh.Sm.CheckSession(cookie.Value)
		if err != nil {
			log.Println(err)
			return 0, err
		}

		if !exists {
			log.Println("cookie does not exist")
			return 0, errors.New("user is not authorized")
		}
		return uid, nil
	}
	log.Println("got no cookie")
	return 0, errors.New("user is not authorized")
}

func (eh EventHandler) GetOneEvent(c echo.Context) error {
	defer c.Request().Body.Close()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	ev, err := eh.UseCase.GetOneEvent(uint64(id))
	if err != nil {
		return err
	}

	if uid, err := eh.GetUserID(c); err == nil {
		if err := eh.UseCase.RecomendSystem(uid, ev.Category); err != nil {
			log.Println(err)
		}
	} else {
		log.Println("cannot get userID")
	}

	if _, err = easyjson.MarshalToWriter(ev, c.Response().Writer); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (eh EventHandler) GetEvents(c echo.Context) error {
	defer c.Request().Body.Close()

	category := c.QueryParam("category")
	events, err := eh.UseCase.GetEventsByCategory(category)

	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if _, err := easyjson.MarshalToWriter(events, c.Response().Writer); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (eh EventHandler) Create(c echo.Context) error {
	defer c.Request().Body.Close()

	newEvent := &models.Event{}

	if err := easyjson.UnmarshalFromReader(c.Request().Body, newEvent); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusTeapot, err.Error())
	}

	if err := eh.UseCase.CreateNewEvent(newEvent); err != nil {
		log.Println(err)
		return err
	}

	return c.JSON(http.StatusOK, *newEvent)
}

func (eh EventHandler) Delete(c echo.Context) error {
	defer c.Request().Body.Close()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = eh.UseCase.Delete(uint64(id))
	if err != nil {
		log.Println(err)
		return err
	}

	return c.String(http.StatusOK, "Event with id "+fmt.Sprint(id)+" successfully deleted \n")
}

func (eh EventHandler) Save(c echo.Context) error {
	defer c.Request().Body.Close()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusTeapot, err.Error())
	}

	img, err := c.FormFile("image")
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusTeapot, err.Error())
	}

	err = eh.UseCase.SaveImage(uint64(id), img)

	if err != nil {
		log.Println(err)
		return err
	}
	return c.JSON(http.StatusOK, "Picture changed successfully")
}

func (eh EventHandler) GetImage(c echo.Context) error {
	defer c.Request().Body.Close()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	file, err := eh.UseCase.GetImage(uint64(id))
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = c.Response().Write(file)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (eh EventHandler) FindEvents(c echo.Context) error {
	defer c.Request().Body.Close()

	str := c.QueryParam("find")
	events, err := eh.UseCase.FindEvents(str)

	if err != nil {
		log.Println(err)
		return err
	}

	if _, err := easyjson.MarshalToWriter(events, c.Response().Writer); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}
