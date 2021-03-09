package auth

import (
	"github.com/go-park-mail-ru/2021_1_Fyvaoldzh/models"
	"github.com/labstack/echo"
	"github.com/mailru/easyjson"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type LoginHandler struct {
	Mu     *sync.Mutex
}

var (
	Store = make(map[string]int)
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)


func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}


func isCorrectUser(user *models.User) (bool, int) {
	for _, value := range models.UserBase {
		if value.Login == (*user).Login && value.Password == (*user).Password {
			return true, value.Id
		}
	}
	return false, 0
}


func (h *LoginHandler) Login(c echo.Context) *echo.HTTPError {
	defer c.Request().Body.Close()
	u := &models.User{}

	cookie, err := c.Cookie("SID")
	if err == nil && Store[cookie.Value] != 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "user is already logged in")
	}


	key := RandStringRunes(32)

	err = easyjson.UnmarshalFromReader(c.Request().Body, u)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	isExisting, uid := isCorrectUser(u)
	if !isExisting {
		return echo.NewHTTPError(http.StatusBadRequest, "incorrect data")
	}

	newCookie := &http.Cookie{
		Name:    "SID",
		Value:   key,
		Expires: time.Now().Add(10 * time.Hour),
		SameSite: http.SameSiteNoneMode,
		Secure: true,
		HttpOnly: true,
	}

	Store[key] = uid
	c.SetCookie(newCookie)

	return nil
}

func (h *LoginHandler) Logout(c echo.Context) *echo.HTTPError {
	defer c.Request().Body.Close()

	cookie, err := c.Cookie("SID")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	uid := Store[cookie.Value]
	if uid == 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	delete(Store, cookie.Value)

	cookie.Expires = time.Now().AddDate(0, 0, -1)
	c.SetCookie(cookie)

	return nil
}

