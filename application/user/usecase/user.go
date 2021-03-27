package usecase

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/labstack/echo"
	"kudago/application/user"
	"kudago/models"
	"net/http"
)

type User struct {
	database user.Repository
}


func NewUser(u user.Repository) user.UseCase {
	return &User{database: u}
}

func (uc User) Login(u *models.User) (uint64, error) {
	hash := sha256.New()
	hash.Write([]byte(u.Password))
	u.Password = base64.URLEncoding.EncodeToString(hash.Sum(nil))

	uid, err := uc.database.IsCorrect(u)

	if err != nil {
		return 0, err
	}

	return uid, nil
}

func (uc User) Add(usr *models.RegData) (uint64, error) {
	hash := sha256.New()
	hash.Write([]byte(usr.Password))
	usr.Password = base64.URLEncoding.EncodeToString(hash.Sum(nil))

	flag, err := uc.database.IsExisting(usr.Login)
	if err != nil {
		return 0, err
	}
	if flag {
		return 0, echo.NewHTTPError(http.StatusBadRequest, "user with this login does exist")
	}

	id, err := uc.database.Add(usr)
	return id, nil
}

func (uc User) Get(id uint) (*models.User, error) {
	panic("implement me")
}

func (uc User) Update(uid uint64, ud *models.UserData) error {
	changeUser, err := uc.database.GetByIdOwn(uid)

	if err != nil {
		return err
	}


	if len(ud.Name) != 0 {
		changeUser.Name = ud.Name
	}

	if len(ud.Password) != 0 {
		changeUser.Password = ud.Password
	}

	if len(ud.Email) != 0 {
		flag, err := uc.database.IsExistingEmail(ud.Email)
		if err != nil {
			return err
		}
		if flag {
			return echo.NewHTTPError(http.StatusBadRequest, "this email does exist")
		}
		changeUser.Email = ud.Email
	}

	if len(ud.About) != 0 {
		changeUser.About = ud.About
	}

	if len(ud.Birthday) != 0 {
		changeUser.Birthday = ud.Birthday
	}

	if len(ud.City) != 0 {
		changeUser.City = ud.City
	}

	err = uc.database.Update(uid, changeUser)

	if err != nil {
		return err
	}

	return nil
}