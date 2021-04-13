package event

import (
	"kudago/application/models"
	"time"
)

type Repository interface {
	AddEvent(newEvent *models.Event) error
	GetAllEvents(now time.Time, page int) ([]models.EventCardWithDateSQL, error)
	GetOneEventByID(eventId uint64) (models.EventSQL, error)
	DeleteById(eventId uint64) error
	GetTags(eventId uint64) (models.Tags, error)
	UpdateEventAvatar(eventId uint64, path string) error
	GetEventsByCategory(typeEvent string, now time.Time, page int) ([]models.EventCardWithDateSQL, error)
	FindEvents(str string, now time.Time) ([]models.EventCardWithDateSQL, error)
	RecomendSystem(uid uint64, category string) error
	GetPreference(uid uint64) (models.Recomend, error)
	GetRecomended(uid uint64, now time.Time) ([]models.EventCardWithDateSQL, error)
	CategorySearch(str string, category string, now time.Time) ([]models.EventCardWithDateSQL, error)
}
