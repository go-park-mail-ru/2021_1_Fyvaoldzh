package usecase

import (
	"errors"
	"fmt"
	"kudago/application/chat"
	"kudago/application/event"
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
			return nil, err
		}
		dialogueCards = append(dialogueCards, models.ConvertDialogueCard(dialogues[i], uid, interlocutor))
	}
	if len(dialogueCards) == 0 {
		c.logger.Debug("page" + fmt.Sprint(page) + "is empty")
		return models.DialogueCards{}, nil
	}

	return dialogueCards, nil
}

//TODO: ОБЯЗАТЕЛЬНО помечать сообщения read!!!!!!!
//Как сделать так, чтобы сообщения выдавались непрочитанные сначала? Ну типа чтоб как в вк, открываешь диалог и тебя
//скролит к последнему сообщению прочитанному, а дальше ты листаешь вниз и читаешь типа
func (c Chat) GetOneDialogue(uid uint64, id uint64, page int) (models.Dialogue, error) {
	_, err := c.repoUser.GetUserByID(id)
	if err != nil {
		return models.Dialogue{}, err
	}
	interlocutor, err := c.repoUser.GetUserByID(id)
	if err != nil {
		c.logger.Warn(err)
		return models.Dialogue{}, err
	}

	isDialogue, d_id, err := c.repo.CheckDialogue(uid, id)
	if err != nil {
		return models.Dialogue{}, err
	}
	if !isDialogue {
		return models.Dialogue{Interlocutor: interlocutor, DialogMessages: models.Messages{}}, nil
	}

	dialogue, err := c.repo.GetEasyDialogue(d_id)
	if err != nil {
		c.logger.Warn(err)
		return models.Dialogue{}, err
	}

	if dialogue.User1 != uid && dialogue.User2 != uid {
		err := errors.New("you are not interlocutor of this dialogue")
		c.logger.Warn(err)
		return models.Dialogue{}, err
	}

	messages, err := c.repo.GetMessages(id, page)
	if err != nil {
		c.logger.Warn(err)
		return models.Dialogue{}, err
	}

	resDialogue := models.ConvertDialogue(dialogue, messages, uid, interlocutor)

	return resDialogue, nil
}

func (c Chat) IsInterlocutor(uid uint64, id uint64, f func(id uint64) (models.EasyDialogueMessageSQL, error)) (bool, error) {
	dialogue, err := f(id)
	if err != nil {
		c.logger.Warn(err)
		return false, err
	}
	if uid != dialogue.User1 && uid != dialogue.User2 {
		return false, nil
	}
	return true, nil
}

func (c Chat) IsSenderMessage(uid uint64, id uint64) (bool, error) {
	message, err := c.repo.GetEasyMessage(id)
	if err != nil {
		c.logger.Warn(err)
		return false, err
	}
	if uid != message.User1 {
		return false, nil
	}
	return true, nil
}

//Не совсем оптимально получается, дважды смотрим один и тот же диалог, сначала на предмет существования, потом на то, являемся ли собеседником,
//можно убрать первую проверку, но тогда будет возвращаться ошибка "0 request error". Далее у функци
func (c Chat) DeleteDialogue(uid uint64, id uint64) error {
	isDialogue, _, err := c.repo.CheckDialogue(uid, id)
	if err != nil {
		return err
	}
	if !isDialogue {
		return errors.New("no dialogue with this id")
	}

	isInterlocutor, err := c.IsInterlocutor(uid, id, c.repo.GetEasyDialogue)
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

//Тут так же 2 запроса, просто устал уже думать(
func (c Chat) DeleteMessage(uid uint64, id uint64) error {
	isMessage, _, err := c.repo.CheckMessage(uid, id)
	if err != nil {
		return err
	}
	if !isMessage {
		return errors.New("no dialogue with this id")
	}

	isInterlocutor, err := c.IsInterlocutor(uid, id, c.repo.GetEasyMessage)
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

//Тут вообще кастомную функцию для проверки существования сообщения по его id надо делать
func (c Chat) EditMessage(uid uint64, newMessage *models.RedactMessage) error {
	/*isMessage, _, err := c.repo.CustomCheckMessage(newMessage.ID)
	if err != nil {
		return err
	}
	if !isMessage {
		return errors.New("no dialogue with this id")
	}*/

	isInterlocutor, err := c.IsSenderMessage(uid, newMessage.ID)
	if err != nil {
		c.logger.Warn(err)
		return err
	}
	if isInterlocutor {
		err := c.repo.EditMessage(newMessage.ID, newMessage.Text)
		if err != nil {
			c.logger.Warn(err)
			return err
		}
		return nil
	}
	return errors.New("user is not interlocutor")
}

func (c Chat) AutoMailingConstructor(to uint64, from, eventName, eventID string) models.NewMessage {
	var mailingMessage models.NewMessage
	mailingMessage.To = to
	mailingMessage.Text = from + constants.MailingText + `"` + eventName + `" ` + constants.MailingAddress + eventID

	return mailingMessage
}

func (c Chat) Mailing(uid uint64, mailing *models.Mailing) error {
	sender, err := c.repoUser.GetUserByID(uid)
	if err != nil {
		return err
	}
	event, err := c.repoEvent.GetOneEventByID(mailing.EventID)
	if err != nil {
		return err
	}

	for _, id := range mailing.To {
		message := c.AutoMailingConstructor(id, sender.Name, event.Title, fmt.Sprint(event.ID))
		err := c.SendMessage(&message, uid)
		if err != nil {
			return err
		}
	}
	return nil
}

//И вот тут 2 раза
func (c Chat) Search(uid uint64, id int, str string, page int) (models.Messages, error) {
	str = strings.ToLower(str)

	var sqlMessages []models.MessageSQL
	var err error
	if id == 0 {
		sqlMessages, err = c.repo.MessagesSearch(uid, str, page)
		if err != nil {
			c.logger.Warn(err)
			return nil, err
		}
	} else {
		isDialogue, _, err := c.repo.CheckMessage(uid, uint64(id))
		if err != nil {
			return nil, err
		}
		if !isDialogue {
			return models.Messages{}, errors.New("no dialogue with this id")
		}

		isInterlocutor, err := c.IsInterlocutor(uid, uint64(id), c.repo.GetEasyDialogue)
		if err != nil {
			c.logger.Warn(err)
			return nil, err
		}
		if isInterlocutor {
			sqlMessages, err = c.repo.DialogueMessagesSearch(uid, uint64(id), str, page)
			if err != nil {
				c.logger.Warn(err)
				return nil, err
			}
		} else {
			return nil, errors.New("user is not interlocutor")
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
