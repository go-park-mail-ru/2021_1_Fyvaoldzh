package event

import "kudago/application/models"

type Repository interface {
	AddEvent(newEvent *models.Event) error
	GetAllEvents() ([]models.EventCardWithDateSQL, error)
	GetOneEventByID(eventId uint64) (models.EventSQL, error)
	DeleteById(eventId uint64) error
	GetTags(eventId uint64) (models.Tags, error)
	UpdateEventAvatar(eventId uint64, path string) error
	GetEventsByCategory(typeEvent string) ([]models.EventCardWithDateSQL, error)
	FindEvents(str string) ([]models.EventCardWithDateSQL, error)
	RecomendSystem(uid uint64, category string) error
}
