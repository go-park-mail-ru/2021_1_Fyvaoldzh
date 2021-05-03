package http

import (
	"errors"
	"github.com/mailru/easyjson"
	clientAuth "kudago/application/microservices/auth/client"
	clientSub "kudago/application/microservices/subscription/client"
	"kudago/application/models"
	"kudago/application/subscription"
	"kudago/pkg/constants"
	"kudago/pkg/custom_sanitizer"
	"kudago/pkg/logger"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

type SubscriptionHandler struct {
	rpcAuth   clientAuth.AuthClient
	rpcSub    clientSub.SubscriptionClient
	usecase   subscription.UseCase
	sanitizer *custom_sanitizer.CustomSanitizer
	Logger    logger.Logger
}

func CreateSubscriptionsHandler(e *echo.Echo,
	rpcA clientAuth.AuthClient,
	rpcS clientSub.SubscriptionClient,
	uc subscription.UseCase,
	sz *custom_sanitizer.CustomSanitizer,
	logger logger.Logger) {
	subscriptionHandler := SubscriptionHandler{
		rpcAuth:   rpcA,
		rpcSub:    rpcS,
		usecase:   uc,
		sanitizer: sz,
		Logger:    logger}

	e.POST("/api/v1/add/planning/:id", subscriptionHandler.AddPlanningEvent)
	e.POST("/api/v1/add/visited/:id", subscriptionHandler.AddVisitedEvent)
	e.DELETE("/api/v1/remove/:id", subscriptionHandler.RemoveEvent)
	e.POST("/api/v1/subscribe/user/:id", subscriptionHandler.Subscribe)
	e.DELETE("/api/v1/unsubscribe/user/:id", subscriptionHandler.Unsubscribe)
	e.GET("/api/v1/followers/:id", subscriptionHandler.GetFollowers)
	e.GET("/api/v1/subscriptions/:id", subscriptionHandler.GetSubscriptions)
	e.GET("/api/v1/event/is_added/:id", subscriptionHandler.IsAdded)

	e.GET("/api/v1/get/planning/:id", subscriptionHandler.GetPlanningEvents)
	e.GET("/api/v1/get/visited/:id", subscriptionHandler.GetVisitedEvents)
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

func (sh SubscriptionHandler) GetFollowers(c echo.Context) error {
	defer c.Request().Body.Close()

	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if userId <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "incorrect data")
	}

	users, err := sh.usecase.GetFollowers(uint64(userId))
	if err != nil {
		return err
	}

	users = sh.sanitizer.SanitizeUserCards(users)
	if _, err = easyjson.MarshalToWriter(users, c.Response().Writer); err != nil {
		sh.Logger.Warn(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return err
}

func (sh SubscriptionHandler) GetSubscriptions(c echo.Context) error {
	defer c.Request().Body.Close()

	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if userId <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "incorrect data")
	}

	users, err := sh.usecase.GetSubscriptions(uint64(userId))
	if err != nil {
		return err
	}

	users = sh.sanitizer.SanitizeUserCards(users)
	if _, err = easyjson.MarshalToWriter(users, c.Response().Writer); err != nil {
		sh.Logger.Warn(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return err
}

func (sh *SubscriptionHandler) IsAdded(c echo.Context) error {
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

	var answer models.IsAddedEvent
	answer.IsAdded, err = sh.usecase.IsAddedEvent(userId, uint64(eventId))
	if err != nil {
		sh.Logger.Warn(err)
		return err
	}

	if _, err = easyjson.MarshalToWriter(answer, c.Response().Writer); err != nil {
		sh.Logger.Warn(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (sh SubscriptionHandler) GetPlanningEvents(c echo.Context) error {
	defer c.Request().Body.Close()

	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if userId <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "incorrect data")
	}

	events, err := sh.usecase.GetPlanningEvents(uint64(userId))
	if err != nil {
		return err
	}

	events = sh.sanitizer.SanitizeEventCards(events)
	if _, err = easyjson.MarshalToWriter(events, c.Response().Writer); err != nil {
		sh.Logger.Warn(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (sh SubscriptionHandler) GetVisitedEvents(c echo.Context) error {
	defer c.Request().Body.Close()

	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if userId <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "incorrect data")
	}

	events, err := sh.usecase.GetVisitedEvents(uint64(userId))
	if err != nil {
		return err
	}

	events = sh.sanitizer.SanitizeEventCards(events)
	if _, err = easyjson.MarshalToWriter(events, c.Response().Writer); err != nil {
		sh.Logger.Warn(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return nil
}
