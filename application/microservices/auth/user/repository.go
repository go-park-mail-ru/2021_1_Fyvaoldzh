package user

import "kudago/application/models"

type Repository interface {
	GetUser(login string) (*models.User, bool, error)
}
