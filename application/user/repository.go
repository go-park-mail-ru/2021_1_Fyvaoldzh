package user

import "kudago/models"

type Repository interface {
	Add(user *models.RegData) (uint64, error)
	Update(id uint64, upUser *models.UserData) error
	GetByIdOther(id uint64) (*models.User, error)
	GetByIdOwn(id uint64) (*models.UserData, error)
	GetByName(login string) (*models.User, error)
	IsExisting(login string) (bool, error)
	IsExistingEmail(login string) (bool, error)
	IsCorrect(user *models.User) (uint64, error)

}