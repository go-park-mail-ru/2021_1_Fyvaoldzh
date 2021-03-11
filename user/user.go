package user

import (
	"errors"
	"fmt"
	"github.com/labstack/echo"
	"github.com/mailru/easyjson"
	"io"
	"io/ioutil"
	"kudago/models"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

// -----------------bases-----------------
type HandlerUser struct {
	UserBase    []*models.User
	ProfileBase []*models.UserOwnProfile
	UserEvent   []*models.UserEvents
	Store       map[string]uint64
	Mu          *sync.Mutex
}

var UserBase = []*models.User{
	{1, "moroz", "123456"},
	{2, "matros", "123456"},
	{3, "mail", "123456"},
}

var ProfileBase = []*models.UserOwnProfile{
	{1, "Анастасия", "6 февраля 2001 г.", "Москва", "moroz@mail.ru",
		12, 2, 36, "люблю котиков", "1default.png", nil},
	{2, "Матрос Матросович Матросов", "7 февраля 1999 г.", "Санкт-Петербург", "matros@mail.ru",
		77, 15, 1000, "главный матрос на корабле", "1.png", nil},
	{3, "Почтальон Печкин", "1 марта 1997 г.", "Москва", "pechkin@mail.ru",
		1000, 99, 123, "ваш любимый почтальон", "2.png", nil},
}

var EventUserBase = []*models.UserEvents{
	{1, 125},
	{1, 126},
	{2, 125},
	{2, 127},
}

// -----------------global-----------------

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	id          = uint64(4)
)

const n = uint8(32)

// -----------------api register-----------------

func (h *HandlerUser) CreateUser(c echo.Context) error {
	defer c.Request().Body.Close()

	cookie, err := c.Cookie("SID")
	if err == nil && h.Store[cookie.Value] != 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "user is already logged in")
	}

	newData := &models.RegData{}

	log.Println(c.Request().Body)
	err = easyjson.UnmarshalFromReader(c.Request().Body, newData)
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	uid, err := h.CreateUserProfile(newData)
	if err != nil {
		return err
	}

	c.SetCookie(CreateCookie(h, n, uid))

	return nil
}

// -----------------api user-profile-----------------

func (h *HandlerUser) CreateUserProfile(data *models.RegData) (uint64, error) {
	newUser := &models.User{}
	newUser.Login = data.Login
	newUser.Password = data.Password

	if IsExistingUser(h, newUser) {
		return 0, echo.NewHTTPError(http.StatusBadRequest, "user with this login does exist")
	}

	newUser.Id = id
	id++

	newProfile := &models.UserOwnProfile{}
	newProfile.Uid = newUser.Id
	newProfile.Name = data.Name

	h.Mu.Lock()
	h.UserBase = append(h.UserBase, newUser)
	h.ProfileBase = append(h.ProfileBase, newProfile)
	h.Mu.Unlock()
	return newUser.Id, nil
}

// -----------------api login-----------------

func (h *HandlerUser) Login(c echo.Context) error {
	defer c.Request().Body.Close()
	u := &models.User{}

	cookie, err := c.Cookie("SID")
	if err == nil && h.Store[cookie.Value] != 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "user is already logged in")
	}

	err = easyjson.UnmarshalFromReader(c.Request().Body, u)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	isCorrect, uid := IsCorrectUser(h, u)
	if !isCorrect {
		return echo.NewHTTPError(http.StatusBadRequest, "incorrect data")
	}

	c.SetCookie(CreateCookie(h, n, uid))

	return nil
}

func (h *HandlerUser) Logout(c echo.Context) error {
	defer c.Request().Body.Close()

	cookie, err := c.Cookie("SID")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	uid := h.Store[cookie.Value]
	if uid == 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	delete(h.Store, cookie.Value)

	cookie.Expires = time.Now().AddDate(0, 0, -1)
	c.SetCookie(cookie)

	return nil
}

// -----------------api profile-----------------

func (h *HandlerUser) GetProfile(c echo.Context) error {
	defer c.Request().Body.Close()

	cookie, err := c.Cookie("SID")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	if h.Store[cookie.Value] == 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	profile := GetProfile(h, h.Store[cookie.Value])
	profile.Event = getUserEvents(h, profile.Uid)

	if _, err = easyjson.MarshalToWriter(profile, c.Response().Writer); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (h *HandlerUser) GetUserProfile(c echo.Context) error {
	defer c.Request().Body.Close()

	uid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user := GetOtherUserProfile(h, uint64(uid))

	if user.Uid == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("user does not exist"))
	}

	user.Event = getUserEvents(h, user.Uid)

	if _, err = easyjson.MarshalToWriter(user, c.Response().Writer); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (h *HandlerUser) UpdateProfile(c echo.Context) error {
	defer c.Request().Body.Close()

	cookie, err := c.Cookie("SID")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	if h.Store[cookie.Value] == 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	ud := &models.UserData{}
	err = easyjson.UnmarshalFromReader(c.Request().Body, ud)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	httperr := changeProfileData(h, ud, h.Store[cookie.Value])
	if httperr != nil {
		return httperr
	}

	return nil
}
func (h *HandlerUser) GetAvatar(c echo.Context) error {
	defer c.Request().Body.Close()

	uid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	profile := GetOtherUserProfile(h, uint64(uid))

	if profile.Uid == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("user does not exist"))
	}

	file, err := ioutil.ReadFile(profile.Avatar)
	if err != nil {
		log.Println("Cannot open file: " + profile.Avatar)
	} else {
		c.Response().Write(file)
	}

	return nil
}

func (h *HandlerUser) UploadAvatar(c echo.Context) error {
	defer c.Request().Body.Close()

	cookie, err := c.Cookie("SID")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	if h.Store[cookie.Value] == 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	user := GetUser(h, h.Store[cookie.Value])
	profile := GetProfile(h, user.Id)

	img, err := c.FormFile("avatar")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if img == nil {
		return nil
	}

	src, err := img.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer src.Close()

	fileName := fmt.Sprint(user.Id) + img.Filename
	dst, err := os.Create(fileName)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	profile.Avatar = fileName
	return nil
}

// -----------------helpful functions-----------------

func RandStringRunes(n uint8) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// -----------------users-----------------

func GetUser(h *HandlerUser, uid uint64) *models.User {
	for _, value := range h.UserBase {
		if value.Id == uid {
			return value
		}
	}
	return &models.User{}
}

func IsExistingUser(h *HandlerUser, user *models.User) bool {
	for _, value := range h.UserBase {
		if value.Login == (*user).Login {
			return true
		}
	}
	return false
}

func IsCorrectUser(h *HandlerUser, user *models.User) (bool, uint64) {
	for _, value := range h.UserBase {
		if value.Login == (*user).Login && value.Password == (*user).Password {
			return true, value.Id
		}
	}
	return false, 0
}

func CreateCookie(h *HandlerUser, n uint8, uid uint64) *http.Cookie {
	key := RandStringRunes(n)

	newCookie := &http.Cookie{
		Name:     "SID",
		Value:    key,
		Expires:  time.Now().Add(10 * time.Hour),
		HttpOnly: true,
	}

	h.Store[key] = uid
	return newCookie
}

// -----------------profile-----------------

func GetProfile(h *HandlerUser, uid uint64) *models.UserOwnProfile {
	for _, value := range h.ProfileBase {
		if value.Uid == uid {
			return value
		}
	}
	return &models.UserOwnProfile{}
}

func changeProfileData(h *HandlerUser, ud *models.UserData, uid uint64) error {
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

func GetOtherUserProfile(h *HandlerUser, uid uint64) *models.UserProfile {
	value := &models.UserOwnProfile{}
	flag := false

	for _, value = range h.ProfileBase {
		if value.Uid == uid {
			flag = true
			break
		}
	}

	if !flag {
		return &models.UserProfile{}
	}

	return models.ConvertOwnOther(*value)
}

func getUserEvents(h *HandlerUser, uid uint64) []uint64 {
	var events []uint64

	for _, value := range h.UserEvent {
		if value.Uid == uid {
			events = append(events, value.Eid)
		}
	}

	return events
}

func IsExistingEMail(h *HandlerUser, email string) bool {
	for _, value := range h.ProfileBase {
		if value.Email == email {
			return true
		}
	}
	return false
}
