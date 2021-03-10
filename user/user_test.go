package user

import (
	"github.com/stretchr/testify/require"
	"kudago/models"
	"sync"
	"testing"
)

// здесь оно вообще все ходит в базу, но без нее чет никак

var h = &HandlerUser{
	UserBase:      UserBase,
	ProfileBase:   ProfileBase,
	PlanningEvent: PlanningEvent,
	Store:         make(map[string]int),
	Mu:            &sync.Mutex{},
}

func TestIsExistingUserTrue(t *testing.T) {
	user := &models.User{Id: 1, Login: "moroz", Password: "123456"}

	require.Equal(t, true, IsExistingUser(h, user), "Existing user does not exist")
}

func TestIsExistingUserFalse(t *testing.T) {
	user := &models.User{Id: 1, Login: "morozik", Password: "123456"}

	require.Equal(t, false, IsExistingUser(h, user), "Existing user does not exist")
}

func TestIsExistingMailTrue(t *testing.T) {
	profile := &models.UserOwnProfile{1, "Анастасия", "6 февраля 2001 г.", "Москва", "moroz@mail.ru",
		12, 2, 36, "люблю котиков", "1.png", nil}

	require.Equal(t, true, IsExistingEMail(h, profile.Email), "Existing user does not exist")
}

func TestIsExistingMailFalse(t *testing.T) {
	profile := &models.UserOwnProfile{1, "Анастасия", "6 февраля 2001 г.", "Москва", "morozik@mail.ru",
		12, 2, 36, "люблю котиков", "1.png", nil}

	require.Equal(t, false, IsExistingEMail(h, profile.Email), "Existing user does not exist")
}

func TestGetUser(t *testing.T) {
	user := GetUser(h, 1)
	expUser := &models.User{Id: 1, Login: "moroz", Password: "123456"}

	require.Equal(t, expUser, user, "Users are not equal")
}
