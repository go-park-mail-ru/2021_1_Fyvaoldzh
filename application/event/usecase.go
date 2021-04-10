package event

import (
	"kudago/application/models"
	"mime/multipart"
)

type UseCase interface {
	GetAllEvents(page int) (models.EventCards, error)
	GetOneEvent(eventId uint64) (models.Event, error)
	Delete(eventId uint64) error
	CreateNewEvent(newEvent *models.Event) error
	SaveImage(eventId uint64, img *multipart.FileHeader) error
	GetEventsByCategory(typeEvent string) (models.EventCards, error)
	GetImage(eventId uint64) ([]byte, error)
	FindEvents(str string) (models.EventCards, error)
	RecomendSystem(uid uint64, category string) error
	GetRecomended(uid uint64) (models.EventCards, error)
}
