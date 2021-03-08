package models

type User struct {
	Id int
	Login  string `json:"login"`
	Name string `json:"name"`
	Password string `json:"password"`
}

type Profile struct {
	Uid int
	Birthday string `json:"birthday"`
	Age int `json:"age"`
	City string `json:"city"`
	Email string `json:"email"`
	Visited int `json:"visited"`
	Planning int `json:"planning"`
	Followers int `json:"followers"`
	About string `json:"about"`
}

type UserData struct {
	Name string `json:"name"`
	Birthday string `json:"birthday"`
	City string `json:"city"`
	Email string `json:"email"`
	About string `json:"about"`
	Password string `json:"password"`
}






