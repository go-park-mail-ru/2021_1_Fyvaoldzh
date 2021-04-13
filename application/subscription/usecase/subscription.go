package usecase

import (
	"kudago/application/subscription"
	"net/http"

	"github.com/labstack/echo"
	"go.uber.org/zap"
)

type Subscription struct {
	repo   subscription.Repository
	logger *zap.SugaredLogger
}

func NewSubscription(subRepo subscription.Repository, logger *zap.SugaredLogger) subscription.UseCase {
	return &Subscription{repo: subRepo, logger: logger}
}

func (s Subscription) SubscribeUser(subscriberId uint64, subscribedToId uint64) error {
	if subscriberId == subscribedToId {
		s.logger.Warn("subscriberId == subscribedToId")
		return echo.NewHTTPError(http.StatusBadRequest, "incorrect data")
	}

	return s.repo.SubscribeUser(subscriberId, subscribedToId)
}

func (s Subscription) UnsubscribeUser(subscriberId uint64, subscribedToId uint64) error {
	if subscriberId == subscribedToId {
		s.logger.Warn("subscriberId == subscribedToId")
		return echo.NewHTTPError(http.StatusBadRequest, "incorrect data")
	}

	return s.repo.UnsubscribeUser(subscriberId, subscribedToId)
}

func (s Subscription) AddPlanning(userId uint64, eventId uint64) error {
	return s.repo.AddPlanning(userId, eventId)
}

func (s Subscription) RemovePlanning(userId uint64, eventId uint64) error {
	return s.repo.RemovePlanning(userId, eventId)
}

func (s Subscription) AddVisited(userId uint64, eventId uint64) error {
	return s.repo.AddVisited(userId, eventId)
}

func (s Subscription) RemoveVisited(userId uint64, eventId uint64) error {
	return s.repo.RemoveVisited(userId, eventId)
}

func (s Subscription) UpdateEventStatus(userId uint64, eventId uint64) error {
	return s.repo.UpdateEventStatus(userId, eventId)
}

func (s Subscription) IsAddedEvent(userId uint64, eventId uint64) (bool, error) {
	return s.repo.IsAddedEvent(userId, eventId)
}