package usecase

import (
	"errors"
	"fmt"
	"kudago/application/chat"
	"kudago/application/models"
	"kudago/application/subscription"
	"kudago/pkg/logger"
	"time"
)

type Chat struct {
	repo    chat.Repository
	repoSub subscription.Repository
	logger  logger.Logger
}

func NewChat(c chat.Repository, repoSubscription subscription.Repository, logger logger.Logger) chat.UseCase {
	return &Chat{repo: c, repoSub: repoSubscription, logger: logger}
}

func (c Chat) GetAllDialogues(uid uint64, page int) (models.DialogueCards, error) {
	dialogues, err := c.repo.GetAllDialogues(time.Now(), uid, page) //return dialogueCardsSQL
	if err != nil {
		c.logger.Warn(err)
		return models.DialogueCards{}, err
	}

	var dialogueCards models.DialogueCards

	for i := range dialogues {
		dialogueCards = append(dialogueCards, models.ConvertDialogueCard(dialogues[i], uid))
	}
	if len(dialogueCards) == 0 {
		c.logger.Debug("page" + fmt.Sprint(page) + "is empty")
		return models.DialogueCards{}, nil
	}

	return dialogueCards, nil
}

func (c Chat) GetOneDialogue(uid uint64, id uint64, page int) (models.Dialogue, error) {
	dialogueSQL, err := c.repo.GetOneDialogue(id, page) //return dialogueSQL
	if err != nil {
		c.logger.Warn(err)
		return models.Dialogue{}, err
	}

	if dialogueSQL.User1 != uid && dialogueSQL.User2 != uid {
		err := errors.New("you are not interlocutor of this dialogue")
		c.logger.Warn(err)
		return models.Dialogue{}, err
	}

	return models.ConvertDialogue(dialogueSQL, uid), nil
}

func (c Chat) DeleteDialogue(uid uint64, id uint64) error {
	isInterlocutor, err := c.repo.IsInterlocutor(uid, id)
	if err != nil {
		c.logger.Warn(err)
		return err
	}
	if isInterlocutor {
		err := c.repo.DeleteDialogue(id)
		if err != nil {
			c.logger.Warn(err)
			return err
		}
		return nil
	}
	return errors.New("user is not interlocutor")
}

func (c Chat) SendMessage(newMessage *models.NewMessage, uid uint64) error {

}
