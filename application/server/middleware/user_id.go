package middleware

import (
	"net/http"

	"kudago/application/microservices/auth/client"
	"kudago/pkg/constants"

	"github.com/labstack/echo"
)

type Auth struct {
	rpcAuth client.IAuthClient
}

func NewAuth(rpcAuth client.IAuthClient) Auth {
	return Auth{rpcAuth: rpcAuth}
}

func (a Auth) GetSession(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		cookie, err := ctx.Cookie(constants.SessionCookieName)
		if err != nil && cookie != nil {
			return echo.NewHTTPError(http.StatusForbidden, "user is not authorized")
		}

		var uid uint64
		var exists bool
		var code int

		if cookie == nil {
			ErrResponse(ctx, http.StatusForbidden)
			return echo.NewHTTPError(http.StatusForbidden, "user is not authorized")
		}
		exists, uid, err, code = a.rpcAuth.Check(cookie.Value)
		if err != nil {
			ErrResponse(ctx, code)
			return echo.NewHTTPError(http.StatusInternalServerError, "cannot check cookie")
		}

		if !exists {
			ErrResponse(ctx, http.StatusForbidden)
			return echo.NewHTTPError(http.StatusForbidden, "user is not authorized")
		}

		ctx.Set(constants.SessionCookieName, cookie.Value)
		ctx.Set(constants.UserIdKey, uid)
		return next(ctx)
	}
}
