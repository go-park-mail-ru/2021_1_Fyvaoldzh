package http

import (
	"errors"
	"fmt"
	"kudago/application/microservices/auth/client"
	client_chat "kudago/application/microservices/chat/client"
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
	"github.com/mailru/easyjson"
)

type ChatHandler struct {
	rpcAuth   client.IAuthClient
	rpcChat   client_chat.IChatClient
	Logger    logger.Logger
	sanitizer *custom_sanitizer.CustomSanitizer
}


func CreateChatHandler(e *echo.Echo, rpcA client.IAuthClient,
	sz *custom_sanitizer.CustomSanitizer, logger logger.Logger, auth middleware.Auth,
	rpcC client_chat.IChatClient) {

	chatHandler := ChatHandler{rpcChat: rpcC, rpcAuth: rpcA, Logger: logger, sanitizer: sz}

	//TODO групповой чат
	e.GET("/api/v1/dialogues", chatHandler.GetDialogues, auth.GetSession, middleware.GetPage)
	e.GET("/api/v1/dialogues/:id", chatHandler.GetOneDialogue, auth.GetSession, middleware.GetPage, middleware.GetId) //Здесь id собеседника, по просьбе Димы
	e.DELETE("/api/v1/dialogues/:id", chatHandler.DeleteDialogue, auth.GetSession, middleware.GetId)                  //Везде дальше и здесь id сообщения/диалога
	e.POST("/api/v1/send", chatHandler.SendMessage, auth.GetSession)
	e.DELETE("/api/v1/message/:id", chatHandler.DeleteMessage, auth.GetSession, middleware.GetId)
	e.POST("/api/v1/message", chatHandler.EditMessage, auth.GetSession)
	e.POST("/api/v1/message/mailing", chatHandler.Mailing, auth.GetSession)
	e.GET("/api/v1/dialogues/search", chatHandler.Search, auth.GetSession, middleware.GetPage)
}

func (ch ChatHandler) GetDialogues(c echo.Context) error {
	defer c.Request().Body.Close()

	start := time.Now()
	requestId := fmt.Sprintf("%016x", rand.Int())

	page := c.Get(constants.PageKey).(int)
	uid := c.Get(constants.UserIdKey).(uint64)

	dialogues, err := ch.rpcChat.GetAllDialogues(uid, page)
	if err != nil {
		ch.Logger.LogError(c, start, requestId, err)
		return err
	}
	dialogues = ch.sanitizer.SanitizeDialogueCards(dialogues)

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

	page := c.Get(constants.PageKey).(int)
	uid := c.Get(constants.UserIdKey).(uint64)
	id := c.Get(constants.IdKey).(int)

	dialogue, err := ch.rpcChat.GetOneDialogue(uid, uint64(id), page)
	if err != nil {
		ch.Logger.LogError(c, start, requestId, err)
		return err
	}
	ch.sanitizer.SanitizeDialogue(&dialogue)

	if _, err = easyjson.MarshalToWriter(dialogue, c.Response().Writer); err != nil {
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

	id := c.Get(constants.IdKey).(int)
	uid := c.Get(constants.UserIdKey).(uint64)

	err := ch.rpcChat.DeleteDialogue(uid, uint64(id))
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
	newMessageJSON := &models.NewMessageJSON{}

	if err := easyjson.UnmarshalFromReader(c.Request().Body, newMessageJSON); err != nil {
		ch.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusTeapot, err.Error())
	}

	newMessage := &models.NewMessage{}
	toInt, err := strconv.Atoi(newMessageJSON.To)
	if err != nil {
		ch.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if toInt <= 0 {
		err := errors.New("user id cannot be less than zero")
		ch.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusTeapot, err)
	}

	newMessage.To = uint64(toInt)
	newMessage.Text = newMessageJSON.Text

	uid := c.Get(constants.UserIdKey).(uint64)

	err = ch.rpcChat.SendMessage(newMessage, uid)
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

	id := c.Get(constants.IdKey).(int)
	uid := c.Get(constants.UserIdKey).(uint64)

	err := ch.rpcChat.DeleteMessage(uid, uint64(id))
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
	redactMessageJSON := &models.RedactMessageJSON{}

	if err := easyjson.UnmarshalFromReader(c.Request().Body, redactMessageJSON); err != nil {
		ch.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusTeapot, err.Error())
	}

	redactMessage := &models.RedactMessage{}
	idInt, err := strconv.Atoi(redactMessageJSON.ID)
	if err != nil {
		ch.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if idInt <= 0 {
		err := errors.New("message id cannot be less than zero")
		ch.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusTeapot, err)
	}

	redactMessage.ID = uint64(idInt)
	redactMessage.Text = redactMessageJSON.Text

	uid := c.Get(constants.UserIdKey).(uint64)

	err = ch.rpcChat.EditMessage(uid, redactMessage)
	if err != nil {
		ch.Logger.LogError(c, start, requestId, err)
		return err
	}

	ch.Logger.LogInfo(c, start, requestId)
	return nil
}

func (ch ChatHandler) Mailing(c echo.Context) error {
	defer c.Request().Body.Close()

	start := time.Now()
	requestId := fmt.Sprintf("%016x", rand.Int())

	mailingJSON := &models.MailingJSON{}
	if err := easyjson.UnmarshalFromReader(c.Request().Body, mailingJSON); err != nil {
		ch.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusTeapot, err.Error())
	}

	mailing := &models.Mailing{}

	eventIdInt, err := strconv.Atoi(mailingJSON.EventID)
	if err != nil {
		ch.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if eventIdInt <= 0 {
		err := errors.New("message id cannot be less than zero")
		ch.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusTeapot, err)
	}
	mailing.EventID = uint64(eventIdInt)

	for i := range mailingJSON.To {
		idInt, err := strconv.Atoi(mailingJSON.To[i])
		if err != nil {
			ch.Logger.LogError(c, start, requestId, err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		if idInt <= 0 {
			err := errors.New("user id cannot be less than zero")
			ch.Logger.LogError(c, start, requestId, err)
			return echo.NewHTTPError(http.StatusTeapot, err)
		}

		mailing.To = append(mailing.To, uint64(idInt))
	}

	uid := c.Get(constants.UserIdKey).(uint64)

	err = ch.rpcChat.Mailing(uid, mailing)
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

	page := c.Get(constants.PageKey).(int)
	uid := c.Get(constants.UserIdKey).(uint64)
	str := c.QueryParam("find")

	idParam := c.QueryParam("id")
	if idParam == "" {
		idParam = "0"
	}
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ch.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if id < 0 {
		err := errors.New("id cannot be less than zero")
		ch.Logger.LogError(c, start, requestId, err)
		return err
	}

	messages, err := ch.rpcChat.Search(uid, id, str, page)
	if err != nil {
		ch.Logger.LogError(c, start, requestId, err)
		return err
	}
	messages = ch.sanitizer.SanitizeMessages(messages)

	if _, err = easyjson.MarshalToWriter(messages, c.Response().Writer); err != nil {
		ch.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	ch.Logger.LogInfo(c, start, requestId)
	return nil
}
