package user

import "kudago/models"

type UseCase interface {
	Add(usr *models.RegData) (uint64, error)
	GetOwnProfile(id uint64) (*models.UserOwnProfile, error)
	GetOtherProfile(id uint64) (*models.OtherUserProfile, error)
	Update(uid uint64, user *models.UserOwnProfile) error
	Login(user *models.User) (uint64, error)
}