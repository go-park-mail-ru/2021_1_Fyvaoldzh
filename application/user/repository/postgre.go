package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo"
	"kudago/application/models"
	"kudago/application/user"
	"net/http"
	"time"
)

type UserDatabase struct {
	pool *pgxpool.Pool
}

func NewUserDatabase(conn *pgxpool.Pool) user.Repository {
	return &UserDatabase{pool: conn}
}

func (ud UserDatabase) ChangeAvatar(uid uint64, path string) error {
	_, err := ud.pool.Exec(context.Background(),
		`UPDATE users SET "avatar" = $1 WHERE id = $2`, path, uid)

	if err != nil {
		return err
	}

	return nil
}

func (ud UserDatabase) GetPlanningEvents(id uint64) ([]uint64, error) {
	var events []uint64
	err := pgxscan.Select(context.Background(), ud.pool, &events, `SELECT eid
		FROM user_event WHERE uid = $1 AND is_p = $2`, id, true)

	if errors.As(err, &pgx.ErrNoRows) || len(events) == 0 {
		return []uint64{}, nil
	}

	if err != nil {
		return nil, err
	}

	// TODO: запихнуть горутины
	var newEvents []uint64
	for _, elem := range events {
		var date sql.NullTime
		err = ud.pool.QueryRow(context.Background(),
			`SELECT date FROM events WHERE id = $1`,
			elem).Scan(&date)
		if err != nil {
			return []uint64{}, err
		}
		if date.Valid && date.Time.Before(time.Now()) {
			// TODO: здесь надо красиво дергать другую репу, но увы
			// TODO: а еще обработать ошибки, офк

			_, _ = ud.pool.Exec(context.Background(),
				`DELETE FROM user_event WHERE uid = $1 AND eid = $2 AND is_p = $3`,
				id, elem, true)

			_, _ = ud.pool.Exec(context.Background(),
				`INSERT INTO user_event (uid, eid, is_p) VALUES ($1, $2, $3)`,
				id, elem, false)
		} else {
			newEvents = append(newEvents, elem)
		}
	}

	if len(newEvents) == 0 {
		return []uint64{}, nil
	}
	return newEvents, nil
}

func (ud UserDatabase) GetVisitedEvents(id uint64) ([]uint64, error) {
	var events []uint64
	err := pgxscan.Select(context.Background(), ud.pool, &events, `SELECT eid
		FROM user_event WHERE uid = $1 AND is_p = $2`, id, false)

	if errors.As(err, &pgx.ErrNoRows) || len(events) == 0 {
		return []uint64{}, nil
	}

	if err != nil {
		return nil, err
	}
	return events, nil
}

func (ud UserDatabase) GetFollowers(id uint64) ([]uint64, error) {
	var users []uint64
	err := pgxscan.Select(context.Background(), ud.pool, &users, `SELECT uid1
		FROM subscriptions WHERE uid2 = $1`, id)
	if errors.As(err, &pgx.ErrNoRows) {
		return []uint64{}, nil
	}

	if err != nil {
		return nil, err
	}
	return users, nil
}

func (ud UserDatabase) Add(user *models.RegData) (id uint64, err error) {
	err = ud.pool.QueryRow(context.Background(),
		`INSERT INTO users (name, login, password) VALUES ($1, $2, $3) RETURNING id`,
		user.Name, user.Login, user.Password).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
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

func (ud UserDatabase) IsCorrect(user *models.User) (uint64, error) {
	var id uint64
	err := ud.pool.
		QueryRow(context.Background(),
			`SELECT id FROM users WHERE login = $1 AND password = $2`,
			user.Login, user.Password).Scan(&id)
	if errors.As(err, &pgx.ErrNoRows) {
		return 0, echo.NewHTTPError(http.StatusBadRequest, "incorrect data")
	}

	if err != nil {
		return 0, err
	}
	return id, nil
}

func (ud UserDatabase) Update(id uint64, us *models.UserData) error {
	_, err := ud.pool.Exec(context.Background(),
		`UPDATE users SET "name" = $1, "email" = $2, "city" = $3, "about" = $4,`+
			`"avatar" = $5, "birthday" = $6, "password" = $7 WHERE id = $8`,
		us.Name, us.Email, us.City, us.About, us.Avatar, us.Birthday, us.Password, id,
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
		return &models.UserData{}, echo.NewHTTPError(http.StatusBadRequest, errors.New("user does not exist"))
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
