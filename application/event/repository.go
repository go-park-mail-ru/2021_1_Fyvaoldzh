package event

import "kudago/application/models"

type Repository interface {
	AddEvent(newEvent *models.Event) error
	GetAllEvents() ([]models.EventCardSQL, error)
	GetOneEventByID(eventId uint64) (models.EventSQL, error)
	DeleteById(eventId uint64) error
	GetCategoryTags(eventId uint64) ([]models.CategoryTagDescription, error)
	UpdateEventAvatar(eventId uint64, path string) error
	GetEventsByType(typeEvent string) ([]models.EventCardSQL, error)
	FindEvents(str string) ([]models.EventCardSQL, error)
}
