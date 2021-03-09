package models

type User struct {
	Id int
	Login  string `json:"login"`
	Password string `json:"password"`
}

type RegData struct {
	Id int
	Name string `json:"name"`
	Login  string `json:"login"`
	Password string `json:"password"`
}

type UserProfile struct {
	Uid int
	Name string `json:"name"`
	Age int `json:"age"`
	City string `json:"city"`
	Followers int `json:"followers"`
	About string `json:"about"`
	Avatar string `json:"avatar"`
}

type UserOwnProfile struct {
	Uid int
	Name string `json:"name"`
	Birthday string `json:"birthday"`
	City string `json:"city"`
	Email string `json:"email"`
	Visited int `json:"visited"`
	Planning int `json:"planning"`
	Followers int `json:"followers"`
	About string `json:"about"`
	Avatar string `json:"avatar"`
}

type UserData struct {
	Name string `json:"name"`
	Birthday string `json:"birthday"`
	City string `json:"city"`
	Email string `json:"email"`
	About string `json:"about"`
	Password string `json:"password"`
	Avatar string `json:"avatar"`
}

type PlanningEvents struct {
	Uid int
	Eid int
}

var UserBase = []*User{
	{1, "moroz", "123456"},
	{2, "matros", "123456"},
	{3, "mail", "123456"},
}

var ProfileBase = []*UserOwnProfile{
	{1, "Анастасия", "6 февраля 2001 г.", "Москва", "moroz@mail.ru",
		12, 2, 36, "люблю котиков", "public/1.png"},
	{2, "Матрос Матросович Матросов", "7 февраля 1999 г.", "Санкт-Петербург", "matros@mail.ru",
		77, 15, 1000, "главный матрос на корабле", "public/1.png"},
	{3, "Почтальон Печкин", "1 марта 1997 г.", "Москва", "pechkin@mail.ru",
		1000, 99, 123, "ваш любимый почтальон", "public/2.png"},
}

var PlanningEvent = []*PlanningEvents{
	{1, 125},
	{1, 126},
	{2, 125},
	{2, 127},
}

func IsExistingUser(user *User) bool {
	for _, value := range UserBase {
		if value.Login == (*user).Login {
			return true
		}
	}
	return false
}

func IsExistingEMail(email string) bool {
	for _, value := range ProfileBase {
		if value.Email == email {
			return true
		}
	}
	return false
}

func GetUser(uid int) *User {
	for _, value := range UserBase {
		if value.Id == uid {
			return value
		}
	}
	return &User{}
}

func GetProfile(uid int) *UserOwnProfile {
	for _, value := range ProfileBase {
		if value.Uid == uid {
			return value
		}
	}
	return &UserOwnProfile{}
}

func GetOtherUserProfile(uid int) *UserProfile {
	value := &UserOwnProfile{}
	flag := false

	for _, value = range ProfileBase {
		if value.Uid == uid {
			flag = true
			break
		}
	}

	if !flag {
		return &UserProfile{}
	}

	otherProfile := &UserProfile{}
	otherProfile.Uid = value.Uid
	otherProfile.Name = value.Name
	otherProfile.City = value.City
	otherProfile.About = value.About
	otherProfile.Followers = value.Followers
	// здесь оно будет по-умному высчитываться, но пока так
	otherProfile.Age = 20
	return otherProfile
}





