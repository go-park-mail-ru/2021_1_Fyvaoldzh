package repository

import (
	"database/sql"
	"kudago/application/user"
	"kudago/models"
)

type UserDatabase struct {
	DB   *sql.DB
}

func NewUserDatabase(conn *sql.DB) user.Repository {
	return &UserDatabase{DB: conn}
}

/*
func GetUser(h *UserHandler, uid uint64) *models.User {
	for _, value := range h.UserBase {
		if value.Id == uid {
			return value
		}
	}
	return &models.User{}
}

func IsExistingUser(h *UserHandler, user *models.User) bool {
	for _, value := range h.UserBase {
		if value.Login == (*user).Login {
			return true
		}
	}
	return false
}

func IsCorrectUser(h *UserHandler, user *models.User) (bool, uint64) {
	for _, value := range h.UserBase {
		if value.Login == (*user).Login && value.Password == (*user).Password {
			return true, value.Id
		}
	}
	return false, 0
}

func GetProfile(h *UserHandler, uid uint64) *models.UserOwnProfile {
	for _, value := range h.ProfileBase {
		if value.Uid == uid {
			return value
		}
	}
	return &models.UserOwnProfile{}
}

func Update(h *UserHandler, ud *models.UserData, uid uint64) error {
	user := GetUser(h, uid)
	profile := GetProfile(h, uid)

	if len(ud.Name) != 0 {
		profile.Name = ud.Name
	}

	if len(ud.Password) != 0 {
		user.Password = ud.Password
	}

	if len(ud.Email) != 0 {
		if IsExistingEMail(h, ud.Email) {
			return echo.NewHTTPError(http.StatusBadRequest, "this email does exist")
		}
		profile.Email = ud.Email
	}

	if len(ud.About) != 0 {
		profile.About = ud.About
	}

	if len(ud.Birthday) != 0 {
		profile.Birthday = ud.Birthday
		// код на изменение age, который будет, когда будет формат даты
	}

	if len(ud.City) != 0 {
		profile.City = ud.City
	}

	return nil
}

func GetByIdOther(h *UserHandler, uid uint64) *models.OtherUserProfile {
	value := &models.UserOwnProfile{}
	flag := false

	for _, value = range h.ProfileBase {
		if value.Uid == uid {
			flag = true
			break
		}
	}

	if !flag {
		return &models.OtherUserProfile{}
	}

	return models.ConvertOwnOther(*value)
}


func IsExistingEMail(h *UserHandler, email string) bool {
	for _, value := range h.ProfileBase {
		if value.Email == email {
			return true
		}
	}
	return false
}

func getUserEvents(h *UserHandler, uid uint64) []uint64 {
	var events []uint64

	for _, value := range h.UserEvent {
		if value.Uid == uid {
			events = append(events, value.Eid)
		}
	}

	return events
}


 */

func (u UserDatabase) Add(user *models.RegData) (err error) {
	u.DB.QueryRow(
		`INSERT INTO users ("name", "login, "password") VALUES ($1, $2, $3)`,
		user.Name, user.Login, user.Password)

	return nil
}

func (u UserDatabase) Update(id uint, upUser *models.User) error {
	panic("implement me")
}

func (u UserDatabase) GetByIdOther(id uint) (user *models.User, err error) {
	panic("implement me")
}

func (u UserDatabase) GetByIdOwn(id uint) (user *models.User, err error) {
	panic("implement me")
}

func (u UserDatabase) GetByName(login string) (user *models.User, err error) {
	panic("implement me")
}
