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

func (s Subscription) SubscribeUser(subscriberId uint64, subscribedToId uint64) (bool, string, error) {
	exists, err := s.repo.CheckSubscription(subscriberId, subscribedToId)
	if err != nil {
		return false, "", err
	}
	if exists {
		return true, "subscription is already added", nil
	}
	return false, "", s.repo.SubscribeUser(subscriberId, subscribedToId)
}

func (s Subscription) UnsubscribeUser(subscriberId uint64, subscribedToId uint64) (bool, string, error) {
	exists, err := s.repo.CheckSubscription(subscriberId, subscribedToId)
	if err != nil {
		return false, "", err
	}
	if !exists {
		return true, "subscription does not exist", nil
	}
	return false, "", s.repo.UnsubscribeUser(subscriberId, subscribedToId)
}

func (s Subscription) AddPlanning(userId uint64, eventId uint64) (bool, string, error) {
	exists, err := s.repo.CheckEvent(userId, eventId)
	if err != nil {
		return false, "", err
	}
	if exists {
		return true, "event is already added", nil
	}
	return false, "", s.repo.AddPlanning(userId, eventId)
}

func (s Subscription) AddVisited(userId uint64, eventId uint64) (bool, string, error) {
	exists, err := s.repo.CheckEvent(userId, eventId)
	if err != nil {
		return false, "", err
	}
	if exists {
		return true, "event is already added", nil
	}
	return false, "", s.repo.AddVisited(userId, eventId)
}

func (s Subscription) RemoveEvent(userId uint64, eventId uint64) (bool, string, error) {
	exists, err := s.repo.CheckEvent(userId, eventId)
	if err != nil {
		return false, "", err
	}
	if !exists {
		return true, "event does not exist in list", nil
	}
	return false, "", s.repo.RemoveEvent(userId, eventId)
}
