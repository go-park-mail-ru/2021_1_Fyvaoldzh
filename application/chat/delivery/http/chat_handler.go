package http

import (
	"errors"
	"fmt"
	"kudago/application/chat"
	"kudago/application/microservices/auth/client"
	"kudago/application/models"
	"kudago/application/server/middleware"
	"kudago/pkg/constants"
	"kudago/pkg/custom_sanitizer"
	"kudago/pkg/logger"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/mailru/easyjson"
)

type ChatHandler struct {
	UseCase   chat.UseCase
	rpcAuth   client.AuthClient
	Logger    logger.Logger
	sanitizer *custom_sanitizer.CustomSanitizer
}

func CreateChatHandler(e *echo.Echo, uc chat.UseCase, rpcA client.AuthClient,
	sz *custom_sanitizer.CustomSanitizer, logger logger.Logger, auth middleware.Auth) {

	chatHandler := ChatHandler{UseCase: uc, rpcAuth: rpcA, Logger: logger, sanitizer: sz}

	//TODO групповой чат
	e.GET("/api/v1/dialogues", chatHandler.GetDialogues, auth.GetSession)
	e.GET("/api/v1/dialogues/:id", chatHandler.GetOneDialogue, auth.GetSession)
	e.DELETE("/api/v1/dialogues/:id", chatHandler.DeleteDialogue, auth.GetSession)
	e.POST("/api/v1/send", chatHandler.SendMessage, auth.GetSession)
	e.DELETE("/api/v1/message/:id", chatHandler.DeleteMessage, auth.GetSession)
	e.POST("/api/v1/message/:id", chatHandler.EditMessage, auth.GetSession)
	e.GET("/api/v1/dialogues/search", chatHandler.Search, auth.GetSession)
}

//Есть ли возможность как-то шаблонизировать функции ниже, везде одно и то же начало, мб создать функцию, принимающую функцию
//и в зависимости от метода кидать туда свой метод юзкейса?

func (ch ChatHandler) GetDialogues(c echo.Context) error {
	defer c.Request().Body.Close()

	start := time.Now()
	requestId := fmt.Sprintf("%016x", rand.Int())
	pageParam := c.QueryParam("page")
	if pageParam == "" {
		pageParam = "1"
	}
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		ch.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if page == 0 {
		page = 1
	}

	uid := c.Get(constants.UserIdKey).(uint64)

	dialogues, err := ch.UseCase.GetAllDialogues(uid, page)
	//dialogues = ch.sanitizer.SanitizeEventCards(dialogues) SanitizeDialogues А нужно ли?
	if err != nil {
		ch.Logger.LogError(c, start, requestId, err)
		return err
	}

	if _, err = easyjson.MarshalToWriter(dialogues, c.Response().Writer); err != nil {
		ch.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	ch.Logger.LogInfo(c, start, requestId)
	return nil
}

func (ch ChatHandler) GetOneDialogue(c echo.Context) error {
	defer c.Request().Body.Close()

	start := time.Now()
	requestId := fmt.Sprintf("%016x", rand.Int())

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ch.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if id <= 0 {
		err := errors.New("user id cannot be less than zero")
		ch.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusTeapot, err)
	}

	pageParam := c.QueryParam("page")
	if pageParam == "" {
		pageParam = "1"
	}
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		ch.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if page == 0 {
		page = 1
	}

	uid := c.Get(constants.UserIdKey).(uint64)
	log.Info(uid)

	messages, err := ch.UseCase.GetOneDialogue(uid, uint64(id), page)
	//dialogues = ch.sanitizer.SanitizeEventCards(dialogues) SanitizeDialogues А нужно ли?
	if err != nil {
		ch.Logger.LogError(c, start, requestId, err)
		return err
	}

	if _, err = easyjson.MarshalToWriter(messages, c.Response().Writer); err != nil {
		ch.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	ch.Logger.LogInfo(c, start, requestId)
	return nil
}

func (ch ChatHandler) DeleteDialogue(c echo.Context) error {
	defer c.Request().Body.Close()

	start := time.Now()
	requestId := fmt.Sprintf("%016x", rand.Int())
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ch.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	uid := c.Get(constants.UserIdKey).(uint64)

	err = ch.UseCase.DeleteDialogue(uid, uint64(id))
	if err != nil {
		ch.Logger.LogError(c, start, requestId, err)
		return err
	}

	ch.Logger.LogInfo(c, start, requestId)
	return nil
}

func (ch ChatHandler) SendMessage(c echo.Context) error {
	defer c.Request().Body.Close()

	start := time.Now()
	requestId := fmt.Sprintf("%016x", rand.Int())
	newMessage := &models.NewMessage{}

	if err := easyjson.UnmarshalFromReader(c.Request().Body, newMessage); err != nil {
		ch.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusTeapot, err.Error())
	}

	if newMessage.To <= 0 {
		err := errors.New("user id cannot be less than zero")
		ch.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusTeapot, err)
	}

	uid := c.Get(constants.UserIdKey).(uint64)

	err := ch.UseCase.SendMessage(newMessage, uid)
	if err != nil {
		ch.Logger.LogError(c, start, requestId, err)
		return err
	}

	ch.Logger.LogInfo(c, start, requestId)
	return nil
}

func (ch ChatHandler) DeleteMessage(c echo.Context) error {
	defer c.Request().Body.Close()

	start := time.Now()
	requestId := fmt.Sprintf("%016x", rand.Int())
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ch.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	uid := c.Get(constants.UserIdKey).(uint64)

	err = ch.UseCase.DeleteMessage(uid, uint64(id))
	if err != nil {
		ch.Logger.LogError(c, start, requestId, err)
		return err
	}

	ch.Logger.LogInfo(c, start, requestId)
	return nil
}

func (ch ChatHandler) EditMessage(c echo.Context) error {
	defer c.Request().Body.Close()

	start := time.Now()
	requestId := fmt.Sprintf("%016x", rand.Int())
	redactMessage := &models.RedactMessage{}

	if err := easyjson.UnmarshalFromReader(c.Request().Body, redactMessage); err != nil {
		ch.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusTeapot, err.Error())
	}

	uid := c.Get(constants.UserIdKey).(uint64)

	err := ch.UseCase.EditMessage(uid, redactMessage)
	if err != nil {
		ch.Logger.LogError(c, start, requestId, err)
		return err
	}

	ch.Logger.LogInfo(c, start, requestId)
	return nil
}

func (ch ChatHandler) Search(c echo.Context) error {
	defer c.Request().Body.Close()

	start := time.Now()
	requestId := fmt.Sprintf("%016x", rand.Int())

	str := c.QueryParam("find")

	pageParam := c.QueryParam("page")
	if pageParam == "" {
		pageParam = "1"
	}
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		ch.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if page == 0 {
		page = 1
	}

	var id int
	if c.QueryParam("id") == "" {
		//вот тут костыль прям некрасивый, нет идей, как можно пофиксить?
		id = -1
	} else {
		idatoi, err := strconv.Atoi(c.QueryParam("id"))
		if err != nil {
			ch.Logger.LogError(c, start, requestId, err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if idatoi < 0 {
			err := errors.New("id cannot be less than zero")
			ch.Logger.LogError(c, start, requestId, err)
			return err
		} else {
			id = idatoi
		}
	}

	uid := c.Get(constants.UserIdKey).(uint64)

	messages, err := ch.UseCase.Search(uid, id, str, page)
	//dialogues = ch.sanitizer.SanitizeEventCards(dialogues) SanitizeDialogues А нужно ли?
	if err != nil {
		ch.Logger.LogError(c, start, requestId, err)
		return err
	}

	if _, err = easyjson.MarshalToWriter(messages, c.Response().Writer); err != nil {
		ch.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	ch.Logger.LogInfo(c, start, requestId)
	return nil
}
