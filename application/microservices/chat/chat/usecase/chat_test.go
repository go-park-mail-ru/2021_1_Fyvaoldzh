package usecase

import (
	"database/sql"
	"errors"
	"fmt"
	mock_event "kudago/application/event/mocks"
	"kudago/application/microservices/chat/chat"
	mock_chat "kudago/application/microservices/chat/chat/mocks"
	"kudago/application/models"
	mock_subscription "kudago/application/subscription/mocks"
	mock_user "kudago/application/user/mocks"
	"kudago/pkg/logger"
	"log"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	test_id        = 1
	test_page      = 1
	test_id2       = 2
	test_wrong_id  = 10
	test_text      = "test text"
	test_text_2    = "test text 2"
	test_from      = "test from"
	test_event     = "test event"
	test_mail_text = `test from приглашает Вас на мероприятие "test event" qdaqda.ru:3000/event1`
)

var testDialogueMessageSQL = models.EasyDialogueMessageSQL{
	ID:    uint64(test_id),
	User1: uint64(test_id),
	User2: uint64(test_id2),
}

var testMessages = models.Messages{
	{
		ID:     uint64(test_id),
		FromMe: true,
		Text:   test_text,
		Date:   testMessagesSQL[0].Date.String(),
		Redact: false,
		Read:   false,
	},
	{
		ID:     uint64(test_id2),
		FromMe: false,
		Text:   test_text_2,
		Date:   testMessagesSQL[1].Date.String(),
		Redact: false,
		Read:   false,
	},
}

var testEventSQL = models.EventSQL{
	ID:          1,
	Title:       test_event,
	Place:       test_text,
	Description: test_text,
	StartDate:   sql.NullTime{Time: time.Now(), Valid: true},
	EndDate:     sql.NullTime{Time: time.Now().Add(15000 * time.Hour), Valid: true},
	Subway:      sql.NullString{String: test_text, Valid: true},
	Street:      sql.NullString{String: test_text, Valid: true},
	Category:    test_text,
	Image:       sql.NullString{String: test_text, Valid: true},
}

var testMail = models.NewMessage{
	To:   uint64(test_id2),
	Text: test_mail_text,
}

var testMailing = models.Mailing{
	EventID: uint64(test_id),
	To:      []uint64{uint64(test_id2)},
}

var testRedactMessage = models.RedactMessage{
	ID:   uint64(test_id),
	Text: test_text,
}

var testNewMessage = models.NewMessage{
	To:   uint64(test_id2),
	Text: test_text,
}

var testMessagesSQL = models.MessagesSQL{
	{
		ID:     uint64(test_id),
		From:   uint64(test_id),
		To:     uint64(test_id2),
		Text:   test_text,
		Date:   time.Now(),
		Redact: false,
		Read:   false,
	},
	{
		ID:     uint64(test_id2),
		From:   uint64(test_id2),
		To:     uint64(test_id),
		Text:   test_text_2,
		Date:   time.Now(),
		Redact: false,
		Read:   false,
	},
}

var testDialogue = models.Dialogue{
	ID:           uint64(test_id),
	Interlocutor: testUser,
	DialogMessages: models.Messages{
		{
			ID:     uint64(test_id),
			FromMe: true,
			Text:   test_text,
			Date:   testMessagesSQL[0].Date.String(),
			Redact: false,
			Read:   false,
		},
		{
			ID:     uint64(test_id2),
			FromMe: false,
			Text:   test_text_2,
			Date:   testMessagesSQL[1].Date.String(),
			Redact: false,
			Read:   false,
		},
	},
}

var testAllDialoguesSQL = models.DialogueCardsSQL{
	{
		ID:     uint64(test_id),
		User1:  uint64(test_id),
		User2:  uint64(test_id2),
		IDMes:  uint64(test_id),
		From:   uint64(test_id),
		To:     uint64(test_id2),
		Text:   test_text,
		Date:   time.Now(),
		Redact: false,
		Read:   false,
	},
	{
		ID:     uint64(test_id),
		User1:  uint64(test_id),
		User2:  uint64(test_id2),
		IDMes:  uint64(test_id2),
		From:   uint64(test_id2),
		To:     uint64(test_id),
		Text:   test_text_2,
		Date:   time.Now(),
		Redact: false,
		Read:   false,
	},
}

var testUser = models.UserOnEvent{
	Id:     uint64(test_id2),
	Name:   test_text,
	Avatar: test_text_2,
}

var testUserMail = models.UserOnEvent{
	Id:     uint64(test_id2),
	Name:   test_from,
	Avatar: test_text_2,
}

var testAllDialogues = models.DialogueCards{
	{
		ID:           uint64(test_id),
		Interlocutor: testUser,
		LastMessage: models.Message{
			ID:     uint64(test_id),
			FromMe: true,
			Text:   test_text,
			Date:   testAllDialoguesSQL[0].Date.String(),
			Redact: false,
			Read:   false,
		},
	},
	{
		ID:           uint64(test_id),
		Interlocutor: testUser,
		LastMessage: models.Message{
			ID:     uint64(test_id2),
			FromMe: false,
			Text:   test_text_2,
			Date:   testAllDialoguesSQL[1].Date.String(),
			Redact: false,
			Read:   false,
		},
	},
}

func setUp(t *testing.T) (*mock_chat.MockRepository, *mock_subscription.MockRepository,
	*mock_user.MockRepository, *mock_event.MockRepository, chat.UseCase) {
	ctrl := gomock.NewController(t)

	rep := mock_chat.NewMockRepository(ctrl)
	repSub := mock_subscription.NewMockRepository(ctrl)
	repUser := mock_user.NewMockRepository(ctrl)
	repEvent := mock_event.NewMockRepository(ctrl)

	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	sugar := l.Sugar()
	zap.NewAtomicLevelAt(zapcore.DebugLevel)

	uc := NewChat(rep, repSub, repUser, repEvent, logger.NewLogger(sugar))
	return rep, repSub, repUser, repEvent, uc
}

func TestChatUseCase_GetAllDialoguesOk(t *testing.T) {
	rep, _, repUser, _, uc := setUp(t)
	rep.EXPECT().GetAllDialogues(uint64(test_id), test_page).Return(testAllDialoguesSQL, nil)
	repUser.EXPECT().GetUserByID(uint64(test_id2)).Return(testUser, nil).Times(2)

	dcards, err := uc.GetAllDialogues(uint64(test_id), test_page)
	assert.Nil(t, err)
	assert.Equal(t, dcards, testAllDialogues)
}

func TestChatUseCase_GetAllDialoguesErrorRepo(t *testing.T) {
	rep, _, _, _, uc := setUp(t)
	rep.EXPECT().GetAllDialogues(uint64(test_id), -1).Return(testAllDialoguesSQL, errors.New("invalid page number"))

	dcards, err := uc.GetAllDialogues(uint64(test_id), -1)
	assert.NotNil(t, err)
	assert.Equal(t, dcards, models.DialogueCards{})
}

func TestChatUseCase_GetAllDialoguesErrorRepoUser(t *testing.T) {
	rep, _, repUser, _, uc := setUp(t)
	rep.EXPECT().GetAllDialogues(uint64(test_id), test_page).Return(testAllDialoguesSQL, nil)
	repUser.EXPECT().GetUserByID(uint64(test_id2)).Return(models.UserOnEvent{}, errors.New("invalid id"))

	dcards, err := uc.GetAllDialogues(uint64(test_id), test_page)
	assert.NotNil(t, err)
	assert.Equal(t, dcards, models.DialogueCards{})
}

func TestChatUseCase_GetAllDialoguesZeroLength(t *testing.T) {
	rep, _, _, _, uc := setUp(t)
	rep.EXPECT().GetAllDialogues(uint64(test_id), test_page).Return(models.DialogueCardsSQL{}, nil)

	dcards, err := uc.GetAllDialogues(uint64(test_id), test_page)
	assert.Nil(t, err)
	assert.Equal(t, dcards, models.DialogueCards{})
}

func TestChatUseCase_GetOneDialogueOk(t *testing.T) {
	rep, _, repUser, _, uc := setUp(t)
	repUser.EXPECT().GetUserByID(uint64(test_id2)).Return(testUser, nil)
	rep.EXPECT().CheckDialogueUsers(uint64(test_id), uint64(test_id2)).Return(true, testDialogueMessageSQL, nil)
	rep.EXPECT().GetMessages(uint64(test_id), test_page).Return(testMessagesSQL, nil)
	rep.EXPECT().ReadMessages(uint64(test_id), test_page, uint64(test_id)).Return(nil)

	dialogues, err := uc.GetOneDialogue(uint64(test_id), uint64(test_id2), test_page)
	assert.Nil(t, err)
	assert.Equal(t, dialogues, testDialogue)
}

func TestChatUseCase_GetOneDialogueErrorRepoUser(t *testing.T) {
	_, _, repUser, _, uc := setUp(t)
	repUser.EXPECT().GetUserByID(uint64(test_id2)).Return(models.UserOnEvent{}, errors.New("invalid id"))

	dialogues, err := uc.GetOneDialogue(uint64(test_id), uint64(test_id2), test_page)
	assert.NotNil(t, err)
	assert.Equal(t, dialogues, models.Dialogue{})
}

func TestChatUseCase_GetOneDialogueErrorCheck(t *testing.T) {
	rep, _, repUser, _, uc := setUp(t)
	repUser.EXPECT().GetUserByID(uint64(test_id2)).Return(testUser, nil)
	rep.EXPECT().CheckDialogueUsers(uint64(test_id), uint64(test_id2)).Return(false, models.EasyDialogueMessageSQL{}, errors.New("invalid id"))

	dialogues, err := uc.GetOneDialogue(uint64(test_id), uint64(test_id2), test_page)
	assert.NotNil(t, err)
	assert.Equal(t, dialogues, models.Dialogue{})
}

func TestChatUseCase_GetOneDialogueNoDialogue(t *testing.T) {
	rep, _, repUser, _, uc := setUp(t)
	repUser.EXPECT().GetUserByID(uint64(test_id2)).Return(testUser, nil)
	rep.EXPECT().CheckDialogueUsers(uint64(test_id), uint64(test_id2)).Return(false, models.EasyDialogueMessageSQL{}, nil)

	dialogues, err := uc.GetOneDialogue(uint64(test_id), uint64(test_id2), test_page)
	assert.Nil(t, err)
	assert.Equal(t, dialogues, models.Dialogue{Interlocutor: testUser, DialogMessages: models.Messages{}})
}

func TestChatUseCase_GetOneDialogueMessagesError(t *testing.T) {
	rep, _, repUser, _, uc := setUp(t)
	repUser.EXPECT().GetUserByID(uint64(test_id2)).Return(testUser, nil)
	rep.EXPECT().CheckDialogueUsers(uint64(test_id), uint64(test_id2)).Return(true, testDialogueMessageSQL, nil)
	rep.EXPECT().GetMessages(uint64(test_id), test_page).Return(models.MessagesSQL{}, errors.New("invalid id"))

	dialogues, err := uc.GetOneDialogue(uint64(test_id), uint64(test_id2), test_page)
	assert.NotNil(t, err)
	assert.Equal(t, dialogues, models.Dialogue{})
}

func TestChatUseCase_GetOneDialogueReadError(t *testing.T) {
	rep, _, repUser, _, uc := setUp(t)
	repUser.EXPECT().GetUserByID(uint64(test_id2)).Return(testUser, nil)
	rep.EXPECT().CheckDialogueUsers(uint64(test_id), uint64(test_id2)).Return(true, testDialogueMessageSQL, nil)
	rep.EXPECT().GetMessages(uint64(test_id), test_page).Return(testMessagesSQL, nil)
	rep.EXPECT().ReadMessages(uint64(test_id), test_page, uint64(test_id)).Return(errors.New("smthing wrong"))

	dialogues, err := uc.GetOneDialogue(uint64(test_id), uint64(test_id2), test_page)
	assert.Nil(t, err)
	assert.Equal(t, dialogues, testDialogue)
}

func TestChatUseCase_IsInterlocutorTrue(t *testing.T) {
	_, _, _, _, uc := setUp(t)

	is := uc.IsInterlocutor(uint64(test_id), testDialogueMessageSQL)
	assert.Equal(t, is, true)
}

func TestChatUseCase_IsInterlocutorFalse(t *testing.T) {
	_, _, _, _, uc := setUp(t)

	is := uc.IsInterlocutor(uint64(test_wrong_id), testDialogueMessageSQL)
	assert.Equal(t, is, false)
}

func TestChatUseCase_IsSenderMessageTrue(t *testing.T) {
	_, _, _, _, uc := setUp(t)

	is := uc.IsSenderMessage(uint64(test_id), testDialogueMessageSQL)
	assert.Equal(t, is, true)
}

func TestChatUseCase_IsSenderMessageFalse(t *testing.T) {
	_, _, _, _, uc := setUp(t)

	is := uc.IsSenderMessage(uint64(test_id2), testDialogueMessageSQL)
	assert.Equal(t, is, false)
}

func TestChatUseCase_DeleteDialogueOk(t *testing.T) {
	rep, _, _, _, uc := setUp(t)
	rep.EXPECT().CheckDialogueID(uint64(test_id)).Return(true, testDialogueMessageSQL, nil)
	rep.EXPECT().DeleteDialogue(uint64(test_id)).Return(nil)

	err := uc.DeleteDialogue(uint64(test_id), uint64(test_id))
	assert.Nil(t, err)
}

func TestChatUseCase_DeleteDialogueErrorDelete(t *testing.T) {
	rep, _, _, _, uc := setUp(t)
	rep.EXPECT().CheckDialogueID(uint64(test_id)).Return(true, testDialogueMessageSQL, nil)
	rep.EXPECT().DeleteDialogue(uint64(test_id)).Return(errors.New("smthing wrong"))

	err := uc.DeleteDialogue(uint64(test_id), uint64(test_id))
	assert.NotNil(t, err)
}

func TestChatUseCase_DeleteDialogueErrorNotInterlocutor(t *testing.T) {
	rep, _, _, _, uc := setUp(t)
	rep.EXPECT().CheckDialogueID(uint64(test_id)).Return(true, testDialogueMessageSQL, nil)

	err := uc.DeleteDialogue(uint64(test_wrong_id), uint64(test_id))
	assert.NotNil(t, err)
}

func TestChatUseCase_DeleteDialogueNoDialogue(t *testing.T) {
	rep, _, _, _, uc := setUp(t)
	rep.EXPECT().CheckDialogueID(uint64(test_id)).Return(false, models.EasyDialogueMessageSQL{}, nil)

	err := uc.DeleteDialogue(uint64(test_id), uint64(test_id))
	assert.NotNil(t, err)
}

func TestChatUseCase_DeleteDialogueErrorCheck(t *testing.T) {
	rep, _, _, _, uc := setUp(t)
	rep.EXPECT().CheckDialogueID(uint64(test_id)).Return(false, models.EasyDialogueMessageSQL{}, errors.New("invalid id"))

	err := uc.DeleteDialogue(uint64(test_id), uint64(test_id))
	assert.NotNil(t, err)
}

func TestChatUseCase_SendMessageOk(t *testing.T) {
	rep, _, repUser, _, uc := setUp(t)
	repUser.EXPECT().GetUserByID(uint64(test_id2)).Return(testUser, nil)
	rep.EXPECT().CheckDialogueUsers(uint64(test_id), uint64(test_id2)).Return(true, testDialogueMessageSQL, nil)
	rep.EXPECT().SendMessage(testDialogueMessageSQL.ID, &testNewMessage, uint64(test_id), gomock.Any()).Return(nil)

	err := uc.SendMessage(&testNewMessage, uint64(test_id))
	assert.Nil(t, err)
}

func TestChatUseCase_SendMessageErrorSend(t *testing.T) {
	rep, _, repUser, _, uc := setUp(t)
	repUser.EXPECT().GetUserByID(uint64(test_id2)).Return(testUser, nil)
	rep.EXPECT().CheckDialogueUsers(uint64(test_id), uint64(test_id2)).Return(true, testDialogueMessageSQL, nil)
	rep.EXPECT().SendMessage(testDialogueMessageSQL.ID, &testNewMessage, uint64(test_id), gomock.Any()).Return(errors.New("smthing wrong"))

	err := uc.SendMessage(&testNewMessage, uint64(test_id))
	assert.NotNil(t, err)
}

func TestChatUseCase_SendMessageNoDialogue(t *testing.T) {
	rep, _, repUser, _, uc := setUp(t)
	repUser.EXPECT().GetUserByID(uint64(test_id2)).Return(testUser, nil)
	rep.EXPECT().CheckDialogueUsers(uint64(test_id), uint64(test_id2)).Return(false, models.EasyDialogueMessageSQL{}, nil)
	rep.EXPECT().NewDialogue(uint64(test_id), testNewMessage.To).Return(uint64(test_id), nil)
	rep.EXPECT().SendMessage(testDialogueMessageSQL.ID, &testNewMessage, uint64(test_id), gomock.Any()).Return(nil)

	err := uc.SendMessage(&testNewMessage, uint64(test_id))
	assert.Nil(t, err)
}

func TestChatUseCase_SendMessageNoDialogueError(t *testing.T) {
	rep, _, repUser, _, uc := setUp(t)
	repUser.EXPECT().GetUserByID(uint64(test_id2)).Return(testUser, nil)
	rep.EXPECT().CheckDialogueUsers(uint64(test_id), uint64(test_id2)).Return(false, models.EasyDialogueMessageSQL{}, nil)
	rep.EXPECT().NewDialogue(uint64(test_id), testNewMessage.To).Return(uint64(test_id), errors.New("smthing wrong"))

	err := uc.SendMessage(&testNewMessage, uint64(test_id))
	assert.NotNil(t, err)
}

func TestChatUseCase_SendMessageCheckError(t *testing.T) {
	rep, _, repUser, _, uc := setUp(t)
	repUser.EXPECT().GetUserByID(uint64(test_id2)).Return(testUser, nil)
	rep.EXPECT().CheckDialogueUsers(uint64(test_id), uint64(test_id2)).Return(false, models.EasyDialogueMessageSQL{}, errors.New("smthing wrong"))

	err := uc.SendMessage(&testNewMessage, uint64(test_id))
	assert.NotNil(t, err)
}

func TestChatUseCase_SendMessageUserError(t *testing.T) {
	_, _, repUser, _, uc := setUp(t)
	repUser.EXPECT().GetUserByID(uint64(test_id2)).Return(testUser, errors.New("invalid id"))

	err := uc.SendMessage(&testNewMessage, uint64(test_id))
	assert.NotNil(t, err)
}

func TestChatUseCase_DeleteMessageOk(t *testing.T) {
	rep, _, _, _, uc := setUp(t)
	rep.EXPECT().CheckMessage(uint64(test_id)).Return(true, testDialogueMessageSQL, nil)
	rep.EXPECT().DeleteMessage(uint64(test_id)).Return(nil)

	err := uc.DeleteMessage(uint64(test_id), uint64(test_id))
	assert.Nil(t, err)
}

func TestChatUseCase_DeleteMessageErrorDelete(t *testing.T) {
	rep, _, _, _, uc := setUp(t)
	rep.EXPECT().CheckMessage(uint64(test_id)).Return(true, testDialogueMessageSQL, nil)
	rep.EXPECT().DeleteMessage(uint64(test_id)).Return(errors.New("smthing wrong"))

	err := uc.DeleteMessage(uint64(test_id), uint64(test_id))
	assert.NotNil(t, err)
}

func TestChatUseCase_DeleteMessageErrorNotInterlocutor(t *testing.T) {
	rep, _, _, _, uc := setUp(t)
	rep.EXPECT().CheckMessage(uint64(test_id)).Return(true, testDialogueMessageSQL, nil)

	err := uc.DeleteMessage(uint64(test_wrong_id), uint64(test_id))
	assert.NotNil(t, err)
}

func TestChatUseCase_DeleteMessageNoMessage(t *testing.T) {
	rep, _, _, _, uc := setUp(t)
	rep.EXPECT().CheckMessage(uint64(test_id)).Return(false, models.EasyDialogueMessageSQL{}, nil)

	err := uc.DeleteMessage(uint64(test_id), uint64(test_id))
	assert.NotNil(t, err)
}

func TestChatUseCase_DeleteMessageErrorCheck(t *testing.T) {
	rep, _, _, _, uc := setUp(t)
	rep.EXPECT().CheckMessage(uint64(test_id)).Return(false, models.EasyDialogueMessageSQL{}, errors.New("invalid id"))

	err := uc.DeleteMessage(uint64(test_id), uint64(test_id))
	assert.NotNil(t, err)
}

func TestChatUseCase_EditMessageOk(t *testing.T) {
	rep, _, _, _, uc := setUp(t)
	rep.EXPECT().CheckMessage(uint64(test_id)).Return(true, testDialogueMessageSQL, nil)
	rep.EXPECT().EditMessage(uint64(test_id), test_text).Return(nil)

	err := uc.EditMessage(uint64(test_id), &testRedactMessage)
	assert.Nil(t, err)
}

func TestChatUseCase_EditMessageError(t *testing.T) {
	rep, _, _, _, uc := setUp(t)
	rep.EXPECT().CheckMessage(uint64(test_id)).Return(true, testDialogueMessageSQL, nil)
	rep.EXPECT().EditMessage(uint64(test_id), test_text).Return(errors.New("smthing wrong"))

	err := uc.EditMessage(uint64(test_id), &testRedactMessage)
	assert.NotNil(t, err)
}

func TestChatUseCase_EditMessageNotSender(t *testing.T) {
	rep, _, _, _, uc := setUp(t)
	rep.EXPECT().CheckMessage(uint64(test_id)).Return(true, testDialogueMessageSQL, nil)
	rep.EXPECT().EditMessage(uint64(test_id), test_text).Return(nil)

	err := uc.EditMessage(uint64(test_id2), &testRedactMessage)
	assert.NotNil(t, err)
}

func TestChatUseCase_EditMessageNoMessage(t *testing.T) {
	rep, _, _, _, uc := setUp(t)
	rep.EXPECT().CheckMessage(uint64(test_id)).Return(false, models.EasyDialogueMessageSQL{}, nil)

	err := uc.EditMessage(uint64(test_id2), &testRedactMessage)
	assert.NotNil(t, err)
}

func TestChatUseCase_EditMessageCheckError(t *testing.T) {
	rep, _, _, _, uc := setUp(t)
	rep.EXPECT().CheckMessage(uint64(test_id)).Return(false, models.EasyDialogueMessageSQL{}, errors.New("invalid id"))

	err := uc.EditMessage(uint64(test_id2), &testRedactMessage)
	assert.NotNil(t, err)
}

func TestChatUseCase_AutoMailingConstructor(t *testing.T) {
	_, _, _, _, uc := setUp(t)

	mail := uc.AutoMailingConstructor(uint64(test_id2), test_from, test_event, fmt.Sprint(test_id))
	assert.Equal(t, mail, testMail)
}

func TestChatUseCase_MailingOk(t *testing.T) {
	rep, _, repUser, repEvent, uc := setUp(t)
	repUser.EXPECT().GetUserByID(uint64(test_id)).Return(testUserMail, nil)
	repEvent.EXPECT().GetOneEventByID(uint64(test_id)).Return(testEventSQL, nil)
	repUser.EXPECT().GetUserByID(uint64(test_id2)).Return(testUser, nil)
	rep.EXPECT().CheckDialogueUsers(uint64(test_id), uint64(test_id2)).Return(true, testDialogueMessageSQL, nil)
	rep.EXPECT().SendMessage(testDialogueMessageSQL.ID, &testMail, uint64(test_id), gomock.Any()).Return(nil)

	err := uc.Mailing(uint64(test_id), &testMailing)
	assert.Nil(t, err)
}

func TestChatUseCase_MailingSendError(t *testing.T) {
	rep, _, repUser, repEvent, uc := setUp(t)
	repUser.EXPECT().GetUserByID(uint64(test_id)).Return(testUserMail, nil)
	repEvent.EXPECT().GetOneEventByID(uint64(test_id)).Return(testEventSQL, nil)
	repUser.EXPECT().GetUserByID(uint64(test_id2)).Return(testUser, nil)
	rep.EXPECT().CheckDialogueUsers(uint64(test_id), uint64(test_id2)).Return(true, testDialogueMessageSQL, nil)
	rep.EXPECT().SendMessage(testDialogueMessageSQL.ID, &testMail, uint64(test_id), gomock.Any()).Return(errors.New("smthing wrong"))

	err := uc.Mailing(uint64(test_id), &testMailing)
	assert.NotNil(t, err)
}

func TestChatUseCase_MailingEventError(t *testing.T) {
	_, _, repUser, repEvent, uc := setUp(t)
	repUser.EXPECT().GetUserByID(uint64(test_id)).Return(testUserMail, nil)
	repEvent.EXPECT().GetOneEventByID(uint64(test_id)).Return(models.EventSQL{}, errors.New("invalid id"))

	err := uc.Mailing(uint64(test_id), &testMailing)
	assert.NotNil(t, err)
}

func TestChatUseCase_MailingUserError(t *testing.T) {
	_, _, repUser, _, uc := setUp(t)
	repUser.EXPECT().GetUserByID(uint64(test_id)).Return(testUserMail, errors.New("invalid id"))

	err := uc.Mailing(uint64(test_id), &testMailing)
	assert.NotNil(t, err)
}

func TestChatUseCase_SearchOkZeroId(t *testing.T) {
	rep, _, _, _, uc := setUp(t)
	rep.EXPECT().MessagesSearch(uint64(test_id), test_text, test_page).Return(testMessagesSQL, nil)

	messages, err := uc.Search(uint64(test_id), 0, test_text, test_page)
	assert.Nil(t, err)
	assert.Equal(t, messages, testMessages)
}

func TestChatUseCase_SearchOkZeroLength(t *testing.T) {
	rep, _, _, _, uc := setUp(t)
	rep.EXPECT().MessagesSearch(uint64(test_id), test_text, test_page).Return(models.MessagesSQL{}, nil)

	messages, err := uc.Search(uint64(test_id), 0, test_text, test_page)
	assert.Nil(t, err)
	assert.Equal(t, messages, models.Messages{})
}

func TestChatUseCase_SearchOkZeroIdError(t *testing.T) {
	rep, _, _, _, uc := setUp(t)
	rep.EXPECT().MessagesSearch(uint64(test_id), test_text, test_page).Return(models.MessagesSQL{}, errors.New("smthing wrong"))

	messages, err := uc.Search(uint64(test_id), 0, test_text, test_page)
	assert.NotNil(t, err)
	assert.Equal(t, messages, models.Messages{})
}

func TestChatUseCase_SearchOk(t *testing.T) {
	rep, _, _, _, uc := setUp(t)
	rep.EXPECT().CheckDialogueID(uint64(test_id)).Return(true, testDialogueMessageSQL, nil)
	rep.EXPECT().DialogueMessagesSearch(uint64(test_id), uint64(test_id), test_text, test_page).Return(testMessagesSQL, nil)

	messages, err := uc.Search(uint64(test_id), test_id, test_text, test_page)
	assert.Nil(t, err)
	assert.Equal(t, messages, testMessages)
}

func TestChatUseCase_SearchNotInterlocutor(t *testing.T) {
	rep, _, _, _, uc := setUp(t)
	rep.EXPECT().CheckDialogueID(uint64(test_id)).Return(true, testDialogueMessageSQL, nil)

	messages, err := uc.Search(uint64(test_wrong_id), test_id, test_text, test_page)
	assert.NotNil(t, err)
	assert.Equal(t, messages, models.Messages{})
}

func TestChatUseCase_SearchError(t *testing.T) {
	rep, _, _, _, uc := setUp(t)
	rep.EXPECT().CheckDialogueID(uint64(test_id)).Return(true, testDialogueMessageSQL, nil)
	rep.EXPECT().DialogueMessagesSearch(uint64(test_id), uint64(test_id), test_text, test_page).Return(models.MessagesSQL{}, errors.New("smthing wrong"))

	messages, err := uc.Search(uint64(test_id), test_id, test_text, test_page)
	assert.NotNil(t, err)
	assert.Equal(t, messages, models.Messages{})
}

func TestChatUseCase_SearchNoDialogue(t *testing.T) {
	rep, _, _, _, uc := setUp(t)
	rep.EXPECT().CheckDialogueID(uint64(test_id)).Return(false, models.EasyDialogueMessageSQL{}, nil)

	messages, err := uc.Search(uint64(test_id), test_id, test_text, test_page)
	assert.NotNil(t, err)
	assert.Equal(t, messages, models.Messages{})
}

func TestChatUseCase_SearchErrorCheck(t *testing.T) {
	rep, _, _, _, uc := setUp(t)
	rep.EXPECT().CheckDialogueID(uint64(test_id)).Return(false, models.EasyDialogueMessageSQL{}, errors.New("smthing wrong"))

	messages, err := uc.Search(uint64(test_id), test_id, test_text, test_page)
	assert.NotNil(t, err)
	assert.Equal(t, messages, models.Messages{})
}
