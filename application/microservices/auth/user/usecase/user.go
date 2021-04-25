package usecase

import (
	"errors"
	"github.com/labstack/echo"
	"kudago/application/microservices/auth/session"
	"kudago/application/microservices/auth/user"
	"kudago/pkg/generator"
	"kudago/pkg/logger"
	"net/http"
)

type UserUseCase struct {
	repo    user.Repository
	repoSession session.Repository
	logger  logger.Logger
}

func NewUserUseCase(u user.Repository, repoS session.Repository, logger logger.Logger) user.UseCase {
	return &UserUseCase{repo: u, repoSession: repoS, logger: logger}
}

func (uc UserUseCase) Login(login string, password string) (uint64, error) {
	return uc.CheckUser(login , password )
}

func (uc UserUseCase) CheckUser(login string, password string) (uint64, error) {
	gotUser, err := uc.repo.IsCorrect(login)
	if err != nil {
		uc.logger.Warn(err)
		return 0, err
	}

	if !generator.CheckHashedPassword(gotUser.Password, password) {
		uc.logger.Warn(errors.New("incorrect data"))
		return 0, echo.NewHTTPError(http.StatusBadRequest, "incorrect data")
	}

	return gotUser.Id, nil
}

func (uc UserUseCase) Logout(value string) error {
	err := uc.repoSession.DeleteSession(value)
	if err != nil {
		return err
	}

	return nil
}
