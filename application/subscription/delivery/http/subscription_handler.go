package http

import (
	"github.com/labstack/echo"
	"kudago/application/subscription"
	"kudago/pkg/constants"
	"kudago/pkg/infrastructure"
	"net/http"
	"strconv"
)

type SubscriptionHandler struct {
	UseCase subscription.UseCase
	Sm      *infrastructure.SessionManager
}

func CreateSubscriptionsHandler(e *echo.Echo, uc subscription.UseCase, sm *infrastructure.SessionManager) {

	subscriptionHandler := SubscriptionHandler{UseCase: uc, Sm: sm}

	e.POST("/api/v1/subscribe/user/:id", subscriptionHandler.Subscribe)
	e.DELETE("/api/v1/unsubscribe/user/:id", subscriptionHandler.Unsubscribe)
	e.POST("/api/v1/add/planning/:id", subscriptionHandler.AddPlanningEvent)
	e.DELETE("/api/v1/remove/planning/:id", subscriptionHandler.RemovePlanningEvent)
	e.POST("/api/v1/add/visited/:id", subscriptionHandler.AddVisitedEvent)
	e.DELETE("/api/v1/remove/visited/:id", subscriptionHandler.RemoveVisitedEvent)
}

func (h SubscriptionHandler) Subscribe(c echo.Context) error {
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
