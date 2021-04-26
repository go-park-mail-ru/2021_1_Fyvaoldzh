package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"kudago/application/models"
	"kudago/application/subscription"
	"kudago/pkg/logger"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo"
)

type SubscriptionDatabase struct {
	pool   *pgxpool.Pool
	logger logger.Logger
}

func NewSubscriptionDatabase(conn *pgxpool.Pool, logger logger.Logger) subscription.Repository {
	return &SubscriptionDatabase{pool: conn, logger: logger}
}

func (sd SubscriptionDatabase) GetFollowers(id uint64) ([]uint64, error) {
	var users []uint64
	err := pgxscan.Select(context.Background(), sd.pool, &users, `SELECT subscriber_id
		FROM subscriptions WHERE subscribed_to_id = $1`, id)
	if errors.As(err, &sql.ErrNoRows) || len(users) == 0 {
		sd.logger.Debug("got no rows in method GetFollowers with id " + fmt.Sprint(id))
		return []uint64{}, nil
	}
	if err != nil {
		sd.logger.Warn(err)
		return []uint64{}, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return users, nil
}

func (sd SubscriptionDatabase) UpdateEventStatus(userId uint64, eventId uint64) error {
	resp, err := sd.pool.Exec(context.Background(),
		`UPDATE user_event SET is_planning = $1 WHERE user_id = $2 AND event_id = $3`,
		false, userId, eventId)
	if err != nil {
		sd.logger.Warn(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if resp.RowsAffected() == 0 {
		sd.logger.Warn(errors.New("event does not exist in list"))
		return echo.NewHTTPError(http.StatusBadRequest,
			"event does not exist in list")
	}

	return nil
}

func (sd SubscriptionDatabase) GetEventFollowers(eventId uint64) (models.UsersOnEvent, error) {
	var users models.UsersOnEvent
	err := pgxscan.Select(context.Background(), sd.pool, &users,
		`SELECT u.id, u.name, u.avatar
		FROM user_event
		JOIN users u ON u.id = user_id
		WHERE event_id = $1`, eventId)
	if errors.As(err, &sql.ErrNoRows) || len(users) == 0 {
		sd.logger.Debug("got no rows in method GetEventFollowers with id " + fmt.Sprint(eventId))
		return models.UsersOnEvent{}, nil
	}
	if err != nil {
		sd.logger.Warn(err)
		return models.UsersOnEvent{}, err
	}

	return users, nil
}

func (sd SubscriptionDatabase) IsAddedEvent(userId uint64, eventId uint64) (bool, error) {
	var id uint64
	err := sd.pool.
		QueryRow(context.Background(),
			`SELECT event_id FROM user_event WHERE event_id = $1 AND user_id = $2`,
			eventId, userId).Scan(&id)

	if errors.As(err, &sql.ErrNoRows) {
		sd.logger.Debug("no rows in method IsAddedEvent")
		return false, nil
	}
	if err != nil {
		sd.logger.Warn(err)
		return false, err
	}

	return true, nil
}