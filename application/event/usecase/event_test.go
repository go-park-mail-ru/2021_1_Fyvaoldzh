package usecase

import (
	"database/sql"
	"kudago/application/event"
	mock_event "kudago/application/event/mocks"
	"kudago/application/models"
	mock_subscription "kudago/application/subscription/mocks"
	"kudago/pkg/logger"
	"log"
	"mime/multipart"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
var testAllEventsWithDateSQL = []models.EventCardWithDateSQL{
	{
		ID:          1,
		Title:       title,
		Description: desc,
		Image:       sql.NullString{String: img, Valid: true},
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(15000 * time.Hour),
	},
	{
		ID:          2,
		Title:       "test_title2",
		Description: "test_description2",
		Image:       sql.NullString{String: "test_img_addr2", Valid: true},
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(15000 * time.Hour),
	},
}

var testAllEvents = models.EventCards{
	{
		ID:          1,
		Title:       title,
		Description: desc,
		Image:       img,
		StartDate:   testAllEventsWithDateSQL[0].StartDate.String(),
		EndDate:     testAllEventsWithDateSQL[0].EndDate.String(),
	},
	{
		ID:          2,
		Title:       "test_title2",
		Description: "test_description2",
		Image:       testAllEventsWithDateSQL[1].Image.String,
		StartDate:   testAllEventsWithDateSQL[1].StartDate.String(),
		EndDate:     testAllEventsWithDateSQL[1].EndDate.String(),
	},
}

var testEventSQL = models.EventSQL{
	ID:          1,
	Title:       title,
	Place:       place,
	Description: desc,
	StartDate:   sql.NullTime{Time: time.Now(), Valid: true},
	EndDate:     sql.NullTime{Time: time.Now().Add(15000 * time.Hour), Valid: true},
	Subway:      sql.NullString{String: subway, Valid: true},
	Street:      sql.NullString{String: street, Valid: true},
	Category:    category,
	Image:       sql.NullString{String: img, Valid: true},
}

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
	StartDate:   testEventSQL.StartDate.Time.String(),
	EndDate:     testEventSQL.EndDate.Time.String(),
	Subway:      subway,
	Street:      street,
	Tags:        testTags,
	Category:    category,
	Image:       img,
	Followers:   testFollowers,
}

var testEventErrTags = models.Event{
	ID:          1,
	Title:       title,
	Place:       place,
	Description: desc,
	StartDate:   testEventSQL.StartDate.Time.String(),
	EndDate:     testEventSQL.EndDate.Time.String(),
	Subway:      subway,
	Street:      street,
	Tags:        models.Tags(nil),
	Category:    category,
	Image:       img,
	Followers:   models.UsersOnEvent(nil),
}

var testEventErrFollowers = models.Event{
	ID:          1,
	Title:       title,
	Place:       place,
	Description: desc,
	StartDate:   testEventSQL.StartDate.Time.String(),
	EndDate:     testEventSQL.EndDate.Time.String(),
	Subway:      subway,
	Street:      street,
	Tags:        testTags,
	Category:    category,
	Image:       img,
	Followers:   models.UsersOnEvent(nil),
}

func setUp(t *testing.T) (*mock_event.MockRepository, *mock_subscription.MockRepository, event.UseCase) {
	ctrl := gomock.NewController(t)

	rep := mock_event.NewMockRepository(ctrl)
	repSub := mock_subscription.NewMockRepository(ctrl)

	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	sugar := l.Sugar()
	zap.NewAtomicLevelAt(zapcore.DebugLevel)

	uc := NewEvent(rep, repSub, logger.NewLogger(sugar))
	return rep, repSub, uc
}

func TestEventUseCase_GetAllEventsOk(t *testing.T) {
	rep, _, uc := setUp(t)
	rep.EXPECT().GetAllEvents(gomock.Any(), 1).Return(testAllEventsWithDateSQL, nil)

	allEv, err := uc.GetAllEvents(1)
	assert.Nil(t, err)
	assert.Equal(t, allEv, testAllEvents)
}

func TestEventUseCase_GetAllEventsError(t *testing.T) {
	rep, _, uc := setUp(t)
	rep.EXPECT().GetAllEvents(gomock.Any(), -1).Return(testAllEventsWithDateSQL, errors.New("invalid page number"))

	allEv, err := uc.GetAllEvents(-1)
	assert.NotNil(t, err)
	assert.Equal(t, allEv, models.EventCards{})
}

func TestEventUseCase_GetAllEventsZeroLength(t *testing.T) {
	rep, _, uc := setUp(t)
	rep.EXPECT().GetAllEvents(gomock.Any(), 1).Return([]models.EventCardWithDateSQL{}, nil)

	allEv, err := uc.GetAllEvents(1)
	assert.Nil(t, err)
	assert.Equal(t, allEv, models.EventCards{})
}

func TestEventUseCase_GetOneEventOk(t *testing.T) {
	rep, srep, uc := setUp(t)
	rep.EXPECT().GetOneEventByID(uint64(1)).Return(testEventSQL, nil)
	rep.EXPECT().GetTags(uint64(1)).Return(testTags, nil)
	srep.EXPECT().GetEventFollowers(uint64(1)).Return(testFollowers, nil)

	oneEv, err := uc.GetOneEvent(uint64(1))
	assert.Nil(t, err)
	assert.Equal(t, oneEv, testEvent)
}

func TestEventUseCase_GetOneEventErrorByID(t *testing.T) {
	rep, _, uc := setUp(t)
	rep.EXPECT().GetOneEventByID(uint64(1)).Return(models.EventSQL{}, errors.New("get event error"))

	oneEv, err := uc.GetOneEvent(uint64(1))
	assert.NotNil(t, err)
	assert.Equal(t, oneEv, models.Event{})
}

func TestEventUseCase_GetOneEventErrorTags(t *testing.T) {
	rep, _, uc := setUp(t)
	rep.EXPECT().GetOneEventByID(uint64(1)).Return(testEventSQL, nil)
	rep.EXPECT().GetTags(uint64(1)).Return(models.Tags{}, errors.New("get tags error"))

	oneEv, err := uc.GetOneEvent(uint64(1))
	assert.NotNil(t, err)
	assert.Equal(t, oneEv, testEventErrTags)
}

func TestEventUseCase_GetOneEventErrorFollowers(t *testing.T) {
	rep, srep, uc := setUp(t)
	rep.EXPECT().GetOneEventByID(uint64(1)).Return(testEventSQL, nil)
	rep.EXPECT().GetTags(uint64(1)).Return(testTags, nil)
	srep.EXPECT().GetEventFollowers(uint64(1)).Return(models.UsersOnEvent{}, errors.New("get followers error"))

	oneEv, err := uc.GetOneEvent(uint64(1))
	assert.NotNil(t, err)
	assert.Equal(t, oneEv, testEventErrFollowers)
}

func TestEventUseCase_Delete(t *testing.T) {
	rep, _, uc := setUp(t)
	rep.EXPECT().DeleteById(uint64(1)).Return(nil)

	err := uc.Delete(uint64(1))
	assert.Nil(t, err)
}

func TestEventUseCase_CreateNewEvent(t *testing.T) {
	rep, _, uc := setUp(t)
	rep.EXPECT().AddEvent(gomock.Any()).Return(nil)

	err := uc.CreateNewEvent(&models.Event{})
	assert.Nil(t, err)
}

func TestEventUseCase_SaveImageError(t *testing.T) {
	rep, _, uc := setUp(t)

	err := uc.SaveImage(uint64(1), &multipart.FileHeader{})
	assert.NotNil(t, err)
}

/*func TestEventUseCase_SaveImageOk(t *testing.T) {
	rep, _, uc := setUp(t)
	rep.EXPECT().UpdateEventAvatar(1, gomock.Any()).Return(nil)

	err := uc.SaveImage(uint64(1), **realfile**)
	assert.Nil(t, err)
}*/

func TestEventUseCase_GetEventsByCategoryOk(t *testing.T) {
	rep, _, uc := setUp(t)
	rep.EXPECT().GetEventsByCategory(ttype, gomock.Any(), 1).Return(testAllEventsWithDateSQL, nil)

	catEv, err := uc.GetEventsByCategory(ttype, 1)
	assert.Nil(t, err)
	assert.Equal(t, catEv, testAllEvents)
}

func TestEventUseCase_GetEventsByCategoryWithoutCategory(t *testing.T) {
	rep, _, uc := setUp(t)
	rep.EXPECT().GetAllEvents(gomock.Any(), 1).Return(testAllEventsWithDateSQL, nil)

	catEv, err := uc.GetEventsByCategory("", 1)
	assert.Nil(t, err)
	assert.Equal(t, catEv, testAllEvents)
}

func TestEventUseCase_GetEventsByCategoryError(t *testing.T) {
	rep, _, uc := setUp(t)
	rep.EXPECT().GetEventsByCategory(ttype, gomock.Any(), 1).Return([]models.EventCardWithDateSQL{}, errors.New("get category error"))

	catEv, err := uc.GetEventsByCategory(ttype, 1)
	assert.NotNil(t, err)
	assert.Equal(t, catEv, models.EventCards{})
}

func TestEventUseCase_GetEventsByCategoryZeroLength(t *testing.T) {
	rep, _, uc := setUp(t)
	rep.EXPECT().GetEventsByCategory(ttype, gomock.Any(), 1).Return([]models.EventCardWithDateSQL{}, nil)

	catEv, err := uc.GetEventsByCategory(ttype, 1)
	assert.Nil(t, err)
	assert.Equal(t, catEv, models.EventCards{})
}

/*func TestEventUseCase_GetImageOk(t *testing.T) {
	rep, _, uc := setUp(t)
	rep.EXPECT().GetOneEventByID(uint64(1)).Return(**eventWithImgExist**, nil)

	slicebyte, err := uc.GetImage(uint64(1))
	assert.Nil(t, err)
	assert.Equal(t, slicebyte, **file**)
}*/

func TestEventUseCase_GetImageUnexistingImg(t *testing.T) {
	rep, _, uc := setUp(t)
	rep.EXPECT().GetOneEventByID(uint64(1)).Return(testEventSQL, nil)

	slicebyte, err := uc.GetImage(uint64(1))
	assert.NotNil(t, err)
	assert.Equal(t, slicebyte, []byte{})
}

func TestEventUseCase_GetImageError(t *testing.T) {
	rep, _, uc := setUp(t)
	rep.EXPECT().GetOneEventByID(uint64(1)).Return(models.EventSQL{}, errors.New(""))

	slicebyte, err := uc.GetImage(uint64(1))
	assert.NotNil(t, err)
	assert.Equal(t, slicebyte, []byte{})
}

func TestEventUseCase_FindEventsOk(t *testing.T) {
	rep, _, uc := setUp(t)
	rep.EXPECT().FindEvents(search, gomock.Any(), 1).Return(testAllEventsWithDateSQL, nil)

	searchEv, err := uc.FindEvents(search, "", 1)
	assert.Nil(t, err)
	assert.Equal(t, searchEv, testAllEvents)
}

func TestEventUseCase_FindEventsZeroLengthWithCategory(t *testing.T) {
	rep, _, uc := setUp(t)
	rep.EXPECT().CategorySearch(search, category, gomock.Any(), 1).Return([]models.EventCardWithDateSQL{}, nil)

	searchEv, err := uc.FindEvents(search, category, 1)
	assert.Nil(t, err)
	assert.Equal(t, searchEv, models.EventCards{})
}

func TestEventUseCase_FindEventsError(t *testing.T) {
	rep, _, uc := setUp(t)
	rep.EXPECT().FindEvents(search, gomock.Any(), 1).Return([]models.EventCardWithDateSQL{}, errors.New("get find error"))

	searchEv, err := uc.FindEvents(search, "", 1)
	assert.NotNil(t, err)
	assert.Equal(t, searchEv, models.EventCards{})
}

func TestEventUseCase_RecomendSystemOk(t *testing.T) {
	rep, _, uc := setUp(t)
	rep.EXPECT().RecomendSystem(uint64(1), category).Return(nil)

	err := uc.RecomendSystem(uint64(1), category)
	assert.Nil(t, err)
}

func TestEventUseCase_RecomendSystemError(t *testing.T) {
	rep, _, uc := setUp(t)
	rep.EXPECT().RecomendSystem(uint64(1), category).Times(2).Return(errors.New("get recommend problem"))

	err := uc.RecomendSystem(uint64(1), category)
	assert.NotNil(t, err)
}

func TestEventUseCase_GetRecommendedOk(t *testing.T) {
	rep, _, uc := setUp(t)
	rep.EXPECT().GetRecommended(uint64(1), gomock.Any(), 1).Return(testAllEventsWithDateSQL, nil)

	recEv, err := uc.GetRecommended(uint64(1), 1)
	assert.Nil(t, err)
	assert.Equal(t, recEv, testAllEvents)
}

func TestEventUseCase_GetRecommendedError(t *testing.T) {
	rep, _, uc := setUp(t)
	rep.EXPECT().GetRecommended(uint64(1), gomock.Any(), 1).Return([]models.EventCardWithDateSQL{}, errors.New("get recommended problem"))

	recEv, err := uc.GetRecommended(uint64(1), 1)
	assert.NotNil(t, err)
	assert.Equal(t, recEv, models.EventCards{})
}

func TestEventUseCase_GetRecommendedZeroLength(t *testing.T) {
	rep, _, uc := setUp(t)
	rep.EXPECT().GetRecommended(uint64(1), gomock.Any(), 1).Return([]models.EventCardWithDateSQL{}, nil)

	recEv, err := uc.GetRecommended(uint64(1), 1)
	assert.Nil(t, err)
	assert.Equal(t, recEv, models.EventCards{})
}
