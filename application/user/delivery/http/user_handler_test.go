package http

import (
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/microcosm-cc/bluemonday"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"kudago/application/models"
	mock_user "kudago/application/user/mocks"
	"kudago/pkg/constants"
	"kudago/pkg/custom_sanitizer"

	mock_infrastructure "kudago/pkg/infrastructure/mocks"
	"kudago/pkg/logger"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var (
	userId        uint64  = 1
	pageNum         = 1
	login           = "userlogin"
	name            = "username"
	frontPassword   = "123456"
	backPassword    = "IvJrQEdIeoTzLsMX_839spM7MzaXS7aJ_b3xTzmYqbotq3HRKAs="
	badBackPassword = "1111IvJrQEdIeoTzLsMX_839spM7MzaXS7aJ_b3xTzmYqbotq3HRKAs="
	email           = "email@mail.ru"
	birthdayStr     = "1999-01-01"
	birthday, err   = time.Parse(constants.TimeFormat, "1999-01-01")
	city            = "City"
	about           = "some personal information"
	avatar          = "public/users/default.png"
	imageName       = "image.png"
	evPlanningSQL   = models.EventCardWithDateSQL{
		ID:        1,
		StartDate: time.Now(),
		EndDate:   time.Now().Add(10 * time.Hour),
	}
	evPlanning = models.EventCard{
		ID:        1,
		StartDate: evPlanningSQL.StartDate.String(),
		EndDate:   evPlanningSQL.EndDate.String(),
	}
	evVisitedSQL = models.EventCardWithDateSQL{
		ID:        2,
		StartDate: time.Now(),
		EndDate:   time.Now(),
	}
	evVisited = models.EventCard{
		ID:        2,
		StartDate: evVisitedSQL.StartDate.String(),
		EndDate:   evVisitedSQL.EndDate.String(),
	}
	eventsPlanningSQL = []models.EventCardWithDateSQL{
		evPlanningSQL, evVisitedSQL,
	}
	eventsVisitedSQL = []models.EventCardWithDateSQL{
		evVisitedSQL,
	}
	eventsPlanning = []models.EventCard{
		evPlanning,
	}
	eventsVisited = []models.EventCard{
		evVisited,
	}

	followers = []uint64{2, 2, 3}
)

var testOwnUserProfile = &models.UserOwnProfile{
	Uid:       userId,
	Login:     login,
	Visited:   eventsVisited,
	Planning:  eventsPlanning,
	Followers: followers,
}

func setUp(t *testing.T, url, method string) (echo.Context,
	UserHandler, *mock_user.MockUseCase,*mock_infrastructure.MockSessionTarantool) {
	e := echo.New()
	r := e.Router()
	r.Add(method, url, func(echo.Context) error { return nil })
	var req *http.Request
	switch method {
	//case http.MethodPost:
		//f, _ := testFilm.MarshalJSON()
		//req = httptest.NewRequest(http.MethodGet, url, bytes.NewBuffer(f))
	case http.MethodGet:
		req = httptest.NewRequest(http.MethodGet, url, nil)
	}
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(url)
	ctrl := gomock.NewController(t)
	_ = mock_user.NewMockRepository(ctrl)
	usecase := mock_user.NewMockUseCase(ctrl)
	sm := mock_infrastructure.NewMockSessionTarantool(ctrl)

	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	sugar := l.Sugar()
	zap.NewAtomicLevelAt(zapcore.DebugLevel)
	cs := custom_sanitizer.NewCustomSanitizer(bluemonday.UGCPolicy())

	handler := UserHandler{
		UseCase:   usecase,
		Sm: sm,
		Logger:    logger.NewLogger(sugar),
		sanitizer: cs,
	}

	return c, handler, usecase, sm
}


func TestUserHandler_GetOwnProfile(t *testing.T) {
	ctx, h, usecase, sm := setUp(t, "/api/v1/profile",  http.MethodGet)
	ctx.SetCookie(h.CreateCookie(constants.CookieLength))
	cookie, _ := ctx.Cookie(constants.SessionCookieName)
	sm.EXPECT().CheckSession(cookie.Value).Return(true, userId, nil)
	usecase.EXPECT().GetOwnProfile(userId).Return(testOwnUserProfile)

	err = h.GetOwnProfile(ctx)

	assert.Nil(t, err)
}
