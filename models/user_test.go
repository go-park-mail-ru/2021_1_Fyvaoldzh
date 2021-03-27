package models

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConvert(t *testing.T) {
	own := UserOwnProfile{1, "Анастасия", "6 февраля 2001 г.", "Москва", "moroz@mail.ru",
		12, 2, 36, "люблю котиков", "1default.png", nil}
	expect := &OtherUserProfile{Uid: 1, Name: "Анастасия", Age: 20, City: "Москва", Followers: 36, About: "люблю котиков", Avatar: "1default.png", Event: nil}

	require.Equal(t, expect, ConvertOwnOther(own), "keys are not different")
}
