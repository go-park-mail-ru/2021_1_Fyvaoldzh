package chat

import (
	"kudago/application/models"
	"time"
)

type Repository interface {
	GetAllDialogues(uid uint64, page int) (models.DialogueCardsSQL, error)
	GetAllNotifications(uid uint64, page int, now time.Time) (models.NotificationsSQL, error)
	GetMessages(id uint64, page int) (models.MessagesSQL, error)
	CheckDialogueID(id uint64) (bool, models.EasyDialogueMessageSQL, error)
	GetEasyMessage(id uint64) (models.EasyDialogueMessageSQL, error)
	DeleteDialogue(id uint64) error
	SendMessage(id uint64, newMessage *models.NewMessage, uid uint64, now time.Time) error
	DeleteMessage(id uint64) error
	EditMessage(id uint64, text string) error
	MessagesSearch(uid uint64, str string, page int) (models.DialogueCardsSQL, error)
	DialogueMessagesSearch(uid uint64, id uint64, str string, page int) (models.DialogueCardsSQL, error)
	CheckDialogueUsers(uid1 uint64, uid2 uint64) (bool, models.EasyDialogueMessageSQL, error)
	CheckMessage(id uint64) (bool, models.EasyDialogueMessageSQL, error)
	NewDialogue(uid1 uint64, uid2 uint64) (uint64, error)
	ReadMessages(id uint64, page int, uid uint64) (int64, error)
	ReadNotifications(uid uint64, page int, now time.Time) error
	AddMailNotification(id uint64, idTo uint64, now time.Time) error
	AddCountNotification(id uint64) error
	SetZeroCountNotifications(id uint64) error
	AddCountMessages(id uint64) error
	DecrementCountMessages(id uint64, count int64) error
}
