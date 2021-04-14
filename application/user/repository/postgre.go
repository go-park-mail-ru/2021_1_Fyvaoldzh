package repository

import (
	"context"
	"database/sql"
	"errors"
	"kudago/application/models"
	"kudago/application/user"
	"kudago/pkg/logger"
	"net/http"

	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo"
)

type UserDatabase struct {
	conn   *pgx.Conn
	logger logger.Logger
}

func NewUserDatabase(conn *pgx.Conn, logger logger.Logger) user.Repository {
	return &UserDatabase{conn: conn, logger: logger}
}

func (ud UserDatabase) ChangeAvatar(uid uint64, path string) error {
	_, err := ud.conn.Exec(context.Background(),
		`UPDATE users SET avatar = $1 WHERE id = $2`, path, uid)
	if err != nil {
		ud.logger.Warn(err)
		return err
	}

	return nil
}

func (ud UserDatabase) Add(usr *models.RegData) (uint64, error) {
	var id uint64
	err := ud.conn.QueryRow(context.Background(),
		`INSERT INTO users (name, login, password) VALUES ($1, $2, $3) RETURNING id`,
		usr.Name, usr.Login, usr.Password).Scan(&id)
	if err != nil {
		ud.logger.Warn(err)
		return 0, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return id, nil
}

func (ud UserDatabase) AddToPreferences(id uint64) error {
	_, err := ud.conn.Query(context.Background(),
		`INSERT INTO user_preference (user_id) VALUES ($1)`, id)
	if err != nil {
		ud.logger.Warn(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}

func (ud UserDatabase) IsExisting(login string) (bool, error) {
	var id uint64
	err := ud.conn.
		QueryRow(context.Background(),
			`SELECT id FROM users WHERE login = $1`, login).Scan(&id)

	if err == sql.ErrNoRows {
		ud.logger.Debug("no rows in method IsExisting")
		return false, nil
	}
	if err != nil {
		ud.logger.Warn(err)
		return false, err
	}

	return true, nil
}

func (ud UserDatabase) IsCorrect(user *models.User) (*models.User, error) {
	var gotUser models.User
	err := ud.conn.
		QueryRow(context.Background(),
			`SELECT id, password FROM users WHERE login = $1`,
			user.Login).Scan(&gotUser.Id, &gotUser.Password)
	if err == sql.ErrNoRows {
		ud.logger.Debug("no rows in method IsCorrect")
		return &gotUser, echo.NewHTTPError(http.StatusBadRequest, "incorrect data")
	}
	if err != nil {
		ud.logger.Warn(err)
		return &gotUser, err
	}

	return &gotUser, nil
}

func (ud UserDatabase) Update(id uint64, us *models.UserData) error {
	_, err := ud.conn.Exec(context.Background(),
		`UPDATE users SET "name" = $1, "email" = $2, "city" = $3, "about" = $4,
			"birthday" = $5, "password" = $6 WHERE id = $7`,
		us.Name, us.Email, us.City, us.About, us.Birthday, us.Password, id,
	)
	if err != nil {
		ud.logger.Warn(err)
		return err
	}

	return nil
}

func (ud UserDatabase) GetByIdOwn(id uint64) (*models.UserData, error) {
	usr := &models.UserData{}
	err := ud.conn.QueryRow(context.Background(),`SELECT id, name, login, birthday, city, email, about, password, avatar 
		FROM users WHERE id = $1`, id).Scan(
		&usr.Id, &usr.Name, &usr.Login,
		&usr.Birthday, &usr.City, &usr.Email,
		&usr.About, &usr.Password,
		&usr.Avatar)

	if err == sql.ErrNoRows {
		return &models.UserData{}, echo.NewHTTPError(http.StatusBadRequest, "user does not exist")
	}
	if err != nil {
		ud.logger.Warn(err)
		return &models.UserData{}, err
	}

	return usr, nil
}

func (ud UserDatabase) IsExistingEmail(email string) (bool, error) {
	var id uint64
	err := ud.conn.
		QueryRow(context.Background(),
			`SELECT id FROM users WHERE email = $1`, email).Scan(&id)

	if err == sql.ErrNoRows {
		ud.logger.Debug("no rows in method IsExistingEmail")
		return false, nil
	}
	if err != nil {
		ud.logger.Warn(err)
		return false, err
	}

	return true, nil
}

func (ud UserDatabase) IsExistingUserId(userId uint64) error {
	_, err := ud.conn.Query(context.Background(),
		`SELECT id FROM users WHERE id = $1`, userId)
	if err == sql.ErrNoRows {
		ud.logger.Debug("no rows in method IsExistingUser")
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("user does not exist"))
	}
	if err != nil {
		ud.logger.Warn(err)
		return err
	}

	return nil
}

func (ud UserDatabase) GetUsers(page int) (models.UsersOnEvent, error) {
	var users models.UsersOnEvent
	rows, err := ud.conn.Query(context.Background(),`SELECT id, name, avatar
		FROM users
		LIMIT 10 OFFSET $1`, (page-1)*10)
	if err == sql.ErrNoRows {
		ud.logger.Debug("no rows in method GetUsers")
		return models.UsersOnEvent{}, nil
	}
	if err != nil {
		ud.logger.Warn(err)
		return models.UsersOnEvent{}, err
	}

	for rows.Next() {
		usr := models.UserOnEvent{}
		err = rows.Scan(&usr.Id, &usr.Name, &usr.Avatar)
		if err != nil {
			return nil, err
		}
		users = append(users, usr)
	}

	return users, nil
}
