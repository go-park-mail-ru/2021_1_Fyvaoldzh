package user

import "kudago/application/models"

type Repository interface {
	Add(user *models.RegData) (uint64, error)
	AddToPreferences(id uint64) error
	Update(id uint64, upUser *models.UserData) error
	ChangeAvatar(id uint64, path string) error
	GetByIdOwn(id uint64) (*models.UserData, error)
	IsExisting(login string) (bool, error)
	IsExistingEmail(login string) (bool, error)
	IsCorrect(user *models.User) (*models.User, error)
	IsExistingUserId(userId uint64) error
	GetUsers(page int) ([]models.UserCardSQL, error)
	UpdateEventStatus(userId uint64, eventId uint64) error
	GetPlanningEvents(id uint64) ([]models.EventCardWithDateSQL, error)
	GetVisitedEvents(id uint64) ([]models.EventCardWithDateSQL, error)
	GetFollowers(id uint64) ([]uint64, error)
	GetEventFollowers(eventId uint64) (models.UsersOnEvent, error)
	IsAddedEvent(userId uint64, eventId uint64) (bool, error)
}
