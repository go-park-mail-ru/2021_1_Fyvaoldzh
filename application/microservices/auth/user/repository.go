package user

import "kudago/application/models"

type Repository interface {
	IsCorrect(login string) (*models.User, error)
}
