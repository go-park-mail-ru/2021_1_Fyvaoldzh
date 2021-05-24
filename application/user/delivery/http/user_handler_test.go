package http

import (
	"bytes"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/microcosm-cc/bluemonday"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"kudago/application/microservices/auth/client"
	"kudago/application/models"
	middleware1 "kudago/application/server/middleware"
	mock_user "kudago/application/user/mocks"
	"kudago/pkg/constants"
	"kudago/pkg/custom_sanitizer"
	"kudago/pkg/generator"
	"kudago/pkg/logger"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	pageNum        = 1
	userId  uint64 = 1

	login         = "userlogin"
	name          = "username"
	frontPassword = "123456"

	followers = uint64(1)
)

var testOwnUserProfile = &models.UserOwnProfile{
	Uid:       userId,
	Login:     login,
	Followers: followers,
}

var testOtherUserProfile = &models.OtherUserProfile{
	Uid:       userId,
	Followers: followers,
}

var testUserFront = &models.User{
	Login:    login,
	Password: frontPassword,
}

var testRegData = &models.RegData{
	Login:    login,
	Password: frontPassword,
}

var testUserCard = &models.UserCard{
	Id:   userId,
	Name: name,
}

var testUserCards = models.UserCards{*testUserCard}

func setUp(t *testing.T, url, method string) (echo.Context,
	UserHandler, *mock_user.MockUseCase, *client.MockIAuthClient) {
	e := echo.New()
	r := e.Router()
	r.Add(method, url, func(echo.Context) error { return nil })

	ctrl := gomock.NewController(t)
	usecase := mock_user.NewMockUseCase(ctrl)
	rpcAuth := client.NewMockIAuthClient(ctrl)

	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	sugar := l.Sugar()
	zap.NewAtomicLevelAt(zapcore.DebugLevel)
	cs := custom_sanitizer.NewCustomSanitizer(bluemonday.UGCPolicy())
	auth := middleware1.NewAuth(rpcAuth)

	handler := UserHandler{
		UseCase:   usecase,
		rpcAuth:   rpcAuth,
		auth:      auth,
		Logger:    logger.NewLogger(sugar),
		sanitizer: cs,
	}

	var req *http.Request
	switch method {
	case http.MethodPost:
		switch url {
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

	return c, handler, usecase, rpcAuth
}

func TestUserHandler_GetOwnProfile(t *testing.T) {
	c, h, usecase, _ := setUp(t, "/api/v1/profile", http.MethodGet)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Set(constants.SessionCookieName, cookie.Value)
	c.Set(constants.UserIdKey, userId)

	usecase.EXPECT().GetOwnProfile(userId).Return(testOwnUserProfile, nil)

	err := h.GetOwnProfile(c)

	assert.Nil(t, err)
}

/*
func TestUserHandler_GetOwnProfileErrorNoCookie(t *testing.T) {
	c, h, _, _ := setUp(t, "/api/v1/profile", http.MethodGet)

	err := h.GetOwnProfile(c)

	assert.Error(t, err)
}

*/

/*
func TestUserHandler_GetOwnProfileRpcAuthFalse(t *testing.T) {
	c, h, _, rpcAuth := setUp(t, "/api/v1/profile", http.MethodGet)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Set(constants.SessionCookieName, cookie.Value)
	c.Set(constants.UserIdKey, userId)
	rpcAuth.EXPECT().Check(cookie.Value).Return(false, userId, nil)

	err := h.GetOwnProfile(c)

	assert.Error(t, err)
}


*/

/*
func TestUserHandler_GetOwnProfileErrorRpcAuth(t *testing.T) {
	c, h, _, rpcAuth := setUp(t, "/api/v1/profile", http.MethodGet)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Set(constants.SessionCookieName, cookie.Value)
	c.Set(constants.UserIdKey, userId)
	rpcAuth.EXPECT().Check(cookie.Value).Return(false, userId, echo.NewHTTPError(http.StatusInternalServerError))

	err := h.GetOwnProfile(c)

	assert.Error(t, err)
}

*/

func TestUserHandler_GetOwnProfileErrorUCGetOwnProfile(t *testing.T) {
	c, h, usecase, _ := setUp(t, "/api/v1/profile", http.MethodGet)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Set(constants.SessionCookieName, cookie.Value)
	c.Set(constants.UserIdKey, userId)
	usecase.EXPECT().GetOwnProfile(userId).Return(testOwnUserProfile, echo.NewHTTPError(http.StatusBadRequest))

	err := h.GetOwnProfile(c)

	assert.Error(t, err)
}

///////////////////////////////////////////////////

func TestUserHandler_GetOtherUserProfile(t *testing.T) {
	c, h, usecase, _ := setUp(t, "/api/v1/profile/:id", http.MethodGet)
	c.Set(constants.IdKey, pageNum)

	usecase.EXPECT().GetOtherProfile(userId).Return(testOtherUserProfile, nil)

	err := h.GetOtherUserProfile(c)

	assert.Nil(t, err)
}

/*
func TestUserHandler_GetOtherUserProfileErrorAtoi(t *testing.T) {
	c, h, _, _ := setUp(t, "/api/v1/profile/:id", http.MethodGet)
	c.SetParamNames("id")
	c.SetParamValues("a")

	err := h.GetOtherUserProfile(c)

	assert.Error(t, err, echo.NewHTTPError(http.StatusBadRequest))
}

*/

/*
func TestUserHandler_GetOtherUserProfileErrorMinus(t *testing.T) {
	c, h, _, _ := setUp(t, "/api/v1/profile/:id", http.MethodGet)
	c.SetParamNames("id")
	c.SetParamValues("-1")

	err := h.GetOtherUserProfile(c)

	assert.Error(t, err)
}

*/

func TestUserHandler_GetOtherUserProfileErrorUC(t *testing.T) {
	c, h, usecase, _ := setUp(t, "/api/v1/profile/:id", http.MethodGet)
	c.Set(constants.IdKey, pageNum)

	usecase.EXPECT().GetOtherProfile(userId).Return(testOtherUserProfile, echo.NewHTTPError(http.StatusInternalServerError))

	err := h.GetOtherUserProfile(c)

	assert.Error(t, err)
}

///////////////////////////////////////////////////

func TestUserHandler_GetUsers(t *testing.T) {
	c, h, usecase, _ := setUp(t, "/api/v1/users", http.MethodGet)
	c.Set(constants.PageKey, pageNum)
	usecase.EXPECT().GetUsers(1).Return(testUserCards, nil)

	err := h.GetUsers(c)

	assert.Nil(t, err)
}

/*
func TestUserHandler_GetUsersErrorAtoi(t *testing.T) {
	c, h, _, _ := setUp(t, "/api/v1/users?page=a", http.MethodGet)

	err := h.GetUsers(c)

	assert.Error(t, err)
}

*/

/*
func TestUserHandler_GetUsersErrorMinus(t *testing.T) {
	c, h, _, _ := setUp(t, "/api/v1/users?page=-1", http.MethodGet)

	err := h.GetUsers(c)

	assert.Error(t, err)
}


*/
func TestUserHandler_GetUsersErrorUC(t *testing.T) {
	c, h, usecase, _ := setUp(t, "/api/v1/users", http.MethodGet)
	c.Set(constants.PageKey, pageNum)

	usecase.EXPECT().GetUsers(1).Return(testUserCards, echo.NewHTTPError(http.StatusInternalServerError))

	err := h.GetUsers(c)

	assert.Error(t, err)
}

///////////////////////////////////////////////////

func TestUserHandler_GetAvatar(t *testing.T) {
	c, h, usecase, _ := setUp(t, "/api/v1/avatar/:id", http.MethodGet)
	c.Set(constants.IdKey, pageNum)

	usecase.EXPECT().GetAvatar(userId).Return([]byte{}, nil)

	err := h.GetAvatar(c)

	assert.Nil(t, err)
}

func TestUserHandler_GetAvatarErrorUC(t *testing.T) {
	c, h, usecase, _ := setUp(t, "/api/v1/avatar/:id", http.MethodGet)
	c.Set(constants.IdKey, pageNum)

	usecase.EXPECT().GetAvatar(userId).Return([]byte{}, echo.NewHTTPError(http.StatusInternalServerError))

	err := h.GetAvatar(c)

	assert.Error(t, err)
}

/*
func TestUserHandler_GetAvatarErrorAtoi(t *testing.T) {
	c, h, _, _ := setUp(t, "/api/v1/avatar/:id", http.MethodGet)
	c.SetParamNames("id")
	c.SetParamValues("tt")

	err := h.GetAvatar(c)

	assert.Error(t, err)
}

*/

///////////////////////////////////////////////////

func TestUserHandler_Login(t *testing.T) {
	c, h, _, rpcAuth := setUp(t, "/api/v1/login", http.MethodPost)

	rpcAuth.EXPECT().Login(gomock.Any(), gomock.Any(), "").Return(userId, "lalala", nil, 200)

	err := h.Login(c)

	assert.Nil(t, err)
}

func TestUserHandler_LoginCookieNotNil(t *testing.T) {
	c, h, _, rpcAuth := setUp(t, "/api/v1/login", http.MethodPost)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)

	rpcAuth.EXPECT().Login(gomock.Any(), gomock.Any(), cookie.Value).Return(userId, "lalala", nil, 200)

	err := h.Login(c)

	assert.Nil(t, err)
}

/*
func TestUserHandler_LoginErrorAlreadyLogin(t *testing.T) {
	c, h, _, _ := setUp(t, "/api/v1/login", http.MethodPost)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Set(constants.SessionCookieName, cookie.Value)
	c.Set(constants.UserIdKey, userId)

	err := h.Login(c)

	assert.Error(t, err)
}

*/

/*


func TestUserHandler_LoginErrorRpcAuthCheckSession(t *testing.T) {
	c, h, usecase, rpcAuth := setUp(t, "/api/v1/login", http.MethodPost)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Set(constants.SessionCookieName, cookie.Value)
	c.Set(constants.UserIdKey, userId)

	usecase.EXPECT().Login(testUserFront).Return(userId, nil)
	rpcAuth.EXPECT().Check(gomock.Any()).Return(true, userId, echo.NewHTTPError(http.StatusInternalServerError))

	err := h.Login(c)

	assert.Error(t, err)
}
*/

func TestUserHandler_LoginErrorRpcAuthLogin(t *testing.T) {
	c, h, _, rpcAuth := setUp(t, "/api/v1/login", http.MethodPost)

	rpcAuth.EXPECT().Login(gomock.Any(), gomock.Any(), "").Return(userId, login, echo.NewHTTPError(http.StatusInternalServerError), 500)

	err := h.Login(c)

	assert.Error(t, err)
}

///////////////////////////////////////////////////

func TestUserHandler_Logout(t *testing.T) {
	c, h, _, rpcAuth := setUp(t, "/api/v1/login", http.MethodDelete)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)

	rpcAuth.EXPECT().Logout(cookie.Value).Return(nil, 200)

	err := h.Logout(c)

	assert.Nil(t, err)
}

func TestUserHandler_LogoutErrorRpcAuthLogout(t *testing.T) {
	c, h, _, rpcAuth := setUp(t, "/api/v1/login", http.MethodDelete)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)

	rpcAuth.EXPECT().Logout(cookie.Value).Return(echo.NewHTTPError(http.StatusInternalServerError), 500)

	err := h.Logout(c)

	assert.Error(t, err)
}

///////////////////////////////////////////////////

func TestUserHandler_Register(t *testing.T) {
	c, h, usecase, rpcAuth := setUp(t, "/api/v1/register", http.MethodPost)

	usecase.EXPECT().Add(testRegData).Return(userId, nil)
	rpcAuth.EXPECT().Login(gomock.Any(), gomock.Any(), "").Return(userId, "lalala", nil, 200)

	err := h.Register(c)

	assert.Nil(t, err)
}

func TestUserHandler_RegisterErrorRpcAuthLogin(t *testing.T) {
	c, h, usecase, rpcAuth := setUp(t, "/api/v1/register", http.MethodPost)

	usecase.EXPECT().Add(testRegData).Return(userId, nil)
	rpcAuth.EXPECT().Login(gomock.Any(), gomock.Any(), "").Return(userId, "", echo.NewHTTPError(http.StatusInternalServerError), 500)

	err := h.Register(c)

	assert.Error(t, err)
}

func TestUserHandler_RegisterErrorUCAdd(t *testing.T) {
	c, h, usecase, _ := setUp(t, "/api/v1/register", http.MethodPost)

	usecase.EXPECT().Add(testRegData).Return(userId, echo.NewHTTPError(http.StatusInternalServerError))

	err := h.Register(c)

	assert.Error(t, err)
}

func TestUserHandler_RegisterLoggedIn(t *testing.T) {
	c, h, _, rpcAuth := setUp(t, "/api/v1/register", http.MethodPost)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)

	rpcAuth.EXPECT().Check(gomock.Any()).Return(true, userId, nil, 200)

	err := h.Register(c)

	assert.Error(t, err)
}

func TestUserHandler_RegisterErrorRpcAuthCheck(t *testing.T) {
	c, h, _, rpcAuth := setUp(t, "/api/v1/register", http.MethodPost)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)

	rpcAuth.EXPECT().Check(gomock.Any()).Return(true, userId, echo.NewHTTPError(http.StatusInternalServerError), 500)

	err := h.Register(c)

	assert.Error(t, err)
}

///////////////////////////////////////////////////

func TestUserHandler_Update(t *testing.T) {
	c, h, usecase, _ := setUp(t, "/api/v1/update", http.MethodPost)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Set(constants.SessionCookieName, cookie.Value)
	c.Set(constants.UserIdKey, userId)

	usecase.EXPECT().Update(userId, testOwnUserProfile).Return(nil)

	err := h.Update(c)

	assert.Nil(t, err)
}

func TestUserHandler_UpdateErrorUC(t *testing.T) {
	c, h, usecase, _ := setUp(t, "/api/v1/update", http.MethodPost)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Set(constants.SessionCookieName, cookie.Value)
	c.Set(constants.UserIdKey, userId)

	usecase.EXPECT().Update(userId, testOwnUserProfile).Return(echo.NewHTTPError(http.StatusInternalServerError))

	err := h.Update(c)

	assert.Error(t, err)
}

/*
func TestUserHandler_UpdateSessionNotExists(t *testing.T) {
	c, h, _, rpcAuth := setUp(t, "/api/v1/update", http.MethodPost)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)

	rpcAuth.EXPECT().Check(gomock.Any()).Return(false, userId, nil)

	err := h.Update(c)

	assert.Error(t, err)
}
*/

/*

func TestUserHandler_UpdateErrorSM(t *testing.T) {
	c, h, _, rpcAuth := setUp(t, "/api/v1/update", http.MethodPost)
	cookie := generator.CreateCookie(constants.CookieLength)
	c.Request().AddCookie(cookie)

	rpcAuth.EXPECT().Check(gomock.Any()).Return(false, userId, echo.NewHTTPError(http.StatusInternalServerError))

	err := h.Update(c)

	assert.Error(t, err)
}
*/

/*
func TestUserHandler_UpdateNoCookie(t *testing.T) {
	c, h, _, _ := setUp(t, "/api/v1/update", http.MethodPost)

	err := h.Update(c)

	assert.Error(t, err)
}

*/

///////////////////////////////////////////////////

func TestUserHandler_FindUsers(t *testing.T) {
	c, h, usecase, _ := setUp(t, "/api/v1/find", http.MethodGet)
	c.Set(constants.PageKey, pageNum)
	usecase.EXPECT().FindUsers(gomock.Any(), pageNum).Return(testUserCards, nil)

	err := h.FindUsers(c)

	assert.Nil(t, err)
}

func TestUserHandler_FindUsersErrorUCFindUsers(t *testing.T) {
	c, h, usecase, _ := setUp(t, "/api/v1/find", http.MethodGet)
	c.Set(constants.PageKey, pageNum)
	usecase.EXPECT().FindUsers(gomock.Any(), pageNum).Return(testUserCards, echo.NewHTTPError(http.StatusInternalServerError))

	err := h.FindUsers(c)

	assert.Error(t, err)
}

///////////////////////////////////////////////////
