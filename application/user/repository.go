package user

import (
	"kudago/application/models"
)

type Repository interface {
	Add(user *models.RegData) (uint64, error)
	AddToPreferences(id uint64) error
	Update(id uint64, upUser *models.UserDataSQL) error
	ChangeAvatar(id uint64, path string) error
	GetByIdOwn(id uint64) (*models.UserDataSQL, error)
	IsExisting(login string) (bool, error)
	IsExistingEmail(login string) (bool, error)
	IsCorrect(user *models.User) (*models.User, error)
	IsExistingUserId(userId uint64) error
	GetUsers(page int) ([]models.UserCardSQL, error)
	FindUsers(str string, page int) ([]models.UserCardSQL, error)
	GetUserByID(id uint64) (models.UserOnEvent, error)
	GetActions(id uint64, page int) ([]*models.ActionCard, error)
}
