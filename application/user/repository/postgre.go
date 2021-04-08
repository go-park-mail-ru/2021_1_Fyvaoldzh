package repository

import (
	"context"
	"errors"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo"
	"kudago/application/models"
	"kudago/application/user"
	"net/http"
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

func (ud UserDatabase) GetPlanningEvents(id uint64) ([]models.EventCardWithDateSQL, error) {
	var events []models.EventCardWithDateSQL
	err := pgxscan.Select(context.Background(), ud.pool, &events,
		`SELECT ue.eid AS id, e.title, e.description, e.image, e.date  
		FROM user_event ue
		JOIN events e ON ue.eid = e.id
		WHERE ue.uid = $1 AND ue.is_p = $2`, id, true)
	if errors.As(err, &pgx.ErrNoRows) || len(events) == 0 {
		return []models.EventCardWithDateSQL{}, nil
	}

	if err != nil {
		return []models.EventCardWithDateSQL{}, err
	}

	return events, nil
}

func (ud UserDatabase) DeletePlanningEvent(userId uint64, eventId uint64) error {
	resp, err := ud.pool.Exec(context.Background(),
		`DELETE FROM user_event WHERE uid = $1 AND eid = $2 AND is_p = $3`,
		userId, eventId, true)
	if err != nil {
		return err
	}
	if resp.RowsAffected() == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "user or event does not exist")
	}

	return nil
}

func (ud UserDatabase) AddVisitedEvent(userId uint64, eventId uint64) error {
	_, err := ud.pool.Exec(context.Background(),
		`INSERT INTO user_event (uid, eid, is_p) VALUES ($1, $2, $3)`,
		userId, eventId, false)
	if err != nil {
		return err
	}

	return nil
}

func (ud UserDatabase) GetVisitedEvents(id uint64) ([]models.EventCardSQL, error) {
	var events []models.EventCardSQL
	err := pgxscan.Select(context.Background(), ud.pool, &events,
		`SELECT ue.eid AS id, e.title, e.description, e.image
		FROM user_event ue
		JOIN events e ON ue.eid = e.id
		WHERE ue.uid = $1 AND ue.is_p = $2`, id, false)
	if errors.As(err, &pgx.ErrNoRows) || len(events) == 0 {
		return []models.EventCardSQL{}, nil
	}
	if err != nil {
		return []models.EventCardSQL{}, err
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
