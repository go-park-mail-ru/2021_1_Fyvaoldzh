package http

import (
	"fmt"
	"kudago/application/microservices/auth/client"
	"kudago/application/models"
	"kudago/application/user"
	"kudago/pkg/constants"
	"kudago/pkg/custom_sanitizer"
	"kudago/pkg/generator"
	"kudago/pkg/logger"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/mailru/easyjson"
)

type UserHandler struct {
	UseCase   user.UseCase
	rpcAuth client.AuthClient
	Logger    logger.Logger
	sanitizer *custom_sanitizer.CustomSanitizer
}

func CreateUserHandler(e *echo.Echo, uc user.UseCase,
	auth client.AuthClient, sz *custom_sanitizer.CustomSanitizer, logger logger.Logger) {
	userHandler := UserHandler{
		UseCase: uc,
		rpcAuth: auth,
		sanitizer: sz,
		Logger: logger}

	e.POST("/api/v1/login", userHandler.Login)
	e.DELETE("/api/v1/logout", userHandler.Logout)
	e.POST("/api/v1/register", userHandler.Register)
	e.GET("/api/v1/profile", userHandler.GetOwnProfile)
	e.GET("/api/v1/profile/:id", userHandler.GetOtherUserProfile)
	e.PUT("/api/v1/profile", userHandler.Update)
	e.POST("/api/v1/upload_avatar", userHandler.UploadAvatar)
	e.GET("/api/v1/avatar/:id", userHandler.GetAvatar)
	e.GET("/api/v1/users", userHandler.GetUsers)
}

func (uh *UserHandler) Login(c echo.Context) error {
	defer c.Request().Body.Close()
	start := time.Now()
	requestId := fmt.Sprintf("%016x", rand.Int())
	u := &models.User{}

	cookie, err := c.Cookie(constants.SessionCookieName)

	// если убрать куки нил, то упадет с no cookie
	if err != nil && cookie != nil {
		uh.Logger.LogError(c, start, requestId, err)
		return err
	}

	err = easyjson.UnmarshalFromReader(c.Request().Body, u)
	if err != nil {
		uh.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var value string
	if cookie != nil {
		_, value, err = uh.rpcAuth.Login(u.Login, u.Password, cookie.Value)
	} else {
		_, value, err = uh.rpcAuth.Login(u.Login, u.Password, "")
	}
	if err != nil {
		uh.Logger.LogError(c, start, requestId, err)
		return err
	}

	cookie = generator.CreateCookieWithValue(value)
	c.SetCookie(cookie)

	uh.Logger.LogInfo(c, start, requestId)
	return nil
}

func (uh *UserHandler) Logout(c echo.Context) error {
	defer c.Request().Body.Close()
	start := time.Now()
	requestId := fmt.Sprintf("%016x", rand.Int())

	cookie, err := c.Cookie(constants.SessionCookieName)
	if err != nil && cookie != nil {
		uh.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusBadRequest, "error getting cookie")
	}

	if cookie == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	err = uh.rpcAuth.Logout(cookie.Value)
	if err != nil {
		uh.Logger.LogError(c, start, requestId, err)
		return err
	}

	cookie.Expires = time.Now().AddDate(0, 0, -1)
	c.SetCookie(cookie)

	return nil
}

func (uh *UserHandler) Register(c echo.Context) error {
	defer c.Request().Body.Close()
	start := time.Now()
	requestId := fmt.Sprintf("%016x", rand.Int())

	cookie, err := c.Cookie(constants.SessionCookieName)
	if err != nil && cookie != nil {
		uh.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusBadRequest, "error getting cookie")
	}

	if cookie != nil {
		exists, _, err := uh.rpcAuth.Check(cookie.Value)
		if err != nil {
			uh.Logger.LogError(c, start, requestId, err)
			return err
		}

		if exists {
			return echo.NewHTTPError(http.StatusBadRequest, "user is already logged in")
		}
	}

	newData := &models.RegData{}

	err = easyjson.UnmarshalFromReader(c.Request().Body, newData)
	if err != nil {
		uh.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	oldPassword := newData.Password
	_, err = uh.UseCase.Add(newData)
	if err != nil {
		uh.Logger.LogError(c, start, requestId, err)
		return err
	}

	_, value, err := uh.rpcAuth.Login(newData.Login, oldPassword, "")
	if err != nil {
		uh.Logger.LogError(c, start, requestId, err)
		return err
	}

	cookie = generator.CreateCookieWithValue(value)
	c.SetCookie(cookie)
	return nil
}

func (uh *UserHandler) Update(c echo.Context) error {
	defer c.Request().Body.Close()
	start := time.Now()
	requestId := fmt.Sprintf("%016x", rand.Int())

	cookie, err := c.Cookie(constants.SessionCookieName)
	if err != nil {
		uh.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	var uid uint64
	var exists bool
	if cookie != nil {
		exists, uid, err = uh.rpcAuth.Check(cookie.Value)
		if err != nil {
			uh.Logger.LogError(c, start, requestId, err)
			return err
		}

		if !exists {
			return echo.NewHTTPError(http.StatusBadRequest, "user is not authorized")
		}
	}

	ud := &models.UserOwnProfile{}
	err = easyjson.UnmarshalFromReader(c.Request().Body, ud)
	if err != nil {
		uh.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = uh.UseCase.Update(uid, ud)
	if err != nil {
		uh.Logger.LogError(c, start, requestId, err)
		return err
	}

	return nil
}

func (uh *UserHandler) GetOwnProfile(c echo.Context) error {
	defer c.Request().Body.Close()
	start := time.Now()
	requestId := fmt.Sprintf("%016x", rand.Int())

	cookie, err := c.Cookie("SID")
	if err != nil {
		uh.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	var uid uint64
	var exists bool
	exists, uid, err = uh.rpcAuth.Check(cookie.Value)
	if err != nil {
		uh.Logger.LogError(c, start, requestId, err)
		return err
	}

	if !exists {
		return echo.NewHTTPError(http.StatusBadRequest, "user is not authorized")
	}

	usr, err := uh.UseCase.GetOwnProfile(uid)
	if err != nil {
		uh.Logger.LogError(c, start, requestId, err)
		return err
	}

	uh.sanitizer.SanitizeOwnProfile(usr)
	if _, err = easyjson.MarshalToWriter(usr, c.Response().Writer); err != nil {
		uh.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (uh *UserHandler) GetOtherUserProfile(c echo.Context) error {
	defer c.Request().Body.Close()

	uid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		uh.Logger.Warn(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if uid <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "incorrect data")
	}

	usr, err := uh.UseCase.GetOtherProfile(uint64(uid))
	if err != nil {
		uh.Logger.Warn(err)
		return err
	}

	uh.sanitizer.SanitizeOtherProfile(usr)
	if _, err = easyjson.MarshalToWriter(usr, c.Response().Writer); err != nil {
		uh.Logger.Warn(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (uh *UserHandler) UploadAvatar(c echo.Context) error {
	defer c.Request().Body.Close()

	cookie, err := c.Cookie("SID")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
	}

	var uid uint64
	var exists bool
	if cookie != nil {
		exists, uid, err = uh.rpcAuth.Check(cookie.Value)
		if err != nil {
			uh.Logger.Warn(err)
			return err
		}

		if !exists {
			return echo.NewHTTPError(http.StatusBadRequest, "user is not authorized")
		}
	}

	img, err := c.FormFile("avatar")
	if err != nil {
		uh.Logger.Warn(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if img == nil {
		return nil
	}

	src, err := img.Open()
	if err != nil {
		uh.Logger.Warn(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer src.Close()
	err = uh.UseCase.UploadAvatar(uid, src, img.Filename)

	return nil
}

func (uh *UserHandler) GetAvatar(c echo.Context) error {
	defer c.Request().Body.Close()
	start := time.Now()
	requestId := fmt.Sprintf("%016x", rand.Int())

	uid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		uh.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	file, err := uh.UseCase.GetAvatar(uint64(uid))
	if err != nil {
		uh.Logger.LogError(c, start, requestId, err)
		return err
	}

	c.Response().Write(file)

	return nil
}

func (uh *UserHandler) GetUsers(c echo.Context) error {
	defer c.Request().Body.Close()
	start := time.Now()
	requestId := fmt.Sprintf("%016x", rand.Int())

	page := c.QueryParam("page")
	pageNum, err := strconv.Atoi(page)
	if err != nil {
		uh.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if pageNum < 1 {
		return echo.NewHTTPError(http.StatusBadRequest, "incorrect data")
	}

	users, err := uh.UseCase.GetUsers(pageNum)
	if err != nil {
		uh.Logger.LogError(c, start, requestId, err)
		return err
	}

	users = uh.sanitizer.SanitizeUserCards(users)
	if _, err = easyjson.MarshalToWriter(users, c.Response().Writer); err != nil {
		uh.Logger.LogError(c, start, requestId, err)
		return err
	}

	return nil
}
