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

//Функции конвертов хотелось бы перенести в модели, но обязательно в них нужно использовать методы userRepository, как быть?
//Просто сейчас часть конвертов тут, часть в модельках
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

func (c Chat) ConvertDialogue(dialogue models.EasyDialogueMessageSQL, messages models.MessagesSQL, uid uint64) (models.Dialogue, error) {
	var newDialogue models.Dialogue
	newDialogue.ID = dialogue.ID
	for i := range messages {
		newDialogue.DialogMessages = append(newDialogue.DialogMessages, models.ConvertMessage(messages[i], uid))
	}
	var err error
	if dialogue.User1 == uid {
		newDialogue.Interlocutor, err = c.repoUser.GetUserByID(dialogue.User2)
	} else {
		newDialogue.Interlocutor, err = c.repoUser.GetUserByID(dialogue.User1)
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

//ОБЯЗАТЕЛЬНО помечать сообщения read
func (c Chat) GetOneDialogue(uid uint64, id uint64, page int) (models.Dialogue, error) {
	//Проверить, существует ли такой диалог вообще, иначе падает поросто реквест еррор с 0 ошибкой
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

	resDialogue, err := c.ConvertDialogue(dialogue, messages, uid)
	if err != nil {
		c.logger.Warn(err)
		return models.Dialogue{}, err
	}

	return resDialogue, nil
}

//Очень похожие функции, дублирование кода, мб есть смысл как-то вынести, хотя тут не уверен, что возможно
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

func (c Chat) DeleteDialogue(uid uint64, id uint64) error {
	//Проверить, существует ли такой диалог вообще, иначе падает поросто реквест еррор с 0 ошибкой
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
	//Проверить, существует ли такое сообщение вообще, иначе падает поросто реквест еррор с 0 ошибкой
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

func (c Chat) EditMessage(uid uint64, newMessage *models.RedactMessage) error {
	//Проверить, существует ли такое сообщение вообще, иначе падает поросто реквест еррор с 0 ошибкой
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
func (c Chat) Search(uid uint64, id int, str string, page int) (models.Messages, error) {
	str = strings.ToLower(str)

	var sqlMessages []models.MessageSQL
	var err error
	if id == -1 {
		sqlMessages, err = c.repo.MessagesSearch(uid, str, page)
		if err != nil {
			c.logger.Warn(err)
			return nil, err
		}
	} else {
		//Проверить, существует ли такой диалог вообще, иначе падает поросто реквест еррор с 0 ошибкой
		isInterlocutor, err := c.IsInterlocutorDialogue(uid, uint64(id))
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
