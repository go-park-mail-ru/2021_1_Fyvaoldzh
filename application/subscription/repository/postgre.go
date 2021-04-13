package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"kudago/application/models"
	"kudago/application/subscription"
	"kudago/pkg/logger"
	"net/http"

	"github.com/georgysavva/scany/pgxscan"
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

func (sd SubscriptionDatabase) SubscribeUser(subscriberId uint64, subscribedToId uint64) error {
	_, err := sd.pool.Exec(context.Background(),
		`INSERT INTO subscriptions (subscriber_id, subscribed_to_id) VALUES ($1, $2)`,
		subscriberId, subscribedToId)
	if err != nil {
		sd.logger.Warn(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}

func (sd SubscriptionDatabase) UnsubscribeUser(subscriberId uint64, subscribedToId uint64) error {
	resp, err := sd.pool.Exec(context.Background(),
		`DELETE FROM subscriptions WHERE subscriber_id = $1 AND subscribed_to_id = $2`,
		subscriberId, subscribedToId)
	if err != nil {
		sd.logger.Warn(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if resp.RowsAffected() == 0 {
		sd.logger.Warn(errors.New("subscription does not exist"))
		return echo.NewHTTPError(http.StatusBadRequest,
			"subscription does not exist")
	}

	return nil
}

func (sd SubscriptionDatabase) AddPlanning(userId uint64, eventId uint64) error {
	_, err := sd.pool.Exec(context.Background(),
		`INSERT INTO user_event (user_id, event_id, is_planning) VALUES ($1, $2, $3)`,
		userId, eventId, true)
	if err != nil {
		sd.logger.Warn(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}

func (sd SubscriptionDatabase) RemovePlanning(userId uint64, eventId uint64) error {
	resp, err := sd.pool.Exec(context.Background(),
		`DELETE FROM user_event WHERE user_id = $1 AND event_id = $2 AND is_planning = $3`,
		userId, eventId, true)
	if err != nil {
		sd.logger.Warn(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if resp.RowsAffected() == 0 {
		sd.logger.Warn(errors.New("subscription with id " + fmt.Sprint(userId) + "in profile with id" + fmt.Sprint(eventId) + "does not exist in planing"))
		return echo.NewHTTPError(http.StatusBadRequest,
			"event does not exist in list")
	}

	return nil
}

func (sd SubscriptionDatabase) AddVisited(userId uint64, eventId uint64) error {
	_, err := sd.pool.Exec(context.Background(),
		`INSERT INTO user_event (user_id, event_id, is_planning) VALUES ($1, $2, $3)`,
		userId, eventId, false)
	if err != nil {
		sd.logger.Warn(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}

func (sd SubscriptionDatabase) RemoveVisited(uid uint64, eid uint64) error {
	resp, err := sd.pool.Exec(context.Background(),
		`DELETE FROM user_event WHERE user_id = $1 AND event_id = $2 AND is_planning = $3`,
		uid, eid, false)
	if err != nil {
		sd.logger.Warn(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if resp.RowsAffected() == 0 {
		sd.logger.Warn(errors.New("subscription with id " + fmt.Sprint(uid) + "in profile with id" + fmt.Sprint(eid) + "does not exist in visited"))
		return echo.NewHTTPError(http.StatusBadRequest,
			"event does not exist in list")
	}

	return nil
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

func (sd SubscriptionDatabase) GetPlanningEvents(id uint64) ([]models.EventCardWithDateSQL, error) {
	var events []models.EventCardWithDateSQL
	err := pgxscan.Select(context.Background(), sd.pool, &events,
		`SELECT e.id, e.title, e.description, e.image, e.start_date, e.end_date
		FROM events e
		JOIN user_event ON user_id = $1
		WHERE event_id = e.id AND is_planning = $2`, id, true)
	if errors.As(err, &sql.ErrNoRows) || len(events) == 0 {
		sd.logger.Debug("got no rows in method GetPlanningEvents with id " + fmt.Sprint(id))
		return []models.EventCardWithDateSQL{}, nil
	}
	if err != nil {
		sd.logger.Warn(err)
		return []models.EventCardWithDateSQL{}, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return events, nil
}

func (sd SubscriptionDatabase) GetVisitedEvents(id uint64) ([]models.EventCardWithDateSQL, error) {
	var events []models.EventCardWithDateSQL
	err := pgxscan.Select(context.Background(), sd.pool, &events,
		`SELECT e.id, e.title, e.description, e.image, e.start_date, e.end_date  
		FROM events e
		JOIN user_event ON user_id = $1
		WHERE event_id = e.id AND is_planning = $2`, id, false)
	if errors.As(err, &sql.ErrNoRows) || len(events) == 0 {
		sd.logger.Debug("got no rows in method GetVisitedEvents with id " + fmt.Sprint(id))
		return []models.EventCardWithDateSQL{}, nil
	}
	if err != nil {
		sd.logger.Warn(err)
		return []models.EventCardWithDateSQL{}, err
	}

	return events, nil
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
