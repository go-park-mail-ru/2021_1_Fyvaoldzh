package http

import (
	"errors"
	clientAuth "kudago/application/microservices/auth/client"
	clientSub "kudago/application/microservices/subscription/client"
	"kudago/pkg/constants"
	"kudago/pkg/logger"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

type SubscriptionHandler struct {
	rpcAuth clientAuth.AuthClient
	rpcSub  clientSub.SubscriptionClient
	Logger  logger.Logger
}

func CreateSubscriptionsHandler(e *echo.Echo, rpcA clientAuth.AuthClient, rpcS clientSub.SubscriptionClient, logger logger.Logger) {
	subscriptionHandler := SubscriptionHandler{
		rpcAuth: rpcA,
		rpcSub:  rpcS,
		Logger:  logger}

	e.POST("/api/v1/add/planning/:id", subscriptionHandler.AddPlanningEvent)
	e.POST("/api/v1/add/visited/:id", subscriptionHandler.AddVisitedEvent)
	e.DELETE("/api/v1/remove/:id", subscriptionHandler.RemoveEvent)
	e.POST("/api/v1/subscribe/user/:id", subscriptionHandler.Subscribe)
	e.DELETE("/api/v1/unsubscribe/user/:id", subscriptionHandler.Unsubscribe)
}



func (sh SubscriptionHandler) Subscribe(c echo.Context) error {
	defer c.Request().Body.Close()
	cookie, err := c.Cookie(constants.SessionCookieName)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	var subscriberId uint64
	var exists bool
	exists, subscriberId, err = sh.rpcAuth.Check(cookie.Value)
	if err != nil {
		sh.Logger.Warn(err)
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
	if subscriberId == uint64(subscribedToId) {
		return echo.NewHTTPError(http.StatusBadRequest, "incorrect data")
	}

	return sh.rpcSub.Subscribe(subscriberId, uint64(subscribedToId))
}

func (sh SubscriptionHandler) Unsubscribe(c echo.Context) error {
	defer c.Request().Body.Close()

	cookie, err := c.Cookie(constants.SessionCookieName)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	var subscriberId uint64
	var exists bool
	exists, subscriberId, err = sh.rpcAuth.Check(cookie.Value)
	if err != nil {
		sh.Logger.Warn(err)
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
	if subscriberId == uint64(subscribedToId) {
		sh.Logger.Warn(errors.New("subscriberId == subscribedToId"))
		return echo.NewHTTPError(http.StatusBadRequest, "incorrect data")
	}

	return sh.rpcSub.Unsubscribe(subscriberId, uint64(subscribedToId))
}

func (sh SubscriptionHandler) AddPlanningEvent(c echo.Context) error {
	defer c.Request().Body.Close()

	cookie, err := c.Cookie(constants.SessionCookieName)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	var userId uint64
	var exists bool
	exists, userId, err = sh.rpcAuth.Check(cookie.Value)
	if err != nil {
		sh.Logger.Warn(err)
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

	return sh.rpcSub.AddPlanningEvent(userId, uint64(eventId))
}

func (sh SubscriptionHandler) AddVisitedEvent(c echo.Context) error {
	defer c.Request().Body.Close()

	cookie, err := c.Cookie(constants.SessionCookieName)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	var userId uint64
	var exists bool
	exists, userId, err = sh.rpcAuth.Check(cookie.Value)
	if err != nil {
		sh.Logger.Warn(err)
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

	return sh.rpcSub.AddVisitedEvent(userId, uint64(eventId))
}

func (sh SubscriptionHandler) RemoveEvent(c echo.Context) error {
	defer c.Request().Body.Close()

	cookie, err := c.Cookie(constants.SessionCookieName)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	var userId uint64
	var exists bool
	exists, userId, err = sh.rpcAuth.Check(cookie.Value)
	if err != nil {
		sh.Logger.Warn(err)
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

	return sh.rpcSub.RemoveEvent(userId, uint64(eventId))
}



