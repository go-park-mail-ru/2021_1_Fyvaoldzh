package user

import "kudago/models"

type Repository interface {
	Add(user *models.RegData) (err error)
	Update(id uint, upUser *models.User) error
	GetByIdOther(id uint) (user *models.User, err error)
	GetByIdOwn(id uint) (user *models.User, err error)
	GetByName(login string) (user *models.User, err error)

}