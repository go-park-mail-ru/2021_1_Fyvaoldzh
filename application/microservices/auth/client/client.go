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
		Value: value,
	}

	if value != "" {
		flag, _, err := a.Check(value)
		if err != nil {
			return 0, "", err
		}
		if flag {
			return 0, "", echo.NewHTTPError(http.StatusBadRequest, "user is already logged in")
		}

	}

	answer, err := a.client.Login(context.Background(), usr)
	if err != nil {
		return 0, "", err
	}

	return answer.UserId, answer.Value, nil
}

func (a *AuthClient) Check(value string) (bool, uint64, error) {
	sessionValue := &proto.Session{Value: value}

	checkAnswer, err := a.client.Check(context.Background(), sessionValue)
	if err != nil {
		return false, 0, err
	}

	return checkAnswer.Answer, checkAnswer.UserId, err
}

func (a *AuthClient) Logout(value string) error {
	sessionValue := &proto.Session{Value: value}

	_, err := a.client.Logout(context.Background(), sessionValue)
	if err != nil {
		if err.Error() == "rpc error: code = InvalidArgument desc = user is not authorized" {
			return echo.NewHTTPError(http.StatusUnauthorized, "user is not authorized")
		}
		return err
	}

	return nil
}

func (a *AuthClient) Close() {
	if err := a.gConn.Close(); err != nil {
		a.logger.Warn(err)
	}
}
