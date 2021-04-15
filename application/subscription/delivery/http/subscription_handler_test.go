package http

import (
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	mock_subscription "kudago/application/subscription/mocks"
	"kudago/pkg/constants"
	"kudago/pkg/generator"

	mock_infrastructure "kudago/pkg/infrastructure/mocks"
	"kudago/pkg/logger"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	id    uint64 = 1
	strId        = "1"
)

func setUp(t *testing.T, url, method string) (echo.Context,
	SubscriptionHandler, *mock_subscription.MockUseCase, *mock_infrastructure.MockSessionTarantool) {
	e := echo.New()
	r := e.Router()
	r.Add(method, url, func(echo.Context) error { return nil })

	ctrl := gomock.NewController(t)
	usecase := mock_subscription.NewMockUseCase(ctrl)
	sm := mock_infrastructure.NewMockSessionTarantool(ctrl)

	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	sugar := l.Sugar()
	zap.NewAtomicLevelAt(zapcore.DebugLevel)

	handler := SubscriptionHandler{
		UseCase: usecase,
		Sm:      sm,
		Logger:  logger.NewLogger(sugar),
	}

	var req *http.Request
	req = httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(url)

	return c, handler, usecase, sm
}

///////////////////////////////////////////////////

func TestSubscriptionHandler_Subscribe(t *testing.T) {
	c, h, usecase, sm := setUp(t, "/api/v1/subscribe/user/:id", http.MethodPost)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	sm.EXPECT().CheckSession(cookie.Value).Return(true, id, nil)
	usecase.EXPECT().SubscribeUser(id, id).Return(nil)

	err := h.Subscribe(c)

	assert.Nil(t, err)
}

func TestSubscriptionHandler_SubscribeErrorMinus(t *testing.T) {
	c, h, _, sm := setUp(t, "/api/v1/subscribe/user/:id", http.MethodPost)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)
	c.SetParamNames("id")
	c.SetParamValues("-1")

	sm.EXPECT().CheckSession(cookie.Value).Return(true, id, nil)

	err := h.Subscribe(c)

	assert.Error(t, err)
}

func TestSubscriptionHandler_SubscribeErrorAtoi(t *testing.T) {
	c, h, _, sm := setUp(t, "/api/v1/subscribe/user/:id", http.MethodPost)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)
	c.SetParamNames("id")
	c.SetParamValues("a")

	sm.EXPECT().CheckSession(cookie.Value).Return(true, id, nil)

	err := h.Subscribe(c)

	assert.Error(t, err)
}

func TestSubscriptionHandler_SubscribeNotAuthorized(t *testing.T) {
	c, h, _, sm := setUp(t, "/api/v1/subscribe/user/:id", http.MethodPost)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	sm.EXPECT().CheckSession(cookie.Value).Return(false, id, nil)

	err := h.Subscribe(c)

	assert.Error(t, err)
}

func TestSubscriptionHandler_SubscribeErrorSM(t *testing.T) {
	c, h, _, sm := setUp(t, "/api/v1/subscribe/user/:id", http.MethodPost)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	sm.EXPECT().CheckSession(cookie.Value).Return(false, id, echo.NewHTTPError(http.StatusInternalServerError))

	err := h.Subscribe(c)

	assert.Error(t, err)
}

func TestSubscriptionHandler_SubscribeNoCookie(t *testing.T) {
	c, h, _, _ := setUp(t, "/api/v1/subscribe/user/:id", http.MethodPost)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	err := h.Subscribe(c)

	assert.Error(t, err)
}

///////////////////////////////////////////////////

func TestSubscriptionHandler_Unsubscribe(t *testing.T) {
	c, h, usecase, sm := setUp(t, "/api/v1/unsubscribe/user/:id", http.MethodDelete)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	sm.EXPECT().CheckSession(cookie.Value).Return(true, id, nil)
	usecase.EXPECT().UnsubscribeUser(id, id).Return(nil)

	err := h.Unsubscribe(c)

	assert.Nil(t, err)
}

func TestSubscriptionHandler_UnsubscribeErrorMinus(t *testing.T) {
	c, h, _, sm := setUp(t, "/api/v1/unsubscribe/user/:id", http.MethodDelete)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)
	c.SetParamNames("id")
	c.SetParamValues("-1")

	sm.EXPECT().CheckSession(cookie.Value).Return(true, id, nil)

	err := h.Unsubscribe(c)

	assert.Error(t, err)
}

func TestSubscriptionHandler_UnsubscribeErrorAtoi(t *testing.T) {
	c, h, _, sm := setUp(t, "/api/v1/unsubscribe/user/:id", http.MethodDelete)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)
	c.SetParamNames("id")
	c.SetParamValues("a")

	sm.EXPECT().CheckSession(cookie.Value).Return(true, id, nil)

	err := h.Unsubscribe(c)

	assert.Error(t, err)
}

func TestSubscriptionHandler_UnsubscribeNotAuthorized(t *testing.T) {
	c, h, _, sm := setUp(t, "/api/v1/unsubscribe/user/:id", http.MethodDelete)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	sm.EXPECT().CheckSession(cookie.Value).Return(false, id, nil)

	err := h.Unsubscribe(c)

	assert.Error(t, err)
}

func TestSubscriptionHandler_UnsubscribeErrorSM(t *testing.T) {
	c, h, _, sm := setUp(t, "/api/v1/unsubscribe/user/:id", http.MethodDelete)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	sm.EXPECT().CheckSession(cookie.Value).Return(false, id, echo.NewHTTPError(http.StatusInternalServerError))

	err := h.Unsubscribe(c)

	assert.Error(t, err)
}

func TestSubscriptionHandler_UnsubscribeNoCookie(t *testing.T) {
	c, h, _, _ := setUp(t, "/api/v1/unsubscribe/user/:id", http.MethodDelete)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	err := h.Unsubscribe(c)

	assert.Error(t, err)
}

///////////////////////////////////////////////////

func TestSubscriptionHandler_AddPlanningEvent(t *testing.T) {
	c, h, usecase, sm := setUp(t, "/api/v1/add/planning/:id", http.MethodPost)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	sm.EXPECT().CheckSession(cookie.Value).Return(true, id, nil)
	usecase.EXPECT().AddPlanning(id, id).Return(nil)

	err := h.AddPlanningEvent(c)

	assert.Nil(t, err)
}

func TestSubscriptionHandler_AddPlanningEventErrorMinus(t *testing.T) {
	c, h, _, sm := setUp(t, "/api/v1/add/planning/:id", http.MethodPost)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)
	c.SetParamNames("id")
	c.SetParamValues("-1")

	sm.EXPECT().CheckSession(cookie.Value).Return(true, id, nil)

	err := h.AddPlanningEvent(c)

	assert.Error(t, err)
}

func TestSubscriptionHandler_AddPlanningEventErrorAtoi(t *testing.T) {
	c, h, _, sm := setUp(t, "/api/v1/add/planning/:id", http.MethodPost)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)
	c.SetParamNames("id")
	c.SetParamValues("a")

	sm.EXPECT().CheckSession(cookie.Value).Return(true, id, nil)

	err := h.AddPlanningEvent(c)

	assert.Error(t, err)
}

func TestSubscriptionHandler_AddPlanningEventNotAuthorized(t *testing.T) {
	c, h, _, sm := setUp(t, "/api/v1/add/planning/:id", http.MethodPost)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	sm.EXPECT().CheckSession(cookie.Value).Return(false, id, nil)

	err := h.AddPlanningEvent(c)

	assert.Error(t, err)
}

func TestSubscriptionHandler_AddPlanningEventErrorSM(t *testing.T) {
	c, h, _, sm := setUp(t, "/api/v1/add/planning/:id", http.MethodPost)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	sm.EXPECT().CheckSession(cookie.Value).Return(false, id, echo.NewHTTPError(http.StatusInternalServerError))

	err := h.AddPlanningEvent(c)

	assert.Error(t, err)
}

func TestSubscriptionHandler_AddPlanningEventNoCookie(t *testing.T) {
	c, h, _, _ := setUp(t, "/api/v1/add/planning/:id", http.MethodPost)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	err := h.AddPlanningEvent(c)

	assert.Error(t, err)
}

///////////////////////////////////////////////////

func TestSubscriptionHandler_RemovePlanningEvent(t *testing.T) {
	c, h, usecase, sm := setUp(t, "/api/v1/remove/planning/:id", http.MethodDelete)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	sm.EXPECT().CheckSession(cookie.Value).Return(true, id, nil)
	usecase.EXPECT().RemovePlanning(id, id).Return(nil)

	err := h.RemovePlanningEvent(c)

	assert.Nil(t, err)
}

func TestSubscriptionHandler_RemovePlanningEventErrorMinus(t *testing.T) {
	c, h, _, sm := setUp(t, "/api/v1/uremove/planning/:id", http.MethodDelete)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)
	c.SetParamNames("id")
	c.SetParamValues("-1")

	sm.EXPECT().CheckSession(cookie.Value).Return(true, id, nil)

	err := h.RemovePlanningEvent(c)

	assert.Error(t, err)
}

func TestSubscriptionHandler_RemovePlanningEventErrorAtoi(t *testing.T) {
	c, h, _, sm := setUp(t, "/api/v1/remove/planning/:id", http.MethodDelete)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)
	c.SetParamNames("id")
	c.SetParamValues("a")

	sm.EXPECT().CheckSession(cookie.Value).Return(true, id, nil)

	err := h.RemovePlanningEvent(c)

	assert.Error(t, err)
}

func TestSubscriptionHandler_RemovePlanningEventNotAuthorized(t *testing.T) {
	c, h, _, sm := setUp(t, "/api/v1/remove/planning/:id", http.MethodDelete)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	sm.EXPECT().CheckSession(cookie.Value).Return(false, id, nil)

	err := h.RemovePlanningEvent(c)

	assert.Error(t, err)
}

func TestSubscriptionHandler_RemovePlanningEventErrorSM(t *testing.T) {
	c, h, _, sm := setUp(t, "/api/v1/remove/planning/:id", http.MethodDelete)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	sm.EXPECT().CheckSession(cookie.Value).Return(false, id, echo.NewHTTPError(http.StatusInternalServerError))

	err := h.RemovePlanningEvent(c)

	assert.Error(t, err)
}

func TestSubscriptionHandler_RemovePlanningEventNoCookie(t *testing.T) {
	c, h, _, _ := setUp(t, "/api/v1/remove/planning/:id", http.MethodDelete)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	err := h.RemovePlanningEvent(c)

	assert.Error(t, err)
}

///////////////////////////////////////////////////

func TestSubscriptionHandler_AddVisitedEvent(t *testing.T) {
	c, h, usecase, sm := setUp(t, "/api/v1/add/visited/:id", http.MethodPost)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	sm.EXPECT().CheckSession(cookie.Value).Return(true, id, nil)
	usecase.EXPECT().AddVisited(id, id).Return(nil)

	err := h.AddVisitedEvent(c)

	assert.Nil(t, err)
}

func TestSubscriptionHandler_AddVisitedEventErrorMinus(t *testing.T) {
	c, h, _, sm := setUp(t, "/api/v1/add/visited/:id", http.MethodPost)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)
	c.SetParamNames("id")
	c.SetParamValues("-1")

	sm.EXPECT().CheckSession(cookie.Value).Return(true, id, nil)

	err := h.AddVisitedEvent(c)

	assert.Error(t, err)
}

func TestSubscriptionHandler_AddVisitedEventErrorAtoi(t *testing.T) {
	c, h, _, sm := setUp(t, "/api/v1/add/visited/:id", http.MethodPost)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)
	c.SetParamNames("id")
	c.SetParamValues("a")

	sm.EXPECT().CheckSession(cookie.Value).Return(true, id, nil)

	err := h.AddVisitedEvent(c)

	assert.Error(t, err)
}

func TestSubscriptionHandler_AddVisitedEventNotAuthorized(t *testing.T) {
	c, h, _, sm := setUp(t, "/api/v1/add/visited/:id", http.MethodPost)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	sm.EXPECT().CheckSession(cookie.Value).Return(false, id, nil)

	err := h.AddVisitedEvent(c)

	assert.Error(t, err)
}

func TestSubscriptionHandler_AddVisitedEventErrorSM(t *testing.T) {
	c, h, _, sm := setUp(t, "/api/v1/add/visited/:id", http.MethodPost)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	sm.EXPECT().CheckSession(cookie.Value).Return(false, id, echo.NewHTTPError(http.StatusInternalServerError))

	err := h.AddVisitedEvent(c)

	assert.Error(t, err)
}

func TestSubscriptionHandler_AddVisitedEventNoCookie(t *testing.T) {
	c, h, _, _ := setUp(t, "/api/v1/add/visited/:id", http.MethodPost)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	err := h.AddVisitedEvent(c)

	assert.Error(t, err)
}

///////////////////////////////////////////////////

func TestSubscriptionHandler_RemoveVisitedEvent(t *testing.T) {
	c, h, usecase, sm := setUp(t, "/api/v1/remove/visited/:id", http.MethodDelete)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	sm.EXPECT().CheckSession(cookie.Value).Return(true, id, nil)
	usecase.EXPECT().RemoveVisited(id, id).Return(nil)

	err := h.RemoveVisitedEvent(c)

	assert.Nil(t, err)
}

func TestSubscriptionHandler_RemoveVisitedEventErrorMinus(t *testing.T) {
	c, h, _, sm := setUp(t, "/api/v1/uremove/visited/:id", http.MethodDelete)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)
	c.SetParamNames("id")
	c.SetParamValues("-1")

	sm.EXPECT().CheckSession(cookie.Value).Return(true, id, nil)

	err := h.RemoveVisitedEvent(c)

	assert.Error(t, err)
}

func TestSubscriptionHandler_RemoveVisitedEventErrorAtoi(t *testing.T) {
	c, h, _, sm := setUp(t, "/api/v1/remove/visited/:id", http.MethodDelete)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)
	c.SetParamNames("id")
	c.SetParamValues("a")

	sm.EXPECT().CheckSession(cookie.Value).Return(true, id, nil)

	err := h.RemoveVisitedEvent(c)

	assert.Error(t, err)
}

func TestSubscriptionHandler_RemoveVisitedEventNotAuthorized(t *testing.T) {
	c, h, _, sm := setUp(t, "/api/v1/remove/visited/:id", http.MethodDelete)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	sm.EXPECT().CheckSession(cookie.Value).Return(false, id, nil)

	err := h.RemoveVisitedEvent(c)

	assert.Error(t, err)
}

func TestSubscriptionHandler_RemoveVisitedEventErrorSM(t *testing.T) {
	c, h, _, sm := setUp(t, "/api/v1/remove/visited/:id", http.MethodDelete)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	sm.EXPECT().CheckSession(cookie.Value).Return(false, id, echo.NewHTTPError(http.StatusInternalServerError))

	err := h.RemoveVisitedEvent(c)

	assert.Error(t, err)
}

func TestSubscriptionHandler_RemoveVisitedEventNoCookie(t *testing.T) {
	c, h, _, _ := setUp(t, "/api/v1/remove/visited/:id", http.MethodDelete)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	err := h.RemoveVisitedEvent(c)

	assert.Error(t, err)
}

///////////////////////////////////////////////////

func TestSubscriptionHandler_IsAdded(t *testing.T) {
	c, h, usecase, sm := setUp(t, "/api/v1/event/is_added/:id", http.MethodGet)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	sm.EXPECT().CheckSession(cookie.Value).Return(true, id, nil)
	usecase.EXPECT().IsAddedEvent(id, id).Return(true, nil)

	err := h.IsAdded(c)

	assert.Nil(t, err)
}

func TestSubscriptionHandler_IsAddedErrorUC(t *testing.T) {
	c, h, usecase, sm := setUp(t, "/api/v1/event/is_added/:id", http.MethodGet)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	sm.EXPECT().CheckSession(cookie.Value).Return(true, id, nil)
	usecase.EXPECT().IsAddedEvent(id, id).Return(true, echo.NewHTTPError(http.StatusInternalServerError))

	err := h.IsAdded(c)

	assert.Error(t, err)
}

func TestSubscriptionHandler_IsAddedErrorSM(t *testing.T) {
	c, h, _, sm := setUp(t, "/api/v1/event/is_added/:id", http.MethodGet)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	sm.EXPECT().CheckSession(cookie.Value).Return(true, id, echo.NewHTTPError(http.StatusInternalServerError))

	err := h.IsAdded(c)

	assert.Error(t, err)
}

func TestSubscriptionHandler_IsAddedNotAuthorized(t *testing.T) {
	c, h, _, sm := setUp(t, "/api/v1/event/is_added/:id", http.MethodGet)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	sm.EXPECT().CheckSession(cookie.Value).Return(false, id, nil)

	err := h.IsAdded(c)

	assert.Error(t, err)
}

func TestSubscriptionHandler_IsAddedErrorMinus(t *testing.T) {
	c, h, _, sm := setUp(t, "/api/v1/event/is_added/:id", http.MethodGet)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)
	c.SetParamNames("id")
	c.SetParamValues("-1")

	sm.EXPECT().CheckSession(cookie.Value).Return(true, id, nil)

	err := h.IsAdded(c)

	assert.Error(t, err)
}

func TestSubscriptionHandler_IsAddedErrorAtoi(t *testing.T) {
	c, h, _, sm := setUp(t, "/api/v1/event/is_added/:id", http.MethodGet)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)
	c.SetParamNames("id")
	c.SetParamValues("a")

	sm.EXPECT().CheckSession(cookie.Value).Return(true, id, nil)

	err := h.IsAdded(c)

	assert.Error(t, err)
}

func TestSubscriptionHandler_IsAddedErrorNoCookie(t *testing.T) {
	c, h, _, _ := setUp(t, "/api/v1/event/is_added/:id", http.MethodGet)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	err := h.IsAdded(c)

	assert.Error(t, err)
}
