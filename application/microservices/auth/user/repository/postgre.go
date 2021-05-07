package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"kudago/application/microservices/auth/user"
	"kudago/application/models"
	"kudago/pkg/logger"
)

type UserDatabase struct {
	pool   *pgxpool.Pool
	logger logger.Logger
}

func NewUserDatabase(conn *pgxpool.Pool, logger logger.Logger) user.Repository {
	return &UserDatabase{pool: conn, logger: logger}
}

func (ud UserDatabase) GetUser(login string) (*models.User, bool, error) {
	var gotUser models.User
	err := ud.pool.
		QueryRow(context.Background(),
			`SELECT id, password FROM users WHERE login = $1`,
			login).Scan(&gotUser.Id, &gotUser.Password)
	if errors.As(err, &pgx.ErrNoRows) {
		return &gotUser, true, nil
	}
	if err != nil {
		ud.logger.Warn(err)
		return &gotUser, false, status.Error(codes.Internal, err.Error())
	}

	return &gotUser, false, nil
}
