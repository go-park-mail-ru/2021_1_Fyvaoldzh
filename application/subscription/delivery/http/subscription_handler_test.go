package http

import (
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	clientAuth "kudago/application/microservices/auth/client"
	clientSub "kudago/application/microservices/subscription/client"
	"kudago/application/models"
	mock_subscription "kudago/application/subscription/mocks"
	"kudago/pkg/constants"
	"kudago/pkg/logger"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var (
	userId        uint64 = 1
	n                    = 1
	id            uint64 = 1
	num                  = 2
	numu          uint64 = 2
	strId                = "1"
	name                 = "name"
	evPlanningSQL        = models.EventCardWithDateSQL{
		ID:        1,
		StartDate: time.Now(),
		EndDate:   time.Now().Add(10 * time.Hour),
	}
	evPlanning = models.EventCard{
		ID:          1,
		Title:       "title",
		Place:       "place",
		Description: "desc",
		StartDate:   evPlanningSQL.StartDate.String(),
		EndDate:     evPlanningSQL.EndDate.String(),
	}
	evVisitedSQL = models.EventCardWithDateSQL{
		ID:        2,
		StartDate: time.Now(),
		EndDate:   time.Now(),
	}
	evVisited = models.EventCard{
		ID:          2,
		Title:       "title",
		Place:       "place",
		Description: "desc",
		StartDate:   evVisitedSQL.StartDate.String(),
		EndDate:     evVisitedSQL.EndDate.String(),
	}
	eventsPlanningSQL = []models.EventCardWithDateSQL{
		evPlanningSQL, evVisitedSQL,
	}
	eventsVisitedSQL = []models.EventCardWithDateSQL{
		evVisitedSQL,
	}
	eventsPlanning = models.EventCards{
		evPlanning,
	}
	eventsVisited = models.EventCards{
		evVisited,
	}
)

var testUserCard = &models.UserCard{
	Id:   userId,
	Name: name,
}

var testUserCards = models.UserCards{*testUserCard}

func setUp(t *testing.T, url, method string) (echo.Context,
	SubscriptionHandler,
	*mock_subscription.MockUseCase,
	*clientAuth.MockIAuthClient,
	*clientSub.MockISubscriptionClient) {
	e := echo.New()
	r := e.Router()
	r.Add(method, url, func(echo.Context) error { return nil })

	ctrl := gomock.NewController(t)
	usecase := mock_subscription.NewMockUseCase(ctrl)
	rpcAuth := clientAuth.NewMockIAuthClient(ctrl)
	rpcSub := clientSub.NewMockISubscriptionClient(ctrl)

	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	sugar := l.Sugar()
	zap.NewAtomicLevelAt(zapcore.DebugLevel)

	handler := SubscriptionHandler{
		usecase: usecase,
		rpcAuth: rpcAuth,
		rpcSub:  rpcSub,
		Logger:  logger.NewLogger(sugar),
	}

	var req *http.Request
	req = httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(url)

	return c, handler, usecase, rpcAuth, rpcSub
}

///////////////////////////////////////////////////

func TestSubscriptionHandler_Subscribe(t *testing.T) {
	c, h, _, _, rpcSub := setUp(t, "/api/v1/subscribe/user/:id", http.MethodPost)
	c.Set(constants.UserIdKey, userId)
	c.Set(constants.IdKey, num)

	rpcSub.EXPECT().Subscribe(userId, numu).Return(nil)

	err := h.Subscribe(c)

	assert.Nil(t, err)
}

func TestSubscriptionHandler_SubscribeErrorSame(t *testing.T) {
	c, h, _, _, _ := setUp(t, "/api/v1/subscribe/user/:id", http.MethodPost)
	c.Set(constants.UserIdKey, userId)
	c.Set(constants.IdKey, n)

	err := h.Subscribe(c)

	assert.Error(t, err)
}

///////////////////////////////////////////////////

func TestSubscriptionHandler_Unsubscribe(t *testing.T) {
	c, h, _, _, rpcSub := setUp(t, "/api/v1/unsubscribe/user/:id", http.MethodDelete)
	c.Set(constants.UserIdKey, userId)
	c.Set(constants.IdKey, num)

	rpcSub.EXPECT().Unsubscribe(userId, numu).Return(nil)

	err := h.Unsubscribe(c)

	assert.Nil(t, err)
}

///////////////////////////////////////////////////

func TestSubscriptionHandler_AddPlanningEvent(t *testing.T) {
	c, h, _, _, rpcSub := setUp(t, "/api/v1/add/planning/:id", http.MethodPost)
	c.Set(constants.UserIdKey, userId)
	c.Set(constants.IdKey, num)

	rpcSub.EXPECT().AddPlanningEvent(userId, numu).Return(nil)

	err := h.AddPlanningEvent(c)

	assert.Nil(t, err)
}

///////////////////////////////////////////////////

func TestSubscriptionHandler_RemoveEvent(t *testing.T) {
	c, h, _, _, rpcSub := setUp(t, "/api/v1/remove/:id", http.MethodDelete)
	c.Set(constants.UserIdKey, userId)
	c.Set(constants.IdKey, num)

	rpcSub.EXPECT().RemoveEvent(userId, numu).Return(nil)

	err := h.RemoveEvent(c)

	assert.Nil(t, err)
}

///////////////////////////////////////////////////

func TestSubscriptionHandler_AddVisitedEvent(t *testing.T) {
	c, h, _, _, rpcSub := setUp(t, "/api/v1/add/visited/:id", http.MethodPost)
	c.Set(constants.UserIdKey, userId)
	c.Set(constants.IdKey, num)

	rpcSub.EXPECT().AddVisitedEvent(userId, numu).Return(nil)

	err := h.AddVisitedEvent(c)

	assert.Nil(t, err)
}

///////////////////////////////////////////////////

func TestSubscriptionHandler_IsAdded(t *testing.T) {
	c, h, usecase, _, _ := setUp(t, "/api/v1/event/is_added/:id", http.MethodGet)
	c.Set(constants.UserIdKey, userId)
	c.Set(constants.IdKey, num)

	usecase.EXPECT().IsAddedEvent(userId, numu).Return(true, nil)

	err := h.IsAdded(c)

	assert.Nil(t, err)
}

func TestSubscriptionHandler_IsAddedErrorUC(t *testing.T) {
	c, h, usecase, _, _ := setUp(t, "/api/v1/event/is_added/:id", http.MethodGet)
	c.Set(constants.UserIdKey, userId)
	c.Set(constants.IdKey, num)

	usecase.EXPECT().IsAddedEvent(userId, numu).Return(true, echo.NewHTTPError(http.StatusInternalServerError))

	err := h.IsAdded(c)

	assert.Error(t, err)
}

///////////////////////////////////////////////////

func TestSubscriptionHandler_GetFollowersErrorUC(t *testing.T) {
	c, h, usecase, _, _ := setUp(t, "/api/v1/followers/:id", http.MethodGet)
	c.Set(constants.IdKey, n)
	c.Set(constants.PageKey, num)

	usecase.EXPECT().GetFollowers(userId, num).Return(testUserCards,
		echo.NewHTTPError(http.StatusInternalServerError))

	err := h.GetFollowers(c)

	assert.Error(t, err)
}

///////////////////////////////////////////////////

func TestSubscriptionHandler_GetSubscriptionsErrorUC(t *testing.T) {
	c, h, usecase, _, _ := setUp(t, "/api/v1/subscriptions/:id", http.MethodGet)
	c.Set(constants.IdKey, n)
	c.Set(constants.PageKey, num)

	usecase.EXPECT().GetSubscriptions(userId, num).Return(testUserCards,
		echo.NewHTTPError(http.StatusInternalServerError))

	err := h.GetSubscriptions(c)

	assert.Error(t, err)
}

///////////////////////////////////////////////////

func TestSubscriptionHandler_GetPlanningEventsErrorUC(t *testing.T) {
	c, h, usecase, _, _ := setUp(t, "/api/v1/get/planning/:id", http.MethodGet)
	c.Set(constants.IdKey, n)
	c.Set(constants.PageKey, num)

	usecase.EXPECT().GetPlanningEvents(userId, num).Return(eventsPlanning,
		echo.NewHTTPError(http.StatusInternalServerError))

	err := h.GetPlanningEvents(c)

	assert.Error(t, err)
}

///////////////////////////////////////////////////

func TestSubscriptionHandler_GetVisitedEventsErrorUC(t *testing.T) {
	c, h, usecase, _, _ := setUp(t, "/api/v1/get/visited/:id", http.MethodGet)
	c.Set(constants.IdKey, n)
	c.Set(constants.PageKey, num)

	usecase.EXPECT().GetVisitedEvents(userId, num).Return(eventsPlanning,
		echo.NewHTTPError(http.StatusInternalServerError))

	err := h.GetVisitedEvents(c)

	assert.Error(t, err)
}

///////////////////////////////////////////////////
