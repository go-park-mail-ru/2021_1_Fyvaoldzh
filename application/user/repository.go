package user

import "kudago/models"

type Repository interface {
	Add(user *models.RegData) (uint64, error)
	Update(id uint64, upUser *models.UserData) error
	GetByIdOther(id uint64) (*models.User, error)
	GetByIdOwn(id uint64) (*models.UserData, error)
	IsExisting(login string) (bool, error)
	IsExistingEmail(login string) (bool, error)
	IsCorrect(user *models.User) (uint64, error)
	GetPlanningEvents(id uint64) ([]uint64, error)
	GetVisitedEvents(id uint64) ([]uint64, error)
	GetFollowers(id uint64) ([]uint64, error)
}