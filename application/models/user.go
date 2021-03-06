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

type UserOnEvent struct {
	Id     uint64 `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type UserCard struct {
	Id        uint64 `json:"id"`
	Name      string `json:"name"`
	Avatar    string `json:"avatar"`
	Age       uint8  `json:"age"`
	City      string `json:"city"`
	Followers uint64 `json:"followers"`
}

//easyjson:json
type UserCards []UserCard

type UserCardSQL struct {
	Id       uint64
	Name     string
	Avatar   string
	Birthday sql.NullTime
	City     sql.NullString
}

func ConvertUserCard(sqlCard UserCardSQL) *UserCard {
	var card UserCard
	card.Id = sqlCard.Id
	card.Name = sqlCard.Name
	card.Avatar = sqlCard.Avatar
	if sqlCard.Birthday.Valid {
		dif := sqlCard.Birthday.Time.Sub(time.Now())
		secdif := math.Abs(dif.Seconds())
		card.Age = uint8(secdif / 31536000)
	}
	card.City = sqlCard.City.String
	return &card
}

//easyjson:json
type UsersOnEvent []UserOnEvent

type OtherUserProfile struct {
	Uid           uint64
	Name          string `json:"name"`
	Age           uint8  `json:"age"`
	City          string `json:"city"`
	About         string `json:"about"`
	Avatar        string `json:"avatar"`
	Followers     uint64 `json:"followers"`
	Subscriptions uint64 `json:"subscriptions"`
}

type UserOwnProfile struct {
	Uid           uint64
	Name          string `json:"name"`
	Login         string `json:"login"`
	Birthday      string `json:"birthday"`
	City          string `json:"city"`
	Email         string `json:"email"`
	Followers     uint64 `json:"followers"`
	Subscriptions uint64 `json:"subscriptions"`
	About         string `json:"about"`
	Avatar        string `json:"avatar"`
	OldPassword   string `json:"old_password"`
	NewPassword   string `json:"new_password"`
}

type UserDataSQL struct {
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

func ConvertToOwn(own UserDataSQL) *UserOwnProfile {
	var usr UserOwnProfile
	usr.Uid = own.Id
	usr.About = own.About.String
	usr.Email = own.Email.String
	//usr.Password = own.Password.String
	usr.Name = own.Name.String
	usr.Login = own.Login
	if own.Birthday.Valid {
		usr.Birthday = own.Birthday.Time.Format(constants.DateFormat)
	}
	usr.Avatar = own.Avatar.String
	usr.City = own.City.String
	return &usr
}

func ConvertToOther(own UserDataSQL) *OtherUserProfile {
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

type ActionCard struct {
	Id1   uint64
	Name1 string
	Id2   uint64
	Name2 string
	Time  time.Time
	Type  string
}

type ActionCardStringTime struct {
	Id1   uint64 `json:"id_1"`
	Name1 string `json:"name_1"`
	Id2   uint64 `json:"id_2"`
	Name2 string `json:"name_2"`
	Time  string `json:"time"`
	Type  string `json:"type"`
}

//easyjson:json
type ActionCards []ActionCardStringTime

func ConvertActionCard(card ActionCard) ActionCardStringTime {
	var newCard ActionCardStringTime
	newCard.Id1 = card.Id1
	newCard.Id2 = card.Id2
	newCard.Name1 = card.Name1
	newCard.Name2 = card.Name2
	newCard.Type = card.Type
	newCard.Time = card.Time.Format(constants.DateTimeFormat)
	return newCard
}
