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

func (a *AuthServer) Login(c context.Context, usr *proto.User) (*proto.LoginAnswer, error) {
	sessionValue := usr.Value
	if len(sessionValue) != 0 {
		flag, _, err := a.usecase.Check(sessionValue)
		if err != nil {
			return &proto.LoginAnswer{Value: sessionValue, Flag: false}, err
		}
		if flag {
			return &proto.LoginAnswer{Value: sessionValue,
				Flag: true,
				Msg:  "user is already logged in"}, nil
		}
	}

	sessionValue, flag, err := a.usecase.Login(usr.Login, usr.Password)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if flag {
		return &proto.LoginAnswer{Value: sessionValue,
			Flag: true,
			Msg:  "incorrect data"}, nil
	}

	return &proto.LoginAnswer{Value: sessionValue, Flag: false}, nil
}

func (a *AuthServer) Check(c context.Context, s *proto.Session) (*proto.CheckAnswer, error) {
	flag, userId, err := a.usecase.Check(s.Value)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &proto.CheckAnswer{Answer: flag, UserId: userId}, nil
}

func (a *AuthServer) Logout(c context.Context, s *proto.Session) (*proto.LogoutAnswer, error) {
	flag, _, err := a.usecase.Check(s.Value)
	if err != nil {
		return &proto.LogoutAnswer{}, status.Error(codes.Internal, err.Error())
	}
	if !flag {
		return &proto.LogoutAnswer{Flag: true, Msg: "user is not authorized"}, nil
	}

	err = a.usecase.Logout(s.Value)
	if err != nil {
		log.Println(err)
		return &proto.LogoutAnswer{}, err
	}

	return &proto.LogoutAnswer{Flag: false}, nil
}
