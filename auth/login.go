package auth

import (
	"github.com/go-park-mail-ru/2021_1_Fyvaoldzh/models"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/mailru/easyjson"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

type LoginHandler struct {
	Mu     *sync.Mutex
}

var UserBase = []models.User {
	{"moroz", "moroz@mail.ru", "123456"},
	{"matros", "matros@mail.ru", "123456"},
	{"mail", "mail@mail.ru", "123456"},
}

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
// заменить как-нибудь потом
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")


func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func isExistingUser(user *models.User) bool {
	for _, value := range UserBase {
		if value == *user {
			return true
		}
	}
	return false
}

func (h *LoginHandler) Login(c echo.Context) error {
	defer c.Request().Body.Close()
	u := &models.User{}

	log.Println(c.Request().Body)
	err := easyjson.UnmarshalFromReader(c.Request().Body, u)
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if !isExistingUser(u) {
		//возврат ошибки
	}

	key := RandStringRunes(32)

	s, _ := store.Get(c.Request(), key)

	err = s.Save(c.Request(), c.Response())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// ставим куку
	// здесь явно что-то не то
	cookie := new(http.Cookie)
	cookie.Name = "SID"
	cookie.Value = key
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(cookie)

	return nil
}