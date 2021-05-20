package event

import (
	"kudago/application/models"
	"mime/multipart"
)

//go:generate mockgen -destination=./mock/usecase_mock.go -package=mock -source=./application/event/usecase.go

type UseCase interface {
	GetAllEvents(page int) (models.EventCards, error)
	GetOneEvent(eventId uint64) (models.Event, error)
	GetOneEventName(eventId uint64) (string, error)
	Delete(eventId uint64) error
	CreateNewEvent(newEvent *models.Event) error
	SaveImage(eventId uint64, img *multipart.FileHeader) error
	GetEventsByCategory(typeEvent string, page int) (models.EventCards, error)
	GetImage(eventId uint64) ([]byte, error)
	FindEvents(str string, category string, page int) (models.EventCards, error)
	RecomendSystem(uid uint64, category string) error
	GetRecommended(uid uint64, page int) (models.EventCards, error)
	GetNear(coord models.Coordinates, page int) (models.EventCards, error)
}
