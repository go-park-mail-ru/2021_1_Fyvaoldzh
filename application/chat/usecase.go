package chat

import "kudago/application/models"

//go:generate mockgen -destination=./mock/usecase_mock.go -package=mock -source=./application/event/usecase.go

type UseCase interface {
	ConvertDialogueCard(old models.DialogueCardSQL, uid uint64) (models.DialogueCard, error)
	ConvertDialogue(old models.DialogueSQL, uid uint64) (models.Dialogue, error)
	GetAllDialogues(uid uint64, page int) (models.DialogueCards, error)
	GetOneDialogue(uid uint64, id uint64, page int) (models.Dialogue, error)
	DeleteDialogue(uid uint64, id uint64) error
	SendMessage(newMessage *models.NewMessage, uid uint64) error
	DeleteMessage(uid uint64, id uint64) error
	EditMessage(uid uint64, id uint64, newMessage *models.NewMessage) error
	Search(uid uint64, id int, str string, page int) (models.Messages, error)
}
