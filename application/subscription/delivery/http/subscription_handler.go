package http

import (
	"github.com/labstack/echo"
	"github.com/mailru/easyjson"
	clientAuth "kudago/application/microservices/auth/client"
	clientSub "kudago/application/microservices/subscription/client"
	"kudago/application/models"
	"kudago/application/server/middleware"
	"kudago/application/subscription"
	"kudago/pkg/constants"
	"kudago/pkg/custom_sanitizer"
	"kudago/pkg/logger"
	"net/http"
)

type SubscriptionHandler struct {
	rpcAuth   clientAuth.AuthClient
	rpcSub    clientSub.SubscriptionClient
	usecase   subscription.UseCase
	sanitizer *custom_sanitizer.CustomSanitizer
	Logger    logger.Logger
	auth      middleware.Auth
}

func CreateSubscriptionsHandler(e *echo.Echo,
	rpcA clientAuth.AuthClient,
	rpcS clientSub.SubscriptionClient,
	uc subscription.UseCase,
	sz *custom_sanitizer.CustomSanitizer,
	logger logger.Logger,
	a middleware.Auth) {
	subscriptionHandler := SubscriptionHandler{
		rpcAuth:   rpcA,
		rpcSub:    rpcS,
		usecase:   uc,
		sanitizer: sz,
		Logger:    logger,
		auth:      a}

	e.POST("/api/v1/add/planning/:id",
		subscriptionHandler.AddPlanningEvent,
		a.GetSession,
		middleware.GetId)
	e.POST("/api/v1/add/visited/:id",
		subscriptionHandler.AddVisitedEvent,
		a.GetSession,
		middleware.GetId)
	e.DELETE("/api/v1/remove/:id",
		subscriptionHandler.RemoveEvent,
		a.GetSession,
		middleware.GetId)
	e.POST("/api/v1/subscribe/user/:id",
		subscriptionHandler.Subscribe,
		a.GetSession,
		middleware.GetId)
	e.DELETE("/api/v1/unsubscribe/user/:id",
		subscriptionHandler.Unsubscribe,
		a.GetSession,
		middleware.GetId)
	e.GET("/api/v1/followers/:id",
		subscriptionHandler.GetFollowers,
		middleware.GetPage,
		middleware.GetId)
	e.GET("/api/v1/subscriptions/:id",
		subscriptionHandler.GetSubscriptions,
		middleware.GetPage,
		middleware.GetId)
	e.GET("/api/v1/event/is_added/:id",
		subscriptionHandler.IsAdded,
		a.GetSession,
		middleware.GetId)

	e.GET("/api/v1/get/planning/:id",
		subscriptionHandler.GetPlanningEvents,
		middleware.GetPage,
		middleware.GetId)
	e.GET("/api/v1/get/visited/:id",
		subscriptionHandler.GetVisitedEvents,
		middleware.GetPage,
		middleware.GetId)
}

func (sh SubscriptionHandler) Subscribe(c echo.Context) error {
	defer c.Request().Body.Close()

	subscribedToId := c.Get(constants.IdKey).(int)
	subscriberId := c.Get(constants.UserIdKey).(uint64)

	if subscriberId == uint64(subscribedToId) {
		return echo.NewHTTPError(http.StatusBadRequest, "incorrect data")
	}

	return sh.rpcSub.Subscribe(subscriberId, uint64(subscribedToId))
}

func (sh SubscriptionHandler) Unsubscribe(c echo.Context) error {
	defer c.Request().Body.Close()

	subscribedToId := c.Get(constants.IdKey).(int)
	subscriberId := c.Get(constants.UserIdKey).(uint64)
	if subscriberId == uint64(subscribedToId) {
		return echo.NewHTTPError(http.StatusBadRequest, "incorrect data")
	}

	return sh.rpcSub.Unsubscribe(subscriberId, uint64(subscribedToId))
}

func (sh SubscriptionHandler) AddPlanningEvent(c echo.Context) error {
	defer c.Request().Body.Close()

	userId := c.Get(constants.UserIdKey).(uint64)
	eventId := c.Get(constants.IdKey).(int)

	return sh.rpcSub.AddPlanningEvent(userId, uint64(eventId))
}

func (sh SubscriptionHandler) AddVisitedEvent(c echo.Context) error {
	defer c.Request().Body.Close()

	userId := c.Get(constants.UserIdKey).(uint64)
	eventId := c.Get(constants.IdKey).(int)

	return sh.rpcSub.AddVisitedEvent(userId, uint64(eventId))
}

func (sh SubscriptionHandler) RemoveEvent(c echo.Context) error {
	defer c.Request().Body.Close()

	userId := c.Get(constants.UserIdKey).(uint64)
	eventId := c.Get(constants.IdKey).(int)

	return sh.rpcSub.RemoveEvent(userId, uint64(eventId))
}

func (sh SubscriptionHandler) GetFollowers(c echo.Context) error {
	defer c.Request().Body.Close()

	userId := c.Get(constants.IdKey).(int)
	page := c.Get(constants.PageKey).(int)

	users, err := sh.usecase.GetFollowers(uint64(userId), page)
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

	userId := c.Get(constants.IdKey).(int)
	page := c.Get(constants.PageKey).(int)

	users, err := sh.usecase.GetSubscriptions(uint64(userId), page)
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

	var err error
	userId := c.Get(constants.UserIdKey).(uint64)
	eventId := c.Get(constants.IdKey).(int)

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

	userId := c.Get(constants.IdKey).(int)
	page := c.Get(constants.PageKey).(int)

	events, err := sh.usecase.GetPlanningEvents(uint64(userId), page)
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

	userId := c.Get(constants.IdKey).(int)
	page := c.Get(constants.PageKey).(int)

	events, err := sh.usecase.GetVisitedEvents(uint64(userId), page)
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
