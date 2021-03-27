package repository

import (
	"database/sql"
	"github.com/labstack/echo"
	"kudago/application/user"
	"kudago/models"
	"log"
	"net/http"
	"time"
)

type UserDatabase struct {
	DB *sql.DB
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

func (ud UserDatabase) Add(user *models.RegData) (id uint64, err error) {
	err = ud.DB.QueryRow(
		`INSERT INTO users ("name", "login", "password") VALUES ($1, $2, $3) RETURNING id`,
		user.Name, user.Login, user.Password).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}

func (ud UserDatabase) IsExisting(login string) (bool, error) {
	var id uint64
	err := ud.DB.
		QueryRow(`SELECT id FROM users WHERE login = $1`, login).Scan(&id)
	switch {
	case err == sql.ErrNoRows:
		return false, nil
	case err != nil:
		return false, err
	}
	return true, nil
}

func (ud UserDatabase) IsCorrect(user *models.User) (uint64, error) {
	var id uint64
	err := ud.DB.
		QueryRow(`SELECT id FROM users WHERE login = $1 AND password = $2`,
			user.Login, user.Password).Scan(&id)
	switch {
	case err == sql.ErrNoRows:
		return 0, echo.NewHTTPError(http.StatusBadRequest, "incorrect login or password")
	case err != nil:
		return 0, err
	}
	return id, nil
}

// TODO: null поля, дату тоже в null
func (ud UserDatabase) Update(id uint64, us *models.UserData) error {
	dt, err := time.Parse("2006-01-02", us.Birthday)
	_, err = ud.DB.Exec(
		`UPDATE users SET "name" = $1, "email" = $2, "city" = $3, "about" = $4,`+
			`"avatar" = $5, "birthday" = $6, "password" = $7 WHERE id = $8`,
		us.Name, us.Email, us.City, us.About, us.Avatar, dt, us.Password, id,
	)
	log.Println(err)
	if err != nil {
		return err
	}

	return nil
}

func (ud UserDatabase) GetByIdOther(id uint64) (user *models.User, err error) {
	panic("implement me")
}

func (ud UserDatabase) GetByIdOwn(id uint64) (*models.UserData, error) {
	var massiv [7]sql.NullString
	usr := &models.UserData{}
	err := ud.DB.
		QueryRow(`SELECT name, birthday, city, email, about, password, avatar 
		FROM users WHERE id = $1`,
			id).Scan(&massiv[0], &massiv[1], &massiv[2], &massiv[3], &massiv[4],
				&massiv[5], &massiv[6])
	log.Println(err)

	usr.Name = massiv[0].String
	usr.Birthday = massiv[1].String
	usr.City = massiv[2].String
	usr.Email = massiv[3].String
	usr.About = massiv[4].String
	usr.Password = massiv[5].String
	usr.Avatar = massiv[6].String

	if err != nil {
		return &models.UserData{}, err
	}
	return usr, nil
}

func (ud UserDatabase) GetByName(login string) (user *models.User, err error) {
	panic("implement me")
}

func (ud UserDatabase) IsExistingEmail(email string) (bool, error) {
	var id uint64
	err := ud.DB.
		QueryRow(`SELECT id FROM users WHERE email = $1`, email).Scan(&id)
	switch {
	case err == sql.ErrNoRows:
		return false, nil
	case err != nil:
		return false, err
	}
	return true, nil
}
