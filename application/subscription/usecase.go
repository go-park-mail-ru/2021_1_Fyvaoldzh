package subscription

import "kudago/application/models"

type UseCase interface {
	UpdateEventStatus(userId uint64, eventId uint64) error
	IsAddedEvent(userId uint64, eventId uint64) (bool, error)
	GetFollowers(id uint64) (models.UserCards, error)
	GetSubscriptions(id uint64) (models.UserCards, error)
	GetPlanningEvents(id uint64) (models.EventCards, error)
	GetVisitedEvents(id uint64) (models.EventCards, error)
}

