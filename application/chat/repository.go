package chat

import (
	"kudago/application/models"
)

type Repository interface {
	GetAllDialogues(uid uint64, page int) (models.DialogueCardsSQL, error)
	GetOneDialogue(id uint64, page int) (models.DialogueSQL, error)
	GetEasyDialogue(id uint64) (models.EasyDialogueMessageSQL, error)
	GetEasyMessage(id uint64) (models.EasyDialogueMessageSQL, error)
	DeleteDialogue(id uint64) error
	SendMessage(newMessage *models.NewMessage, uid uint64) error
	DeleteMessage(id uint64) error
	EditMessage(id uint64) error
	FindFollowers(str string, uid uint64) (models.UsersOnEvent, error)
	MessagesSearch(uid uint64, str string, page int) (models.MessagesSQL, error)
	DialogueMessagesSearch(uid uint64, id uint64, str string, page int) (models.MessagesSQL, error)
}
