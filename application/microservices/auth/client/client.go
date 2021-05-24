package client

/*
protoc --go_out=plugins=grpc:. *.proto
*/

import (
	"context"
	"github.com/labstack/echo"
	traceutils "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"kudago/application/microservices/auth/proto"
	"kudago/pkg/constants"
	"kudago/pkg/logger"
	"net/http"
)

type AuthClient struct {
	client proto.AuthClient
	gConn  *grpc.ClientConn
	logger logger.Logger
}

func NewAuthClient(port string, logger logger.Logger, tracer opentracing.Tracer) (IAuthClient, error) {
	gConn, err := grpc.Dial(
		constants.Localhost+port,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(traceutils.OpenTracingClientInterceptor(tracer)),
	)
	if err != nil {
		return nil, err
	}

	return &AuthClient{client: proto.NewAuthClient(gConn), gConn: gConn, logger: logger}, nil
}

func (a *AuthClient) Login(login string, password string, value string) (uint64, string, error, int) {
	usr := &proto.User{
		Login:    login,
		Password: password,
		Value:    value,
	}

	answer, err := a.client.Login(context.Background(), usr)
	if err != nil {
		return 0, "", err, http.StatusInternalServerError
	}
	if answer.Flag {
		return 0, "", echo.NewHTTPError(http.StatusBadRequest, answer.Msg), http.StatusBadRequest
	}

	return answer.UserId, answer.Value, nil, http.StatusOK
}

func (a *AuthClient) Check(value string) (bool, uint64, error, int) {
	sessionValue := &proto.Session{Value: value}

	answer, err := a.client.Check(context.Background(), sessionValue)
	if err != nil {
		return false, 0, err, http.StatusInternalServerError
	}

	return answer.Answer, answer.UserId, err, http.StatusOK
}

func (a *AuthClient) Logout(value string) (error, int) {
	sessionValue := &proto.Session{Value: value}

	answer, err := a.client.Logout(context.Background(), sessionValue)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	if answer.Flag {
		return echo.NewHTTPError(http.StatusBadRequest, answer.Msg), http.StatusBadRequest
	}

	return nil, http.StatusOK
}

func (a *AuthClient) Close() {
	if err := a.gConn.Close(); err != nil {
		a.logger.Warn(err)
	}
}
