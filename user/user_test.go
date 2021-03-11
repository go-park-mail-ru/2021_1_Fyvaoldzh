package user

import (
	"github.com/stretchr/testify/require"
	"kudago/models"
	"sync"
	"testing"
)

var h = &HandlerUser{
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
	newData := &models.RegData{Login: "morozik_jr", Password: "123456", Name: "Morozik Jr."}
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
	user := GetUser(h, 1)
	expUser := &models.User{Id: 1, Login: "moroz", Password: "123456"}

	require.Equal(t, expUser, user, "users are not equal")
}

func TestGetUserFalse(t *testing.T) {
	user := GetUser(h, 6)
	expUser := &models.User{}

	require.Equal(t, expUser, user, "got user that does not exist")
}

func TestIsExistingUserTrue(t *testing.T) {
	user := &models.User{Id: 1, Login: "moroz", Password: "123456"}

	require.Equal(t, true, IsExistingUser(h, user), "existing user does not exist")
}

func TestIsExistingUserFalse(t *testing.T) {
	user := &models.User{Id: 1, Login: "morozik", Password: "123456"}

	require.Equal(t, false, IsExistingUser(h, user), "not existing user does exist")
}

func TestIsCorrectUserTrue(t *testing.T) {
	user := &models.User{Id: 1, Login: "moroz", Password: "123456"}

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
	profile := &models.UserOwnProfile{Uid: 1, Name: "Анастасия", Birthday: "6 февраля 2001 г.", City: "Москва", Email: "moroz@mail.ru",
		Visited: 12, Planning: 2, Followers: 36, About: "люблю котиков", Avatar: "1default.png"}

	require.Equal(t, profile, GetProfile(h, profile.Uid), "existing mail does not exist")
}

func TestGetProfileFalse(t *testing.T) {
	require.Equal(t, &models.UserOwnProfile{}, GetProfile(h, 6), "returns not existing profile")
}

func TestChangesProfileDataTrue(t *testing.T) {
	profile := &models.UserData{Name: "Матрос", Birthday: "6 марта 1999 г.", City: "Санкт-Петербург", Email: "matrosik@mail.ru",
		About: "люблю котиков", Avatar: "1default.png"}

	require.Equal(t, nil, changeProfileData(h, profile, 2), "returns error on correct data")
}

func TestChangesProfileDataFalse(t *testing.T) {
	profile := &models.UserData{Email: "moroz@mail.ru"}

	require.NotEqual(t, nil, changeProfileData(h, profile, 2), "does not catch changed repeated mail")
}

func TestGetOtherUserProfileTrue(t *testing.T) {
	profile := &models.UserProfile{Uid: 1, Name: "Анастасия", Age: 20, City: "Москва", Followers: 36, About: "люблю котиков", Avatar: "1default.png"}

	require.Equal(t, profile, GetOtherUserProfile(h, 1), "gets incorrect profile")
}

func TestGetOtherUserProfileFalse(t *testing.T) {
	require.Equal(t, &models.UserProfile{}, GetOtherUserProfile(h, 6), "does not catch changed repeated mail")
}

func TestGetUserEvents(t *testing.T) {
	require.Equal(t, []uint64{125, 126}, getUserEvents(h, 1), "gets incorrect events")
}

func TestIsExistingMailTrue(t *testing.T) {
	profile := &models.UserOwnProfile{Uid: 1, Name: "Анастасия", Birthday: "6 февраля 2001 г.", City: "Москва", Email: "moroz@mail.ru",
		Visited: 12, Planning: 2, Followers: 36, About: "люблю котиков", Avatar: "1default.png"}

	require.Equal(t, true, IsExistingEMail(h, profile.Email), "existing mail does not exist")
}

func TestIsExistingMailFalse(t *testing.T) {
	profile := &models.UserOwnProfile{Uid: 1, Name: "Анастасия", Birthday: "6 февраля 2001 г.", City: "Москва", Email: "morozik@mail.ru",
		Visited: 12, Planning: 2, Followers: 36, About: "люблю котиков", Avatar: "1default.png"}

	require.Equal(t, false, IsExistingEMail(h, profile.Email), "not existing mail does exist")
}
