package http

import (
	"github.com/labstack/echo"
	"kudago/application/subscriptions"
)

type SubscriptionHandler struct {
	UseCase subscriptions.UseCase
}

func CreateSubscriptionsHandler(e *echo.Echo, uc subscriptions.UseCase){

	subscriptionsHandler := SubscriptionHandler{UseCase: uc}

	e.POST("/api/v1/subscribe/user/:id", subscriptionsHandler.Subscribe)
	e.DELETE("/api/v1/subscribe/user/:id", subscriptionsHandler.Unsubscribe)
	e.POST("/api/v1/subscribe/planning/:id", subscriptionsHandler.AddPlanningEvent)
	e.DELETE("/api/v1/subscribe/planning/:id", subscriptionsHandler.RemovePlanningEvent)
	e.POST("/api/v1/subscribe/visited/:id", subscriptionsHandler.AddVisitedEvent)
	e.DELETE("/api/v1/subscribe/visited/:id", subscriptionsHandler.RemoveVisitedEvent)
}

func (h SubscriptionHandler) Subscribe(c echo.Context) error {
	// TODO: проверка на куку
	/*


	uid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	h.UseCase.SubscribeUser()

	 */

	return nil
}

func (h SubscriptionHandler) Unsubscribe(c echo.Context) error {


	return nil
}

func (h SubscriptionHandler) AddPlanningEvent(c echo.Context) error {


	return nil
}

func (h SubscriptionHandler) RemovePlanningEvent(c echo.Context) error {


	return nil
}

func (h SubscriptionHandler) AddVisitedEvent(c echo.Context) error {


	return nil
}

func (h SubscriptionHandler) RemoveVisitedEvent(c echo.Context) error {


	return nil
}
