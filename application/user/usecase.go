package user

import (
	"kudago/application/models"
	"mime/multipart"
)

type UseCase interface {
	Add(usr *models.RegData) (uint64, error)
	GetOwnProfile(id uint64) (*models.UserOwnProfile, error)
	GetOtherProfile(id uint64) (*models.OtherUserProfile, error)
	Update(id uint64, user *models.UserOwnProfile) error
	Login(usr *models.User) (uint64, error)
	CheckUser(usr *models.User) (uint64, error)
	UploadAvatar(id uint64, img multipart.File, filename string) error
	GetAvatar(id uint64) ([]byte, error)
	GetUsers(page int) (models.UserCards, error)
	UpdateEventStatus(userId uint64, eventId uint64) error
	IsAddedEvent(userId uint64, eventId uint64) (bool, error)
}
