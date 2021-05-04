package usecase

import (
	"kudago/application/microservices/auth/user"
	"kudago/pkg/generator"
	"kudago/pkg/logger"
)

type UserUseCase struct {
	repo   user.Repository
	logger logger.Logger
}

func NewUserUseCase(u user.Repository, logger logger.Logger) user.UseCase {
	return &UserUseCase{repo: u, logger: logger}
}

func (uc UserUseCase) CheckUser(login string, password string) (uint64, bool, error) {
	gotUser, flag, err := uc.repo.GetUser(login)
	if err != nil {
		uc.logger.Warn(err)
		return 0, false, err
	}
	if flag {
		return 0, true, nil
	}

	if !generator.CheckHashedPassword(gotUser.Password, password) {
		return 0, true, nil
	}

	return gotUser.Id, false, nil
}
