package chat

import "kudago/application/models"

//go:generate mockgen -destination=./mock/usecase_mock.go -package=mock -source=./application/event/usecase.go

type UseCase interface {
	GetAllDialogues(uid uint64, page int) (models.DialogueCards, error)
	GetOneDialogue(uid uint64, id uint64, page int) (models.Dialogue, error)
	DeleteDialogue(uid uint64, id uint64) error
	SendMessage(newMessage *models.NewMessage) error
	DeleteMessage(uid uint64, id uint64) error
	EditMessage()
}
