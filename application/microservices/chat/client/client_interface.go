package client

import "kudago/application/models"

type IChatClient interface {
	GetAllDialogues(uid uint64, page int) (models.DialogueCards, error)
	GetOneDialogue(uid1 uint64, uid2 uint64, page int) (models.Dialogue, error)
	DeleteDialogue(uid uint64, id uint64) error
	SendMessage(newMessage *models.NewMessage, uid uint64) error
	EditMessage(uid uint64, newMessage *models.RedactMessage) error
	DeleteMessage(uid uint64, id uint64) error
	Mailing(uid uint64, mailing *models.Mailing) error
	Search(uid uint64, id int, str string, page int) (models.Messages, error)
	Close()
}