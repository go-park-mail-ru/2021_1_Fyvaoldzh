package usecase

import (
	"database/sql"
	"image"
	"image/png"
	"io/ioutil"
	"kudago/application/models"
	mock_subscription "kudago/application/subscription/mocks"
	"kudago/application/user"
	mock_user "kudago/application/user/mocks"
	"kudago/pkg/constants"
	"kudago/pkg/logger"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	userId          uint64 = 1
	pageNum                = 1
	login                  = "userlogin"
	name                   = "username"
	frontPassword          = "123456"
	backPassword           = "IvJrQEdIeoTzLsMX_839spM7MzaXS7aJ_b3xTzmYqbotq3HRKAs="
	badBackPassword        = "1111IvJrQEdIeoTzLsMX_839spM7MzaXS7aJ_b3xTzmYqbotq3HRKAs="
	email                  = "email@mail.ru"
	birthdayStr            = "1999-01-01"
	birthday, err          = time.Parse(constants.TimeFormat, "1999-01-01")
	city                   = "City"
	about                  = "some personal information"
	avatar                 = "public/users/default.png"
	imageName              = "image.png"
	evPlanningSQL          = models.EventCardWithDateSQL{
		ID:        1,
		StartDate: time.Now(),
		EndDate:   time.Now().Add(10 * time.Hour),
	}
	evPlanning = models.EventCard{
		ID:        1,
		StartDate: evPlanningSQL.StartDate.String(),
		EndDate:   evPlanningSQL.EndDate.String(),
	}
	evVisitedSQL = models.EventCardWithDateSQL{
		ID:        2,
		StartDate: time.Now(),
		EndDate:   time.Now(),
	}
	evVisited = models.EventCard{
		ID:        2,
		StartDate: evVisitedSQL.StartDate.String(),
		EndDate:   evVisitedSQL.EndDate.String(),
	}
	eventsPlanningSQL = []models.EventCardWithDateSQL{
		evPlanningSQL, evVisitedSQL,
	}
	eventsVisitedSQL = []models.EventCardWithDateSQL{
		evVisitedSQL,
	}
	eventsPlanning = []models.EventCard{
		evPlanning,
	}
	eventsVisited = []models.EventCard{
		evVisited,
	}

	followers = []uint64{2, 2, 3}
)

var testUserFront = &models.User{
	Login:    login,
	Password: frontPassword,
}

var testRegData = &models.RegData{
	Login:    login,
	Password: frontPassword,
}

var testUserBack = &models.User{
	Id:       userId,
	Login:    login,
	Password: backPassword,
}

var testUserBadBack = &models.User{
	Id:       userId,
	Login:    login,
	Password: badBackPassword,
}

var testUserData = &models.UserDataSQL{
	Id:    userId,
	Login: login,
}

var testUserDataWithAvatar = &models.UserDataSQL{
	Id:     userId,
	Login:  login,
	Avatar: sql.NullString{String: imageName, Valid: true},
}

var testOtherUserProfile = &models.OtherUserProfile{
	Uid:       userId,
	Visited:   eventsVisited,
	Planning:  eventsPlanning,
	Followers: followers,
}

var testOwnUserProfile = &models.UserOwnProfile{
	Uid:       userId,
	Login:     login,
	Visited:   eventsVisited,
	Planning:  eventsPlanning,
	Followers: followers,
}
var testOwnUserProfileToUpdate = &models.UserOwnProfile{
	Uid:      userId,
	Name:     name,
	Login:    login,
	Birthday: birthdayStr,
	City:     city,
	Email:    email,
	About:    about,
	Avatar:   avatar,
}

var testNewUserData = &models.UserDataSQL{
	Id:       userId,
	Name:     sql.NullString{String: name, Valid: true},
	Login:    login,
	Birthday: sql.NullTime{Time: birthday, Valid: true},
	City:     sql.NullString{String: city, Valid: true},
	Email:    sql.NullString{String: email, Valid: true},
	About:    sql.NullString{String: about, Valid: true},
}

var testUserOnEvent = &models.UserOnEvent{
	Id:   userId,
	Name: name,
}

var testUserCardSQL = &models.UserCardSQL{
	Id:   userId,
	Name: name,
}

var testUserCardsSQL = []models.UserCardSQL{*testUserCardSQL}


func setUp(t *testing.T) (*mock_user.MockRepository, *mock_subscription.MockRepository, user.UseCase) {
	ctrl := gomock.NewController(t)

	rep := mock_user.NewMockRepository(ctrl)
	repSub := mock_subscription.NewMockRepository(ctrl)

	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	sugar := l.Sugar()
	zap.NewAtomicLevelAt(zapcore.DebugLevel)

	uc := NewUser(rep, repSub, logger.NewLogger(sugar))
	return rep, repSub, uc
}

func createImage() {
	width := 200
	height := 100

	upLeft := image.Point{}
	lowRight := image.Point{X: width, Y: height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	f, _ := os.Create(imageName)
	png.Encode(f, img)
}

func deleteImage() {
	os.Remove(imageName)
}

///////////////////////////////////////////////////

func TestUserUseCase_Login(t *testing.T) {
	rep, _, uc := setUp(t)

	rep.EXPECT().IsCorrect(testUserFront).Return(testUserBack, nil)

	actual, err := uc.CheckUser(testUserFront)

	assert.Nil(t, err)
	assert.Equal(t, testUserBack.Id, actual)
}

func TestUserUseCase_CheckUserOK(t *testing.T) {
	rep, _, uc := setUp(t)

	rep.EXPECT().IsCorrect(testUserFront).Return(testUserBack, nil)

	actual, err := uc.CheckUser(testUserFront)

	assert.Nil(t, err)
	assert.Equal(t, testUserBack.Id, actual)
}

func TestUserUseCase_CheckUserIncorrectData(t *testing.T) {
	rep, _, uc := setUp(t)

	rep.EXPECT().IsCorrect(testUserFront).Return(testUserBadBack, nil)

	_, err := uc.CheckUser(testUserFront)

	assert.Equal(t, err, echo.NewHTTPError(http.StatusBadRequest, "incorrect data"))
}

func TestUserUseCase_CheckUserIncorrectDBError(t *testing.T) {
	rep, _, uc := setUp(t)

	rep.EXPECT().IsCorrect(testUserFront).Return(&models.User{},
		echo.NewHTTPError(http.StatusBadRequest, "incorrect data"))

	_, err := uc.CheckUser(testUserFront)

	assert.Equal(t, err, echo.NewHTTPError(http.StatusBadRequest, "incorrect data"))
}

///////////////////////////////////////////////////

func TestUserUseCase_AddOK(t *testing.T) {
	rep, _, uc := setUp(t)

	rep.EXPECT().IsExisting(testUserFront.Login).Return(false, nil)
	rep.EXPECT().Add(testRegData).Return(userId, nil)
	rep.EXPECT().AddToPreferences(userId).Return(nil)

	actual, err := uc.Add(testRegData)

	assert.Nil(t, err)
	assert.Equal(t, testUserBack.Id, actual)
}

func TestUserUseCase_AddExistingLogin(t *testing.T) {
	rep, _, uc := setUp(t)

	rep.EXPECT().IsExisting(testUserFront.Login).Return(true, nil)

	_, err := uc.Add(testRegData)

	assert.Equal(t, err, echo.NewHTTPError(http.StatusBadRequest, "user with this login does exist"))
}

func TestUserUseCase_AddDBErrorAdd(t *testing.T) {
	rep, _, uc := setUp(t)

	rep.EXPECT().IsExisting(testUserFront.Login).Return(false, nil)
	rep.EXPECT().Add(testRegData).Return(userId, nil)
	rep.EXPECT().AddToPreferences(userId).Return(echo.NewHTTPError(http.StatusInternalServerError))

	_, err := uc.Add(testRegData)

	assert.Error(t, err)
}

func TestUserUseCase_AddDBErrorAddToPreferences(t *testing.T) {
	rep, _, uc := setUp(t)

	rep.EXPECT().IsExisting(testUserFront.Login).Return(false, nil)
	rep.EXPECT().Add(testRegData).Return(userId, echo.NewHTTPError(http.StatusInternalServerError))

	_, err := uc.Add(testRegData)

	assert.Error(t, err)
}

func TestUserUseCase_AddDBErrorIsExisting(t *testing.T) {
	rep, _, uc := setUp(t)

	rep.EXPECT().IsExisting(testUserFront.Login).Return(false, echo.NewHTTPError(http.StatusInternalServerError))

	_, err := uc.Add(testRegData)

	assert.Error(t, err)
}

///////////////////////////////////////////////////

func TestUserUseCase_GetOtherProfileOK(t *testing.T) {
	rep, repSub, uc := setUp(t)

	rep.EXPECT().GetByIdOwn(userId).Return(testUserData, nil)
	repSub.EXPECT().GetPlanningEvents(userId).Return(eventsPlanningSQL, nil)
	repSub.EXPECT().GetVisitedEvents(userId).Return(eventsVisitedSQL, nil)
	repSub.EXPECT().GetFollowers(userId).Return(followers, nil)
	repSub.EXPECT().UpdateEventStatus(userId, eventsPlanningSQL[1].ID).Return(nil)

	other, err := uc.GetOtherProfile(userId)

	assert.Nil(t, err)
	assert.Equal(t, testOtherUserProfile, other)
}

func TestUserUseCase_GetOtherProfileDBErrorGetByID(t *testing.T) {
	rep, _, uc := setUp(t)

	rep.EXPECT().GetByIdOwn(userId).Return(&models.UserDataSQL{}, echo.NewHTTPError(http.StatusInternalServerError))

	_, err := uc.GetOtherProfile(userId)

	assert.Error(t, err)
}

func TestUserUseCase_GetOtherProfileDBErrorGetPlanningEvents(t *testing.T) {
	rep, repSub, uc := setUp(t)

	rep.EXPECT().GetByIdOwn(userId).Return(testUserData, nil)
	repSub.EXPECT().GetPlanningEvents(userId).Return([]models.EventCardWithDateSQL{},
		echo.NewHTTPError(http.StatusInternalServerError))

	_, err := uc.GetOtherProfile(userId)

	assert.Error(t, err)
}

func TestUserUseCase_GetOtherProfileDBErrorGetVisitedEvents(t *testing.T) {
	rep, repSub, uc := setUp(t)

	rep.EXPECT().GetByIdOwn(userId).Return(testUserData, nil)
	repSub.EXPECT().GetPlanningEvents(userId).Return(eventsPlanningSQL, nil)
	repSub.EXPECT().GetVisitedEvents(userId).Return([]models.EventCardWithDateSQL{},
		echo.NewHTTPError(http.StatusInternalServerError))
	repSub.EXPECT().UpdateEventStatus(userId, eventsPlanningSQL[1].ID).Return(nil)

	_, err := uc.GetOtherProfile(userId)

	assert.Error(t, err)
}

func TestUserUseCase_GetOtherProfileDBErrorUpdateEventStatus(t *testing.T) {
	rep, repSub, uc := setUp(t)

	rep.EXPECT().GetByIdOwn(userId).Return(testUserData, nil)
	repSub.EXPECT().GetPlanningEvents(userId).Return(eventsPlanningSQL, nil)
	repSub.EXPECT().UpdateEventStatus(userId, eventsPlanningSQL[1].ID).Return(echo.NewHTTPError(http.StatusInternalServerError))

	_, err := uc.GetOtherProfile(userId)

	assert.Error(t, err)
}

func TestUserUseCase_GetOtherProfileDBErrorGetFollowers(t *testing.T) {
	rep, repSub, uc := setUp(t)

	rep.EXPECT().GetByIdOwn(userId).Return(testUserData, nil)
	repSub.EXPECT().GetPlanningEvents(userId).Return(eventsPlanningSQL, nil)
	repSub.EXPECT().GetVisitedEvents(userId).Return(eventsVisitedSQL, nil)
	repSub.EXPECT().GetFollowers(userId).Return([]uint64{},
		echo.NewHTTPError(http.StatusInternalServerError))
	repSub.EXPECT().UpdateEventStatus(userId, eventsPlanningSQL[1].ID).Return(nil)

	_, err := uc.GetOtherProfile(userId)

	assert.Error(t, err)
}

///////////////////////////////////////////////////

func TestUserUseCase_GetOwnProfileOK(t *testing.T) {
	rep, repSub, uc := setUp(t)

	rep.EXPECT().GetByIdOwn(userId).Return(testUserData, nil)
	repSub.EXPECT().GetPlanningEvents(userId).Return(eventsPlanningSQL, nil)
	repSub.EXPECT().GetVisitedEvents(userId).Return(eventsVisitedSQL, nil)
	repSub.EXPECT().GetFollowers(userId).Return(followers, nil)
	repSub.EXPECT().UpdateEventStatus(userId, eventsPlanningSQL[1].ID).Return(nil)

	own, err := uc.GetOwnProfile(testUserData.Id)

	assert.Nil(t, err)
	assert.Equal(t, testOwnUserProfile, own)
}

func TestUserUseCase_GetOwnProfileDBErrorGetByID(t *testing.T) {
	rep, _, uc := setUp(t)

	rep.EXPECT().GetByIdOwn(userId).Return(&models.UserDataSQL{}, echo.NewHTTPError(http.StatusInternalServerError))

	_, err := uc.GetOwnProfile(userId)

	assert.Error(t, err)
}

func TestUserUseCase_GetOwnProfileDBErrorGetPlanningEvents(t *testing.T) {
	rep, repSub, uc := setUp(t)

	rep.EXPECT().GetByIdOwn(userId).Return(testUserData, nil)
	repSub.EXPECT().GetPlanningEvents(userId).Return([]models.EventCardWithDateSQL{},
		echo.NewHTTPError(http.StatusInternalServerError))

	_, err := uc.GetOwnProfile(userId)

	assert.Error(t, err)
}

func TestUserUseCase_GetOwnProfileDBErrorGetVisitedEvents(t *testing.T) {
	rep, repSub, uc := setUp(t)

	rep.EXPECT().GetByIdOwn(userId).Return(testUserData, nil)
	repSub.EXPECT().GetPlanningEvents(userId).Return(eventsPlanningSQL, nil)
	repSub.EXPECT().GetVisitedEvents(userId).Return([]models.EventCardWithDateSQL{},
		echo.NewHTTPError(http.StatusInternalServerError))
	repSub.EXPECT().UpdateEventStatus(userId, eventsPlanningSQL[1].ID).Return(nil)

	_, err := uc.GetOwnProfile(userId)

	assert.Error(t, err)
}

func TestUserUseCase_GetOwnProfileDBErrorUpdateEventStatus(t *testing.T) {
	rep, repSub, uc := setUp(t)

	rep.EXPECT().GetByIdOwn(userId).Return(testUserData, nil)
	repSub.EXPECT().GetPlanningEvents(userId).Return(eventsPlanningSQL, nil)
	repSub.EXPECT().UpdateEventStatus(userId, eventsPlanningSQL[1].ID).Return(echo.NewHTTPError(http.StatusInternalServerError))

	_, err := uc.GetOwnProfile(userId)

	assert.Error(t, err)
}

func TestUserUseCase_GetOwnProfileDBErrorGetFollowers(t *testing.T) {
	rep, repSub, uc := setUp(t)

	rep.EXPECT().GetByIdOwn(userId).Return(testUserData, nil)
	repSub.EXPECT().GetPlanningEvents(userId).Return(eventsPlanningSQL, nil)
	repSub.EXPECT().GetVisitedEvents(userId).Return(eventsVisitedSQL, nil)
	repSub.EXPECT().GetFollowers(userId).Return([]uint64{},
		echo.NewHTTPError(http.StatusInternalServerError))
	repSub.EXPECT().UpdateEventStatus(userId, eventsPlanningSQL[1].ID).Return(nil)

	_, err := uc.GetOwnProfile(userId)

	assert.Error(t, err)
}

///////////////////////////////////////////////////

func TestUserUseCase_UpdateOK(t *testing.T) {
	rep, _, uc := setUp(t)

	rep.EXPECT().GetByIdOwn(testOwnUserProfileToUpdate.Uid).Return(testUserData, nil)
	rep.EXPECT().IsExistingEmail(testOwnUserProfileToUpdate.Email).Return(false, nil)
	rep.EXPECT().Update(testOwnUserProfileToUpdate.Uid, testNewUserData).Return(nil)

	err := uc.Update(testOwnUserProfileToUpdate.Uid, testOwnUserProfileToUpdate)

	assert.Nil(t, err)
}

func TestUserUseCase_UpdateBDErrorGetByIdOwn(t *testing.T) {
	rep, _, uc := setUp(t)

	rep.EXPECT().GetByIdOwn(testOwnUserProfileToUpdate.Uid).Return(&models.UserDataSQL{}, echo.NewHTTPError(http.StatusInternalServerError))

	err := uc.Update(testOwnUserProfileToUpdate.Uid, testOwnUserProfileToUpdate)

	assert.Error(t, err)
}

func TestUserUseCase_UpdateBDErrorIsExistingEmail(t *testing.T) {
	rep, _, uc := setUp(t)

	rep.EXPECT().GetByIdOwn(testOwnUserProfileToUpdate.Uid).Return(testUserData, nil)
	rep.EXPECT().IsExistingEmail(testOwnUserProfileToUpdate.Email).Return(false, echo.NewHTTPError(http.StatusInternalServerError))

	err := uc.Update(testOwnUserProfileToUpdate.Uid, testOwnUserProfileToUpdate)

	assert.Error(t, err)
}

func TestUserUseCase_UpdateBDErrorUpdate(t *testing.T) {
	rep, _, uc := setUp(t)

	rep.EXPECT().GetByIdOwn(testOwnUserProfileToUpdate.Uid).Return(testUserData, nil)
	rep.EXPECT().IsExistingEmail(testOwnUserProfileToUpdate.Email).Return(false, nil)
	rep.EXPECT().Update(testOwnUserProfileToUpdate.Uid, testNewUserData).Return(echo.NewHTTPError(http.StatusInternalServerError))

	err := uc.Update(testOwnUserProfileToUpdate.Uid, testOwnUserProfileToUpdate)

	assert.Error(t, err)
}

///////////////////////////////////////////////////

func TestUserUseCase_GetAvatar(t *testing.T) {
	rep, _, uc := setUp(t)
	createImage()
	file, err := ioutil.ReadFile(imageName)
	if err != nil {
		t.Error(err)
	}

	rep.EXPECT().GetByIdOwn(userId).Return(testUserDataWithAvatar, nil)

	gotFile, err := uc.GetAvatar(userId)

	assert.Nil(t, err)
	assert.Equal(t, file, gotFile)
	deleteImage()
}

func TestUserUseCase_GetAvatarDBErrorGetByIdOwn(t *testing.T) {
	rep, _, uc := setUp(t)

	rep.EXPECT().GetByIdOwn(userId).Return(&models.UserDataSQL{}, echo.NewHTTPError(http.StatusInternalServerError))

	_, err := uc.GetAvatar(userId)

	assert.Error(t, err)
}

func TestUserUseCase_GetAvatarErrorReadFile(t *testing.T) {
	rep, _, uc := setUp(t)

	rep.EXPECT().GetByIdOwn(userId).Return(testUserDataWithAvatar, nil)

	_, err := uc.GetAvatar(userId)

	assert.Error(t, err)
}

///////////////////////////////////////////////////

func TestUserUseCase_GetUsers(t *testing.T) {
	rep, repSub, uc := setUp(t)

	rep.EXPECT().GetUsers(pageNum).Return(testUserCardsSQL, nil)
	repSub.EXPECT().GetFollowers(userId).Return([]uint64{userId}, nil)

	_, err := uc.GetUsers(pageNum)

	assert.Nil(t, err)
}

/*
func TestUserUseCase_UploadAvatar(t *testing.T) {
	rep, _, uc := setUp(t)
	createImage()
	file, err := os.Open(imageName)
	if err != nil {
		t.Error(err)
	}

	rep.EXPECT().GetByIdOwn(userId).Return(testUserDataWithAvatar, nil)

	err = uc.UploadAvatar(userId, file, imageName)

	assert.Nil(t, err)
	deleteImage()
}


*/
