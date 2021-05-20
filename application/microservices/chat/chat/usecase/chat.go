package usecase

import (
	"errors"
	"fmt"
	"kudago/application/event"
	"kudago/application/microservices/chat/chat"
	"kudago/application/models"
	"kudago/application/subscription"
	"kudago/application/user"
	"kudago/pkg/constants"
	"kudago/pkg/logger"
	"strings"
	"time"
)

type Chat struct {
	repo      chat.Repository
	repoSub   subscription.Repository
	repoUser  user.Repository
	repoEvent event.Repository
	logger    logger.Logger
}

func NewChat(c chat.Repository, repoSubscription subscription.Repository,
	repoUser user.Repository, repoEvent event.Repository, logger logger.Logger) chat.UseCase {
	return &Chat{repo: c, repoSub: repoSubscription, repoUser: repoUser, repoEvent: repoEvent, logger: logger}
}

func (c Chat) GetAllDialogues(uid uint64, page int) (models.DialogueCards, error) {
	dialogues, err := c.repo.GetAllDialogues(uid, page)
	if err != nil {
		c.logger.Warn(err)
		return models.DialogueCards{}, err
	}

	var dialogueCards models.DialogueCards

	for i := range dialogues {
		var interlocutor models.UserOnEvent
		if dialogues[i].User1 == uid {
			interlocutor, err = c.repoUser.GetUserByID(dialogues[i].User2)
		} else {
			interlocutor, err = c.repoUser.GetUserByID(dialogues[i].User1)
		}
		if err != nil {
			c.logger.Warn(err)
			return models.DialogueCards{}, err
		}
		dialogueCards = append(dialogueCards, models.ConvertDialogueCard(dialogues[i], uid, interlocutor))
	}
	if len(dialogueCards) == 0 {
		c.logger.Debug("page" + fmt.Sprint(page) + "is empty")
		return models.DialogueCards{}, nil
	}

	return dialogueCards, nil
}

func (c Chat) GetOneDialogue(uid1 uint64, uid2 uint64, page int) (models.Dialogue, error) {
	interlocutor, err := c.repoUser.GetUserByID(uid2)
	if err != nil {
		c.logger.Warn(err)
		return models.Dialogue{}, err
	}

	isDialogue, dialogue, err := c.repo.CheckDialogueUsers(uid1, uid2)
	if err != nil {
		return models.Dialogue{}, err
	}
	if !isDialogue {
		return models.Dialogue{Interlocutor: interlocutor, DialogMessages: models.Messages{}}, nil
	}

	messages, err := c.repo.GetMessages(dialogue.ID, page)
	if err != nil {
		c.logger.Warn(err)
		return models.Dialogue{}, err
	}

	resDialogue := models.ConvertDialogue(dialogue, messages, uid1, interlocutor)

	err = c.repo.ReadMessages(dialogue.ID, page, uid1)
	if err != nil {
		c.logger.Warn(err)
		return resDialogue, nil
	}

	return resDialogue, nil
}

func (c Chat) IsInterlocutor(uid uint64, elem models.EasyDialogueMessageSQL) bool {
	if uid != elem.User1 && uid != elem.User2 {
		return false
	}
	return true
}

func (c Chat) IsSenderMessage(uid uint64, elem models.EasyDialogueMessageSQL) bool {
	if uid != elem.User1 {
		return false
	}
	return true
}

func (c Chat) DeleteDialogue(uid uint64, id uint64) error {
	isDialogue, dialogue, err := c.repo.CheckDialogueID(id)
	if err != nil {
		return err
	}
	if !isDialogue {
		return errors.New("no dialogue with this id")
	}

	isInterlocutor := c.IsInterlocutor(uid, dialogue)
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

	isDialogue, dialogue, err := c.repo.CheckDialogueUsers(uid, newMessage.To)
	if err != nil {
		return err
	}
	if !isDialogue {
		dialogue.ID, err = c.repo.NewDialogue(uid, newMessage.To)
		if err != nil {
			return err
		}
	}

	err = c.repo.SendMessage(dialogue.ID, newMessage, uid, time.Now())
	if err != nil {
		c.logger.Warn(err)
		return err
	}
	return nil
}

func (c Chat) DeleteMessage(uid uint64, id uint64) error {
	isMessage, message, err := c.repo.CheckMessage(id)
	if err != nil {
		return err
	}
	if !isMessage {
		return errors.New("no dialogue with this id")
	}

	isInterlocutor := c.IsInterlocutor(uid, message)
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

//Неплохо бы проверять, как давно было написано сообщение, чтобы нельзя было менять сообщения, написанные пару часов назад
func (c Chat) EditMessage(uid uint64, newMessage *models.RedactMessage) error {
	isMessage, message, err := c.repo.CheckMessage(newMessage.ID)
	if err != nil {
		return err
	}
	if !isMessage {
		return errors.New("no message with this id")
	}

	isSender := c.IsSenderMessage(uid, message)
	if !isSender {
		return errors.New("user is not sender")
	}
	err = c.repo.EditMessage(newMessage.ID, newMessage.Text)
	if err != nil {
		c.logger.Warn(err)
		return err
	}

	return nil
}

func (c Chat) AutoMailingConstructor(to uint64, from, eventName, eventID string) models.NewMessage {
	var mailingMessage models.NewMessage
	mailingMessage.To = to
	mailingMessage.Text = from + constants.MailingText + eventName + " " + constants.MailingAddress + eventID

	return mailingMessage
}

func (c Chat) Mailing(uid uint64, mailing *models.Mailing) error {
	sender, err := c.repoUser.GetUserByID(uid)
	if err != nil {
		return err
	}
	ev, err := c.repoEvent.GetOneEventByID(mailing.EventID)
	if err != nil {
		return err
	}

	for _, id := range mailing.To {
		message := c.AutoMailingConstructor(id, sender.Name, ev.Title, fmt.Sprint(ev.ID))
		err := c.SendMessage(&message, uid)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c Chat) Search(uid uint64, id int, str string, page int) (models.Messages, error) {
	str = strings.ToLower(str)

	var sqlMessages []models.MessageSQL
	var err error
	if id == 0 {
		sqlMessages, err = c.repo.MessagesSearch(uid, str, page)
		if err != nil {
			c.logger.Warn(err)
			return models.Messages{}, err
		}
	} else {
		isDialogue, dialogue, err := c.repo.CheckDialogueID(uint64(id))
		if err != nil {
			return models.Messages{}, err
		}
		if !isDialogue {
			return models.Messages{}, errors.New("no dialogue with this id")
		}

		isInterlocutor := c.IsInterlocutor(uid, dialogue)
		if isInterlocutor {
			sqlMessages, err = c.repo.DialogueMessagesSearch(uid, uint64(id), str, page)
			if err != nil {
				c.logger.Warn(err)
				return models.Messages{}, err
			}
		} else {
			return models.Messages{}, errors.New("user is not interlocutor")
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
