package event

import "kudago/application/models"

type Repository interface {
	GetAllEvents() ([]models.EventCardSQL, error)
	GetOneEventByID(eventId uint64) (models.EventSQL, error)
	GetCategoryTags(eventId uint64) ([]models.CategoryTagDescription, error)
}