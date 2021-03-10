package models

import (
	"github.com/stretchr/testify/require"
	"testing"
)

// здесь оно вообще все ходит в базу, но без нее чет никак

func TestIsExistingUserTrue(t *testing.T) {
	user := &User{Id: 1, Login: "moroz", Password: "123456"}

	require.Equal(t, true, IsExistingUser(user), "Existing user does not exist")
}

func TestIsExistingUserFalse(t *testing.T) {
	user := &User{Id: 1, Login: "morozik", Password: "123456"}

	require.Equal(t, false, IsExistingUser(user),  "Existing user does not exist")
}

func TestIsExistingMailTrue(t *testing.T) {
	profile := &UserOwnProfile{1, "Анастасия", "6 февраля 2001 г.", "Москва", "moroz@mail.ru",
		12, 2, 36, "люблю котиков", "1.png", nil}

	require.Equal(t, true, IsExistingEMail(profile.Email),  "Existing user does not exist")
}

func TestIsExistingMailFalse(t *testing.T) {
	profile := &UserOwnProfile{1, "Анастасия", "6 февраля 2001 г.", "Москва", "morozik@mail.ru",
		12, 2, 36, "люблю котиков", "1.png", nil}

	require.Equal(t, false, IsExistingEMail(profile.Email),  "Existing user does not exist")
}

func TestGetUser(t *testing.T) {
	user := GetUser(1)
	expUser := &User{Id: 1, Login: "morozik", Password: "123456"}

	require.Equal(t, expUser, user, "Users are not equal")
}