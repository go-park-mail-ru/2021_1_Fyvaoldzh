package usecase

import (
	"kudago/application/subscription"
	"kudago/pkg/logger"
)

type Subscription struct {
	repo   subscription.Repository
	Logger logger.Logger
}

func NewSubscription(subRepo subscription.Repository, logger logger.Logger) subscription.UseCase {
	return &Subscription{repo: subRepo, Logger: logger}
}

func (s Subscription) UpdateEventStatus(userId uint64, eventId uint64) error {
	return s.repo.UpdateEventStatus(userId, eventId)
}

func (s Subscription) IsAddedEvent(userId uint64, eventId uint64) (bool, error) {
	return s.repo.IsAddedEvent(userId, eventId)
}




