package usecase

import (
	"errors"
	"fmt"
	"kudago/application/chat"
	"kudago/application/models"
	"kudago/application/subscription"
	"kudago/application/user"
	"kudago/pkg/logger"
	"strings"
	"time"
)

type Chat struct {
	repo     chat.Repository
	repoSub  subscription.Repository
	repoUser user.Repository
	logger   logger.Logger
}

func NewChat(c chat.Repository, repoSubscription subscription.Repository, repoUser user.Repository, logger logger.Logger) chat.UseCase {
	return &Chat{repo: c, repoSub: repoSubscription, repoUser: repoUser, logger: logger}
}

func (c Chat) ConvertDialogueCard(old models.DialogueCardSQL, uid uint64) (models.DialogueCard, error) {
	var newDialogueCard models.DialogueCard
	newDialogueCard.ID = old.ID
	newDialogueCard.LastMessage = models.ConvertMessageFromCard(old, uid)
	var err error
	if old.User_1 == uid {
		newDialogueCard.Interlocutor, err = c.repoUser.GetUserByID(old.User_2)
	} else {
		newDialogueCard.Interlocutor, err = c.repoUser.GetUserByID(old.User_1)
	}
	if err != nil {
		c.logger.Warn(err)
		return models.DialogueCard{}, err
	}
	return newDialogueCard, nil
}

func (c Chat) ConvertDialogue(old models.DialogueSQL, uid uint64) (models.Dialogue, error) {
	var newDialogue models.Dialogue
	newDialogue.ID = old.ID
	for i := range old.DialogMessages {
		newDialogue.DialogMessages = append(newDialogue.DialogMessages, models.ConvertMessage(old.DialogMessages[i], uid))
	}
	var err error
	if old.User1 == uid {
		newDialogue.Interlocutor, err = c.repoUser.GetUserByID(old.User2)
	} else {
		newDialogue.Interlocutor, err = c.repoUser.GetUserByID(old.User1)
	}
	if err != nil {
		c.logger.Warn(err)
		return models.Dialogue{}, err
	}
	return newDialogue, nil
}

func (c Chat) GetAllDialogues(uid uint64, page int) (models.DialogueCards, error) {
	dialogues, err := c.repo.GetAllDialogues(uid, page)
	if err != nil {
		c.logger.Warn(err)
		return models.DialogueCards{}, err
	}

	var dialogueCards models.DialogueCards

	for i := range dialogues {
		dialogueCard, err := c.ConvertDialogueCard(dialogues[i], uid)
		if err != nil {
			c.logger.Warn(err)
			return models.DialogueCards{}, err
		}
		dialogueCards = append(dialogueCards, dialogueCard)
	}
	if len(dialogueCards) == 0 {
		c.logger.Debug("page" + fmt.Sprint(page) + "is empty")
		return models.DialogueCards{}, nil
	}

	return dialogueCards, nil
}

func (c Chat) GetOneDialogue(uid uint64, id uint64, page int) (models.Dialogue, error) {
	dialogue, err := c.repo.GetEasyDialogue(id)
	if err != nil {
		c.logger.Warn(err)
		return models.Dialogue{}, err
	}

	if dialogue.User1 != uid && dialogue.User2 != uid {
		err := errors.New("you are not interlocutor of this dialogue")
		c.logger.Warn(err)
		return models.Dialogue{}, err
	}

	messages, err := c.repo.GetMessages(id)
	if err != nil {
		c.logger.Warn(err)
		return models.Dialogue{}, err
	}

	//TODO : Вынести в конверт эту штуку
	var dialogueSQL models.DialogueSQL
	dialogueSQL.DialogMessages = messages
	dialogueSQL.ID = dialogue.ID
	dialogueSQL.User1 = dialogue.User1
	dialogueSQL.User2 = dialogue.User2

	resDialogue, err := c.ConvertDialogue(dialogueSQL, uid)
	if err != nil {
		c.logger.Warn(err)
		return models.Dialogue{}, err
	}

	return resDialogue, nil
}

func (c Chat) IsInterlocutorDialogue(uid uint64, id uint64) (bool, error) {
	dialogue, err := c.repo.GetEasyDialogue(id)
	if err != nil {
		c.logger.Warn(err)
		return false, err
	}
	if uid != dialogue.User1 && uid != dialogue.User2 {
		return false, nil
	}
	return true, nil
}

func (c Chat) IsInterlocutorMessage(uid uint64, id uint64) (bool, error) {
	message, err := c.repo.GetEasyMessage(id)
	if err != nil {
		c.logger.Warn(err)
		return false, err
	}
	if uid != message.User1 && uid != message.User2 {
		return false, nil
	}
	return true, nil
}

func (c Chat) DeleteDialogue(uid uint64, id uint64) error {
	isInterlocutor, err := c.IsInterlocutorDialogue(uid, id)
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
	_, err := c.repoUser.GetUserByID(newMessage.To)
	if err != nil {
		return err
	}

	isDialogue, id, err := c.repo.CheckDialogue(uid, newMessage.To)
	if err != nil {
		return err
	}
	if !isDialogue {
		id, err = c.repo.NewDialogue(uid, newMessage.To)
		if err != nil {
			return err
		}
	}

	err = c.repo.SendMessage(id, newMessage, uid, time.Now())
	if err != nil {
		c.logger.Warn(err)
		return err
	}
	return nil
}

func (c Chat) DeleteMessage(uid uint64, id uint64) error {
	isInterlocutor, err := c.IsInterlocutorMessage(uid, id)
	if err != nil {
		c.logger.Warn(err)
		return err
	}
	if isInterlocutor {
		err := c.repo.DeleteMessage(id)
		if err != nil {
			c.logger.Warn(err)
			return err
		}
		return nil
	}
	return errors.New("user is not interlocutor")
}

func (c Chat) EditMessage(uid uint64, id uint64, newMessage *models.NewMessage) error {
	//нужна проверка на то, что исправляемое сообщение ИМЕННО пользователя, а не собеседника
	isInterlocutor, err := c.IsInterlocutorMessage(uid, id)
	if err != nil {
		c.logger.Warn(err)
		return err
	}
	if isInterlocutor {
		err := c.repo.EditMessage(id, newMessage.Text, time.Now())
		if err != nil {
			c.logger.Warn(err)
			return err
		}
		return nil
	}
	return errors.New("user is not interlocutor")
}
func (c Chat) Search(uid uint64, id int, str string, page int) (models.Messages, error) {
	str = strings.ToLower(str)

	var sqlMessages []models.MessageSQL
	var err error
	if id == -1 {
		sqlMessages, err = c.repo.MessagesSearch(uid, str, page)
		if err != nil {
			c.logger.Warn(err)
			return models.Messages{}, err
		}
	} else {
		//TODO: проверка на собеседника!!!
		sqlMessages, err = c.repo.DialogueMessagesSearch(uid, uint64(id), str, page)
		if err != nil {
			c.logger.Warn(err)
			return models.Messages{}, err
		}
	}
	var messages models.Messages

	for i := range sqlMessages {
		messages = append(messages, models.ConvertMessage(sqlMessages[i], uid))
	}
	if len(messages) == 0 {
		c.logger.Debug("page" + fmt.Sprint(page) + "is empty")
		return models.Messages{}, nil
	}

	return messages, nil
}
