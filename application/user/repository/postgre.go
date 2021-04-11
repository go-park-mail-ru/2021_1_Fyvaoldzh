package repository

import (
	"context"
	"database/sql"
	"errors"
	"kudago/application/models"
	"kudago/application/user"
	"net/http"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo"
)

type UserDatabase struct {
	pool *pgxpool.Pool
}

func NewUserDatabase(conn *pgxpool.Pool) user.Repository {
	return &UserDatabase{pool: conn}
}

func (ud UserDatabase) ChangeAvatar(uid uint64, path string) error {
	_, err := ud.pool.Exec(context.Background(),
		`UPDATE users SET avatar = $1 WHERE id = $2`, path, uid)
	if err != nil {
		return err
	}

	return nil
}

func (ud UserDatabase) Add(usr *models.RegData) (uint64, error) {
	var id uint64
	err := ud.pool.QueryRow(context.Background(),
		`INSERT INTO users (name, login, password) VALUES ($1, $2, $3) RETURNING id`,
		usr.Name, usr.Login, usr.Password).Scan(&id)
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return id, nil
}

func (ud UserDatabase) AddToPreferences(id uint64) error {
	_, err := ud.pool.Query(context.Background(),
		`INSERT INTO user_preference (user_id) VALUES ($1)`, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}


func (ud UserDatabase) IsExisting(login string) (bool, error) {
	var id uint64
	err := ud.pool.
		QueryRow(context.Background(),
			`SELECT id FROM users WHERE login = $1`, login).Scan(&id)

	if errors.As(err, &pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (ud UserDatabase) IsCorrect(user *models.User) (*models.User, error) {
	var gotUser models.User
	err := ud.pool.
		QueryRow(context.Background(),
			`SELECT id, password FROM users WHERE login = $1`,
			user.Login).Scan(&gotUser.Id, &gotUser.Password)
	if errors.As(err, &pgx.ErrNoRows) {
		return &gotUser, echo.NewHTTPError(http.StatusBadRequest, "incorrect data")
	}
	if err != nil {
		return &gotUser, err
	}

	return &gotUser, nil
}

func (ud UserDatabase) Update(id uint64, us *models.UserData) error {
	_, err := ud.pool.Exec(context.Background(),
		`UPDATE users SET "name" = $1, "email" = $2, "city" = $3, "about" = $4,
			"birthday" = $5, "password" = $6 WHERE id = $7`,
		us.Name, us.Email, us.City, us.About, us.Birthday, us.Password, id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (ud UserDatabase) GetByIdOwn(id uint64) (*models.UserData, error) {
	var usr []*models.UserData
	err := pgxscan.Select(context.Background(), ud.pool, &usr, `SELECT id, name, login, birthday, city, email, about, password, avatar 
		FROM users WHERE id = $1`, id)

	if len(usr) == 0 {
		return &models.UserData{}, echo.NewHTTPError(http.StatusBadRequest, "user does not exist")
	}
	if err != nil {
		return &models.UserData{}, err
	}

	return usr[0], nil
}

func (ud UserDatabase) IsExistingEmail(email string) (bool, error) {
	var id uint64
	err := ud.pool.
		QueryRow(context.Background(),
			`SELECT id FROM users WHERE email = $1`, email).Scan(&id)

	if errors.As(err, &pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (ud UserDatabase) IsExistingUserId(userId uint64) error {
	_, err := ud.pool.Query(context.Background(),
			`SELECT id FROM users WHERE id = $1`, userId)
	if err == sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("user does not exist"))
	}
	if err != nil {
		return err
	}

	return nil
}

func (ud UserDatabase) GetUsers(page int) (models.UsersOnEvent, error) {
	var users models.UsersOnEvent
	err := pgxscan.Select(context.Background(), ud.pool, &users,
		`SELECT id, name, avatar
		FROM users
		LIMIT 10 OFFSET $1`, (page-1)*10)
	if errors.As(err, &sql.ErrNoRows) || len(users) == 0 {
		return models.UsersOnEvent{}, nil
	}
	if err != nil {
		return models.UsersOnEvent{}, err
	}

	return users, nil
}



