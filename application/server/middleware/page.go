package middleware

import (
	"errors"
	"kudago/pkg/constants"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

func GetPage(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pageParam := ctx.QueryParam("page")
		if pageParam == "" {
			pageParam = "1"
		}
		page, err := strconv.Atoi(pageParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusTeapot, errors.New("page must be a number"))
		}
		if page == 0 {
			page = 1
		}

		ctx.Set(constants.PageKey, page)

		return next(ctx)
	}
}
