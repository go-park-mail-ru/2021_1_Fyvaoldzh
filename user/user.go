package user

import (
	"kudago/models"
	"sync"
)

// -----------------bases-----------------
type UserHandler struct {
	UserBase    []*models.User
	ProfileBase []*models.UserOwnProfile
	UserEvent   []*models.UserEvents
	Store       map[string]uint64
	Mu          *sync.Mutex
}

var UserBase = []*models.User{
	{1, "moroz", "123456"},
	{2, "matros", "123456"},
	{3, "mail", "123456"},
}

/*
var ProfileBase = []*models.UserOwnProfile{
	{1, "Анастасия", "6 февраля 2001 г.", "Москва", "moroz@mail.ru",
		12, 2, 36, "люблю котиков", "1default.png", nil},
	{2, "Матрос Матросович Матросов", "7 февраля 1999 г.", "Санкт-Петербург", "matros@mail.ru",
		77, 15, 1000, "главный матрос на корабле", "1.png", nil},
	{3, "Почтальон Печкин", "1 марта 1997 г.", "Москва", "pechkin@mail.ru",
		1000, 99, 123, "ваш любимый почтальон", "2.png", nil},
}*/

var EventUserBase = []*models.UserEvents{
	{1, 125},
	{1, 126},
	{2, 125},
	{2, 127},
}







