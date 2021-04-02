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

	subscriptionsHandler := SubscriptionHandler{UseCase: uc, Sm: sm}

	e.POST("/api/v1/subscribe/user/:id", subscriptionsHandler.Subscribe)
	e.DELETE("/api/v1/unsubscribe/user/:id", subscriptionsHandler.Unsubscribe)
	e.POST("/api/v1/add/planning/:id", subscriptionsHandler.AddPlanningEvent)
	e.DELETE("/api/v1/remove/planning/:id", subscriptionsHandler.RemovePlanningEvent)
	e.POST("/api/v1/add/visited/:id", subscriptionsHandler.AddVisitedEvent)
	e.DELETE("/api/v1/remove/visited/:id", subscriptionsHandler.RemoveVisitedEvent)
}

// TODO: общую часть с кукой и проверкой вынести в общий метод

func (h SubscriptionHandler) Subscribe(c echo.Context) error {
	defer c.Request().Body.Close()

	cookie, err := c.Cookie(constants.SessionCookieName)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	var uid1 uint64
	var exists bool
	exists, uid1, err = h.Sm.CheckSession(cookie.Value)
	if err != nil {
		return err
	}

	if !exists {
		return echo.NewHTTPError(http.StatusBadRequest, "user is not authorized")
	}

	uid2, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = h.UseCase.SubscribeUser(uid1, uint64(uid2))

	if err != nil {
		return err
	}
	return nil
}

func (h SubscriptionHandler) Unsubscribe(c echo.Context) error {
	defer c.Request().Body.Close()

	cookie, err := c.Cookie(constants.SessionCookieName)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	var uid1 uint64
	var exists bool
	exists, uid1, err = h.Sm.CheckSession(cookie.Value)
	if err != nil {
		return err
	}

	if !exists {
		return echo.NewHTTPError(http.StatusBadRequest, "user is not authorized")
	}

	uid2, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return h.UseCase.UnsubscribeUser(uid1, uint64(uid2))
}

func (h SubscriptionHandler) AddPlanningEvent(c echo.Context) error {
	defer c.Request().Body.Close()

	cookie, err := c.Cookie(constants.SessionCookieName)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	var uid uint64
	var exists bool
	exists, uid, err = h.Sm.CheckSession(cookie.Value)
	if err != nil {
		return err
	}

	if !exists {
		return echo.NewHTTPError(http.StatusBadRequest, "user is not authorized")
	}

	eid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return h.UseCase.AddPlanning(uid, uint64(eid))
}

func (h SubscriptionHandler) RemovePlanningEvent(c echo.Context) error {
	defer c.Request().Body.Close()

	cookie, err := c.Cookie(constants.SessionCookieName)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	var uid uint64
	var exists bool
	exists, uid, err = h.Sm.CheckSession(cookie.Value)
	if err != nil {
		return err
	}

	if !exists {
		return echo.NewHTTPError(http.StatusBadRequest, "user is not authorized")
	}

	eid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return h.UseCase.RemovePlanning(uid, uint64(eid))
}

func (h SubscriptionHandler) AddVisitedEvent(c echo.Context) error {
	defer c.Request().Body.Close()

	cookie, err := c.Cookie(constants.SessionCookieName)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	var uid uint64
	var exists bool
	exists, uid, err = h.Sm.CheckSession(cookie.Value)
	if err != nil {
		return err
	}

	if !exists {
		return echo.NewHTTPError(http.StatusBadRequest, "user is not authorized")
	}

	eid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return h.UseCase.AddVisited(uid, uint64(eid))
}

func (h SubscriptionHandler) RemoveVisitedEvent(c echo.Context) error {
	defer c.Request().Body.Close()

	cookie, err := c.Cookie(constants.SessionCookieName)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	var uid uint64
	var exists bool
	exists, uid, err = h.Sm.CheckSession(cookie.Value)
	if err != nil {
		return err
	}

	if !exists {
		return echo.NewHTTPError(http.StatusBadRequest, "user is not authorized")
	}

	eid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return h.UseCase.RemoveVisited(uid, uint64(eid))
}
