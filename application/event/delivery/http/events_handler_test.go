package http

import (
	"bytes"
	"errors"
	mock_event "kudago/application/event/mocks"
	"kudago/application/models"
	"kudago/pkg/constants"
	"kudago/pkg/custom_sanitizer"
	"kudago/pkg/generator"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/microcosm-cc/bluemonday"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	mock_infrastructure "kudago/pkg/infrastructure/mocks"
	"kudago/pkg/logger"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	title    = "test_title"
	desc     = "test_description"
	img      = "test_img_addr"
	place    = "test_place"
	subway   = "test_subway"
	street   = "test_street"
	category = "test_category"
	name     = "test_name"
	ttype    = "test_type"
	search   = "test_search"
)

var testTags = models.Tags{
	{ID: 1,
		Name: name},
	{ID: 2,
		Name: "test_name2"},
}

var testFollowers = models.UsersOnEvent{
	{Id: 1,
		Name:   name,
		Avatar: img},
}

var testEvent = models.Event{
	ID:          1,
	Title:       title,
	Place:       place,
	Description: desc,
	StartDate:   time.Now().String(),
	EndDate:     time.Now().Add(15000 * time.Hour).String(),
	Subway:      subway,
	Street:      street,
	Tags:        testTags,
	Category:    category,
	Image:       img,
	Followers:   testFollowers,
}

var testAllEvents = models.EventCards{
	{
		ID:          1,
		Title:       title,
		Description: desc,
		Place:       img,
		StartDate:   time.Now().String(),
		EndDate:     time.Now().Add(15000 * time.Hour).String(),
	},
	{
		ID:          2,
		Title:       "test_title2",
		Description: "test_description2",
		Place:       img,
		StartDate:   time.Now().String(),
		EndDate:     time.Now().Add(15000 * time.Hour).String(),
	},
}

func setUp(t *testing.T, url, method string) (echo.Context,
	EventHandler, *mock_event.MockUseCase, *mock_infrastructure.MockSessionTarantool) {
	e := echo.New()
	r := e.Router()
	r.Add(method, url, func(echo.Context) error { return nil })

	ctrl := gomock.NewController(t)
	_ = mock_event.NewMockRepository(ctrl)
	usecase := mock_event.NewMockUseCase(ctrl)
	sm := mock_infrastructure.NewMockSessionTarantool(ctrl)

	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	sugar := l.Sugar()
	zap.NewAtomicLevelAt(zapcore.DebugLevel)
	cs := custom_sanitizer.NewCustomSanitizer(bluemonday.UGCPolicy())

	handler := EventHandler{
		UseCase:   usecase,
		Sm:        sm,
		Logger:    logger.NewLogger(sugar),
		sanitizer: cs,
	}

	var req *http.Request
	switch method {
	case http.MethodPost:
		switch url {
		case "/api/v1/create":
			f, _ := testEvent.MarshalJSON()
			req = httptest.NewRequest(http.MethodGet, url, bytes.NewBuffer(f))
		case "/api/v1/save/:id":
			req = httptest.NewRequest(http.MethodGet, url, bytes.NewBuffer([]byte{}))
			/*case "/api/v1/update":
			f, _ := testOwnUserProfile.MarshalJSON()
			req = httptest.NewRequest(http.MethodGet, url, bytes.NewBuffer(f))*/
		}
	case http.MethodGet:
		req = httptest.NewRequest(http.MethodGet, url, nil)
	case http.MethodDelete:
		req = httptest.NewRequest(http.MethodGet, url, nil)
	}
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(url)

	return c, handler, usecase, sm
}

func TestEventsHandler_GetAllEventsOk(t *testing.T) {
	c, h, usecase, _ := setUp(t, "/api/v1/?page=1", http.MethodGet)
	usecase.EXPECT().GetAllEvents(1).Return(testAllEvents, nil)

	err := h.GetAllEvents(c)

	assert.Nil(t, err)
}

func TestEventsHandler_GetAllEventsAtoiError(t *testing.T) {
	c, h, _, _ := setUp(t, "/api/v1/?page='aaaa'", http.MethodGet)

	err := h.GetAllEvents(c)

	assert.NotNil(t, err)
}

func TestEventsHandler_GetAllEventsWithoutPage(t *testing.T) {
	c, h, usecase, _ := setUp(t, "/api/v1/", http.MethodGet)
	usecase.EXPECT().GetAllEvents(1).Return(testAllEvents, nil)

	err := h.GetAllEvents(c)

	assert.Nil(t, err)
}

func TestEventsHandler_GetAllEventsWithZero(t *testing.T) {
	c, h, usecase, _ := setUp(t, "/api/v1/?page=0", http.MethodGet)
	usecase.EXPECT().GetAllEvents(1).Return(testAllEvents, nil)

	err := h.GetAllEvents(c)

	assert.Nil(t, err)
}

func TestEventsHandler_GetAllEventsError(t *testing.T) {
	c, h, usecase, _ := setUp(t, "/api/v1/?page=1", http.MethodGet)
	usecase.EXPECT().GetAllEvents(1).Return(models.EventCards{}, errors.New("get all error"))

	err := h.GetAllEvents(c)

	assert.NotNil(t, err)
}

func TestEventsHandler_GetOneEventOk(t *testing.T) {
	c, h, usecase, sm := setUp(t, "/api/v1/event/:id", http.MethodGet)
	c.SetParamNames("id")
	c.SetParamValues("1")
	key := generator.RandStringRunes(constants.CookieLength)
	newCookie := &http.Cookie{
		Name:     constants.SessionCookieName,
		Value:    key,
		Expires:  time.Now().Add(10 * time.Hour),
		HttpOnly: true,
	}
	c.Request().AddCookie(newCookie)
	usecase.EXPECT().GetOneEvent(uint64(1)).Return(testEvent, nil)
	sm.EXPECT().CheckSession(key).Return(true, uint64(1), nil)
	usecase.EXPECT().RecomendSystem(uint64(1), testEvent.Category).Return(nil)

	err := h.GetOneEvent(c)

	assert.Nil(t, err)
}

func TestEventsHandler_GetOneEventWrongCookie(t *testing.T) {
	c, h, usecase, _ := setUp(t, "/api/v1/event/:id", http.MethodGet)
	c.SetParamNames("id")
	c.SetParamValues("1")
	key := generator.RandStringRunes(constants.CookieLength)
	newCookie := &http.Cookie{
		Name:     "test_name",
		Value:    key,
		Expires:  time.Now().Add(10 * time.Hour),
		HttpOnly: true,
	}
	c.Request().AddCookie(newCookie)
	usecase.EXPECT().GetOneEvent(uint64(1)).Return(testEvent, nil)

	err := h.GetOneEvent(c)

	assert.Nil(t, err)
}

func TestEventsHandler_GetOneEventCookieNotExist(t *testing.T) {
	c, h, usecase, sm := setUp(t, "/api/v1/event/:id", http.MethodGet)
	c.SetParamNames("id")
	c.SetParamValues("1")
	key := generator.RandStringRunes(constants.CookieLength)
	newCookie := &http.Cookie{
		Name:     constants.SessionCookieName,
		Value:    key,
		Expires:  time.Now().Add(10 * time.Hour),
		HttpOnly: true,
	}
	c.Request().AddCookie(newCookie)
	usecase.EXPECT().GetOneEvent(uint64(1)).Return(testEvent, nil)
	sm.EXPECT().CheckSession(key).Return(false, uint64(1), nil)

	err := h.GetOneEvent(c)

	assert.Nil(t, err)
}

func TestEventsHandler_GetOneEventErrorCookie(t *testing.T) {
	c, h, usecase, sm := setUp(t, "/api/v1/event/:id", http.MethodGet)
	c.SetParamNames("id")
	c.SetParamValues("1")
	key := generator.RandStringRunes(constants.CookieLength)
	newCookie := &http.Cookie{
		Name:     constants.SessionCookieName,
		Value:    key,
		Expires:  time.Now().Add(10 * time.Hour),
		HttpOnly: true,
	}
	c.Request().AddCookie(newCookie)
	usecase.EXPECT().GetOneEvent(uint64(1)).Return(testEvent, nil)
	sm.EXPECT().CheckSession(key).Return(true, uint64(1), errors.New("error cookie"))

	err := h.GetOneEvent(c)

	assert.Nil(t, err)
}

func TestEventsHandler_GetOneEventAtoiError(t *testing.T) {
	c, h, _, _ := setUp(t, "/api/v1/event/:id", http.MethodGet)
	c.SetParamNames("id")
	c.SetParamValues("aaaa")

	err := h.GetOneEvent(c)

	assert.NotNil(t, err)
}

func TestEventsHandler_GetOneEventError(t *testing.T) {
	c, h, usecase, _ := setUp(t, "/api/v1/event/:id", http.MethodGet)
	c.SetParamNames("id")
	c.SetParamValues("1")

	usecase.EXPECT().GetOneEvent(uint64(1)).Return(models.Event{}, errors.New("get oneEv error"))

	err := h.GetOneEvent(c)

	assert.NotNil(t, err)
}

func TestEventsHandler_GetEventsOk(t *testing.T) {
	c, h, usecase, _ := setUp(t, "/api/v1/event?page=1", http.MethodGet)
	usecase.EXPECT().GetEventsByCategory("", 1).Return(testAllEvents, nil)

	err := h.GetEvents(c)

	assert.Nil(t, err)
}

func TestEventsHandler_GetEventsAtoiErros(t *testing.T) {
	c, h, _, _ := setUp(t, "/api/v1/event?page='aaaa'", http.MethodGet)

	err := h.GetEvents(c)

	assert.NotNil(t, err)
}

func TestEventsHandler_GetEventsWithoutPage(t *testing.T) {
	c, h, usecase, _ := setUp(t, "/api/v1/event", http.MethodGet)
	usecase.EXPECT().GetEventsByCategory("", 1).Return(testAllEvents, nil)

	err := h.GetEvents(c)

	assert.Nil(t, err)
}

func TestEventsHandler_GetEventsWithZero(t *testing.T) {
	c, h, usecase, _ := setUp(t, "/api/v1/event?page=0", http.MethodGet)
	usecase.EXPECT().GetEventsByCategory("", 1).Return(testAllEvents, nil)

	err := h.GetEvents(c)

	assert.Nil(t, err)
}

func TestEventsHandler_GetEventsError(t *testing.T) {
	c, h, usecase, _ := setUp(t, "/api/v1/event?page=1", http.MethodGet)
	usecase.EXPECT().GetEventsByCategory("", 1).Return(models.EventCards{}, errors.New("get events error"))

	err := h.GetEvents(c)

	assert.NotNil(t, err)
}

func TestEventsHandler_FindEventsOk(t *testing.T) {
	c, h, usecase, _ := setUp(t, "/api/v1/search?page=1", http.MethodGet)
	usecase.EXPECT().FindEvents("", "", 1).Return(testAllEvents, nil)

	err := h.FindEvents(c)

	assert.Nil(t, err)
}

func TestEventsHandler_FindEventsAtoiErros(t *testing.T) {
	c, h, _, _ := setUp(t, "/api/v1/search?page='aaaa'", http.MethodGet)

	err := h.FindEvents(c)

	assert.NotNil(t, err)
}

func TestEventsHandler_FindEventsWithoutPage(t *testing.T) {
	c, h, usecase, _ := setUp(t, "/api/v1/search", http.MethodGet)
	usecase.EXPECT().FindEvents("", "", 1).Return(testAllEvents, nil)

	err := h.FindEvents(c)
	assert.Nil(t, err)
}

func TestEventsHandler_FindEventsWithZero(t *testing.T) {
	c, h, usecase, _ := setUp(t, "/api/v1/search?page=0", http.MethodGet)
	usecase.EXPECT().FindEvents("", "", 1).Return(testAllEvents, nil)

	err := h.FindEvents(c)

	assert.Nil(t, err)
}

func TestEventsHandler_FindEventsError(t *testing.T) {
	c, h, usecase, _ := setUp(t, "/api/v1/search?page=1", http.MethodGet)
	usecase.EXPECT().FindEvents("", "", 1).Return(models.EventCards{}, errors.New("get events error"))

	err := h.FindEvents(c)

	assert.NotNil(t, err)
}

func TestEventsHandler_GetImageOk(t *testing.T) {
	c, h, usecase, _ := setUp(t, "/api/v1/event/:id", http.MethodGet)
	c.SetParamNames("id")
	c.SetParamValues("1")
	usecase.EXPECT().GetImage(uint64(1)).Return([]byte{}, nil)

	err := h.GetImage(c)

	assert.Nil(t, err)
}

func TestEventsHandler_GetImageAtoiError(t *testing.T) {
	c, h, _, _ := setUp(t, "/api/v1/event/:id", http.MethodGet)
	c.SetParamNames("id")
	c.SetParamValues("aaaa")

	err := h.GetImage(c)

	assert.NotNil(t, err)
}

func TestEventsHandler_GetImageError(t *testing.T) {
	c, h, usecase, _ := setUp(t, "/api/v1/event/:id", http.MethodGet)
	c.SetParamNames("id")
	c.SetParamValues("1")

	usecase.EXPECT().GetImage(uint64(1)).Return([]byte{}, errors.New("get image error"))

	err := h.GetImage(c)

	assert.NotNil(t, err)
}

func TestEventsHandler_RecommendNonAuthorize(t *testing.T) {
	c, h, _, _ := setUp(t, "/api/v1/recomend?page=1", http.MethodGet)

	err := h.Recommend(c)
	assert.NotNil(t, err)
}

func TestEventsHandler_RecommendAtoiError(t *testing.T) {
	c, h, _, _ := setUp(t, "/api/v1/recomend?page='aaaa'", http.MethodGet)

	err := h.Recommend(c)

	assert.NotNil(t, err)
}

func TestEventsHandler_RecommendWithoutPageOk(t *testing.T) {
	c, h, usecase, sm := setUp(t, "/api/v1/recomend", http.MethodGet)
	key := generator.RandStringRunes(constants.CookieLength)
	newCookie := &http.Cookie{
		Name:     constants.SessionCookieName,
		Value:    key,
		Expires:  time.Now().Add(10 * time.Hour),
		HttpOnly: true,
	}
	c.Request().AddCookie(newCookie)

	sm.EXPECT().CheckSession(key).Return(true, uint64(1), nil)
	usecase.EXPECT().GetRecommended(uint64(1), 1).Return(testAllEvents, nil)

	err := h.Recommend(c)

	assert.Nil(t, err)
}

func TestEventsHandler_RecommendWithZeroError(t *testing.T) {
	c, h, usecase, sm := setUp(t, "/api/v1/recomend?page=0", http.MethodGet)
	key := generator.RandStringRunes(constants.CookieLength)
	newCookie := &http.Cookie{
		Name:     constants.SessionCookieName,
		Value:    key,
		Expires:  time.Now().Add(10 * time.Hour),
		HttpOnly: true,
	}
	c.Request().AddCookie(newCookie)

	sm.EXPECT().CheckSession(key).Return(true, uint64(1), nil)
	usecase.EXPECT().GetRecommended(uint64(1), 1).Return(models.EventCards{}, errors.New("get recommend error"))

	err := h.Recommend(c)

	assert.NotNil(t, err)
}

func TestEventsHandler_CreateOk(t *testing.T) {
	c, h, usecase, _ := setUp(t, "/api/v1/create", http.MethodPost)

	usecase.EXPECT().CreateNewEvent(&testEvent).Return(nil)

	err := h.Create(c)

	assert.Nil(t, err)
}

func TestEventsHandler_CreateError(t *testing.T) {
	c, h, usecase, _ := setUp(t, "/api/v1/create", http.MethodPost)

	usecase.EXPECT().CreateNewEvent(&testEvent).Return(errors.New("get create error"))

	err := h.Create(c)

	assert.NotNil(t, err)
}

func TestEventsHandler_SaveError(t *testing.T) {
	c, h, _, _ := setUp(t, "/api/v1/save/:id", http.MethodPost)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := h.Save(c)

	assert.NotNil(t, err)
}

func TestEventsHandler_SaveAtoiError(t *testing.T) {
	c, h, _, _ := setUp(t, "/api/v1/save/:id", http.MethodPost)
	c.SetParamNames("id")
	c.SetParamValues("aaa")

	err := h.Save(c)

	assert.NotNil(t, err)
}

func TestEventsHandler_DeleteOk(t *testing.T) {
	c, h, usecase, _ := setUp(t, "/api/v1/event/:id", http.MethodDelete)
	c.SetParamNames("id")
	c.SetParamValues("1")

	usecase.EXPECT().Delete(uint64(1)).Return(nil)

	err := h.Delete(c)

	assert.Nil(t, err)
}

func TestEventsHandler_DeleteAtoiError(t *testing.T) {
	c, h, _, _ := setUp(t, "/api/v1/event/:id", http.MethodDelete)
	c.SetParamNames("id")
	c.SetParamValues("aaaa")

	err := h.Delete(c)

	assert.NotNil(t, err)
}

func TestEventsHandler_DeleteError(t *testing.T) {
	c, h, usecase, _ := setUp(t, "/api/v1/event/:id", http.MethodDelete)
	c.SetParamNames("id")
	c.SetParamValues("1")

	usecase.EXPECT().Delete(uint64(1)).Return(errors.New("get delete error"))

	err := h.Delete(c)

	assert.NotNil(t, err)
}
