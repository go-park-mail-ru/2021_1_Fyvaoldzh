package http

import (
	"github.com/labstack/echo"
	"github.com/mailru/easyjson"
	"kudago/application/user"
	"kudago/models"
	"log"
	"net/http"
)

type UserHandler struct {
	UseCase user.UseCase
}

func CreateUserHandler(e *echo.Echo, uc user.UseCase){

	userHandler := UserHandler{UseCase: uc}

	e.POST("/api/v1/login", userHandler.Login)
	e.DELETE("/api/v1/logout", userHandler.Logout)
	e.POST("/api/v1/register", userHandler.Register)
	e.PUT("/api/v1/profile", userHandler.Update)
	e.PUT("/api/v1/upload_avatar", userHandler.UploadAvatar)

}


func (h *UserHandler) Login(c echo.Context) error {
	/*
	defer c.Request().Body.Close()
	u := &models.User{}

	cookie, err := c.Cookie(constants.SessionCookieName)
	if err == nil && h.Store[cookie.Value] != 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "user is already logged in")
	}

	err = easyjson.UnmarshalFromReader(c.Request().Body, u)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	hash := sha256.New()
	hash.Write([]byte(u.Password))
	u.Password = base64.URLEncoding.EncodeToString(hash.Sum(nil))

	isCorrect, uid := IsCorrectUser(h, u)
	if !isCorrect {
		return echo.NewHTTPError(http.StatusBadRequest, "incorrect data")
	}

	c.SetCookie(h.CreateCookie(constants.CookieLength, uid))
	*/


	return nil
}

func (h *UserHandler) Logout(c echo.Context) error {
	defer c.Request().Body.Close()

	/*
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

	 */

	return nil
}

func (h *UserHandler) Register(c echo.Context) error {
	defer c.Request().Body.Close()

	newData := &models.RegData{}

	err := easyjson.UnmarshalFromReader(c.Request().Body, newData)
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	h.UseCase.Add(newData)
	/*
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

	c.SetCookie(h.CreateCookie(constants.CookieLength, uid))

	 */

	return nil
}

func (h *UserHandler) Update(c echo.Context) error {
	defer c.Request().Body.Close()

	/*

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

	 */

	return nil
}

/*
func (h *UserHandler) CreateCookie(n uint8, uid uint64) *http.Cookie {

	key := generator.RandStringRunes(n)

	newCookie := &http.Cookie{
		Name:     constants.SessionCookieName,
		Value:    key,
		Expires:  time.Now().Add(10 * time.Hour),
		HttpOnly: true,
	}

	h.Store[key] = uid
	return newCookie
}

func CreateProfileHandler(e *echo.Echo, uc profile.UseCase) error {

	profileHandler := ProfileHandler{useCase: uc}

	e.PUT("/api/v1/upload_avatar", profileHandler.UploadAvatar)
	return nil
}

func (h *UserHandler) CreateProfile(data *models.RegData) (uint64, error) {
	newUser := &models.User{}
	newUser.Login = data.Login

	hash := sha256.New()
	hash.Write([]byte(data.Password))
	newUser.Password = base64.URLEncoding.EncodeToString(hash.Sum(nil))

	if IsExistingUser(h, newUser) {
		return 0, echo.NewHTTPError(http.StatusBadRequest, "user with this login does exist")
	}


	h.Mu.Lock()
	newUser.Id = id
	id++

	newProfile := &models.UserOwnProfile{}
	newProfile.Uid = newUser.Id
	newProfile.Name = data.Name
	h.UserBase = append(h.UserBase, newUser)
	h.ProfileBase = append(h.ProfileBase, newProfile)
	h.Mu.Unlock()
	return newUser.Id, nil
}

 */

func (h *UserHandler) GetOwnProfile(c echo.Context) error {
	defer c.Request().Body.Close()
	/*

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


	 */
	return nil
}

func (h *UserHandler) GetOtherUserProfile(c echo.Context) error {
	defer c.Request().Body.Close()

	/*
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


	 */
	return nil
}


func (h *UserHandler) UploadAvatar(c echo.Context) error {
	defer c.Request().Body.Close()

	/*
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

	 */
	return nil
}

func (h *UserHandler) GetAvatar(c echo.Context) error {
	defer c.Request().Body.Close()

	/*
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

	 */

	return nil
}

