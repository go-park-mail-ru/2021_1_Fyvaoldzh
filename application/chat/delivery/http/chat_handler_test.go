package http

import (
	"bytes"
	"errors"
	mock_chat "kudago/application/microservices/chat/client"
	"kudago/application/models"
	"kudago/pkg/constants"
	"kudago/pkg/custom_sanitizer"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/microcosm-cc/bluemonday"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"kudago/application/microservices/auth/client"
	"kudago/pkg/logger"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	test_id       = uint64(1)
	test_id2      = uint64(2)
	test_id_param = 1
	test_page     = 1
	test_text     = "test text"
	test_text_2   = "test text 2"
	test_err      = errors.New("smthing wrong")
)

var testUser = models.UserOnEvent{
	Id:     uint64(test_id2),
	Name:   test_text,
	Avatar: test_text_2,
}

var testDialogue = models.Dialogue{
	ID:           uint64(test_id),
	Interlocutor: testUser,
	DialogMessages: models.Messages{
		{
			ID:     uint64(test_id),
			FromMe: true,
			Text:   test_text,
			Date:   time.Now().String(),
			Redact: false,
			Read:   false,
		},
		{
			ID:     uint64(test_id2),
			FromMe: false,
			Text:   test_text_2,
			Date:   time.Now().String(),
			Redact: false,
			Read:   false,
		},
	},
}

var testAllDialogues = models.DialogueCards{
	{
		ID:           uint64(test_id),
		Interlocutor: testUser,
		LastMessage: models.Message{
			ID:     uint64(test_id),
			FromMe: true,
			Text:   test_text,
			Date:   time.Now().String(),
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
			Date:   time.Now().String(),
			Redact: false,
			Read:   false,
		},
	},
}
var testNewMessageJSON = models.NewMessageJSON{
	To:   "1",
	Text: test_text,
}

var testNewMessageJSONAtoi = models.NewMessageJSON{
	To:   "ooo",
	Text: test_text,
}

var testNewMessageJSONLessZero = models.NewMessageJSON{
	To:   "-2",
	Text: test_text,
}

var testEditMessageJSON = models.RedactMessageJSON{
	ID:   "1",
	Text: test_text,
}

var testEditMessageJSONAtoi = models.RedactMessageJSON{
	ID:   "ooo",
	Text: test_text,
}

var testEditMessageJSONLessZero = models.RedactMessageJSON{
	ID:   "-2",
	Text: test_text,
}

var testNewMessage = models.NewMessage{
	To:   test_id,
	Text: test_text,
}

var testEditMessage = models.RedactMessage{
	ID:   uint64(1),
	Text: test_text,
}

var testMailingJSON = models.MailingJSON{
	EventID: test_id,
	To:      []string{"1"},
}

var testMailingJSONZeroEvent = models.MailingJSON{
	EventID: 0,
	To:      []string{"1"},
}

var testMailingJSONAtoi = models.MailingJSON{
	EventID: 0,
	To:      []string{"ooo"},
}

var testMailingJSONLessZero = models.MailingJSON{
	EventID: 0,
	To:      []string{"-1"},
}

var testMailing = models.Mailing{
	EventID: test_id,
	To:      []uint64{test_id},
}

/*
var testMessages = models.Messages{
	{
		ID:     uint64(test_id),
		FromMe: true,
		Text:   test_text,
		Date:   time.Now().String(),
		Redact: false,
		Read:   false,
	},
	{
		ID:     uint64(test_id2),
		FromMe: false,
		Text:   test_text_2,
		Date:   time.Now().String(),
		Redact: false,
		Read:   false,
	},
}


 */
func setUp(t *testing.T, url, method string) (echo.Context,
	ChatHandler, *mock_chat.MockIChatClient, *client.MockIAuthClient) {
	e := echo.New()
	r := e.Router()
	r.Add(method, url, func(echo.Context) error { return nil })

	ctrl := gomock.NewController(t)
	rpcChat := mock_chat.NewMockIChatClient(ctrl)
	rpcAuth := client.NewMockIAuthClient(ctrl)

	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	sugar := l.Sugar()
	zap.NewAtomicLevelAt(zapcore.DebugLevel)
	cs := custom_sanitizer.NewCustomSanitizer(bluemonday.UGCPolicy())

	handler := ChatHandler{
		rpcAuth:   rpcAuth,
		rpcChat:   rpcChat,
		Logger:    logger.NewLogger(sugar),
		sanitizer: cs,
	}

	var req *http.Request
	switch method {
	case http.MethodPost:
		switch url {
		case "/api/v1/send":
			f, _ := testNewMessageJSON.MarshalJSON()
			req = httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(f))

		case "/api/v1/message":
			f, _ := testEditMessageJSON.MarshalJSON()
			req = httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(f))

		case "/api/v1/message/mailing":
			f, _ := testMailingJSON.MarshalJSON()
			req = httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(f))
		}
	case http.MethodGet:
		req = httptest.NewRequest(http.MethodGet, url, nil)
	case http.MethodDelete:
		req = httptest.NewRequest(http.MethodGet, url, nil)
	}
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(url)

	return c, handler, rpcChat, rpcAuth
}

func TestEventsHandler_GetDialoguesOk(t *testing.T) {
	c, h, rpcChat, _ := setUp(t, "/api/v1/dialogues", http.MethodGet)
	rpcChat.EXPECT().GetAllDialogues(test_id, test_page).Return(testAllDialogues, nil, 200)
	c.Set(constants.PageKey, test_page)
	c.Set(constants.UserIdKey, test_id)

	err := h.GetDialogues(c)

	assert.Nil(t, err)
}

func TestEventsHandler_GetDialoguesError(t *testing.T) {
	c, h, rpcChat, _ := setUp(t, "/api/v1/dialogues", http.MethodGet)
	rpcChat.EXPECT().GetAllDialogues(test_id, test_page).Return(models.DialogueCards{}, test_err, 500)
	c.Set(constants.PageKey, test_page)
	c.Set(constants.UserIdKey, test_id)

	err := h.GetDialogues(c)

	assert.Equal(t, err, test_err)
}

func TestEventsHandler_GetOneDialogueOk(t *testing.T) {
	c, h, rpcChat, _ := setUp(t, "/api/v1/dialogues/:id", http.MethodGet)
	rpcChat.EXPECT().GetOneDialogue(test_id, uint64(test_id_param), test_page).Return(testDialogue, nil, 200)
	c.Set(constants.PageKey, test_page)
	c.Set(constants.UserIdKey, test_id)
	c.Set(constants.IdKey, test_id_param)

	err := h.GetOneDialogue(c)

	assert.Nil(t, err)
}

func TestEventsHandler_GetOneDialogueError(t *testing.T) {
	c, h, rpcChat, _ := setUp(t, "/api/v1/dialogues/:id", http.MethodGet)
	rpcChat.EXPECT().GetOneDialogue(test_id, uint64(test_id_param), test_page).Return(models.Dialogue{}, test_err, 500)
	c.Set(constants.PageKey, test_page)
	c.Set(constants.UserIdKey, test_id)
	c.Set(constants.IdKey, test_id_param)

	err := h.GetOneDialogue(c)

	assert.Equal(t, err, test_err)
}

func TestEventsHandler_DeleteDialogueOk(t *testing.T) {
	c, h, rpcChat, _ := setUp(t, "/api/v1/dialogues/:id", http.MethodDelete)
	rpcChat.EXPECT().DeleteDialogue(test_id, uint64(test_id_param)).Return(nil, 200)
	c.Set(constants.UserIdKey, test_id)
	c.Set(constants.IdKey, test_id_param)

	err := h.DeleteDialogue(c)

	assert.Nil(t, err)
}

func TestEventsHandler_DeleteDialogueError(t *testing.T) {
	c, h, rpcChat, _ := setUp(t, "/api/v1/dialogues/:id", http.MethodDelete)
	rpcChat.EXPECT().DeleteDialogue(test_id, uint64(test_id_param)).Return(test_err, 500)
	c.Set(constants.UserIdKey, test_id)
	c.Set(constants.IdKey, test_id_param)

	err := h.DeleteDialogue(c)

	assert.Equal(t, err, test_err)
}

func TestEventsHandler_SendMessageOk(t *testing.T) {
	c, h, rpcChat, _ := setUp(t, "/api/v1/send", http.MethodPost)
	rpcChat.EXPECT().SendMessage(&testNewMessage, test_id).Return(nil, 200)
	c.Set(constants.UserIdKey, test_id)

	err := h.SendMessage(c)

	assert.Nil(t, err)
}

func TestEventsHandler_SendMessageAtoi(t *testing.T) {
	f, _ := testNewMessageJSONAtoi.MarshalJSON()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/send", bytes.NewBuffer(f))
	rec := httptest.NewRecorder()
	_, h, _, _ := setUp(t, "/api/v1/send", http.MethodPost)
	c := echo.New().NewContext(req, rec)
	c.Set(constants.UserIdKey, test_id)

	err := h.SendMessage(c)

	assert.NotNil(t, err)
}

func TestEventsHandler_SendMessageLessZero(t *testing.T) {
	f, _ := testNewMessageJSONLessZero.MarshalJSON()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/send", bytes.NewBuffer(f))
	rec := httptest.NewRecorder()
	_, h, _, _ := setUp(t, "/api/v1/send", http.MethodPost)
	c := echo.New().NewContext(req, rec)
	c.Set(constants.UserIdKey, test_id)

	err := h.SendMessage(c)

	assert.NotNil(t, err)
}

func TestEventsHandler_SendMessageError(t *testing.T) {
	c, h, rpcChat, _ := setUp(t, "/api/v1/send", http.MethodPost)
	rpcChat.EXPECT().SendMessage(&testNewMessage, test_id).Return(test_err, 500)
	c.Set(constants.UserIdKey, test_id)

	err := h.SendMessage(c)

	assert.Equal(t, err, test_err)
}

func TestEventsHandler_DeleteMessageOk(t *testing.T) {
	c, h, rpcChat, _ := setUp(t, "/api/v1/message/:id", http.MethodDelete)
	rpcChat.EXPECT().DeleteMessage(test_id, uint64(test_id_param)).Return(nil, 200)
	c.Set(constants.UserIdKey, test_id)
	c.Set(constants.IdKey, test_id_param)

	err := h.DeleteMessage(c)

	assert.Nil(t, err)
}

func TestEventsHandler_DeleteMessageError(t *testing.T) {
	c, h, rpcChat, _ := setUp(t, "/api/v1/message/:id", http.MethodDelete)
	rpcChat.EXPECT().DeleteMessage(test_id, uint64(test_id_param)).Return(test_err, 500)
	c.Set(constants.UserIdKey, test_id)
	c.Set(constants.IdKey, test_id_param)

	err := h.DeleteMessage(c)

	assert.Equal(t, err, test_err)
}

func TestEventsHandler_EditMessageOk(t *testing.T) {
	c, h, rpcChat, _ := setUp(t, "/api/v1/message", http.MethodPost)
	rpcChat.EXPECT().EditMessage(test_id, &testEditMessage).Return(nil, 200)
	c.Set(constants.UserIdKey, test_id)

	err := h.EditMessage(c)

	assert.Nil(t, err)
}

func TestEventsHandler_EditMessageAtoi(t *testing.T) {
	f, _ := testEditMessageJSONAtoi.MarshalJSON()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/message", bytes.NewBuffer(f))
	rec := httptest.NewRecorder()
	_, h, _, _ := setUp(t, "/api/v1/message", http.MethodPost)
	c := echo.New().NewContext(req, rec)
	c.Set(constants.UserIdKey, test_id)

	err := h.EditMessage(c)

	assert.NotNil(t, err)
}

func TestEventsHandler_EditMessageLessZero(t *testing.T) {
	f, _ := testEditMessageJSONLessZero.MarshalJSON()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/message", bytes.NewBuffer(f))
	rec := httptest.NewRecorder()
	_, h, _, _ := setUp(t, "/api/v1/message", http.MethodPost)
	c := echo.New().NewContext(req, rec)
	c.Set(constants.UserIdKey, test_id)

	err := h.EditMessage(c)

	assert.NotNil(t, err)
}

func TestEventsHandler_EditMessageError(t *testing.T) {
	c, h, rpcChat, _ := setUp(t, "/api/v1/message", http.MethodPost)
	rpcChat.EXPECT().EditMessage(test_id, &testEditMessage).Return(test_err, 500)
	c.Set(constants.UserIdKey, test_id)

	err := h.EditMessage(c)

	assert.Equal(t, err, test_err)
}

func TestEventsHandler_MailingOk(t *testing.T) {
	c, h, rpcChat, _ := setUp(t, "/api/v1/message/mailing", http.MethodPost)
	rpcChat.EXPECT().Mailing(test_id, &testMailing).Return(nil, 200)
	c.Set(constants.UserIdKey, test_id)

	err := h.Mailing(c)

	assert.Nil(t, err)
}

func TestEventsHandler_MailingZeroEvent(t *testing.T) {
	f, _ := testMailingJSONZeroEvent.MarshalJSON()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/message/mailing", bytes.NewBuffer(f))
	rec := httptest.NewRecorder()
	_, h, _, _ := setUp(t, "/api/v1/message/mailing", http.MethodPost)
	c := echo.New().NewContext(req, rec)
	c.Set(constants.UserIdKey, test_id)

	err := h.Mailing(c)

	assert.NotNil(t, err)
}

func TestEventsHandler_MailingAtoi(t *testing.T) {
	f, _ := testMailingJSONAtoi.MarshalJSON()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/message/mailing", bytes.NewBuffer(f))
	rec := httptest.NewRecorder()
	_, h, _, _ := setUp(t, "/api/v1/message/mailing", http.MethodPost)
	c := echo.New().NewContext(req, rec)
	c.Set(constants.UserIdKey, test_id)

	err := h.Mailing(c)

	assert.NotNil(t, err)
}

func TestEventsHandler_MailingLessZero(t *testing.T) {
	f, _ := testMailingJSONLessZero.MarshalJSON()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/message/mailing", bytes.NewBuffer(f))
	rec := httptest.NewRecorder()
	_, h, _, _ := setUp(t, "/api/v1/message/mailing", http.MethodPost)
	c := echo.New().NewContext(req, rec)
	c.Set(constants.UserIdKey, test_id)

	err := h.Mailing(c)

	assert.NotNil(t, err)
}

func TestEventsHandler_MailingError(t *testing.T) {
	c, h, rpcChat, _ := setUp(t, "/api/v1/message/mailing", http.MethodPost)
	rpcChat.EXPECT().Mailing(test_id, &testMailing).Return(test_err, 500)
	c.Set(constants.UserIdKey, test_id)

	err := h.Mailing(c)

	assert.Equal(t, err, test_err)
}

/*
func TestEventsHandler_SearchOk(t *testing.T) {
	c, h, rpcChat, _ := setUp(t, "/api/v1/dialogues", http.MethodGet)
	rpcChat.EXPECT().Search(test_id, 0, "", test_page).Return(testMessages, nil, 200)
	c.Set(constants.PageKey, test_page)
	c.Set(constants.UserIdKey, test_id)
	c.SetParamNames("find")
	c.SetParamValues("")
	c.SetParamNames("id")
	c.SetParamValues("")

	err := h.Search(c)

	assert.Nil(t, err)
}

 */

func TestEventsHandler_SearchAtoi(t *testing.T) {
	c, h, _, _ := setUp(t, "/api/v1/dialogues?id='aaa'", http.MethodGet)

	c.Set(constants.PageKey, test_page)
	c.Set(constants.UserIdKey, test_id)
	c.SetParamNames("find")
	c.SetParamValues("")
	c.SetParamNames("id")
	c.SetParamValues("aaa")

	err := h.Search(c)

	assert.NotNil(t, err)
}

func TestEventsHandler_SearchLessZero(t *testing.T) {
	c, h, _, _ := setUp(t, "/api/v1/dialogues?id=-5", http.MethodGet)
	c.Set(constants.PageKey, test_page)
	c.Set(constants.UserIdKey, test_id)
	c.SetParamNames("find")
	c.SetParamValues("")
	c.SetParamNames("id")
	c.SetParamValues("-5")

	err := h.Search(c)

	assert.NotNil(t, err)
}

/*
func TestEventsHandler_SearchError(t *testing.T) {
	c, h, rpcChat, _ := setUp(t, "/api/v1/dialogues", http.MethodGet)
	rpcChat.EXPECT().Search(test_id, 0, "", test_page).Return(models.Messages{}, test_err, 500)
	c.Set(constants.PageKey, test_page)
	c.Set(constants.UserIdKey, test_id)
	c.SetParamNames("find")
	c.SetParamValues("")
	c.SetParamNames("id")
	c.SetParamValues("")

	err := h.Search(c)

	assert.Equal(t, err, test_err)
}


 */