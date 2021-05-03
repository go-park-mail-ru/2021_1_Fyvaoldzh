package client

/*
protoc --go_out=plugins=grpc:. *.proto
*/

import (
	"context"
	"github.com/labstack/echo"
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

func NewAuthClient(port string, logger logger.Logger) (*AuthClient, error) {
	gConn, err := grpc.Dial(
		constants.Localhost+port,
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	return &AuthClient{client: proto.NewAuthClient(gConn), gConn: gConn, logger: logger}, nil
}

func (a *AuthClient) Login(login string, password string, value string) (uint64, string, error) {
	usr := &proto.User{
		Login:    login,
		Password: password,
		Value:    value,
	}

	answer, err := a.client.Login(context.Background(), usr)
	if err != nil {
		return 0, "", err
	}
	if answer.Flag {
		return 0, "", echo.NewHTTPError(http.StatusBadRequest, answer.Msg)
	}

	return answer.UserId, answer.Value, nil
}

func (a *AuthClient) Check(value string) (bool, uint64, error) {
	sessionValue := &proto.Session{Value: value}

	answer, err := a.client.Check(context.Background(), sessionValue)
	if err != nil {
		return false, 0, err
	}

	return answer.Answer, answer.UserId, err
}

func (a *AuthClient) Logout(value string) error {
	sessionValue := &proto.Session{Value: value}

	answer, err := a.client.Logout(context.Background(), sessionValue)
	if err != nil {
		return err
	}
	if answer.Flag {
		return echo.NewHTTPError(http.StatusBadRequest, answer.Msg)
	}

	return nil
}

func (a *AuthClient) Close() {
	if err := a.gConn.Close(); err != nil {
		a.logger.Warn(err)
	}
}
