package user

import "kudago/models"

type UseCase interface {
	Add(usr *models.RegData) (uint64, error)
	Get(id uint) (*models.User, error)
	Update(uid uint64, user *models.UserData) error
	Login(user *models.User) (uint64, error)
}