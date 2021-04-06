package http

import (
	"github.com/labstack/echo"
	"github.com/mailru/easyjson"
	"kudago/application/models"
	"kudago/application/user"
	"kudago/pkg/constants"
	"kudago/pkg/generator"
	"kudago/pkg/infrastructure"
	"log"
	"net/http"
	"strconv"
	"time"
)

type UserHandler struct {
	UseCase user.UseCase
	Sm      *infrastructure.SessionManager
}

func CreateUserHandler(e *echo.Echo, uc user.UseCase, sm *infrastructure.SessionManager) {

	userHandler := UserHandler{UseCase: uc, Sm: sm}

	e.POST("/api/v1/login", userHandler.Login)
	e.DELETE("/api/v1/logout", userHandler.Logout)
	e.POST("/api/v1/register", userHandler.Register)
	e.GET("/api/v1/profile", userHandler.GetOwnProfile)
	e.GET("/api/v1/profile/:id", userHandler.GetOtherUserProfile)
	e.PUT("/api/v1/profile", userHandler.Update)
	e.POST("/api/v1/check_user", userHandler.CheckUser)
	e.POST("/api/v1/upload_avatar", userHandler.UploadAvatar)
	e.GET("/api/v1/avatar/:id", userHandler.GetAvatar)
}

func (h *UserHandler) Login(c echo.Context) error {
	defer c.Request().Body.Close()
	u := &models.User{}

	cookie, err := c.Cookie(constants.SessionCookieName)

	if err != nil && cookie != nil {
		log.Println(err)
		return err
	}

	err = easyjson.UnmarshalFromReader(c.Request().Body, u)
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	uid, err := h.UseCase.Login(u)

	if err != nil {
		log.Println(err)
		return err
	}

	if cookie != nil {
		exists, id, err := h.Sm.CheckSession(cookie.Value)
		if err != nil {
			log.Println(err)
			return err
		}

		if exists && id == uid {
			return echo.NewHTTPError(http.StatusBadRequest, "user is already logged in")
		}
	}

	cookie = h.CreateCookie(constants.CookieLength)
	err = h.Sm.InsertSession(uid, cookie.Value)

	if err != nil {
		log.Println(err)
		return err
	}

	c.SetCookie(cookie)

	return nil
}

func (h *UserHandler) Logout(c echo.Context) error {
	defer c.Request().Body.Close()

	cookie, err := c.Cookie(constants.SessionCookieName)
	if err != nil && cookie != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, "error getting cookie")
	}

	if cookie == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	flag, _, err := h.Sm.CheckSession(cookie.Value)

	if err != nil {
		log.Println(err)
		return err
	}

	if !flag {
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	err = h.Sm.DeleteSession(cookie.Value)

	cookie.Expires = time.Now().AddDate(0, 0, -1)
	c.SetCookie(cookie)

	return nil
}

func (h *UserHandler) Register(c echo.Context) error {
	defer c.Request().Body.Close()

	cookie, err := c.Cookie(constants.SessionCookieName)
	if err != nil && cookie != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, "error getting cookie")
	}

	// вынести эту штуку
	if cookie != nil {
		exists, _, err := h.Sm.CheckSession(cookie.Value)
		if err != nil {
			log.Println(err)
			return err
		}

		if exists {
			return echo.NewHTTPError(http.StatusBadRequest, "user is already logged in")
		}
	}

	newData := &models.RegData{}

	err = easyjson.UnmarshalFromReader(c.Request().Body, newData)
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	uid, err := h.UseCase.Add(newData)

	if err != nil {
		log.Println(err)
		return err
	}

	cookie = h.CreateCookie(constants.CookieLength)
	err = h.Sm.InsertSession(uid, cookie.Value)

	if err != nil {
		log.Println(err)
		return err
	}

	c.SetCookie(cookie)
	return nil
}

func (h *UserHandler) Update(c echo.Context) error {
	defer c.Request().Body.Close()

	cookie, err := c.Cookie(constants.SessionCookieName)
	if err != nil && cookie != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	var uid uint64
	var exists bool
	if cookie != nil {
		exists, uid, err = h.Sm.CheckSession(cookie.Value)
		if err != nil {
			log.Println(err)
			return err
		}

		if !exists {
			return echo.NewHTTPError(http.StatusBadRequest, "user is not authorized")
		}
	}

	ud := &models.UserOwnProfile{}
	err = easyjson.UnmarshalFromReader(c.Request().Body, ud)
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = h.UseCase.Update(uid, ud)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// TODO: унести в session manager
func (h *UserHandler) CreateCookie(n uint8) *http.Cookie {

	key := generator.RandStringRunes(n)

	newCookie := &http.Cookie{
		Name:     constants.SessionCookieName,
		Value:    key,
		Expires:  time.Now().Add(10 * time.Hour),
		HttpOnly: true,
	}

	return newCookie
}

func (h *UserHandler) GetOwnProfile(c echo.Context) error {
	defer c.Request().Body.Close()

	cookie, err := c.Cookie("SID")
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	var uid uint64
	var exists bool

	exists, uid, err = h.Sm.CheckSession(cookie.Value)
	if err != nil {
		log.Println(err)
		return err
	}

	if !exists {
		return echo.NewHTTPError(http.StatusBadRequest, "user is not authorized")
	}

	usr, err := h.UseCase.GetOwnProfile(uid)
	if err != nil {
		log.Println(err)
		return err
	}

	if _, err = easyjson.MarshalToWriter(usr, c.Response().Writer); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (h *UserHandler) GetOtherUserProfile(c echo.Context) error {
	defer c.Request().Body.Close()

	uid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	usr, err := h.UseCase.GetOtherProfile(uint64(uid))

	if err != nil {
		log.Println(err)
		return err
	}

	if _, err = easyjson.MarshalToWriter(usr, c.Response().Writer); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (h *UserHandler) UploadAvatar(c echo.Context) error {
	defer c.Request().Body.Close()

	cookie, err := c.Cookie("SID")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	var uid uint64
	var exists bool
	if cookie != nil {
		exists, uid, err = h.Sm.CheckSession(cookie.Value)
		if err != nil {
			log.Println(err)
			return err
		}

		if !exists {
			return echo.NewHTTPError(http.StatusBadRequest, "user is not authorized")
		}
	}

	img, err := c.FormFile("avatar")
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if img == nil {
		return nil
	}

	err = h.UseCase.UploadAvatar(uid, img)

	return nil
}

func (h *UserHandler) GetAvatar(c echo.Context) error {
	defer c.Request().Body.Close()

	uid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	file, err := h.UseCase.GetAvatar(uint64(uid))

	if err != nil {
		log.Println(err)
		return err
	}
	c.Response().Write(file)

	return nil
}

func (h *UserHandler) CheckUser(c echo.Context) error {
	cookie, err := c.Cookie(constants.SessionCookieName)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	var exists bool
	if cookie != nil {
		exists, _, err = h.Sm.CheckSession(cookie.Value)
		if err != nil {
			log.Println(err)
			return err
		}

		if !exists {
			return echo.NewHTTPError(http.StatusBadRequest, "user is not authorized")
		}
	}

	u := &models.User{}

	err = easyjson.UnmarshalFromReader(c.Request().Body, u)
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	_, err = h.UseCase.CheckUser(u)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
