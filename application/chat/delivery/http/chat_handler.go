package http

import (
	"errors"
	"fmt"
	"kudago/application/chat"
	"kudago/application/models"
	"kudago/pkg/constants"
	"kudago/pkg/custom_sanitizer"
	"kudago/pkg/infrastructure"
	"kudago/pkg/logger"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/mailru/easyjson"
)

type ChatHandler struct {
	UseCase   chat.UseCase
	Sm        infrastructure.SessionTarantool
	Logger    logger.Logger
	sanitizer *custom_sanitizer.CustomSanitizer
}

func CreateChatHandler(e *echo.Echo, uc chat.UseCase, sm infrastructure.SessionTarantool, sz *custom_sanitizer.CustomSanitizer, logger logger.Logger) {

	chatHandler := ChatHandler{UseCase: uc, Sm: sm, Logger: logger, sanitizer: sz}

	//TODO групповой чат
	e.GET("/api/v1/dialogues", chatHandler.GetDialogues)
	e.GET("/api/v1/dialogues/:id", chatHandler.GetOneDialogue)
	e.DELETE("/api/v1/dialogues/:id", chatHandler.DeleteDialogue) //Удаляем диалог для обеих сторон или для одной?
	e.POST("/api/v1/send/:id", chatHandler.SendMessage)
	e.DELETE("/api/v1/message/:id", chatHandler.DeleteMessage) //Удаляем сообщения для обеих сторон или для одной?
	e.POST("/api/v1/message/:id", chatHandler.EditMessage)
	e.GET("/api/v1/dialogues/search", chatHandler.Search) //По фолловерам тоже, если без id
}

func (ch ChatHandler) GetUserID(c echo.Context) (uint64, error) {
	start := time.Now()
	requestId := fmt.Sprintf("%016x", rand.Int())
	cookie, err := c.Cookie(constants.SessionCookieName)
	if err != nil && cookie != nil {
		ch.Logger.LogWarn(c, start, requestId, err)
		return 0, errors.New("user is not authorized")
	}

	var uid uint64
	var exists bool

	if cookie != nil {
		exists, uid, err = ch.Sm.CheckSession(cookie.Value)
		if err != nil {
			ch.Logger.LogWarn(c, start, requestId, err)
			return 0, err
		}

		if !exists {
			ch.Logger.LogWarn(c, start, requestId, err)
			return 0, errors.New("user is not authorized")
		}

		return uid, nil
	}
	ch.Logger.LogWarn(c, start, requestId, err)
	return 0, errors.New("user is not authorized")
}

func (ch ChatHandler) GetDialogues(c echo.Context) error {
	defer c.Request().Body.Close()

	start := time.Now()
	requestId := fmt.Sprintf("%016x", rand.Int())
	var page int
	if c.QueryParam("page") == "" {
		page = 1
	} else {
		pageatoi, err := strconv.Atoi(c.QueryParam("page"))
		if err != nil {
			ch.Logger.LogError(c, start, requestId, err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if pageatoi == 0 {
			page = 1
		} else {
			page = pageatoi
		}
	}

	if uid, err := ch.GetUserID(c); err == nil {
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
	} else {
		ch.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
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
	var page int
	if c.QueryParam("page") == "" {
		page = 1
	} else {
		pageatoi, err := strconv.Atoi(c.QueryParam("page"))
		if err != nil {
			ch.Logger.LogError(c, start, requestId, err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if pageatoi == 0 {
			page = 1
		} else {
			page = pageatoi
		}
	}

	if uid, err := ch.GetUserID(c); err == nil {
		messages, err := ch.UseCase.GetOneDialogue(uid, uint64(id), page) //Должна быть проверка на то, является ли чел собеседником
		//dialogues = ch.sanitizer.SanitizeEventCards(dialogues) SanitizeMessages(mes)
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
	} else {
		ch.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
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

	if uid, err := ch.GetUserID(c); err == nil {
		err := ch.UseCase.DeleteDialogue(uid, uint64(id)) //Должна быть проверка на то, является ли чел собеседником
		if err != nil {
			ch.Logger.LogError(c, start, requestId, err)
			return err
		}

		ch.Logger.LogInfo(c, start, requestId)
		return nil
	} else {
		ch.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
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

	if uid, err := ch.GetUserID(c); err == nil {
		newMessage.From = uid
		err := ch.UseCase.SendMessage(newMessage) //Должна быть проверка на то, является ли чел собеседником
		if err != nil {
			ch.Logger.LogError(c, start, requestId, err)
			return err
		}

		ch.Logger.LogInfo(c, start, requestId)
		return nil
	} else {
		ch.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
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

	if uid, err := ch.GetUserID(c); err == nil {
		err := ch.UseCase.DeleteMessage(uid, uint64(id)) //Должна быть проверка на то, является ли чел собеседником
		if err != nil {
			ch.Logger.LogError(c, start, requestId, err)
			return err
		}

		ch.Logger.LogInfo(c, start, requestId)
		return nil
	} else {
		ch.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
}

func (ch ChatHandler) EditMessage(c echo.Context) error {
	defer c.Request().Body.Close()

	start := time.Now()
	requestId := fmt.Sprintf("%016x", rand.Int())
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ch.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	newMessage := &models.NewMessage{}

	//newMessage должно быть отправлено с полем from!!!
	if err := easyjson.UnmarshalFromReader(c.Request().Body, newMessage); err != nil {
		ch.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusTeapot, err.Error())
	}

	//Здесь не опрокидываем uid в качестве from поля, чтобы не было возможности редактировать чужие сообщения(проверка на соответсвие uid и from)
	if uid, err := ch.GetUserID(c); err == nil {
		err := ch.UseCase.EditMessage(uid, id, newMessage) //Должна быть проверка на то, является ли чел собеседником
		if err != nil {
			ch.Logger.LogError(c, start, requestId, err)
			return err
		}

		ch.Logger.LogInfo(c, start, requestId)
		return nil
	} else {
		ch.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
}

func (ch ChatHandler) Search(c echo.Context) error {
	defer c.Request().Body.Close()

	start := time.Now()
	requestId := fmt.Sprintf("%016x", rand.Int())
	str := c.QueryParam("find")
	var page int
	if c.QueryParam("page") == "" {
		page = 1
	} else {
		pageatoi, err := strconv.Atoi(c.QueryParam("page"))
		if err != nil {
			ch.Logger.LogError(c, start, requestId, err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if pageatoi == 0 {
			page = 1
		} else {
			page = pageatoi
		}
	}

	var id int
	if c.QueryParam("id") == "" {
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

	if uid, err := ch.GetUserID(c); err == nil {
		err := ch.UseCase.Search(uid, id, str, page) //Должна быть проверка на то, является ли чел собеседником

		if err != nil {
			ch.Logger.LogError(c, start, requestId, err)
			return err
		}

		ch.Logger.LogInfo(c, start, requestId)
		return nil
	} else {
		ch.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
}
