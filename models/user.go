package models

import (
	"database/sql"
	"kudago/pkg/constants"
	"math"
	"time"
)

type User struct {
	Id       uint64
	Login    string `json:"login"`
	Password string `json:"password"`
}

type OtherUserProfile struct {
	Uid       uint64
	Name      string   `json:"name"`
	Age       uint8    `json:"age"`
	City      string   `json:"city"`
	About     string   `json:"about"`
	Avatar    string   `json:"avatar"`
	Visited   []uint64   `json:"visited"`
	Planning  []uint64   `json:"planning"`
	Followers []uint64   `json:"followers"`
}

type UserOwnProfile struct {
	Uid       uint64
	Name      string   `json:"name"`
	Birthday  string   `json:"birthday"`
	City      string   `json:"city"`
	Email     string   `json:"email"`
	Visited   []uint64   `json:"visited"`
	Planning  []uint64   `json:"planning"`
	Followers []uint64   `json:"followers"`
	About     string   `json:"about"`
	Avatar    string   `json:"avatar"`
	Password  string   `json:"password"`
}

type UserData struct {
	Name     sql.NullString
	Birthday sql.NullTime
	City     sql.NullString
	Email    sql.NullString
	About    sql.NullString
	Password sql.NullString
	Avatar   sql.NullString
}

type UserEvents struct {
	Uid uint64
	Eid uint64
}

func ConvertToOwn(own UserData) *UserOwnProfile {
	var usr UserOwnProfile
	usr.About = own.About.String
	usr.Email = own.Email.String
	//usr.Password = own.Password.String
	usr.Name = own.Name.String
	usr.Birthday = own.Birthday.Time.Format(constants.TimeFormat)
	usr.Avatar = own.Avatar.String
	usr.City = own.City.String
	return &usr
}

// TODO: возраст считается неправильно
func ConvertToOther(own UserData) *OtherUserProfile {
	var usr OtherUserProfile
	usr.About = own.About.String
	//usr.Password = own.Password.String
	usr.Name = own.Name.String
	if own.Birthday.Valid {
		usr.Age = uint8(math.Round(own.Birthday.Time.Sub(time.Now()).Seconds() / 31207680))
	}
	usr.Avatar = own.Avatar.String
	usr.City = own.City.String
	return &usr
}

type RegData struct {
	Id       uint64
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
}


