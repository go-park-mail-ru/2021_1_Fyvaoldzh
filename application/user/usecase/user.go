package usecase

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"kudago/application/models"
	"kudago/application/subscription"
	"kudago/application/user"
	"kudago/pkg/constants"
	"kudago/pkg/generator"
	"kudago/pkg/logger"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo"
)

type UserUseCase struct {
	repo    user.Repository
	repoSub subscription.Repository
	logger  logger.Logger
}

func NewUser(u user.Repository, repoSubscription subscription.Repository, logger logger.Logger) user.UseCase {
	return &UserUseCase{repo: u, repoSub: repoSubscription, logger: logger}
}

func (uc UserUseCase) Login(user *models.User) (uint64, error) {
	return uc.CheckUser(user)
}

func (uc UserUseCase) CheckUser(u *models.User) (uint64, error) {
	gotUser, err := uc.repo.IsCorrect(u)
	if err != nil {
		uc.logger.Warn(err)
		return 0, err
	}

	if !generator.CheckHashedPassword(gotUser.Password, u.Password) {
		uc.logger.Warn(errors.New("incorrect data"))
		return 0, echo.NewHTTPError(http.StatusBadRequest, "incorrect data")
	}

	return gotUser.Id, nil
}

func (uc UserUseCase) Add(usr *models.RegData) (uint64, error) {
	// TODO: добавить проверки
	usr.Password = generator.HashPassword(usr.Password)

	flag, err := uc.repo.IsExisting(usr.Login)
	if err != nil {
		uc.logger.Warn(err)
		return 0, err
	}
	if flag {
		uc.logger.Warn(errors.New("user with this login does exist"))
		return 0, echo.NewHTTPError(http.StatusBadRequest, "user with this login does exist")
	}

	id, err := uc.repo.Add(usr)
	if err != nil {
		uc.logger.Warn(err)
		return 0, err
	}
	err = uc.repo.AddToPreferences(id)
	if err != nil {
		uc.logger.Warn(err)
		return 0, err
	}

	return id, nil
}

func (uc UserUseCase) GetOtherProfile(id uint64) (*models.OtherUserProfile, error) {
	usr, err := uc.repo.GetByIdOwn(id)

	if err != nil {
		uc.logger.Warn(err)
		return &models.OtherUserProfile{}, err
	}

	other := models.ConvertToOther(*usr)
	var newEvents []models.EventCard

	oldEvents, err := uc.repo.GetPlanningEvents(id)
	if err != nil {
		uc.logger.Warn(err)
		return &models.OtherUserProfile{}, err
	}

	for _, elem := range oldEvents {
		if elem.EndDate.Before(time.Now()) {
			err := uc.repo.UpdateEventStatus(id, elem.ID)
			if err != nil {
				uc.logger.Warn(err)
				return &models.OtherUserProfile{}, err
			}
		} else {
			newEvents = append(newEvents, models.ConvertDateCard(elem))
		}
	}
	other.Planning = newEvents

	visitedEventsSQL, err := uc.repo.GetVisitedEvents(id)
	if err != nil {
		uc.logger.Warn(err)
		return &models.OtherUserProfile{}, err
	}
	for _, elem := range visitedEventsSQL {
		other.Visited = append(other.Visited, models.ConvertDateCard(elem))
	}

	other.Followers, err = uc.repo.GetFollowers(id)
	if err != nil {
		uc.logger.Warn(err)
		return &models.OtherUserProfile{}, err
	}

	return other, err

}

func (uc UserUseCase) GetOwnProfile(id uint64) (*models.UserOwnProfile, error) {
	usr, err := uc.repo.GetByIdOwn(id)

	if err != nil {
		uc.logger.Warn(err)
		return &models.UserOwnProfile{}, err
	}

	own := models.ConvertToOwn(*usr)
	var newEvents []models.EventCard

	oldEvents, err := uc.repo.GetPlanningEvents(id)
	if err != nil {
		uc.logger.Warn(err)
		return &models.UserOwnProfile{}, err
	}

	for _, elem := range oldEvents {
		if elem.EndDate.Before(time.Now()) {
			err := uc.repo.UpdateEventStatus(id, elem.ID)
			if err != nil {
				uc.logger.Warn(err)
				return &models.UserOwnProfile{}, err
			}
		} else {
			newEvents = append(newEvents, models.ConvertDateCard(elem))
		}
	}
	own.Planning = newEvents

	visitedEventsSQL, err := uc.repo.GetVisitedEvents(id)
	if err != nil {
		uc.logger.Warn(err)
		return &models.UserOwnProfile{}, err
	}
	for _, elem := range visitedEventsSQL {
		own.Visited = append(own.Visited, models.ConvertDateCard(elem))
	}

	own.Followers, err = uc.repo.GetFollowers(id)
	if err != nil {
		uc.logger.Warn(err)
		return &models.UserOwnProfile{}, err
	}

	return own, nil
}

func (uc UserUseCase) Update(uid uint64, ud *models.UserOwnProfile) error {
	changeUser, err := uc.repo.GetByIdOwn(uid)
	if err != nil {
		uc.logger.Warn(err)
		return err
	}

	if len(ud.Name) != 0 {
		changeUser.Name.String = ud.Name
		changeUser.Name.Valid = true
	}

	if len(ud.OldPassword) != 0 {
		if !generator.CheckHashedPassword(changeUser.Password.String, ud.OldPassword) {
			uc.logger.Warn(errors.New("passwords are not same"))
			return echo.NewHTTPError(http.StatusBadRequest, "passwords are not same")
		}
		changeUser.Password.String = generator.HashPassword(ud.NewPassword)
		changeUser.Password.Valid = true
	}

	if len(ud.Email) != 0 {
		flag, err := uc.repo.IsExistingEmail(ud.Email)
		if err != nil {
			uc.logger.Warn(err)
			return err
		}
		if flag {
			uc.logger.Warn(errors.New("this email does exist"))
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
			uc.logger.Warn(err)
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
		uc.logger.Warn(err)
		return err
	}

	return nil
}

func (uc UserUseCase) UploadAvatar(uid uint64, src multipart.File, filename string) error {
	fileName := constants.UserPicDir + fmt.Sprint(uid) + generator.RandStringRunes(6) + filename
	dst, err := os.Create(fileName)
	if err != nil {
		uc.logger.Warn(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		uc.logger.Warn(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = uc.repo.ChangeAvatar(uid, fileName)

	if err != nil {
		uc.logger.Warn(err)
		return err
	}

	return nil
}

func (uc UserUseCase) GetAvatar(uid uint64) ([]byte, error) {
	usr, err := uc.repo.GetByIdOwn(uid)
	if err != nil {
		uc.logger.Warn(err)
		return []byte{}, err
	}

	if usr.Avatar.Valid {
		file, err := ioutil.ReadFile(usr.Avatar.String)
		if err != nil {
			uc.logger.Warn(err)
			return []byte{}, echo.NewHTTPError(http.StatusInternalServerError, "Cannot open file")
		}
		return file, nil
	}

	return []byte{}, err
}

func (uc UserUseCase) GetUsers(page int) (models.UserCards, error) {
	var cards models.UserCards
	users, err := uc.repo.GetUsers(page)
	if err != nil {
		return models.UserCards{}, err
	}
	for _, elem := range users {
		followers, err := uc.repo.GetFollowers(elem.Id)
		if err != nil {
			return models.UserCards{}, err
		}
		newCard := *models.ConvertUserCard(elem)
		newCard.Followers = uint64(len(followers))
		cards = append(cards, newCard)
	}
	return cards, nil
}

func (uc UserUseCase) UpdateEventStatus(userId uint64, eventId uint64) error {
	return uc.repo.UpdateEventStatus(userId, eventId)
}

func (uc UserUseCase) IsAddedEvent(userId uint64, eventId uint64) (bool, error) {
	return uc.repo.IsAddedEvent(userId, eventId)
}