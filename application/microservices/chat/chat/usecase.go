package chat

import "kudago/application/models"

//go:generate mockgen -destination=./mock/usecase_mock.go -package=mock -source=./application/event/usecase.go

type UseCase interface {
	IsInterlocutor(uid uint64, elem models.EasyDialogueMessageSQL) bool
	IsSenderMessage(uid uint64, elem models.EasyDialogueMessageSQL) bool
	GetAllDialogues(uid uint64, page int) (models.DialogueCards, error)
	GetAllNotifications(uid uint64, page int) (models.Notifications, error)
	GetOneDialogue(uid uint64, id uint64, page int) (models.Dialogue, error)
	DeleteDialogue(uid uint64, id uint64) error
	SendMessage(newMessage *models.NewMessage, uid uint64) error
	DeleteMessage(uid uint64, id uint64) error
	EditMessage(uid uint64, newMessage *models.RedactMessage) error
	AutoMailingConstructor(to uint64, from, eventName, eventID string) models.NewMessage
	Mailing(uid uint64, mailing *models.Mailing) error
	Search(uid uint64, id int, str string, page int) (models.DialogueCards, error)
}
