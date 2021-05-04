package chat

import (
	"kudago/application/models"
	"time"
)

type Repository interface {
	GetAllDialogues(uid uint64, page int) (models.DialogueCardsSQL, error)
	GetMessages(id uint64, page int) (models.MessagesSQL, error)
	GetEasyDialogue(id uint64) (models.EasyDialogueMessageSQL, error)
	GetEasyMessage(id uint64) (models.EasyDialogueMessageSQL, error)
	DeleteDialogue(id uint64) error
	SendMessage(id uint64, newMessage *models.NewMessage, uid uint64, now time.Time) error
	DeleteMessage(id uint64) error
	EditMessage(id uint64, text string) error
	MessagesSearch(uid uint64, str string, page int) (models.MessagesSQL, error)
	DialogueMessagesSearch(uid uint64, id uint64, str string, page int) (models.MessagesSQL, error)
	CheckDialogue(uid1 uint64, uid2 uint64) (bool, uint64, error)
	CheckMessage(uid1 uint64, uid2 uint64) (bool, uint64, error)
	NewDialogue(uid1 uint64, uid2 uint64) (uint64, error)
	ReadMessages(id uint64, page int, uid uint64) error
}
