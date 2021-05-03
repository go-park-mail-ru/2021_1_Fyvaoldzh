package usecase

import (
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

func (s SessionUseCase) Login(login string, password string) (string, bool, error) {
	userId, flag, err := s.userUseCase.CheckUser(login, password)
	if err != nil {
		return "", false, err
	}
	if flag {
		return "", true, nil
	}
	var cookie string
	// хз, как так получилось, но генератор плохо генерит )0)
	for {
		cookie = generator.CreateCookieValue(constants.CookieLength)
		err = s.repo.InsertSession(userId, cookie)
		if !(err != nil && err.Error() == "Duplicate key exists in unique index 'primary' in space 'qdago' (0x3)") {
			break
		}
	}
	if err != nil {
		return "", false, err
	}
	return cookie, false, nil
}

func (s SessionUseCase) Check(value string) (bool, uint64, error) {
	return s.repo.CheckSession(value)
}

func (s SessionUseCase) Logout(value string) error {
	return s.repo.DeleteSession(value)
}
