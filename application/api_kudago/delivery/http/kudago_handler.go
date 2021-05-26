package kudago_http

import (
	"context"
	"errors"
	"github.com/labstack/echo"
	kudago_client "kudago/application/microservices/api_kudago/client"
	"kudago/application/server/middleware"
	"kudago/pkg/logger"
	"net/http"
	"strconv"
)

type KudagoHandler struct {
	Logger    logger.Logger
	rpcKudago *kudago_client.KudagoClient
}

func CreateKudagoHandler(e *echo.Echo, rKudago *kudago_client.KudagoClient, lg logger.Logger) {
	kudagoHandler := KudagoHandler{
		Logger:    lg,
		rpcKudago: rKudago,
	}
	e.GET("/api/v1/kudago/basic",
		kudagoHandler.AddBasic)
	e.GET("/api/v1/kudago/today",
		kudagoHandler.AddToday)
}

func (kh KudagoHandler) AddBasic(c echo.Context) error {
	defer c.Request().Body.Close()

	num := c.QueryParam("num")
	if num == "" {
		num = "1"
	}
	intNum, err := strconv.Atoi(num)
	if err != nil {
		return echo.NewHTTPError(http.StatusTeapot, errors.New("num must be a number"))
	}
	if intNum == 0 {
		intNum = 1
	}

	err, code := kh.rpcKudago.AddBasic(context.Background(), uint64(intNum))
	if err != nil {
		kh.Logger.Warn(err)
		middleware.ErrResponse(c, code)
		return echo.NewHTTPError(code, err.Error())
	}
	middleware.OkResponse(c)
	return nil
}

func (kh KudagoHandler) AddToday(c echo.Context) error {
	defer c.Request().Body.Close()

	err, code := kh.rpcKudago.AddToday(context.Background())
	if err != nil {
		kh.Logger.Warn(err)
		middleware.ErrResponse(c, code)
		return echo.NewHTTPError(code, err.Error())
	}
	middleware.OkResponse(c)
	return nil
}
