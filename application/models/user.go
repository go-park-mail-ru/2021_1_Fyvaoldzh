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
	Name      string     `json:"name"`
	Age       uint8      `json:"age"`
	City      string     `json:"city"`
	About     string     `json:"about"`
	Avatar    string     `json:"avatar"`
	Visited   EventCards `json:"visited"`
	Planning  EventCards `json:"planning"`
	Followers []uint64   `json:"followers"`
}

type UserOwnProfile struct {
	Uid       uint64
	Name      string     `json:"name"`
	Login     string     `json:"login"`
	Birthday  string     `json:"birthday"`
	City      string     `json:"city"`
	Email     string     `json:"email"`
	Visited   EventCards `json:"visited"`
	Planning  EventCards `json:"planning"`
	Followers []uint64   `json:"followers"`
	About     string     `json:"about"`
	Avatar    string     `json:"avatar"`
	Password  string     `json:"password"`
}

type UserData struct {
	Id       uint64
	Name     sql.NullString
	Login    string
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
	usr.Uid = own.Id
	usr.About = own.About.String
	usr.Email = own.Email.String
	//usr.Password = own.Password.String
	usr.Name = own.Name.String
	usr.Login = own.Login
	usr.Birthday = own.Birthday.Time.Format(constants.TimeFormat)
	usr.Avatar = own.Avatar.String
	usr.City = own.City.String
	return &usr
}

func ConvertToOther(own UserData) *OtherUserProfile {
	var usr OtherUserProfile
	usr.Uid = own.Id
	usr.About = own.About.String
	//usr.Password = own.Password.String
	usr.Name = own.Name.String
	if own.Birthday.Valid {
		dif := own.Birthday.Time.Sub(time.Now())
		secdif := math.Abs(dif.Seconds())
		usr.Age = uint8(secdif / 31536000)
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
