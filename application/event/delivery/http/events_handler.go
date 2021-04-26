package http

import (
	"errors"
	"fmt"
	"kudago/application/event"
	"kudago/application/microservices/auth/client"
	"kudago/application/models"
	"kudago/pkg/constants"
	"kudago/pkg/custom_sanitizer"
	"kudago/pkg/logger"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/mailru/easyjson"
)

type EventHandler struct {
	UseCase   event.UseCase
	rpcAuth client.AuthClient
	Logger    logger.Logger
	sanitizer *custom_sanitizer.CustomSanitizer
}

func CreateEventHandler(e *echo.Echo, uc event.UseCase, rpcA client.AuthClient, sz *custom_sanitizer.CustomSanitizer, logger logger.Logger) {
	eventHandler := EventHandler{UseCase: uc, rpcAuth: rpcA, Logger: logger, sanitizer: sz}

	e.GET("/api/v1/", eventHandler.GetAllEvents)
	e.GET("/api/v1/event/:id", eventHandler.GetOneEvent)
	e.GET("/api/v1/event", eventHandler.GetEvents)
	e.GET("/api/v1/search", eventHandler.FindEvents)
	//create & delete & save вообще не должно быть, пользователь НИКАК не может создавать и удалять что-либо, только админ работает с БД
	e.POST("/api/v1/create", eventHandler.Create)
	e.DELETE("/api/v1/event/:id", eventHandler.Delete)
	e.POST("/api/v1/save/:id", eventHandler.Save)
	e.GET("api/v1/event/:id/image", eventHandler.GetImage)
	// TODO фикс названия ручки был
	e.GET("/api/v1/recommend", eventHandler.Recommend)
}

func (eh EventHandler) Recommend(c echo.Context) error {
	defer c.Request().Body.Close()

	start := time.Now()
	requestId := fmt.Sprintf("%016x", rand.Int())
	var page int
	if c.QueryParam("page") == "" {
		page = 1
	} else {
		pageatoi, err := strconv.Atoi(c.QueryParam("page"))
		if err != nil {
			eh.Logger.LogError(c, start, requestId, err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if pageatoi == 0 {
			page = 1
		} else {
			page = pageatoi
		}
	}

	if uid, err := eh.GetUserID(c); err == nil {
		events, err := eh.UseCase.GetRecommended(uid, page)
		events = eh.sanitizer.SanitizeEventCards(events)
		if err != nil {
			eh.Logger.LogError(c, start, requestId, err)
			return err
		}

		if _, err = easyjson.MarshalToWriter(events, c.Response().Writer); err != nil {
			eh.Logger.LogError(c, start, requestId, err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		eh.Logger.LogInfo(c, start, requestId)
		return nil
	} else {
		eh.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
}

func (eh EventHandler) GetAllEvents(c echo.Context) error {
	defer c.Request().Body.Close()

	start := time.Now()
	requestId := fmt.Sprintf("%016x", rand.Int())
	var page int
	if c.QueryParam("page") == "" {
		page = 1
	} else {
		pageatoi, err := strconv.Atoi(c.QueryParam("page"))
		if err != nil {
			eh.Logger.LogError(c, start, requestId, err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if pageatoi == 0 {
			page = 1
		} else {
			page = pageatoi
		}
	}

	events, err := eh.UseCase.GetAllEvents(page)
	events = eh.sanitizer.SanitizeEventCards(events)
	if err != nil {
		eh.Logger.LogError(c, start, requestId, err)
		return err
	}

	if _, err = easyjson.MarshalToWriter(events, c.Response().Writer); err != nil {
		eh.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	eh.Logger.LogInfo(c, start, requestId)

	return nil
}

func (eh EventHandler) GetUserID(c echo.Context) (uint64, error) {
	start := time.Now()
	requestId := fmt.Sprintf("%016x", rand.Int())
	cookie, err := c.Cookie(constants.SessionCookieName)
	if err != nil && cookie != nil {
		eh.Logger.LogWarn(c, start, requestId, err)
		return 0, errors.New("user is not authorized")
	}

	var uid uint64
	var exists bool

	if cookie != nil {
		exists, uid, err = eh.rpcAuth.Check(cookie.Value)
		if err != nil {
			eh.Logger.LogWarn(c, start, requestId, err)
			return 0, err
		}

		if !exists {
			eh.Logger.LogWarn(c, start, requestId, err)
			return 0, errors.New("user is not authorized")
		}

		return uid, nil
	}
	eh.Logger.LogWarn(c, start, requestId, err)
	return 0, errors.New("user is not authorized")
}

func (eh EventHandler) GetOneEvent(c echo.Context) error {
	defer c.Request().Body.Close()

	start := time.Now()
	requestId := fmt.Sprintf("%016x", rand.Int())
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		eh.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	ev, err := eh.UseCase.GetOneEvent(uint64(id))
	if err != nil {
		eh.Logger.LogError(c, start, requestId, err)
		return err
	}
	eh.sanitizer.SanitizeEvent(&ev)
	if uid, err := eh.GetUserID(c); err == nil {
		if err := eh.UseCase.RecomendSystem(uid, ev.Category); err != nil {
			eh.Logger.LogWarn(c, start, requestId, err)
		}
	} else {
		eh.Logger.LogWarn(c, start, requestId, err)
	}

	if _, err = easyjson.MarshalToWriter(ev, c.Response().Writer); err != nil {
		eh.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	eh.Logger.LogInfo(c, start, requestId)
	return nil
}

func (eh EventHandler) GetEvents(c echo.Context) error {
	defer c.Request().Body.Close()

	start := time.Now()
	requestId := fmt.Sprintf("%016x", rand.Int())
	category := c.QueryParam("category")
	var page int
	if c.QueryParam("page") == "" {
		page = 1
	} else {
		pageatoi, err := strconv.Atoi(c.QueryParam("page"))
		if err != nil {
			eh.Logger.LogError(c, start, requestId, err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if pageatoi == 0 {
			page = 1
		} else {
			page = pageatoi
		}
	}
	events, err := eh.UseCase.GetEventsByCategory(category, page)
	events = eh.sanitizer.SanitizeEventCards(events)

	if err != nil {
		eh.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if _, err := easyjson.MarshalToWriter(events, c.Response().Writer); err != nil {
		eh.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	eh.Logger.LogInfo(c, start, requestId)
	return nil
}

//Эти функции будут удалены, поэтому почти не изменялись с переноса архитектуры
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

	start := time.Now()
	requestId := fmt.Sprintf("%016x", rand.Int())
	str := c.QueryParam("find")
	category := c.QueryParam("category")
	var page int
	if c.QueryParam("page") == "" {
		page = 1
	} else {
		pageatoi, err := strconv.Atoi(c.QueryParam("page"))
		if err != nil {
			eh.Logger.LogError(c, start, requestId, err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if pageatoi == 0 {
			page = 1
		} else {
			page = pageatoi
		}
	}

	events, err := eh.UseCase.FindEvents(str, category, page)
	events = eh.sanitizer.SanitizeEventCards(events)
	if err != nil {
		eh.Logger.LogError(c, start, requestId, err)
		return err
	}

	if _, err := easyjson.MarshalToWriter(events, c.Response().Writer); err != nil {
		eh.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	eh.Logger.LogInfo(c, start, requestId)
	return nil
}

