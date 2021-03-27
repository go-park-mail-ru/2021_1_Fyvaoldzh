package user

import "kudago/models"

type UseCase interface {
	Add(usr *models.RegData) (*models.RegData, error)
	Get(id uint) (*models.User, error)
	Update(user *models.User) error
}