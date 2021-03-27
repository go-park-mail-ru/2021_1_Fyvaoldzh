package user

import (
	"bytes"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/require"
	"kudago/models"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

var h = &UserHandler{
	UserBase:    UserBase,
	ProfileBase: ProfileBase,
	UserEvent:   EventUserBase,
	Store:       make(map[string]uint64),
	Mu:          &sync.Mutex{},
}

// -----------------handler-----------------

func TestHandlerUser_CreateUserProfileTrue(t *testing.T) {
	newData := &models.RegData{Login: "morozik_jr", Password: "123456", Name: "Morozik Jr."}
	uid, err := h.CreateUserProfile(newData)
	require.Equal(t, uint64(4), uid, "id is not equal with correct reg data")
	require.Equal(t, nil, err, "got error with correct reg data")
}

func TestHandlerUser_CreateUserProfileFalse(t *testing.T) {
	newData := &models.RegData{Login: h.UserBase[0].Login, Password: "123456", Name: "Morozik Jr."}
	uid, err := h.CreateUserProfile(newData)
	require.Equal(t, uint64(0), uid, "got id for new user with incorrect reg data")
	require.NotEqual(t, nil, err, "did not get error with incorrect reg data")
}

// -----------------helpful functions-----------------

func TestDifferentKeys(t *testing.T) {
	key1 := RandStringRunes(32)
	key2 := RandStringRunes(32)

	require.NotEqual(t, key1, key2, "keys are not different")
}

// -----------------users-----------------

func TestGetUserTrue(t *testing.T) {
	user := GetUser(h, h.UserBase[0].Id)

	require.Equal(t, h.UserBase[0], user, "users are not equal")
}

func TestGetUserFalse(t *testing.T) {
	user := GetUser(h, h.UserBase[len(h.UserBase) - 1].Id + 1)

	require.Equal(t, &models.User{}, user, "got user that does not exist")
}

func TestIsExistingUserTrue(t *testing.T) {
	user := h.UserBase[0]

	require.Equal(t, true, IsExistingUser(h, user), "existing user does not exist")
}

func TestIsExistingUserFalse(t *testing.T) {
	user := &models.User{Id: 1, Login: "morozik", Password: "123456"}

	require.Equal(t, false, IsExistingUser(h, user), "not existing user does exist")
}

func TestIsCorrectUserTrue(t *testing.T) {
	user := h.UserBase[0]

	flag, uid := IsCorrectUser(h, user)
	require.Equal(t, true, flag, "correct user is not correct")
	require.Equal(t, uint64(1), uid, "got incorrect user id")
}

func TestIsCorrectUserFalse(t *testing.T) {
	user := &models.User{Id: 1, Login: "moroz", Password: "1234567"}

	flag, uid := IsCorrectUser(h, user)
	require.Equal(t, false, flag, "incorrect user data is correct")
	require.Equal(t, uint64(0), uid, "incorrect user data is correct")
}

func TestCreateCookie(t *testing.T) {
	uid := uint64(1)
	_ = CreateCookie(h, 32, uid)

	flag := false

	for _, value := range h.Store {
		if value == uid {
			flag = true
			break
		}
	}

	require.Equal(t, true, flag, "cookie is not saved in the base")
}

// -----------------profiles-----------------

func TestGetProfileTrue(t *testing.T) {
	profile := h.ProfileBase[0]

	require.Equal(t, profile, GetProfile(h, profile.Uid), "existing mail does not exist")
}

func TestGetProfileFalse(t *testing.T) {
	require.Equal(t, &models.UserOwnProfile{}, GetProfile(h, h.ProfileBase[len(h.ProfileBase) - 1].Uid + 1), "returns not existing profile")
}

func TestChangesProfileDataTrue(t *testing.T) {
	profile := &models.UserData{Name: "Матрос", Birthday: "6 марта 1999 г.", City: "Санкт-Петербург", Email: "matrosik@mail.ru",
		About: "люблю котиков", Avatar: "1default.png"}

	require.Equal(t, nil, changeProfileData(h, profile, h.UserBase[0].Id), "returns error on correct data")
}

func TestChangesProfileDataFalse(t *testing.T) {
	profile := &models.UserData{Email: h.ProfileBase[0].Email}

	require.NotEqual(t, nil, changeProfileData(h, profile, 2), "does not catch changed repeated mail")
}

func TestGetOtherUserProfileTrue(t *testing.T) {
	profile := &models.OtherUserProfile{Uid: 2, Name: "Матрос Матросович Матросов", Age: 20, City: "Санкт-Петербург", Followers: 1000, About: "главный матрос на корабле", Avatar: "1.png"}

	require.Equal(t, profile, GetOtherUserProfile(h, 2), "gets incorrect profile")
}

func TestGetOtherUserProfileFalse(t *testing.T) {
	require.Equal(t, &models.OtherUserProfile{}, GetOtherUserProfile(h, 6), "does not catch changed repeated mail")
}

func TestGetUserEvents(t *testing.T) {
	require.Equal(t, []uint64{125, 126}, getUserEvents(h, 1), "gets incorrect events")
}

func TestIsExistingMailTrue(t *testing.T) {
	profile := h.ProfileBase[0]

	require.Equal(t, true, IsExistingEMail(h, profile.Email), "existing mail does not exist")
}

func TestIsExistingMailFalse(t *testing.T) {
	profile := &models.UserOwnProfile{Uid: 1, Name: "Анастасия", Birthday: "6 февраля 2001 г.", City: "Москва", Email: "morozik@mail.ru",
		Visited: 12, Planning: 2, Followers: 36, About: "люблю котиков", Avatar: "1default.png"}

	require.Equal(t, false, IsExistingEMail(h, profile.Email), "not existing mail does exist")
}

// -----------------network-----------------

func setupEcho(t *testing.T, url, method string) (echo.Context,
	UserHandler) {
	user := models.User{Login: "moroz", Password: "123456"}
	e := echo.New()
	var req *http.Request
	switch method {
	case http.MethodPost:
		f, _ := user.MarshalJSON()
		req = httptest.NewRequest(http.MethodGet, url, bytes.NewBuffer(f))
	case http.MethodGet:
		req = httptest.NewRequest(http.MethodGet, url, nil)
	}
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(url)

	uh := UserHandler{
		UserBase:    UserBase,
		ProfileBase: ProfileBase,
		UserEvent:   EventUserBase,
		Store:       make(map[string]uint64),
		Mu:          &sync.Mutex{},
	}
	return c, uh
}

func TestHandlerUser_Logout(t *testing.T) {
	c, uh := setupEcho(t, "/api/v1/profile/logout", http.MethodGet)

	err := uh.Logout(c)
	require.NotEqual(t, nil, err)
}

func TestHandlerUser_GetProfileTrue(t *testing.T) {
	c, uh := setupEcho(t, "/api/v1/profile/:id", http.MethodGet)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := uh.GetOtherUserProfile(c)
	require.Equal(t, nil, err)
}

func TestHandlerUser_GetProfileFalse(t *testing.T) {
	c, uh := setupEcho(t, "/api/v1/profile/:id", http.MethodGet)
	c.SetParamNames("id")
	c.SetParamValues("aaa")

	err := uh.GetOtherUserProfile(c)
	require.NotEqual(t, nil, err)
}

func TestHandlerUser_GetAvatarTrue(t *testing.T) {
	c, uh := setupEcho(t, "/api/v1/avatar/:id", http.MethodGet)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := uh.GetAvatar(c)
	require.Equal(t, nil, err)
}

func TestHandlerUser_GetAvatarFalse(t *testing.T) {
	c, uh := setupEcho(t, "/api/v1/avatar/:id", http.MethodGet)
	c.SetParamNames("id")
	c.SetParamValues("u")

	err := uh.GetOtherUserProfile(c)
	require.NotEqual(t, nil, err)
}

func TestHandlerUser_LoginTrue(t *testing.T) {
	c, uh := setupEcho(t, "/api/v1/login", http.MethodPost)

	err := uh.Login(c)
	require.Equal(t, nil, err)
}
