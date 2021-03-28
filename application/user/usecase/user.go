package usecase

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/labstack/echo"
	"kudago/application/user"
	"kudago/models"
	"kudago/pkg/constants"
	"net/http"
	"time"
)

type User struct {
	repo user.Repository
}

func NewUser(u user.Repository) user.UseCase {
	return &User{repo: u}
}

func (uc User) Login(u *models.User) (uint64, error) {
	hash := sha256.New()
	hash.Write([]byte(u.Password))
	u.Password = base64.URLEncoding.EncodeToString(hash.Sum(nil))

	uid, err := uc.repo.IsCorrect(u)

	if err != nil {
		return 0, err
	}

	return uid, nil
}

func (uc User) Add(usr *models.RegData) (uint64, error) {
	hash := sha256.New()
	hash.Write([]byte(usr.Password))
	usr.Password = base64.URLEncoding.EncodeToString(hash.Sum(nil))

	flag, err := uc.repo.IsExisting(usr.Login)
	if err != nil {
		return 0, err
	}
	if flag {
		return 0, echo.NewHTTPError(http.StatusBadRequest, "user with this login does exist")
	}

	id, err := uc.repo.Add(usr)
	return id, nil
}

func (uc User) GetOtherProfile(id uint64) (*models.OtherUserProfile, error) {
	usr, err := uc.repo.GetByIdOwn(id)

	if err != nil {
		return &models.OtherUserProfile{}, err
	}

	other := models.ConvertToOther(*usr)
	other.Visited, err = uc.repo.GetVisitedEvents(id)
	if err != nil {
		return &models.OtherUserProfile{}, err
	}
	other.Planning, err = uc.repo.GetPlanningEvents(id)
	if err != nil {
		return &models.OtherUserProfile{}, err
	}
	other.Followers, err = uc.repo.GetFollowers(id)
	if err != nil {
		return &models.OtherUserProfile{}, err
	}

	return other, err

}

func (uc User) GetOwnProfile(id uint64) (*models.UserOwnProfile, error) {
	usr, err := uc.repo.GetByIdOwn(id)

	if err != nil {
		return &models.UserOwnProfile{}, err
	}

	own := models.ConvertToOwn(*usr)

	own.Visited, err = uc.repo.GetVisitedEvents(id)
	if err != nil {
		return &models.UserOwnProfile{}, err
	}
	own.Planning, err = uc.repo.GetPlanningEvents(id)
	if err != nil {
		return &models.UserOwnProfile{}, err
	}
	own.Followers, err = uc.repo.GetFollowers(id)
	if err != nil {
		return &models.UserOwnProfile{}, err
	}

	return own, nil
}

func (uc User) Update(uid uint64, ud *models.UserOwnProfile) error {
	changeUser, err := uc.repo.GetByIdOwn(uid)

	if err != nil {
		return err
	}

	if len(ud.Name) != 0 {
		changeUser.Name.String = ud.Name
		changeUser.Name.Valid = true
	}

	if len(ud.Password) != 0 {
		hash := sha256.New()
		hash.Write([]byte(ud.Password))
		changeUser.Password.String = base64.URLEncoding.EncodeToString(hash.Sum(nil))
		changeUser.Password.Valid = true
	}

	if len(ud.Email) != 0 {
		flag, err := uc.repo.IsExistingEmail(ud.Email)
		if err != nil {
			return err
		}
		if flag {
			return echo.NewHTTPError(http.StatusBadRequest, "this email does exist")
		}
		changeUser.Email.String = ud.Email
		changeUser.Email.Valid = true
	}

	if len(ud.About) != 0 {
		changeUser.About.String = ud.About
		changeUser.About.Valid = true
	}

	if len(ud.Birthday) != 0 {
		dt, err := time.Parse(constants.TimeFormat, ud.Birthday)
		if err != nil {
			return err
		}
		changeUser.Birthday.Time = dt
		changeUser.Birthday.Valid = true
	}

	if len(ud.City) != 0 {
		changeUser.City.String = ud.City
		changeUser.City.Valid = true
	}

	err = uc.repo.Update(uid, changeUser)

	if err != nil {
		return err
	}

	return nil
}