package http

import (
	"fmt"
	"kudago/application/models"
	"kudago/application/subscription"
	"kudago/pkg/constants"
	"kudago/pkg/infrastructure"
	"kudago/pkg/logger"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/mailru/easyjson"

	"github.com/labstack/echo"
)

type SubscriptionHandler struct {
	UseCase subscription.UseCase
	Sm      infrastructure.SessionTarantool
	Logger  logger.Logger
}

func CreateSubscriptionsHandler(e *echo.Echo, uc subscription.UseCase, sm infrastructure.SessionTarantool, logger logger.Logger) {

	subscriptionHandler := SubscriptionHandler{UseCase: uc, Sm: sm, Logger: logger}

	e.POST("/api/v1/subscribe/user/:id", subscriptionHandler.Subscribe)
	e.DELETE("/api/v1/unsubscribe/user/:id", subscriptionHandler.Unsubscribe)
	e.POST("/api/v1/add/planning/:id", subscriptionHandler.AddPlanningEvent)
	e.DELETE("/api/v1/remove/planning/:id", subscriptionHandler.RemovePlanningEvent)
	e.POST("/api/v1/add/visited/:id", subscriptionHandler.AddVisitedEvent)
	e.DELETE("/api/v1/remove/visited/:id", subscriptionHandler.RemoveVisitedEvent)
	e.GET("/api/v1/event/is_added/:id", subscriptionHandler.IsAdded)
}

func (h SubscriptionHandler) Subscribe(c echo.Context) error {
	defer c.Request().Body.Close()
	start := time.Now()
	requestId := fmt.Sprintf("%016x", rand.Int())
	cookie, err := c.Cookie(constants.SessionCookieName)
	if err != nil {
		h.Logger.LogWarn(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}
	log.Println(cookie.Value)

	var subscriberId uint64
	var exists bool
	exists, subscriberId, err = h.Sm.CheckSession(cookie.Value)
	if err != nil {
		return err
	}
	if !exists {
		return echo.NewHTTPError(http.StatusBadRequest, "user is not authorized")
	}

	subscribedToId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if subscribedToId <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "incorrect data")
	}

	return h.UseCase.SubscribeUser(subscriberId, uint64(subscribedToId))
}

func (h SubscriptionHandler) Unsubscribe(c echo.Context) error {
	defer c.Request().Body.Close()

	cookie, err := c.Cookie(constants.SessionCookieName)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	var subscriberId uint64
	var exists bool
	exists, subscriberId, err = h.Sm.CheckSession(cookie.Value)
	if err != nil {
		return err
	}

	if !exists {
		return echo.NewHTTPError(http.StatusBadRequest, "user is not authorized")
	}

	subscribedToId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if subscribedToId <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "incorrect data")
	}

	return h.UseCase.UnsubscribeUser(subscriberId, uint64(subscribedToId))
}

func (h SubscriptionHandler) AddPlanningEvent(c echo.Context) error {
	defer c.Request().Body.Close()

	cookie, err := c.Cookie(constants.SessionCookieName)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	var userId uint64
	var exists bool
	exists, userId, err = h.Sm.CheckSession(cookie.Value)
	if err != nil {
		return err
	}
	if !exists {
		return echo.NewHTTPError(http.StatusBadRequest, "user is not authorized")
	}

	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if eventId <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "incorrect data")
	}

	return h.UseCase.AddPlanning(userId, uint64(eventId))
}

func (h SubscriptionHandler) RemovePlanningEvent(c echo.Context) error {
	defer c.Request().Body.Close()

	cookie, err := c.Cookie(constants.SessionCookieName)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	var userId uint64
	var exists bool
	exists, userId, err = h.Sm.CheckSession(cookie.Value)
	if err != nil {
		return err
	}
	if !exists {
		return echo.NewHTTPError(http.StatusBadRequest, "user is not authorized")
	}

	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if eventId <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "incorrect data")
	}

	return h.UseCase.RemovePlanning(userId, uint64(eventId))
}

func (h SubscriptionHandler) AddVisitedEvent(c echo.Context) error {
	defer c.Request().Body.Close()

	cookie, err := c.Cookie(constants.SessionCookieName)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	var userId uint64
	var exists bool
	exists, userId, err = h.Sm.CheckSession(cookie.Value)
	if err != nil {
		return err
	}
	if !exists {
		return echo.NewHTTPError(http.StatusBadRequest, "user is not authorized")
	}

	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if eventId <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "incorrect data")
	}

	return h.UseCase.AddVisited(userId, uint64(eventId))
}

func (h SubscriptionHandler) RemoveVisitedEvent(c echo.Context) error {
	defer c.Request().Body.Close()

	cookie, err := c.Cookie(constants.SessionCookieName)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	var userId uint64
	var exists bool
	exists, userId, err = h.Sm.CheckSession(cookie.Value)
	if err != nil {
		return err
	}
	if !exists {
		return echo.NewHTTPError(http.StatusBadRequest, "user is not authorized")
	}

	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if eventId <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "incorrect data")
	}

	return h.UseCase.RemoveVisited(userId, uint64(eventId))
}

func (h SubscriptionHandler) IsAdded(c echo.Context) error {
	defer c.Request().Body.Close()

	cookie, err := c.Cookie(constants.SessionCookieName)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	var userId uint64
	var exists bool
	exists, userId, err = h.Sm.CheckSession(cookie.Value)
	if err != nil {
		return err
	}
	if !exists {
		return echo.NewHTTPError(http.StatusBadRequest, "user is not authorized")
	}

	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if eventId <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "incorrect data")
	}

	var answer models.IsAddedEvent
	answer.EventId = uint64(eventId)
	answer.UserId = userId
	answer.IsAdded, err = h.UseCase.IsAddedEvent(userId, uint64(eventId))
	if err != nil {
		log.Println(err)
	}

	if _, err = easyjson.MarshalToWriter(answer, c.Response().Writer); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}
