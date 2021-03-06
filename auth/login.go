package auth

import (
	"github.com/go-park-mail-ru/2021_1_Fyvaoldzh/models"
	"sync"
)

type LoginHandler struct {
	Mu     *sync.Mutex
}

var userBase = []models.User {
	{"moroz", "moroz@mail.ru"},
	{"danya", "danya@mail.ru"},
	{"mail", "mail@mail.ru"},
}