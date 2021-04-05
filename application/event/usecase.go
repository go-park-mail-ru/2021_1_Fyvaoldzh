package event

import "kudago/application/models"

type UseCase interface {
	GetAllEvents() (models.EventCards, error)
	GetOneEvent(eventId uint64) (models.Event, error)
}