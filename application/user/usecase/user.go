package usecase

import (
	"kudago/application/user"
	"kudago/models"
)

type User struct {
	database user.Repository
}

func NewUser(u user.Repository) user.UseCase {
	return &User{database: u}
}

func (u User) Add(usr *models.RegData) (*models.RegData, error) {
	u.database.Add(usr)
	return usr, nil
}

func (u User) Get(id uint) (*models.User, error) {
	panic("implement me")
}

func (u User) Update(user *models.User) error {
	panic("implement me")
}