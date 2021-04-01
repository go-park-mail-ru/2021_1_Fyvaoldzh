package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo"
	"kudago/application/subscription"
	"net/http"
)

type SubscriptionDatabase struct {
	pool *pgxpool.Pool
}

func NewSubscriptionDatabase(conn *pgxpool.Pool) subscription.Repository {
	return &SubscriptionDatabase{pool: conn}
}

func (ud SubscriptionDatabase) SubscribeUser(uid1 uint64, uid2 uint64) error {
	var ui1, ui2 uint64
	// TODO: херня запрос, исправить по аналогии с тем, что ниже
	err := ud.pool.
		QueryRow(context.Background(),
			`SELECT uid1, uid2 FROM subscriptions WHERE uid1 = $1 AND uid2 = $2`,
			uid1, uid2).Scan(&ui1, &ui2)
	if !errors.As(err, &pgx.ErrNoRows) || err == nil {
		if err != nil {
			return err
		}
		return echo.NewHTTPError(http.StatusBadRequest, "the subscription exists")
	}

	err = ud.pool.
		QueryRow(context.Background(),
			`SELECT id FROM users WHERE id = $1`,
			uid2).Scan(&ui2)
	if errors.As(err, &pgx.ErrNoRows) {
		return echo.NewHTTPError(http.StatusBadRequest, "user with this id does not exist")
	}

	if err != nil {
		return err
	}

	_, err = ud.pool.Exec(context.Background(),
		`INSERT INTO subscriptions (uid1, uid2) VALUES ($1, $2)`,
		uid1, uid2)

	if err != nil {
		return err
	}
	return nil
}

func (ud SubscriptionDatabase) UnsubscribeUser(uid1 uint64, uid2 uint64) error {
	resp, err := ud.pool.Exec(context.Background(),
		`DELETE FROM subscriptions WHERE uid1 = $1 AND uid2 = $2`,
		uid1, uid2)

	if err != nil {
		return err
	}

	if resp.RowsAffected() == 0 {
		return echo.NewHTTPError(http.StatusBadRequest,
			"subscription does not exist")
	}

	return nil
}

func (ud SubscriptionDatabase) AddPlanning(uid uint64, eid uint64) error {
	resp, err := ud.pool.Exec(context.Background(),
		`SELECT * FROM user_event WHERE uid = $1 AND eid = $2 AND is_p = $3`,
		uid, eid, true)

	if err != nil {
		return err
	}

	if resp.RowsAffected() > 0 {
		return echo.NewHTTPError(http.StatusBadRequest,
			"event is already added")
	}

	resp, err = ud.pool.Exec(context.Background(),
		`SELECT id FROM events WHERE id = $1`,
		eid)

	if err != nil {
		return err
	}

	if resp.RowsAffected() == 0 {
		return echo.NewHTTPError(http.StatusBadRequest,
			"event does not exist")
	}

	_, err = ud.pool.Exec(context.Background(),
		`INSERT INTO user_event (uid, eid, is_p) VALUES ($1, $2, $3)`,
		uid, eid, true)

	if err != nil {
		return err
	}

	return nil
}

func (ud SubscriptionDatabase) RemovePlanning(uid uint64, eid uint64) error {
	resp, err := ud.pool.Exec(context.Background(),
		`DELETE FROM user_event WHERE uid = $1 AND eid = $2 AND is_p = $3`,
		uid, eid, true)

	if err != nil {
		return err
	}

	if resp.RowsAffected() == 0 {
		return echo.NewHTTPError(http.StatusBadRequest,
			"do not have this event in list")
	}

	return nil
}

func (ud SubscriptionDatabase) AddVisited(uid uint64, eid uint64) error {
	resp, err := ud.pool.Exec(context.Background(),
		`SELECT FROM user_event WHERE uid = $1 AND eid = $2 AND is_p = $3`,
		uid, eid, false)

	if err != nil {
		return err
	}

	if resp.RowsAffected() > 0 {
		return echo.NewHTTPError(http.StatusBadRequest,
			"event is already added")
	}

	resp, err = ud.pool.Exec(context.Background(),
		`SELECT id FROM events WHERE id = $1`,
		eid)

	if err != nil {
		return err
	}

	if resp.RowsAffected() == 0 {
		return echo.NewHTTPError(http.StatusBadRequest,
			"event does not exist")
	}

	_, err = ud.pool.Exec(context.Background(),
		`INSERT INTO user_event (uid, eid, is_p) VALUES ($1, $2, $3)`,
		uid, eid, false)

	if err != nil {
		return err
	}

	return nil
}

func (ud SubscriptionDatabase) RemoveVisited(uid uint64, eid uint64) error {
	resp, err := ud.pool.Exec(context.Background(),
		`DELETE FROM user_event WHERE uid = $1 AND eid = $2 AND is_p = $3`,
		uid, eid, false)

	if err != nil {
		return err
	}

	if resp.RowsAffected() == 0 {
		return echo.NewHTTPError(http.StatusBadRequest,
			"do not have this event in list")
	}

	return nil
}

