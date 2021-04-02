package user

import (
	"kudago/application/models"
	"mime/multipart"
)

type UseCase interface {
	Add(usr *models.RegData) (uint64, error)
	GetOwnProfile(id uint64) (*models.UserOwnProfile, error)
	GetOtherProfile(id uint64) (*models.OtherUserProfile, error)
	Update(uid uint64, user *models.UserOwnProfile) error
	CheckUser(user *models.User) (uint64, error)
	UploadAvatar(uid uint64, img *multipart.FileHeader) error
	GetAvatar(uid uint64) ([]byte, error)
}
