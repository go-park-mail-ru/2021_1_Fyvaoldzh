package usecase

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"kudago/application/subscription"
	mock_subscription "kudago/application/subscription/mocks"

	"kudago/pkg/logger"
	"log"
	"testing"
)

var (
	userId  uint64 = 1
	eventId uint64 = 1
)

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

func TestSubscription_SubscribeUser(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().SubscribeUser(userId, userId+1).Return(nil)

	err := uc.SubscribeUser(userId, userId+1)

	assert.Nil(t, err)
}

func TestSubscription_SubscribeUserErrorSameId(t *testing.T) {
	_, uc := setUp(t)

	err := uc.SubscribeUser(userId, userId)

	assert.Error(t, err)
}

///////////////////////////////////////////////////

func TestSubscription_UnsubscribeUser(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().UnsubscribeUser(userId, userId+1).Return(nil)

	err := uc.UnsubscribeUser(userId, userId+1)

	assert.Nil(t, err)
}

func TestSubscription_UnsubscribeUserErrorSameId(t *testing.T) {
	_, uc := setUp(t)

	err := uc.UnsubscribeUser(userId, userId)

	assert.Error(t, err)
}

///////////////////////////////////////////////////

func TestSubscription_AddPlanning(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().AddPlanning(userId, eventId).Return(nil)

	err := uc.AddPlanning(userId, eventId)

	assert.Nil(t, err)
}

///////////////////////////////////////////////////

func TestSubscription_RemovePlanning(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().RemovePlanning(userId, eventId).Return(nil)

	err := uc.RemovePlanning(userId, eventId)

	assert.Nil(t, err)
}

///////////////////////////////////////////////////

func TestSubscription_AddVisited(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().AddVisited(userId, eventId).Return(nil)

	err := uc.AddVisited(userId, eventId)

	assert.Nil(t, err)
}

///////////////////////////////////////////////////

func TestSubscription_RemoveVisited(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().RemoveVisited(userId, eventId).Return(nil)

	err := uc.RemoveVisited(userId, eventId)

	assert.Nil(t, err)
}

///////////////////////////////////////////////////

func TestSubscription_UpdateEventStatus(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().UpdateEventStatus(userId, eventId).Return(nil)

	err := uc.UpdateEventStatus(userId, eventId)

	assert.Nil(t, err)
}

func TestSubscription_IsAddedEvent(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().IsAddedEvent(userId, eventId).Return(true, nil)

	_, err := uc.IsAddedEvent(userId, eventId)

	assert.Nil(t, err)
}