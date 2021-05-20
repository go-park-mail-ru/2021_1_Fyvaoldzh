package http

import (
	"fmt"
	"kudago/application/microservices/auth/client"
	"kudago/application/models"
	"kudago/application/server/middleware"
	"kudago/application/user"
	"kudago/pkg/constants"
	"kudago/pkg/custom_sanitizer"
	"kudago/pkg/generator"
	"kudago/pkg/logger"
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/mailru/easyjson"
)

type UserHandler struct {
	UseCase   user.UseCase
	rpcAuth   client.IAuthClient
	Logger    logger.Logger
	sanitizer *custom_sanitizer.CustomSanitizer
	auth      middleware.Auth
}

func CreateUserHandler(e *echo.Echo,
	uc user.UseCase,
	auth client.IAuthClient,
	sz *custom_sanitizer.CustomSanitizer,
	logger logger.Logger,
	am middleware.Auth) {
	userHandler := UserHandler{
		UseCase:   uc,
		rpcAuth:   auth,
		sanitizer: sz,
		Logger:    logger,
		auth:      am}

	e.POST("/api/v1/login", userHandler.Login)
	e.DELETE("/api/v1/logout",
		userHandler.Logout,
		am.GetSession)
	e.POST("/api/v1/register", userHandler.Register)
	e.GET("/api/v1/profile",
		userHandler.GetOwnProfile,
		am.GetSession)
	e.GET("/api/v1/profile/:id",
		userHandler.GetOtherUserProfile,
		middleware.GetId)
	e.PUT("/api/v1/profile",
		userHandler.Update,
		am.GetSession)
	e.POST("/api/v1/upload_avatar",
		userHandler.UploadAvatar,
		am.GetSession)
	e.GET("/api/v1/avatar/:id",
		userHandler.GetAvatar,
		middleware.GetId)
	e.GET("/api/v1/users",
		userHandler.GetUsers,
		middleware.GetPage)
	e.GET("/api/v1/find",
		userHandler.FindUsers,
		middleware.GetPage)
	e.GET("/api/v1/actions",
		userHandler.GetActions,
		am.GetSession,
		middleware.GetPage)
}

func (uh *UserHandler) Login(c echo.Context) error {
	defer c.Request().Body.Close()
	//start := time.Now()
	//requestId := fmt.Sprintf("%016x", rand.Int())
	u := &models.User{}

	cookie, err := c.Cookie(constants.SessionCookieName)

	// если убрать куки нил, то упадет с no cookie
	if err != nil && cookie != nil {
		uh.Logger.Warn(err)
		return err
	}

	err = easyjson.UnmarshalFromReader(c.Request().Body, u)
	if err != nil {
		uh.Logger.Warn(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var value string
	var code int
	if cookie != nil {
		_, value, err, code = uh.rpcAuth.Login(u.Login, u.Password, cookie.Value)
	} else {
		_, value, err, code = uh.rpcAuth.Login(u.Login, u.Password, "")
	}
	if err != nil {
		uh.Logger.Warn(err)
		middleware.ErrResponse(c, code)
		return err
	}

	cookie = generator.CreateCookieWithValue(value)
	c.SetCookie(cookie)

	middleware.OkResponse(c)
	return nil
}

func (uh *UserHandler) Logout(c echo.Context) error {
	defer c.Request().Body.Close()
	//start := time.Now()
	//requestId := fmt.Sprintf("%016x", rand.Int())

	cookie, _ := c.Cookie(constants.SessionCookieName)

	err, code := uh.rpcAuth.Logout(cookie.Value)
	if err != nil {
		uh.Logger.Warn(err)
		middleware.ErrResponse(c, code)
		return err
	}

	cookie.Expires = time.Now().AddDate(0, 0, -1)
	c.SetCookie(cookie)

	middleware.OkResponse(c)
	return nil
}

func (uh *UserHandler) Register(c echo.Context) error {
	defer c.Request().Body.Close()
	//start := time.Now()
	//requestId := fmt.Sprintf("%016x", rand.Int())

	cookie, err := c.Cookie(constants.SessionCookieName)
	if err != nil && cookie != nil {
		uh.Logger.Warn(err)
		return echo.NewHTTPError(http.StatusBadRequest, "error getting cookie")
	}

	if cookie != nil {
		exists, _, err, code := uh.rpcAuth.Check(cookie.Value)
		if err != nil {
			uh.Logger.Warn(err)
			middleware.ErrResponse(c, code)
			return err
		}

		if exists {
			middleware.ErrResponse(c, http.StatusBadRequest)
			return echo.NewHTTPError(http.StatusBadRequest, "user is already logged in")
		}
	}

	newData := &models.RegData{}

	err = easyjson.UnmarshalFromReader(c.Request().Body, newData)
	if err != nil {
		uh.Logger.Warn(err)
		middleware.ErrResponse(c, http.StatusBadRequest)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	oldPassword := newData.Password
	_, err = uh.UseCase.Add(newData)
	if err != nil {
		middleware.ErrResponse(c, http.StatusInternalServerError)
		uh.Logger.Warn(err)
		return err
	}

	_, value, err, code := uh.rpcAuth.Login(newData.Login, oldPassword, "")
	if err != nil {
		middleware.ErrResponse(c, code)
		uh.Logger.Warn(err)
		return err
	}

	cookie = generator.CreateCookieWithValue(value)
	c.SetCookie(cookie)
	middleware.OkResponse(c)
	return nil
}

func (uh *UserHandler) Update(c echo.Context) error {
	defer c.Request().Body.Close()
	start := time.Now()
	requestId := fmt.Sprintf("%016x", rand.Int())

	uid := c.Get(constants.UserIdKey).(uint64)

	ud := &models.UserOwnProfile{}
	err := easyjson.UnmarshalFromReader(c.Request().Body, ud)
	if err != nil {
		uh.Logger.LogError(c, start, requestId, err)
		middleware.ErrResponse(c, http.StatusBadRequest)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = uh.UseCase.Update(uid, ud)
	if err != nil {
		middleware.ErrResponse(c, http.StatusInternalServerError)
		uh.Logger.LogError(c, start, requestId, err)
		return err
	}
	middleware.OkResponse(c)
	return nil
}

func (uh *UserHandler) GetOwnProfile(c echo.Context) error {
	defer c.Request().Body.Close()
	start := time.Now()
	requestId := fmt.Sprintf("%016x", rand.Int())

	uid := c.Get(constants.UserIdKey).(uint64)

	usr, err := uh.UseCase.GetOwnProfile(uid)
	if err != nil {
		uh.Logger.LogError(c, start, requestId, err)
		middleware.ErrResponse(c, http.StatusInternalServerError)
		return err
	}

	uh.sanitizer.SanitizeOwnProfile(usr)
	if _, err = easyjson.MarshalToWriter(usr, c.Response().Writer); err != nil {
		uh.Logger.LogError(c, start, requestId, err)
		middleware.ErrResponse(c, http.StatusInternalServerError)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	middleware.OkResponse(c)
	return nil
}

func (uh *UserHandler) GetOtherUserProfile(c echo.Context) error {
	defer c.Request().Body.Close()

	uid := c.Get(constants.IdKey).(int)

	usr, err := uh.UseCase.GetOtherProfile(uint64(uid))
	if err != nil {
		uh.Logger.Warn(err)
		middleware.ErrResponse(c, http.StatusInternalServerError)
		return err
	}

	uh.sanitizer.SanitizeOtherProfile(usr)
	if _, err = easyjson.MarshalToWriter(usr, c.Response().Writer); err != nil {
		uh.Logger.Warn(err)
		middleware.ErrResponse(c, http.StatusInternalServerError)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	middleware.OkResponse(c)
	return nil
}

func (uh *UserHandler) UploadAvatar(c echo.Context) error {
	defer c.Request().Body.Close()

	uid := c.Get(constants.UserIdKey).(uint64)

	img, err := c.FormFile("avatar")
	if err != nil {
		uh.Logger.Warn(err)
		middleware.ErrResponse(c, http.StatusBadRequest)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if img == nil {
		return nil
	}

	src, err := img.Open()
	if err != nil {
		uh.Logger.Warn(err)
		middleware.ErrResponse(c, http.StatusInternalServerError)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer src.Close()
	_ = uh.UseCase.UploadAvatar(uid, src, img.Filename)
	middleware.OkResponse(c)
	return nil
}

func (uh *UserHandler) GetAvatar(c echo.Context) error {
	defer c.Request().Body.Close()
	start := time.Now()
	requestId := fmt.Sprintf("%016x", rand.Int())

	uid := c.Get(constants.IdKey).(int)

	file, err := uh.UseCase.GetAvatar(uint64(uid))
	if err != nil {
		uh.Logger.LogError(c, start, requestId, err)
		middleware.ErrResponse(c, http.StatusInternalServerError)
		return err
	}

	_, _ = c.Response().Write(file)
	middleware.OkResponse(c)
	return nil
}

func (uh *UserHandler) GetUsers(c echo.Context) error {
	defer c.Request().Body.Close()
	start := time.Now()
	requestId := fmt.Sprintf("%016x", rand.Int())

	page := c.Get(constants.PageKey).(int)

	users, err := uh.UseCase.GetUsers(page)
	if err != nil {
		uh.Logger.LogError(c, start, requestId, err)
		return err
	}

	users = uh.sanitizer.SanitizeUserCards(users)
	if _, err = easyjson.MarshalToWriter(users, c.Response().Writer); err != nil {
		uh.Logger.LogError(c, start, requestId, err)
		return err
	}
	middleware.OkResponse(c)
	return nil
}

func (uh UserHandler) FindUsers(c echo.Context) error {
	defer c.Request().Body.Close()

	//start := time.Now()
	//requestId := fmt.Sprintf("%016x", rand.Int())
	str := c.QueryParam("search")
	page := c.Get(constants.PageKey).(int)

	users, err := uh.UseCase.FindUsers(str, page)
	if err != nil {
		middleware.ErrResponse(c, http.StatusInternalServerError)
		return err
	}
	users = uh.sanitizer.SanitizeUserCards(users)

	if _, err := easyjson.MarshalToWriter(users, c.Response().Writer); err != nil {
		uh.Logger.Warn(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	middleware.OkResponse(c)
	return nil
}

func (uh *UserHandler) GetActions(c echo.Context) error {
	defer c.Request().Body.Close()

	uid := c.Get(constants.UserIdKey).(uint64)
	page := c.Get(constants.PageKey).(int)

	actions, err := uh.UseCase.GetActions(uid, page)
	if err != nil {
		middleware.ErrResponse(c, http.StatusInternalServerError)
		uh.Logger.Warn(err)
		return err
	}

	actions = uh.sanitizer.SanitizeActions(actions)
	if _, err = easyjson.MarshalToWriter(actions, c.Response().Writer); err != nil {
		uh.Logger.Warn(err)
		middleware.ErrResponse(c, http.StatusInternalServerError)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	middleware.OkResponse(c)
	return nil
}
