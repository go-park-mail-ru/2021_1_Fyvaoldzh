package server

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"kudago/application/microservices/auth/proto"
	"kudago/application/microservices/auth/session"
	"log"
)

type AuthServer struct {
	usecase session.UseCase
}

func NewAuthServer(usecase session.UseCase) *AuthServer {
	return &AuthServer{usecase: usecase}
}

func (a *AuthServer) Login(ctx context.Context, usr *proto.User) (*proto.LoginAnswer, error) {
	sessionValue := usr.Value
	if len(sessionValue) != 0 {
		flag, _, err := a.usecase.Check(sessionValue)
		if err != nil {
			return &proto.LoginAnswer{Value: sessionValue}, err
		}
		if flag {
			return &proto.LoginAnswer{Value: sessionValue},
			status.Error(codes.AlreadyExists, "user is already logged in")
		}
	}

	sessionValue,  err := a.usecase.Login(usr.Login, usr.Password)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &proto.LoginAnswer{
		Value: sessionValue,
	}, nil
}

func (a *AuthServer) Check(ctx context.Context, s *proto.Session) (*proto.CheckAnswer, error) {
	flag, userId, err := a.usecase.Check(s.Value)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &proto.CheckAnswer{Answer: flag, UserId: userId}, nil
}

func (a *AuthServer) Logout(ctx context.Context, s *proto.Session) (*proto.Empty, error) {
	err := a.usecase.Logout(s.Value)
	if err != nil {
		log.Println(err)
		return &proto.Empty{}, err
	}

	return &proto.Empty{}, nil
}
