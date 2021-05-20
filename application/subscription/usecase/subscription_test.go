package usecase

import (
	"database/sql"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"kudago/application/models"
	"kudago/application/subscription"
	mock_subscription "kudago/application/subscription/mocks"
	"net/http"
	"time"

	"kudago/pkg/logger"
	"log"
	"testing"
)

var (
	userId     uint64 = 1
	eventId    uint64 = 1
	eventId2   uint64 = 2
	page              = 1
	num        uint64 = 3
	evVisitedSQL = models.EventCardWithDateSQL{
		ID:        eventId2,
		StartDate: time.Now(),
		EndDate:   time.Now(),
	}
	evPlanningSQL = models.EventCardWithDateSQL{
		ID:        eventId,
		StartDate: time.Now(),
		EndDate:   time.Now().Add(10 * time.Hour),
	}
	eventsPlanningSQL = []models.EventCardWithDateSQL{
		evPlanningSQL, evVisitedSQL,
	}
	eventsVisitedSQL = []models.EventCardWithDateSQL{
		evVisitedSQL,
	}
)

var userCardSql = models.UserCardSQL{
	Id:     userId,
	Name:   "name",
	Avatar: "avatar",
	Birthday: sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	},
	City: sql.NullString{
		String: "city",
		Valid:  true,
	},
}

var userCards = []models.UserCardSQL{userCardSql}

func setUp(t *testing.T) (*mock_subscription.MockRepository, subscription.UseCase) {
	ctrl := gomock.NewController(t)

	rep := mock_subscription.NewMockRepository(ctrl)

	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	sugar := l.Sugar()
	zap.NewAtomicLevelAt(zapcore.DebugLevel)

	uc := NewSubscription(rep, logger.NewLogger(sugar))
	return rep, uc
}

///////////////////////////////////////////////////

func TestSubscription_UpdateEventStatus(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().UpdateEventStatus(userId, eventId).Return(nil)

	err := uc.UpdateEventStatus(userId, eventId)

	assert.Nil(t, err)
}

///////////////////////////////////////////////////

func TestSubscription_IsAddedEvent(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().IsAddedEvent(userId, eventId).Return(true, nil)

	_, err := uc.IsAddedEvent(userId, eventId)

	assert.Nil(t, err)
}

///////////////////////////////////////////////////

func TestSubscription_GetFollowers(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().GetFollowers(userId, page).Return(userCards, nil)
	rep.EXPECT().CountUserFollowers(userId).Return(num, nil)

	_, err := uc.GetFollowers(userId, page)

	assert.Nil(t, err)
}

func TestSubscription_GetFollowersErrorRepGF(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().GetFollowers(userId, page).Return(nil, echo.NewHTTPError(http.StatusInternalServerError))

	_, err := uc.GetFollowers(userId, page)

	assert.Error(t, err)
}

///////////////////////////////////////////////////

func TestSubscription_GetSubscriptions(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().GetSubscriptions(userId, page).Return(userCards, nil)
	rep.EXPECT().CountUserFollowers(userId).Return(num, nil)

	_, err := uc.GetSubscriptions(userId, page)

	assert.Nil(t, err)
}

func TestSubscription_GetSubscriptionsErrorRepGS(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().GetSubscriptions(userId, page).Return(nil, echo.NewHTTPError(http.StatusInternalServerError))

	_, err := uc.GetSubscriptions(userId, page)

	assert.Error(t, err)
}

///////////////////////////////////////////////////

func TestSubscription_GetPlanningEvents(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().GetPlanningEvents(userId, page).Return(eventsPlanningSQL, nil)
	rep.EXPECT().UpdateEventStatus(userId, eventId2).Return(nil)

	_, err := uc.GetPlanningEvents(userId, page)

	assert.Nil(t, err)
}

func TestSubscription_GetPlanningEventsErrorRepUES(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().GetPlanningEvents(userId, page).Return(eventsPlanningSQL, nil)
	rep.EXPECT().UpdateEventStatus(userId, eventId2).Return(
		echo.NewHTTPError(http.StatusInternalServerError))

	_, err := uc.GetPlanningEvents(userId, page)

	assert.Error(t, err)
}

func TestSubscription_GetPlanningEventsErrorRepGPE(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().GetPlanningEvents(userId, page).Return(nil,
		echo.NewHTTPError(http.StatusInternalServerError))

	_, err := uc.GetPlanningEvents(userId, page)

	assert.Error(t, err)
}

///////////////////////////////////////////////////

func TestSubscription_GetVisitedEvents(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().GetVisitedEvents(userId, page).Return(eventsVisitedSQL, nil)

	_, err := uc.GetVisitedEvents(userId, page)

	assert.Nil(t, err)
}

func TestSubscription_GetVisitedEventsErrorRepGVE(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().GetVisitedEvents(userId, page).Return(nil,
		echo.NewHTTPError(http.StatusInternalServerError))

	_, err := uc.GetVisitedEvents(userId, page)

	assert.Error(t, err)
}
