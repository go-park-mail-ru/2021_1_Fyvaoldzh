package http

import (
	"bytes"
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
	strUserId = "1"
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

var testOtherUserProfile = &models.OtherUserProfile{
	Uid:       userId,
	Visited:   eventsVisited,
	Planning:  eventsPlanning,
	Followers: followers,
}

var testUserOnEvent = &models.UserOnEvent{
	Id:   userId,
	Name: name,
}

var testUsersOnEvent = &models.UsersOnEvent{*testUserOnEvent}

var testUserFront = &models.User{
	Login:    login,
	Password: frontPassword,
}

var testRegData = &models.RegData{
	Login:    login,
	Password: frontPassword,
}

func setUp(t *testing.T, url, method string) (echo.Context,
	UserHandler, *mock_user.MockUseCase,*mock_infrastructure.MockSessionTarantool) {
	e := echo.New()
	r := e.Router()
	r.Add(method, url, func(echo.Context) error { return nil })

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

	var req *http.Request
	switch method {
	case http.MethodPost:
		switch url{
		case "/api/v1/login":
			f, _ := testUserFront.MarshalJSON()
			req = httptest.NewRequest(http.MethodGet, url, bytes.NewBuffer(f))
		case "/api/v1/register":
			f, _ := testRegData.MarshalJSON()
			req = httptest.NewRequest(http.MethodGet, url, bytes.NewBuffer(f))
		case "/api/v1/update":
			f, _ := testOwnUserProfile.MarshalJSON()
			req = httptest.NewRequest(http.MethodGet, url, bytes.NewBuffer(f))
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


func TestUserHandler_GetOwnProfile(t *testing.T) {
	c, h, usecase, sm := setUp(t, "/api/v1/profile",  http.MethodGet)
	cookie := h.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)
	sm.EXPECT().CheckSession(cookie.Value).Return(true, userId, nil)
	usecase.EXPECT().GetOwnProfile(userId).Return(testOwnUserProfile, nil)

	err = h.GetOwnProfile(c)

	assert.Nil(t, err)
}

func TestUserHandler_GetOwnProfileErrorNoCookie(t *testing.T) {
	c, h, _, _ := setUp(t, "/api/v1/profile",  http.MethodGet)

	err = h.GetOwnProfile(c)

	assert.Error(t, err)
}

func TestUserHandler_GetOwnProfileSMFalse(t *testing.T) {
	c, h, _, sm := setUp(t, "/api/v1/profile",  http.MethodGet)
	cookie := h.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)
	sm.EXPECT().CheckSession(cookie.Value).Return(false, userId, nil)

	err = h.GetOwnProfile(c)

	assert.Error(t, err)
}

func TestUserHandler_GetOwnProfileErrorSM(t *testing.T) {
	c, h, _, sm := setUp(t, "/api/v1/profile",  http.MethodGet)
	cookie := h.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)
	sm.EXPECT().CheckSession(cookie.Value).Return(false, userId, echo.NewHTTPError(http.StatusInternalServerError))

	err = h.GetOwnProfile(c)

	assert.Error(t, err)
}

func TestUserHandler_GetOwnProfileErrorUCGetOwnProfile(t *testing.T) {
	c, h, usecase, sm := setUp(t, "/api/v1/profile",  http.MethodGet)
	cookie := h.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)
	sm.EXPECT().CheckSession(cookie.Value).Return(true, userId, nil)
	usecase.EXPECT().GetOwnProfile(userId).Return(testOwnUserProfile, echo.NewHTTPError(http.StatusInternalServerError))

	err = h.GetOwnProfile(c)

	assert.Error(t, err)
}

///////////////////////////////////////////////////

func TestUserHandler_GetOtherUserProfile(t *testing.T) {
	c, h, usecase, _ := setUp(t, "/api/v1/profile/:id",  http.MethodGet)
	c.SetParamNames("id")
	c.SetParamValues(strUserId)

	usecase.EXPECT().GetOtherProfile(userId).Return(testOtherUserProfile, nil)

	err = h.GetOtherUserProfile(c)

	assert.Nil(t, err)
}

func TestUserHandler_GetOtherUserProfileErrorAtoi(t *testing.T) {
	c, h, _, _ := setUp(t, "/api/v1/profile/:id",  http.MethodGet)
	c.SetParamNames("id")
	c.SetParamValues("a")

	err = h.GetOtherUserProfile(c)

	assert.Error(t, err, echo.NewHTTPError(http.StatusBadRequest))
}

func TestUserHandler_GetOtherUserProfileErrorMinus(t *testing.T) {
	c, h, _, _ := setUp(t, "/api/v1/profile/:id",  http.MethodGet)
	c.SetParamNames("id")
	c.SetParamValues("-1")

	err = h.GetOtherUserProfile(c)

	assert.Error(t, err)
}

func TestUserHandler_GetOtherUserProfileErrorUC(t *testing.T) {
	c, h, usecase, _ := setUp(t, "/api/v1/profile/:id",  http.MethodGet)
	c.SetParamNames("id")
	c.SetParamValues(strUserId)

	usecase.EXPECT().GetOtherProfile(userId).Return(testOtherUserProfile, echo.NewHTTPError(http.StatusInternalServerError))

	err = h.GetOtherUserProfile(c)

	assert.Error(t, err)
}


///////////////////////////////////////////////////

func TestUserHandler_GetUsers(t *testing.T) {
	c, h, usecase, _ := setUp(t, "/api/v1/users?page=1",  http.MethodGet)

	usecase.EXPECT().GetUsers(1).Return(*testUsersOnEvent, nil)

	err = h.GetUsers(c)

	assert.Nil(t, err)
}

func TestUserHandler_GetUsersErrorAtoi(t *testing.T) {
	c, h, _, _ := setUp(t, "/api/v1/users?page=a",  http.MethodGet)

	err = h.GetUsers(c)

	assert.Error(t, err)
}

func TestUserHandler_GetUsersErrorMinus(t *testing.T) {
	c, h, _, _ := setUp(t, "/api/v1/users?page=-1",  http.MethodGet)

	err = h.GetUsers(c)

	assert.Error(t, err)
}

func TestUserHandler_GetUsersErrorUC(t *testing.T) {
	c, h, usecase, _ := setUp(t, "/api/v1/users?page=1",  http.MethodGet)

	usecase.EXPECT().GetUsers(1).Return(*testUsersOnEvent, echo.NewHTTPError(http.StatusInternalServerError))

	err = h.GetUsers(c)

	assert.Error(t, err)
}

///////////////////////////////////////////////////

func TestUserHandler_GetAvatar(t *testing.T) {
	c, h, usecase, _ := setUp(t, "/api/v1/avatar/:id",  http.MethodGet)
	c.SetParamNames("id")
	c.SetParamValues(strUserId)

	usecase.EXPECT().GetAvatar(userId).Return([]byte{}, nil)

	err = h.GetAvatar(c)

	assert.Nil(t, err)
}

func TestUserHandler_GetAvatarErrorUC(t *testing.T) {
	c, h, usecase, _ := setUp(t, "/api/v1/avatar/:id",  http.MethodGet)
	c.SetParamNames("id")
	c.SetParamValues(strUserId)

	usecase.EXPECT().GetAvatar(userId).Return([]byte{}, echo.NewHTTPError(http.StatusInternalServerError))

	err = h.GetAvatar(c)

	assert.Error(t, err)
}

func TestUserHandler_GetAvatarErrorAtoi(t *testing.T) {
	c, h, _, _ := setUp(t, "/api/v1/avatar/:id",  http.MethodGet)
	c.SetParamNames("id")
	c.SetParamValues("tt")

	err = h.GetAvatar(c)

	assert.Error(t, err)
}

///////////////////////////////////////////////////

func TestUserHandler_Login(t *testing.T) {
	c, h, usecase, sm := setUp(t, "/api/v1/login",  http.MethodPost)

	usecase.EXPECT().Login(testUserFront).Return(userId, nil)
	sm.EXPECT().InsertSession(userId, gomock.Any()).Return(nil)

	err = h.Login(c)

	assert.Nil(t, err)
}

func TestUserHandler_LoginErrorAlreadyLogin(t *testing.T) {
	c, h, usecase, sm := setUp(t, "/api/v1/login",  http.MethodPost)
	cookie := h.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)

	usecase.EXPECT().Login(testUserFront).Return(userId, nil)
	sm.EXPECT().CheckSession(gomock.Any()).Return(true, userId, nil)

	err = h.Login(c)

	assert.Error(t, err)
}

func TestUserHandler_LoginErrorSMCheckSession(t *testing.T) {
	c, h, usecase, sm := setUp(t, "/api/v1/login",  http.MethodPost)
	cookie := h.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)

	usecase.EXPECT().Login(testUserFront).Return(userId, nil)
	sm.EXPECT().CheckSession(gomock.Any()).Return(true, userId, echo.NewHTTPError(http.StatusInternalServerError))

	err = h.Login(c)

	assert.Error(t, err)
}

func TestUserHandler_LoginErrorUC(t *testing.T) {
	c, h, usecase, _ := setUp(t, "/api/v1/login",  http.MethodPost)
	cookie := h.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)

	usecase.EXPECT().Login(testUserFront).Return(userId, echo.NewHTTPError(http.StatusInternalServerError))

	err = h.Login(c)

	assert.Error(t, err)
}

func TestUserHandler_LoginErrorSMInsertSession(t *testing.T) {
	c, h, usecase, sm := setUp(t, "/api/v1/login",  http.MethodPost)

	usecase.EXPECT().Login(testUserFront).Return(userId, nil)
	sm.EXPECT().InsertSession(userId, gomock.Any()).Return(echo.NewHTTPError(http.StatusInternalServerError))

	err = h.Login(c)

	assert.Error(t, err)
}

///////////////////////////////////////////////////

func TestUserHandler_Logout(t *testing.T) {
	c, h, _, sm := setUp(t, "/api/v1/login",  http.MethodDelete)
	cookie := h.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)

	sm.EXPECT().CheckSession(cookie.Value).Return(true, userId, nil)
	sm.EXPECT().DeleteSession(cookie.Value).Return( nil)

	err = h.Logout(c)

	assert.Nil(t, err)
}

func TestUserHandler_LogoutErrorSMDeleteSession(t *testing.T) {
	c, h, _, sm := setUp(t, "/api/v1/login",  http.MethodDelete)
	cookie := h.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)

	sm.EXPECT().CheckSession(cookie.Value).Return(true, userId, nil)
	sm.EXPECT().DeleteSession(cookie.Value).Return( echo.NewHTTPError(http.StatusInternalServerError))

	err = h.Logout(c)

	assert.Error(t, err)
}

func TestUserHandler_LogoutUnauthorized(t *testing.T) {
	c, h, _, sm := setUp(t, "/api/v1/login",  http.MethodDelete)
	cookie := h.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)

	sm.EXPECT().CheckSession(cookie.Value).Return(false, userId, nil)

	err = h.Logout(c)

	assert.Error(t, err)
}

func TestUserHandler_LogoutErrorSMCheckSession(t *testing.T) {
	c, h, _, sm := setUp(t, "/api/v1/login",  http.MethodDelete)
	cookie := h.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)

	sm.EXPECT().CheckSession(cookie.Value).Return(true, userId, echo.NewHTTPError(http.StatusInternalServerError))

	err = h.Logout(c)

	assert.Error(t, err)
}

func TestUserHandler_LogoutNoCookie(t *testing.T) {
	c, h, _, _ := setUp(t, "/api/v1/login",  http.MethodDelete)

	err = h.Logout(c)

	assert.Error(t, err)
}

///////////////////////////////////////////////////

func TestUserHandler_Register(t *testing.T) {
	c, h, usecase, sm := setUp(t, "/api/v1/register",  http.MethodPost)

	usecase.EXPECT().Add(testRegData).Return(userId, nil)
	sm.EXPECT().InsertSession(userId, gomock.Any()).Return(nil)

	err = h.Register(c)

	assert.Nil(t, err)
}

func TestUserHandler_RegisterErrorSMInsertSession(t *testing.T) {
	c, h, usecase, sm := setUp(t, "/api/v1/register",  http.MethodPost)

	usecase.EXPECT().Add(testRegData).Return(userId, nil)
	sm.EXPECT().InsertSession(userId, gomock.Any()).Return(echo.NewHTTPError(http.StatusInternalServerError))

	err = h.Register(c)

	assert.Error(t, err)
}


func TestUserHandler_RegisterErrorUCAdd(t *testing.T) {
	c, h, usecase, _ := setUp(t, "/api/v1/register",  http.MethodPost)

	usecase.EXPECT().Add(testRegData).Return(userId, echo.NewHTTPError(http.StatusInternalServerError))

	err = h.Register(c)

	assert.Error(t, err)
}

func TestUserHandler_RegisterLoggedIn(t *testing.T) {
	c, h, _, sm := setUp(t, "/api/v1/register",  http.MethodPost)
	cookie := h.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)

	sm.EXPECT().CheckSession(gomock.Any()).Return(true, userId, nil)

	err = h.Register(c)

	assert.Error(t, err)
}

func TestUserHandler_RegisterErrorSMCheckSession(t *testing.T) {
	c, h, _, sm := setUp(t, "/api/v1/register",  http.MethodPost)
	cookie := h.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)

	sm.EXPECT().CheckSession(gomock.Any()).Return(true, userId, echo.NewHTTPError(http.StatusInternalServerError))

	err = h.Register(c)

	assert.Error(t, err)
}

///////////////////////////////////////////////////

func TestUserHandler_Update(t *testing.T) {
	c, h, usecase, sm := setUp(t, "/api/v1/update",  http.MethodPost)
	cookie := h.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)

	usecase.EXPECT().Update(userId, testOwnUserProfile).Return(nil)
	sm.EXPECT().CheckSession(gomock.Any()).Return(true, userId, nil)

	err = h.Update(c)

	assert.Nil(t, err)
}

func TestUserHandler_UpdateErrorUC(t *testing.T) {
	c, h, usecase, sm := setUp(t, "/api/v1/update",  http.MethodPost)
	cookie := h.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)

	usecase.EXPECT().Update(userId, testOwnUserProfile).Return(echo.NewHTTPError(http.StatusInternalServerError))
	sm.EXPECT().CheckSession(gomock.Any()).Return(true, userId, nil)

	err = h.Update(c)

	assert.Error(t, err)
}

func TestUserHandler_UpdateSessionNotExists(t *testing.T) {
	c, h, _, sm := setUp(t, "/api/v1/update",  http.MethodPost)
	cookie := h.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)

	sm.EXPECT().CheckSession(gomock.Any()).Return(false, userId, nil)

	err = h.Update(c)

	assert.Error(t, err)
}

func TestUserHandler_UpdateErrorSM(t *testing.T) {
	c, h, _, sm := setUp(t, "/api/v1/update",  http.MethodPost)
	cookie := h.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)

	sm.EXPECT().CheckSession(gomock.Any()).Return(false, userId, echo.NewHTTPError(http.StatusInternalServerError))

	err = h.Update(c)

	assert.Error(t, err)
}

func TestUserHandler_UpdateNoCookie(t *testing.T) {
	c, h, _, _ := setUp(t, "/api/v1/update",  http.MethodPost)

	err = h.Update(c)

	assert.Error(t, err)
}