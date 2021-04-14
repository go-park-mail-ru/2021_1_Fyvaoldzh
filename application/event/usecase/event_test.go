package usecase

import (
	"database/sql"
	"kudago/application/event"
	mock_event "kudago/application/event/mocks"
	"kudago/application/models"
	mock_subscription "kudago/application/subscription/mocks"
	"kudago/pkg/logger"
	"log"
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
