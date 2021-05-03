package usecase

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"kudago/application/microservices/auth/session"
	"kudago/application/microservices/auth/user"
	"kudago/pkg/constants"
	"kudago/pkg/generator"
	"kudago/pkg/logger"
)

type SessionUseCase struct {
	repo    session.Repository
	userUseCase user.UseCase
	logger  logger.Logger
}

func NewSessionUseCase(s session.Repository, u user.UseCase, logger logger.Logger) session.UseCase {
	return &SessionUseCase{repo: s, userUseCase: u, logger: logger}
}

func (s SessionUseCase) Login(login string, password string) (string, error) {
	userId, err := s.userUseCase.CheckUser(login, password)
	if err != nil {
		return "", err
	}
	cookie := generator.CreateCookieValue(constants.CookieLength)
	err = s.repo.InsertSession(userId, cookie)
	if err != nil {
		return "", err
	}
	return cookie, nil
}

func (s SessionUseCase) Check(value string) (bool, uint64, error) {
	return s.repo.CheckSession(value)
}

func (s SessionUseCase) Logout(value string) error {
	flag, _, err := s.Check(value)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	if !flag {
		return status.Error(codes.InvalidArgument, "user is not authorized")
	}

	return s.repo.DeleteSession(value)
}
