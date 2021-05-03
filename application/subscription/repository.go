package subscription

import "kudago/application/models"

type Repository interface {
	UpdateEventStatus(userId uint64, eventId uint64) error
	GetUserFollowers(id uint64) ([]uint64, error)
	GetEventFollowers(eventId uint64) (models.UsersOnEvent, error)
	IsAddedEvent(userId uint64, eventId uint64) (bool, error)
	GetFollowers(id uint64) ([]models.UserCardSQL, error)
	GetSubscriptions(id uint64) ([]models.UserCardSQL, error)
	GetPlanningEvents(id uint64) ([]models.EventCardWithDateSQL, error)
	GetVisitedEvents(id uint64) ([]models.EventCardWithDateSQL, error)
}
