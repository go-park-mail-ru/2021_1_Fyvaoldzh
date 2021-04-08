package event

import "kudago/application/models"

type Repository interface {
	AddEvent(newEvent *models.Event) error
	GetAllEvents() ([]models.EventCardWithDateSQL, error)
	GetOneEventByID(eventId uint64) (models.EventSQL, error)
	DeleteById(eventId uint64) error
	//GetCategoryTags(eventId uint64) ([]models.CategoryTagDescription, error)
	UpdateEventAvatar(eventId uint64, path string) error
	GetEventsByType(typeEvent string) ([]models.EventCard, error)
	FindEvents(str string) ([]models.EventCard, error)
}
