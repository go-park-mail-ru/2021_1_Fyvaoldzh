package usecase

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"kudago/application/microservices/subscription/subscription"
	mock_subscription "kudago/application/microservices/subscription/subscription/mocks"
	"kudago/pkg/logger"
	"log"
	"testing"
)

var (
	userId  uint64 = 1
	userId2 uint64 = 2
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

func TestSubscription_SubscribeUser(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().SubscribeUser(userId, userId2).Return(nil)
	rep.EXPECT().CheckSubscription(userId, userId2).Return(false, nil)
	rep.EXPECT().AddSubscriptionAction(userId, userId2).Return(nil)

	_, _, err := uc.SubscribeUser(userId, userId2)

	assert.Nil(t, err)
}

func TestSubscription_SubscribeUserErrorUCASA(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().SubscribeUser(userId, userId2).Return(nil)
	rep.EXPECT().CheckSubscription(userId, userId2).Return(false, nil)
	rep.EXPECT().AddSubscriptionAction(userId, userId2).Return(
		status.Error(codes.Internal, ""))

	_, _, err := uc.SubscribeUser(userId, userId2)

	assert.Error(t, err)
}

func TestSubscription_SubscribeUserErrorUCSU(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().CheckSubscription(userId, userId2).Return(false, nil)
	rep.EXPECT().SubscribeUser(userId, userId2).Return(status.Error(codes.Internal, ""))

	_, _, err := uc.SubscribeUser(userId, userId2)

	assert.Error(t, err)
}

func TestSubscription_SubscribeUserErrorUCCS(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().CheckSubscription(userId, userId2).Return(false,
		status.Error(codes.Internal, ""))

	_, _, err := uc.SubscribeUser(userId, userId2)

	assert.Error(t, err)
}

func TestSubscription_SubscribeUserSubscriptionExists(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().CheckSubscription(userId, userId2).Return(true, nil)

	flag, _, err := uc.SubscribeUser(userId, userId2)

	assert.Nil(t, err)
	assert.Equal(t, true, flag)
}

///////////////////////////////////////////////////

func TestSubscription_UnsubscribeUser(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().CheckSubscription(userId, userId2).Return(true, nil)
	rep.EXPECT().UnsubscribeUser(userId, userId2).Return(nil)
	rep.EXPECT().RemoveSubscriptionAction(userId, userId2).Return(nil)

	_, _, err := uc.UnsubscribeUser(userId, userId2)

	assert.Nil(t, err)
}

func TestSubscription_UnsubscribeUserErrorUCRSA(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().CheckSubscription(userId, userId2).Return(true, nil)
	rep.EXPECT().UnsubscribeUser(userId, userId2).Return(nil)
	rep.EXPECT().RemoveSubscriptionAction(userId, userId2).Return(
		status.Error(codes.Internal, ""))

	_, _, err := uc.UnsubscribeUser(userId, userId2)

	assert.Error(t, err)
}

func TestSubscription_UnsubscribeUserErrorUCUU(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().CheckSubscription(userId, userId2).Return(true, nil)
	rep.EXPECT().UnsubscribeUser(userId, userId2).Return(status.Error(codes.Internal, ""))

	_, _, err := uc.UnsubscribeUser(userId, userId2)

	assert.Error(t, err)
}

func TestSubscription_UnsubscribeUserDoesNotExist(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().CheckSubscription(userId, userId2).Return(false, nil)

	flag, _, err := uc.UnsubscribeUser(userId, userId2)

	assert.Equal(t, true, flag)
	assert.Nil(t, err)
}

func TestSubscription_UnsubscribeUserUCErrorCS(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().CheckSubscription(userId, userId2).Return(false,
		status.Error(codes.Internal, ""))

	_, _, err := uc.UnsubscribeUser(userId, userId2)

	assert.Error(t, err)
}

///////////////////////////////////////////////////

func TestSubscription_AddPlanning(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().CheckEventInList(eventId).Return(true, nil)
	rep.EXPECT().CheckEventAdded(userId, eventId).Return(false, nil)
	rep.EXPECT().AddPlanning(userId, eventId).Return(nil)
	rep.EXPECT().AddUserEventAction(userId, eventId).Return(nil)

	_, _, err := uc.AddPlanning(userId, eventId)

	assert.Nil(t, err)
}

func TestSubscription_AddPlanningErrorUCAUEA(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().CheckEventInList(eventId).Return(true, nil)
	rep.EXPECT().CheckEventAdded(userId, eventId).Return(false, nil)
	rep.EXPECT().AddPlanning(userId, eventId).Return(nil)
	rep.EXPECT().AddUserEventAction(userId, eventId).Return(
		status.Error(codes.Internal, ""))

	_, _, err := uc.AddPlanning(userId, eventId)

	assert.Error(t, err)
}

func TestSubscription_AddPlanningErrorUCAP(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().CheckEventInList(eventId).Return(true, nil)
	rep.EXPECT().CheckEventAdded(userId, eventId).Return(false, nil)
	rep.EXPECT().AddPlanning(userId, eventId).Return(
		status.Error(codes.Internal, ""))

	_, _, err := uc.AddPlanning(userId, eventId)

	assert.Error(t, err)
}

func TestSubscription_AddPlanningExists(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().CheckEventInList(eventId).Return(true, nil)
	rep.EXPECT().CheckEventAdded(userId, eventId).Return(true, nil)

	flag, _, err := uc.AddPlanning(userId, eventId)

	assert.Nil(t, err)
	assert.Equal(t, true, flag)
}

func TestSubscription_AddPlanningErrorUCCEA(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().CheckEventInList(eventId).Return(true, nil)
	rep.EXPECT().CheckEventAdded(userId, eventId).Return(false,
		status.Error(codes.Internal, ""))

	_, _, err := uc.AddPlanning(userId, eventId)

	assert.Error(t, err)
}

func TestSubscription_AddPlanningNotExistingEvent(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().CheckEventInList(eventId).Return(false, nil)

	flag, _, err := uc.AddPlanning(userId, eventId)

	assert.Nil(t, err)
	assert.Equal(t, true, flag)
}

func TestSubscription_AddPlanningErrorUCCEIL(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().CheckEventInList(eventId).Return(false, status.Error(codes.Internal, ""))

	_, _, err := uc.AddPlanning(userId, eventId)

	assert.Error(t, err)
}

///////////////////////////////////////////////////

func TestSubscription_AddVisited(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().CheckEventInList(eventId).Return(true, nil)
	rep.EXPECT().CheckEventAdded(userId, eventId).Return(false, nil)
	rep.EXPECT().AddVisited(userId, eventId).Return(nil)
	rep.EXPECT().AddUserEventAction(userId, eventId).Return(nil)

	_, _, err := uc.AddVisited(userId, eventId)

	assert.Nil(t, err)
}

func TestSubscription_AddVisitedErrorUCAUEA(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().CheckEventInList(eventId).Return(true, nil)
	rep.EXPECT().CheckEventAdded(userId, eventId).Return(false, nil)
	rep.EXPECT().AddVisited(userId, eventId).Return(nil)
	rep.EXPECT().AddUserEventAction(userId, eventId).Return(
		status.Error(codes.Internal, ""))

	_, _, err := uc.AddVisited(userId, eventId)

	assert.Error(t, err)
}

func TestSubscription_AddVisitedErrorUCAP(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().CheckEventInList(eventId).Return(true, nil)
	rep.EXPECT().CheckEventAdded(userId, eventId).Return(false, nil)
	rep.EXPECT().AddVisited(userId, eventId).Return(
		status.Error(codes.Internal, ""))

	_, _, err := uc.AddVisited(userId, eventId)

	assert.Error(t, err)
}

func TestSubscription_AddVisitedExists(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().CheckEventInList(eventId).Return(true, nil)
	rep.EXPECT().CheckEventAdded(userId, eventId).Return(true, nil)

	flag, _, err := uc.AddVisited(userId, eventId)

	assert.Nil(t, err)
	assert.Equal(t, true, flag)
}

func TestSubscription_AddVisitedErrorUCCEA(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().CheckEventInList(eventId).Return(true, nil)
	rep.EXPECT().CheckEventAdded(userId, eventId).Return(false,
		status.Error(codes.Internal, ""))

	_, _, err := uc.AddVisited(userId, eventId)

	assert.Error(t, err)
}

func TestSubscription_AddVisitedNotExistingEvent(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().CheckEventInList(eventId).Return(false, nil)

	flag, _, err := uc.AddVisited(userId, eventId)

	assert.Nil(t, err)
	assert.Equal(t, true, flag)
}

func TestSubscription_AddVisitedErrorUCCEIL(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().CheckEventInList(eventId).Return(false, status.Error(codes.Internal, ""))

	_, _, err := uc.AddVisited(userId, eventId)

	assert.Error(t, err)
}

///////////////////////////////////////////////////

func TestSubscription_RemoveEvent(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().CheckEventAdded(userId, eventId).Return(true, nil)
	rep.EXPECT().RemoveEvent(userId, eventId).Return(nil)
	rep.EXPECT().RemoveUserEventAction(userId, eventId).Return(nil)

	_, _, err := uc.RemoveEvent(userId, eventId)

	assert.Nil(t, err)
}

func TestSubscription_RemoveEventErrorUCRUEA(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().CheckEventAdded(userId, eventId).Return(true, nil)
	rep.EXPECT().RemoveEvent(userId, eventId).Return(nil)
	rep.EXPECT().RemoveUserEventAction(userId, eventId).Return(
		status.Error(codes.Internal, ""))

	_, _, err := uc.RemoveEvent(userId, eventId)

	assert.Error(t, err)
}

func TestSubscription_RemoveEventErrorUCRE(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().CheckEventAdded(userId, eventId).Return(true, nil)
	rep.EXPECT().RemoveEvent(userId, eventId).Return(
		status.Error(codes.Internal, ""))

	_, _, err := uc.RemoveEvent(userId, eventId)

	assert.Error(t, err)
}

func TestSubscription_RemoveEventNotExists(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().CheckEventAdded(userId, eventId).Return(false, nil)

	flag, _, err := uc.RemoveEvent(userId, eventId)

	assert.Nil(t, err)
	assert.Equal(t, true, flag)
}

func TestSubscription_RemoveEventErrorUCCEA(t *testing.T) {
	rep, uc := setUp(t)

	rep.EXPECT().CheckEventAdded(userId, eventId).Return(false,
		status.Error(codes.Internal, ""))

	_, _, err := uc.RemoveEvent(userId, eventId)

	assert.Error(t, err)
}
