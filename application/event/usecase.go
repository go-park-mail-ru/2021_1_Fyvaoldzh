package event

import (
	"kudago/application/models"
	"mime/multipart"
)

type UseCase interface {
	GetAllEvents() (models.EventCards, error)
	GetOneEvent(eventId uint64) (models.Event, error)
	Delete(eventId uint64) error
	CreateNewEvent(newEvent *models.Event) error
	SaveImage(eventId uint64, img *multipart.FileHeader) error
	GetEventsByType(typeEvent string) (models.EventCards, error)
	GetImage(eventId uint64) ([]byte, error)
}