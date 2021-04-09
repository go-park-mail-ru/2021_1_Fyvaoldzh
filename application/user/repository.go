package user

import "kudago/application/models"

type Repository interface {
	Add(user *models.RegData) (uint64, error)
	Update(id uint64, upUser *models.UserData) error
	ChangeAvatar(uid uint64, path string) error
	GetByIdOwn(id uint64) (*models.UserData, error)
	IsExisting(login string) (bool, error)
	IsExistingEmail(login string) (bool, error)
	IsCorrect(user *models.User) (*models.User, error)
	// TODO: можно вынести в подписки и тыкать на уровне usecase в подписки?
	GetPlanningEvents(id uint64) ([]models.EventCardWithDateSQL, error)
	GetVisitedEvents(id uint64) ([]models.EventCard, error)
	GetFollowers(id uint64) ([]uint64, error)
	DeletePlanningEvent(userId uint64, eventId uint64) error
	AddVisitedEvent(userId uint64, eventId uint64) error
}
