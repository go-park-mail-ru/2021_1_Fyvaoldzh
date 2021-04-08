package usecase

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"kudago/application/models"
	"kudago/application/user"
	"kudago/pkg/constants"
	"kudago/pkg/generator"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo"
)

type User struct {
	repo user.Repository
}

func (uc User) Login(user *models.User) (uint64, error) {
	return uc.CheckUser(user)
}

func NewUser(u user.Repository) user.UseCase {
	return &User{repo: u}
}

func (uc User) CheckUser(u *models.User) (uint64, error) {
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
	var newEvents []models.EventCard

	oldEvents, err := uc.repo.GetPlanningEvents(id)
	if err != nil {
		return &models.OtherUserProfile{}, err
	}

	for _, elem := range oldEvents {
		if elem.StartDate.Before(time.Now()) {
			err := uc.repo.DeletePlanningEvent(id, elem.ID)
			if err != nil {
				return &models.OtherUserProfile{}, err
			}

			err = uc.repo.AddVisitedEvent(id, elem.ID)
			if err != nil {
				return &models.OtherUserProfile{}, err
			}
		} else {
			newEvents = append(newEvents, models.ConvertDateCard(elem))
		}
	}
	other.Planning = newEvents

	visitedEventsSQL, err := uc.repo.GetVisitedEvents(id)
	if err != nil {
		return &models.OtherUserProfile{}, err
	}
	for _, elem := range visitedEventsSQL {
		other.Visited = append(other.Visited, elem)
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
	var newEvents []models.EventCard

	oldEvents, err := uc.repo.GetPlanningEvents(id)
	if err != nil {
		return &models.UserOwnProfile{}, err
	}

	for _, elem := range oldEvents {
		if elem.StartDate.Before(time.Now()) {
			err := uc.repo.DeletePlanningEvent(id, elem.ID)
			if err != nil {
				return &models.UserOwnProfile{}, err
			}

			err = uc.repo.AddVisitedEvent(id, elem.ID)
			if err != nil {
				return &models.UserOwnProfile{}, err
			}
		} else {
			newEvents = append(newEvents, models.ConvertDateCard(elem))
		}
	}
	own.Planning = newEvents

	visitedEventsSQL, err := uc.repo.GetVisitedEvents(id)
	if err != nil {
		return &models.UserOwnProfile{}, err
	}
	for _, elem := range visitedEventsSQL {
		own.Visited = append(own.Visited, elem)
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

func (uc User) UploadAvatar(uid uint64, img *multipart.FileHeader) error {
	src, err := img.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer src.Close()

	fileName := constants.UserPicDir + fmt.Sprint(uid) + generator.RandStringRunes(6) + img.Filename
	dst, err := os.Create(fileName)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = uc.repo.ChangeAvatar(uid, fileName)

	if err != nil {
		return err
	}

	return nil
}

func (uc User) GetAvatar(uid uint64) ([]byte, error) {
	usr, err := uc.repo.GetByIdOwn(uid)

	if err != nil {
		return []byte{}, err
	}

	if usr.Avatar.Valid {
		file, err := ioutil.ReadFile(usr.Avatar.String)
		if err != nil {
			return []byte{}, echo.NewHTTPError(http.StatusInternalServerError, "Cannot open file")
		}
		return file, nil
	}

	return []byte{}, err
}
