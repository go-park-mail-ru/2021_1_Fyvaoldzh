package subscription

import "kudago/application/models"

type Repository interface {
	UpdateEventStatus(userId uint64, eventId uint64) error
	GetFollowers(id uint64) ([]uint64, error)
	GetEventFollowers(eventId uint64) (models.UsersOnEvent, error)
	IsAddedEvent(userId uint64, eventId uint64) (bool, error)
}
