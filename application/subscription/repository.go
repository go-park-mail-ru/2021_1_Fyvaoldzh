package subscription

import "kudago/application/models"

type Repository interface {
	SubscribeUser(subscriberId uint64, subscribedToId uint64) error
	UnsubscribeUser(subscriberId uint64, subscribedToId uint64) error
	AddPlanning(userId uint64, eventId uint64) error
	RemovePlanning(userId uint64, eventId uint64) error
	AddVisited(userId uint64, eventId uint64) error
	RemoveVisited(userId uint64, eventId uint64) error
	UpdateEventStatus(userId uint64, eventId uint64) error
	GetPlanningEvents(id uint64) ([]models.EventCardWithDateSQL, error)
	GetVisitedEvents(id uint64) ([]models.EventCardWithDateSQL, error)
	GetFollowers(id uint64) ([]uint64, error)
	GetEventFollowers(eventId uint64) (models.UsersOnEvent, error)
	IsAddedEvent(userId uint64, eventId uint64) (bool, error)
}
