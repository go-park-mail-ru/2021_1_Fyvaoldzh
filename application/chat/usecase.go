package chat

import "kudago/application/models"

//go:generate mockgen -destination=./mock/usecase_mock.go -package=mock -source=./application/event/usecase.go

type UseCase interface {
	ConvertDialogueCard(old models.DialogueCardSQL, uid uint64) (models.DialogueCard, error)
	ConvertDialogue(dialogue models.EasyDialogueMessageSQL, messages models.MessagesSQL, uid uint64) (models.Dialogue, error)
	IsInterlocutor(uid uint64, id uint64, f func(id uint64) (models.EasyDialogueMessageSQL, error)) (bool, error)
	IsSenderMessage(uid uint64, id uint64) (bool, error)
	GetAllDialogues(uid uint64, page int) (models.DialogueCards, error)
	GetOneDialogue(uid uint64, id uint64, page int) (models.Dialogue, error)
	DeleteDialogue(uid uint64, id uint64) error
	SendMessage(newMessage *models.NewMessage, uid uint64) error
	DeleteMessage(uid uint64, id uint64) error
	EditMessage(uid uint64, newMessage *models.RedactMessage) error
	Search(uid uint64, id int, str string, page int) (models.Messages, error)
}
