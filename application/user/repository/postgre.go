package repository

import (
	"context"
	"database/sql"
	"errors"
	"kudago/application/models"
	"kudago/application/user"
	"kudago/pkg/logger"
	"net/http"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo"
)

type UserDatabase struct {
	pool   *pgxpool.Pool
	logger logger.Logger
}

func NewUserDatabase(conn *pgxpool.Pool, logger logger.Logger) user.Repository {
	return &UserDatabase{pool: conn, logger: logger}
}

func (ud UserDatabase) ChangeAvatar(uid uint64, path string) error {
	_, err := ud.pool.Exec(context.Background(),
		`UPDATE users SET avatar = $1 WHERE id = $2`, path, uid)
	if err != nil {
		ud.logger.Warn(err)
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
		ud.logger.Warn(err)
		return 0, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return id, nil
}

func (ud UserDatabase) AddToPreferences(id uint64) error {
	_, err := ud.pool.Query(context.Background(),
		`INSERT INTO user_preference (user_id) VALUES ($1)`, id)
	if err != nil {
		ud.logger.Warn(err)
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
	err := ud.pool.
		QueryRow(context.Background(),
			`SELECT id, password FROM users WHERE login = $1`,
			user.Login).Scan(&gotUser.Id, &gotUser.Password)
	if errors.As(err, &pgx.ErrNoRows) {
		ud.logger.Debug("no rows in method GetUser")
		return &gotUser, echo.NewHTTPError(http.StatusBadRequest, "incorrect data")
	}
	if err != nil {
		ud.logger.Warn(err)
		return &gotUser, err
	}

	return &gotUser, nil
}

func (ud UserDatabase) Update(id uint64, us *models.UserDataSQL) error {
	_, err := ud.pool.Exec(context.Background(),
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

func (ud UserDatabase) GetByIdOwn(id uint64) (*models.UserDataSQL, error) {
	var usr []*models.UserDataSQL
	err := pgxscan.Select(context.Background(), ud.pool, &usr, `SELECT id, name, login, birthday, city, email, about, password, avatar 
		FROM users WHERE id = $1`, id)

	if len(usr) == 0 {
		return &models.UserDataSQL{}, echo.NewHTTPError(http.StatusBadRequest, "user does not exist")
	}
	if err != nil {
		ud.logger.Warn(err)
		return &models.UserDataSQL{}, err
	}

	return usr[0], nil
}

func (ud UserDatabase) IsExistingEmail(email string) (bool, error) {
	var id uint64
	err := ud.pool.
		QueryRow(context.Background(),
			`SELECT id FROM users WHERE email = $1`, email).Scan(&id)

	if errors.As(err, &pgx.ErrNoRows) {
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
	_, err := ud.pool.Query(context.Background(),
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

func (ud UserDatabase) GetUsers(page int) ([]models.UserCardSQL, error) {
	var users []models.UserCardSQL
	err := pgxscan.Select(context.Background(), ud.pool, &users,
		`SELECT id, name, avatar, birthday, city
		FROM users
		LIMIT 10 OFFSET $1`, (page-1)*10)
	if errors.As(err, &sql.ErrNoRows) || len(users) == 0 {
		ud.logger.Debug("no rows in method GetUsers")
		return []models.UserCardSQL{}, nil
	}
	if err != nil {
		ud.logger.Warn(err)
		return []models.UserCardSQL{}, err
	}

	return users, nil
}

func (ud UserDatabase) FindUsers(str string, page int) ([]models.UserCardSQL, error) {
	var users []models.UserCardSQL
	err := pgxscan.Select(context.Background(), ud.pool, &users,
		`SELECT DISTINCT ON(id) id, name, avatar, birthday, city
		FROM users
		WHERE (LOWER(name) LIKE '%' || $1 || '%' OR LOWER(about) LIKE '%' || $1 || '%')
		LIMIT 10 OFFSET $2`, str, (page-1)*10)
	if errors.As(err, &sql.ErrNoRows) || len(users) == 0 {
		return []models.UserCardSQL{}, nil
	}
	if err != nil {
		ud.logger.Warn(err)
		return []models.UserCardSQL{}, err
	}

	return users, nil
}

func (ud UserDatabase) GetUserByID(id uint64) (models.UserOnEvent, error) {
	var users models.UsersOnEvent
	err := pgxscan.Select(context.Background(), ud.pool, &users,
		`SELECT id, name, avatar
		FROM users
		WHERE id = $1`, id)

	if errors.As(err, &sql.ErrNoRows) {
		err := errors.New("no user with this id")
		ud.logger.Warn(err)
		return models.UserOnEvent{}, err
	}
	if err != nil {
		ud.logger.Warn(err)
		return models.UserOnEvent{}, err
	}

	return users[0], nil
}

func (ud UserDatabase) GetActions(id uint64, page int) ([]*models.ActionCard, error) {
	var actions []*models.ActionCard
	err := pgxscan.Select(context.Background(), ud.pool, &actions,
		`SELECT user_id as Id1, u.name as Name1, event_id as Id2, e.title as Name2, time as Time, 'user_event' as Type
		FROM actions_user_event
		JOIN users u on actions_user_event.user_id = u.id
		JOIN events e on actions_user_event.event_id = e.id
		JOIN subscriptions s2 on subscriber_id = $1 AND user_id = s2.subscribed_to_id
		UNION ALL
		select a_s.subscriber_id, u1.name, a_s.subscribed_to_id, u2.name, a_s.time, 'subscription'
		FROM actions_subscription a_s
		JOIN subscriptions s on s.subscriber_id = $1 AND s.subscribed_to_id = a_s.subscriber_id AND a_s.subscribed_to_id <> $1
		JOIN users u1 on a_s.subscriber_id = u1.id
		JOIN users u2 on a_s.subscribed_to_id = u2.id
		UNION ALL
		SELECT ss.subscriber_id, u3.name, ss.subscribed_to_id, '', ss.time, 'new_follower'
		FROM actions_subscription ss
		JOIN users u3 on ss.subscriber_id = u3.id
		WHERE subscribed_to_id = $1
		ORDER BY Time DESC
		LIMIT 10 OFFSET $2`, id, (page-1)*10)
	if errors.As(err, &sql.ErrNoRows) {
		return []*models.ActionCard{}, nil
	}
	if err != nil {
		ud.logger.Warn(err)
		return []*models.ActionCard{}, err
	}

	return actions, nil
}
