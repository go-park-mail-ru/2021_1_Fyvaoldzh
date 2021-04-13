package http

import (
	"errors"
	"fmt"
	"kudago/application/event"
	"kudago/application/models"
	"kudago/pkg/constants"
	"kudago/pkg/custom_sanitizer"
	"kudago/pkg/infrastructure"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
)

type EventHandler struct {
	UseCase   event.UseCase
	Sm        *infrastructure.SessionManager
	Logger    *zap.SugaredLogger
	sanitizer *custom_sanitizer.CustomSanitizer
}

func CreateEventHandler(e *echo.Echo, uc event.UseCase, sm *infrastructure.SessionManager, sz *custom_sanitizer.CustomSanitizer, logger *zap.SugaredLogger) {

	eventHandler := EventHandler{UseCase: uc, Sm: sm, Logger: logger, sanitizer: sz}

	e.GET("/api/v1/", eventHandler.GetAllEvents)
	e.GET("/api/v1/event/:id", eventHandler.GetOneEvent)
	e.GET("/api/v1/event", eventHandler.GetEvents)
	e.GET("/api/v1/search", eventHandler.FindEvents)
	//create & delete & save вообще не должно быть, пользователь НИКАК не может создавать и удалять что-либо, только админ работает с БД
	e.POST("/api/v1/create", eventHandler.Create)
	e.DELETE("/api/v1/event/:id", eventHandler.Delete)
	e.POST("/api/v1/save/:id", eventHandler.Save)
	e.GET("api/v1/event/:id/image", eventHandler.GetImage)
	e.GET("/api/v1/recomend", eventHandler.Recomend)
}

func (eh EventHandler) LogInfo(c echo.Context, start time.Time, request_id string) {
	eh.Logger.Info(c.Request().URL.Path,
		zap.String("method", c.Request().Method),
		zap.String("remote_addr", c.Request().RemoteAddr),
		zap.String("url", c.Request().URL.Path),
		zap.Duration("work_time", time.Since(start)),
		zap.String("request_id", request_id),
	)
}

func (eh EventHandler) LogWarn(c echo.Context, start time.Time, request_id string, err error) {
	eh.Logger.Warn(c.Request().URL.Path,
		zap.String("method", c.Request().Method),
		zap.String("remote_addr", c.Request().RemoteAddr),
		zap.String("url", c.Request().URL.Path),
		zap.Duration("work_time", time.Since(start)),
		zap.String("request_id", request_id),
		zap.Errors("error", []error{err}),
	)
}

func (eh EventHandler) LogError(c echo.Context, start time.Time, request_id string, err error) {
	eh.Logger.Error(c.Request().URL.Path,
		zap.String("method", c.Request().Method),
		zap.String("remote_addr", c.Request().RemoteAddr),
		zap.String("url", c.Request().URL.Path),
		zap.Duration("work_time", time.Since(start)),
		zap.String("request_id", request_id),
		zap.Errors("error", []error{err}),
	)
}

func (eh EventHandler) Recomend(c echo.Context) error {
	defer c.Request().Body.Close()

	start := time.Now()
	request_id := fmt.Sprintf("%016x", rand.Int())
	page, err := strconv.Atoi(c.QueryParam("page"))
	if c.QueryParam("page") == "" {
		page = 1
	} else {
		if err != nil {
			eh.LogError(c, start, request_id, err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	if page == 0 {
		page = 1
	}

	if uid, err := eh.GetUserID(c); err == nil {
		events, err := eh.UseCase.GetRecomended(uid, page)
		events = eh.sanitizer.SanitizeEventCards(events)
		if err != nil {
			eh.LogError(c, start, request_id, err)
			return err
		}

		if _, err = easyjson.MarshalToWriter(events, c.Response().Writer); err != nil {
			eh.LogError(c, start, request_id, err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		eh.LogInfo(c, start, request_id)
		return nil
	} else {
		eh.LogError(c, start, request_id, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
}

func (eh EventHandler) GetAllEvents(c echo.Context) error {
	defer c.Request().Body.Close()

	start := time.Now()
	request_id := fmt.Sprintf("%016x", rand.Int())
	page, err := strconv.Atoi(c.QueryParam("page"))
	if c.QueryParam("page") == "" {
		page = 1
	} else {
		if err != nil {
			eh.LogError(c, start, request_id, err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	if page == 0 {
		page = 1
	}

	events, err := eh.UseCase.GetAllEvents(page)
	events = eh.sanitizer.SanitizeEventCards(events)
	if err != nil {
		eh.LogError(c, start, request_id, err)
		return err
	}

	if _, err = easyjson.MarshalToWriter(events, c.Response().Writer); err != nil {
		eh.LogError(c, start, request_id, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	eh.LogInfo(c, start, request_id)

	return nil
}

func (eh EventHandler) GetUserID(c echo.Context) (uint64, error) {
	start := time.Now()
	request_id := fmt.Sprintf("%016x", rand.Int())
	cookie, err := c.Cookie(constants.SessionCookieName)
	if err != nil && cookie != nil {
		eh.LogWarn(c, start, request_id, err)
		return 0, errors.New("user is not authorized")
	}

	var uid uint64
	var exists bool

	if cookie != nil {
		exists, uid, err = eh.Sm.CheckSession(cookie.Value)
		if err != nil {
			eh.LogWarn(c, start, request_id, err)
			return 0, err
		}

		if !exists {
			eh.LogWarn(c, start, request_id, err)
			return 0, errors.New("user is not authorized")
		}

		return uid, nil
	}
	eh.LogWarn(c, start, request_id, err)
	return 0, errors.New("user is not authorized")
}

func (eh EventHandler) GetOneEvent(c echo.Context) error {
	defer c.Request().Body.Close()

	start := time.Now()
	request_id := fmt.Sprintf("%016x", rand.Int())
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		eh.LogError(c, start, request_id, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	ev, err := eh.UseCase.GetOneEvent(uint64(id))
	if err != nil {
		eh.LogError(c, start, request_id, err)
		return err
	}
	eh.sanitizer.SanitizeEvent(&ev)
	if uid, err := eh.GetUserID(c); err == nil {
		if err := eh.UseCase.RecomendSystem(uid, ev.Category); err != nil {
			eh.LogWarn(c, start, request_id, err)
		}
	} else {
		eh.LogWarn(c, start, request_id, err)
	}

	if _, err = easyjson.MarshalToWriter(ev, c.Response().Writer); err != nil {
		eh.LogError(c, start, request_id, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	eh.LogInfo(c, start, request_id)
	return nil
}

func (eh EventHandler) GetEvents(c echo.Context) error {
	defer c.Request().Body.Close()

	start := time.Now()
	request_id := fmt.Sprintf("%016x", rand.Int())
	category := c.QueryParam("category")
	page, err := strconv.Atoi(c.QueryParam("page"))
	if c.QueryParam("page") == "" {
		page = 1
	} else {
		if err != nil {
			eh.LogError(c, start, request_id, err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	if page == 0 {
		page = 1
	}
	events, err := eh.UseCase.GetEventsByCategory(category, page)
	events = eh.sanitizer.SanitizeEventCards(events)

	if err != nil {
		eh.LogError(c, start, request_id, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if _, err := easyjson.MarshalToWriter(events, c.Response().Writer); err != nil {
		eh.LogError(c, start, request_id, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	eh.LogInfo(c, start, request_id)
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
	request_id := fmt.Sprintf("%016x", rand.Int())
	str := c.QueryParam("find")
	category := c.QueryParam("category")
	page, err := strconv.Atoi(c.QueryParam("page"))
	if c.QueryParam("page") == "" {
		page = 1
	} else {
		if err != nil {
			eh.LogError(c, start, request_id, err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	if page == 0 {
		page = 1
	}

	events, err := eh.UseCase.FindEvents(str, category, page)
	events = eh.sanitizer.SanitizeEventCards(events)
	if err != nil {
		eh.LogError(c, start, request_id, err)
		return err
	}

	if _, err := easyjson.MarshalToWriter(events, c.Response().Writer); err != nil {
		eh.LogError(c, start, request_id, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	eh.LogInfo(c, start, request_id)
	return nil
}
