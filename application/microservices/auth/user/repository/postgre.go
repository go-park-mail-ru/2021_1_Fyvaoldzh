package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo"
	"kudago/application/microservices/auth/user"
	"kudago/application/models"
	"kudago/pkg/logger"
	"net/http"
)

type UserDatabase struct {
	pool   *pgxpool.Pool
	logger logger.Logger
}

func NewUserDatabase(conn *pgxpool.Pool, logger logger.Logger) user.Repository {
	return &UserDatabase{pool: conn, logger: logger}
}

func (ud UserDatabase) IsCorrect(login string) (*models.User, error) {
	var gotUser models.User
	err := ud.pool.
		QueryRow(context.Background(),
			`SELECT id, password FROM users WHERE login = $1`,
			login).Scan(&gotUser.Id, &gotUser.Password)
	if errors.As(err, &pgx.ErrNoRows) {
		ud.logger.Debug("no rows in method IsCorrect")
		return &gotUser, echo.NewHTTPError(http.StatusBadRequest, "incorrect data")
	}
	if err != nil {
		ud.logger.Warn(err)
		return &gotUser, err
	}

	return &gotUser, nil
}