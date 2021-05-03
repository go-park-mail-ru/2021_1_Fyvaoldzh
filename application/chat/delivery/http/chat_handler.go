package http

import (
	"errors"
	"fmt"
	"kudago/application/chat"
	"kudago/application/microservices/auth/client"
	"kudago/application/models"
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
	UseCase   chat.UseCase
	rpcAuth   client.AuthClient
	Logger    logger.Logger
	sanitizer *custom_sanitizer.CustomSanitizer
}

func CreateChatHandler(e *echo.Echo, uc chat.UseCase, rpcA client.AuthClient, sz *custom_sanitizer.CustomSanitizer, logger logger.Logger) {

	chatHandler := ChatHandler{UseCase: uc, rpcAuth: rpcA, Logger: logger, sanitizer: sz}

	//TODO групповой чат
	e.GET("/api/v1/dialogues", chatHandler.GetDialogues)
	e.GET("/api/v1/dialogues/:id", chatHandler.GetOneDialogue)
	e.DELETE("/api/v1/dialogues/:id", chatHandler.DeleteDialogue)
	e.POST("/api/v1/send", chatHandler.SendMessage)
	e.DELETE("/api/v1/message/:id", chatHandler.DeleteMessage)
	e.POST("/api/v1/message/:id", chatHandler.EditMessage)
	e.GET("/api/v1/dialogues/search", chatHandler.Search)
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
		exists, uid, err = ch.rpcAuth.Check(cookie.Value)
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

//Есть ли возможность как-то шаблонизировать функции ниже, везде одно и то же начало, мб создать функцию, принимающую функцию
//и в зависимости от метода кидать туда свой метод юзкейса?

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
		messages, err := ch.UseCase.GetOneDialogue(uid, uint64(id), page)
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
		err := ch.UseCase.DeleteDialogue(uid, uint64(id))
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
		err := ch.UseCase.SendMessage(newMessage, uid)
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
		err := ch.UseCase.DeleteMessage(uid, uint64(id))
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
	redactMessage := &models.RedactMessage{}

	if err := easyjson.UnmarshalFromReader(c.Request().Body, redactMessage); err != nil {
		ch.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusTeapot, err.Error())
	}

	if uid, err := ch.GetUserID(c); err == nil {
		err := ch.UseCase.EditMessage(uid, redactMessage)
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

//Пришла в голову идея не делать поиск по сообщениям и фолловерам, если id диалога не передан, а фронту кидать запрос сначала на поиск
//по фолловерам(у Насти уже реализовано, кажется), а затем по сообщениям и отрисовывать, как им удобно, пойдет?
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

	if uid, err := ch.GetUserID(c); err == nil {
		messages, err := ch.UseCase.Search(uid, id, str, page)

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
