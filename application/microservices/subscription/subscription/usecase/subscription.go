package usecase

import (
	"kudago/application/microservices/subscription/subscription"
	"kudago/pkg/logger"
)

type Subscription struct {
	repo   subscription.Repository
	logger logger.Logger
}

func NewSubscription(subRepo subscription.Repository, logger logger.Logger) subscription.UseCase {
	return &Subscription{repo: subRepo, logger: logger}
}

func (s Subscription) SubscribeUser(subscriberId uint64, subscribedToId uint64) error {
	return s.repo.SubscribeUser(subscriberId, subscribedToId)
}

func (s Subscription) UnsubscribeUser(subscriberId uint64, subscribedToId uint64) error {
	return s.repo.UnsubscribeUser(subscriberId, subscribedToId)
}

func (s Subscription) AddPlanning(userId uint64, eventId uint64) error {
	return s.repo.AddPlanning(userId, eventId)
}

func (s Subscription) AddVisited(userId uint64, eventId uint64) error {
	return s.repo.AddVisited(userId, eventId)
}

func (s Subscription) RemoveEvent(userId uint64, eventId uint64) error {
	return s.repo.RemoveEvent(userId, eventId)
}

