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

// ----------bases------------
type HandlerUser struct {
	UserBase []*models.User
	ProfileBase []*models.UserOwnProfile
	PlanningEvent []*models.PlanningEvents
	Store map[string]int
	Mu *sync.Mutex
}

var UserBase = []*models.User{
	{1, "moroz", "123456"},
	{2, "matros", "123456"},
	{3, "mail", "123456"},
}

var ProfileBase = []*models.UserOwnProfile{
	{1, "Анастасия", "6 февраля 2001 г.", "Москва", "moroz@mail.ru",
		12, 2, 36, "люблю котиков", "1.png", nil},
	{2, "Матрос Матросович Матросов", "7 февраля 1999 г.", "Санкт-Петербург", "matros@mail.ru",
		77, 15, 1000, "главный матрос на корабле", "1.png", nil},
	{3, "Почтальон Печкин", "1 марта 1997 г.", "Москва", "pechkin@mail.ru",
		1000, 99, 123, "ваш любимый почтальон", "2.png", nil},
}

var PlanningEvent = []*models.PlanningEvents{
	{1, 125},
	{1, 126},
	{2, 125},
	{2, 127},
}

// -----------global----------

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	id = 4
)


// ---------api register----------

func (h *HandlerUser) CreateUser(c echo.Context) error {
	defer c.Request().Body.Close()
	newData := &models.RegData{}

	log.Println(c.Request().Body)
	err := easyjson.UnmarshalFromReader(c.Request().Body, newData)
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	newUser := &models.User{}
	newUser.Login = newData.Login
	newUser.Password = newData.Password

	if IsExistingUser(h, newUser) {
		return echo.NewHTTPError(http.StatusBadRequest, "user with this login does exist")
	}

	newUser.Id = id
	id++

	newProfile := &models.UserOwnProfile{}
	newProfile.Uid = newUser.Id
	newProfile.Name = newData.Name

	h.Mu.Lock()
	UserBase = append(UserBase, newUser)
	ProfileBase = append(ProfileBase, newProfile)
	h.Mu.Unlock()

	key := RandStringRunes(32)

	newCookie := &http.Cookie{
		Name:     "SID",
		Value:    key,
		Expires:  time.Now().Add(10 * time.Hour),
		SameSite: http.SameSiteNoneMode,
		HttpOnly: true,
	}

	h.Store[key] = newUser.Id
	c.SetCookie(newCookie)

	return nil
}

// -----------api login----------

func (h *HandlerUser) Login(c echo.Context) error {
	defer c.Request().Body.Close()
	u := &models.User{}

	cookie, err := c.Cookie("SID")
	if err == nil && h.Store[cookie.Value] != 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "user is already logged in")
	}

	key := RandStringRunes(32)

	err = easyjson.UnmarshalFromReader(c.Request().Body, u)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	isExisting, uid := isCorrectUser(h, u)
	if !isExisting {
		return echo.NewHTTPError(http.StatusBadRequest, "incorrect data")
	}

	newCookie := &http.Cookie{
		Name:     "SID",
		Value:    key,
		Expires:  time.Now().Add(10 * time.Hour),
		SameSite: http.SameSiteNoneMode,
		HttpOnly: true,
	}

	h.Store[key] = uid
	c.SetCookie(newCookie)

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

// ---------api profile-----------------

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

	// некрасиво, но пока
	for _, value := range h.PlanningEvent {
		if value.Uid == profile.Uid {
			profile.Event = append(profile.Event, value.Eid)
		}
	}

	log.Println(profile)

	if _, err = easyjson.MarshalToWriter(profile, c.Response().Writer); err != nil {
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

	if httperr := changeProfileData(h, ud, h.Store[cookie.Value]); httperr != nil {
		return httperr
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


// ---------helpful functions-----------

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// --------------users-------------- //
func IsExistingEMail(h *HandlerUser, email string) bool {
	for _, value := range h.ProfileBase {
		if value.Email == email {
			return true
		}
	}
	return false
}

func GetUser(h *HandlerUser, uid int) *models.User {
	for _, value := range h.UserBase {
		if value.Id == uid {
			return value
		}
	}
	return &models.User{}
}

func isCorrectUser(h *HandlerUser, user *models.User) (bool, int) {
	for _, value := range h.UserBase {
		if value.Login == (*user).Login && value.Password == (*user).Password {
			return true, value.Id
		}
	}
	return false, 0
}

func IsExistingUser(h *HandlerUser, user *models.User) bool {
	for _, value := range h.UserBase {
		if value.Login == (*user).Login {
			return true
		}
	}
	return false
}

// --------------profile------------

func GetProfile(h *HandlerUser, uid int) *models.UserOwnProfile {
	for _, value := range h.ProfileBase {
		if value.Uid == uid {
			return value
		}
	}
	return &models.UserOwnProfile{}
}

func changeProfileData(h *HandlerUser, ud *models.UserData, uid int) error {
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
		//проверка на повторяемость почты
	}

	if len(ud.About) != 0 {
		profile.About = ud.About
	}

	if len(ud.Birthday) != 0 {
		profile.Birthday = ud.Birthday
		// код на изменение age
	}

	if len(ud.City) != 0 {
		profile.City = ud.City
	}

	return nil
}

func (h *HandlerUser) GetUserProfile(c echo.Context) error {
	defer c.Request().Body.Close()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	user := GetOtherUserProfile(h, id)

	if user.Uid == 0 {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New("user does not exist"))
	}

	for _, value := range h.PlanningEvent {
		if value.Uid == user.Uid {
			user.Event = append(user.Event, value.Eid)
		}
	}

	log.Println(user)

	if _, err = easyjson.MarshalToWriter(user, c.Response().Writer); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func GetOtherUserProfile(h *HandlerUser, uid int) *models.UserProfile {
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

	otherProfile := &models.UserProfile{}
	otherProfile.Uid = value.Uid
	otherProfile.Name = value.Name
	otherProfile.City = value.City
	otherProfile.About = value.About
	otherProfile.Followers = value.Followers
	otherProfile.Avatar = value.Avatar
	otherProfile.Event = value.Event
	// здесь оно будет по-умному высчитываться, но пока так
	otherProfile.Age = 20
	log.Println(otherProfile)
	return otherProfile
}

func (h *HandlerUser) GetAvatar(c echo.Context) error {
	defer c.Request().Body.Close()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	profile := GetOtherUserProfile(h, id)

	if profile.Uid == 0 {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New("user does not exist"))
	}

	file, err := ioutil.ReadFile(profile.Avatar)
	if err != nil {
		log.Println("Cannot open file: " + profile.Avatar)
	} else {
		c.Response().Write(file)
	}

	return nil
}
