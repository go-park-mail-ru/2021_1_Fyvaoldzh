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
	Event []uint64 `json:"events"`
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
	Event []uint64 `json:"events"`
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
	Eid uint64
}






