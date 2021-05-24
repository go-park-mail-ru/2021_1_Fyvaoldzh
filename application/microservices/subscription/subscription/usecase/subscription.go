package usecase

import (
	"github.com/labstack/gommon/log"
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

	err = s.repo.SubscribeUser(subscriberId, subscribedToId)
	if err != nil {
		return false, "", err
	}
	err = s.repo.AddSubscriptionAction(subscriberId, subscribedToId)
	if err != nil {
		return false, "", err
	}
	return false, "", nil
}

func (s Subscription) UnsubscribeUser(subscriberId uint64, subscribedToId uint64) (bool, string, error) {
	exists, err := s.repo.CheckSubscription(subscriberId, subscribedToId)
	if err != nil {
		return false, "", err
	}
	if !exists {
		return true, "subscription does not exist", nil
	}

	err = s.repo.UnsubscribeUser(subscriberId, subscribedToId)
	if err != nil {
		return false, "", err
	}
	err = s.repo.RemoveSubscriptionAction(subscriberId, subscribedToId)
	if err != nil {
		return false, "", err
	}
	return false, "", nil
}

func (s Subscription) AddPlanning(userId uint64, eventId uint64) (bool, string, error) {
	existsEvent, err := s.repo.CheckEventInList(eventId)
	log.Info("here1" + err.Error())
	if err != nil {
		return false, "", err
	}
	if !existsEvent {
		return true, "event does not exist", nil
	}
	exists, err := s.repo.CheckEventAdded(userId, eventId)
	log.Info("here2" + err.Error())
	if err != nil {
		return false, "", err
	}
	if exists {
		return true, "event is already added", nil
	}

	err = s.repo.AddPlanning(userId, eventId)
	log.Info("here3" + err.Error())
	if err != nil {
		return false, "", err
	}
	err = s.repo.AddUserEventAction(userId, eventId)
	log.Info("here4" + err.Error())
	if err != nil {
		return false, "", err
	}

	eventDate, err := s.repo.GetTimeEvent(eventId)
	log.Info("here5" + err.Error())
	if err != nil {
		s.logger.Warn(err)
	}

	err = s.repo.AddPlanningNotification(eventId, userId, eventDate)
	log.Info("here6" + err.Error())
	if err != nil {
		s.logger.Warn(err)
	}

	return false, "", nil
}

func (s Subscription) AddVisited(userId uint64, eventId uint64) (bool, string, error) {
	existsEvent, err := s.repo.CheckEventInList(eventId)
	if err != nil {
		return false, "", err
	}
	if !existsEvent {
		return true, "event does not exist", nil
	}
	exists, err := s.repo.CheckEventAdded(userId, eventId)
	if err != nil {
		return false, "", err
	}
	if exists {
		return true, "event is already added", nil
	}

	err = s.repo.AddVisited(userId, eventId)
	if err != nil {
		return false, "", err
	}

	err = s.repo.AddUserEventAction(userId, eventId)
	if err != nil {
		return false, "", err
	}
	return false, "", nil
}

func (s Subscription) RemoveEvent(userId uint64, eventId uint64) (bool, string, error) {
	exists, err := s.repo.CheckEventAdded(userId, eventId)
	if err != nil {
		return false, "", err
	}
	if !exists {
		return true, "event does not exist in list", nil
	}

	err = s.repo.RemoveEvent(userId, eventId)
	if err != nil {
		return false, "", err
	}
	err = s.repo.RemoveUserEventAction(userId, eventId)
	if err != nil {
		return false, "", err
	}

	err = s.repo.RemovePlanningNotification(eventId, userId)
	if err != nil {
		s.logger.Warn(err)
	}

	return false, "", nil

}
