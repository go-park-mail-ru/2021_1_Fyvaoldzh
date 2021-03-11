package models

type User struct {
	Id       uint64
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RegData struct {
	Id       uint64
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserProfile struct {
	Uid       uint64
	Name      string   `json:"name"`
	Age       uint8    `json:"age"`
	City      string   `json:"city"`
	Followers uint64   `json:"followers"`
	About     string   `json:"about"`
	Avatar    string   `json:"avatar"`
	Event     []uint64 `json:"events"`
}

type UserOwnProfile struct {
	Uid       uint64
	Name      string   `json:"name"`
	Birthday  string   `json:"birthday"`
	City      string   `json:"city"`
	Email     string   `json:"email"`
	Visited   uint64   `json:"visited"`
	Planning  uint64   `json:"planning"`
	Followers uint64   `json:"followers"`
	About     string   `json:"about"`
	Avatar    string   `json:"avatar"`
	Event     []uint64 `json:"events"`
}

type UserData struct {
	Name     string `json:"name"`
	Birthday string `json:"birthday"`
	City     string `json:"city"`
	Email    string `json:"email"`
	About    string `json:"about"`
	Password string `json:"password"`
	Avatar   string `json:"avatar"`
}

type UserEvents struct {
	Uid uint64
	Eid uint64
}

func ConvertOwnOther(own UserOwnProfile) *UserProfile {
	otherProfile := &UserProfile{}
	otherProfile.Uid = own.Uid
	otherProfile.Name = own.Name
	otherProfile.City = own.City
	otherProfile.About = own.About
	otherProfile.Followers = own.Followers
	otherProfile.Avatar = own.Avatar
	otherProfile.Event = own.Event
	// здесь оно будет по-умному высчитываться, но пока так
	otherProfile.Age = 20

	return otherProfile
}
