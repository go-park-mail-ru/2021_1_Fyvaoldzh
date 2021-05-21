package client

import "kudago/application/models"

type IChatClient interface {
	GetAllDialogues(uid uint64, page int) (models.DialogueCards, error, int)
	GetOneDialogue(uid1 uint64, uid2 uint64, page int) (models.Dialogue, error, int)
	DeleteDialogue(uid uint64, id uint64) (error, int)
	SendMessage(newMessage *models.NewMessage, uid uint64) (error, int)
	EditMessage(uid uint64, newMessage *models.RedactMessage) (error, int)
	DeleteMessage(uid uint64, id uint64) (error, int)
	Mailing(uid uint64, mailing *models.Mailing) (error, int)
	Search(uid uint64, id int, str string, page int) (models.DialogueCards, error, int)
	GetNotifications(uid uint64, page int) (models.Notifications, error, int)
	Close()
}
