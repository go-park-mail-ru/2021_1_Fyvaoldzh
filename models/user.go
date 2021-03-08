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
}

type UserData struct {
	Name string `json:"name"`
	Birthday string `json:"birthday"`
	City string `json:"city"`
	Email string `json:"email"`
	About string `json:"about"`
	Password string `json:"password"`
}






