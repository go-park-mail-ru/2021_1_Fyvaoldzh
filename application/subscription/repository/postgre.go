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

func (sd SubscriptionDatabase) CountUserFollowers(id uint64) (uint64, error) {
	var num uint64
	err := sd.pool.QueryRow(context.Background(), `SELECT COUNT(subscriber_id)
		FROM subscriptions WHERE subscribed_to_id = $1`, id).Scan(&num)
	if err != nil {
		sd.logger.Warn(err)
		return 0, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return num, nil
}

func (sd SubscriptionDatabase) CountUserSubscriptions(id uint64) (uint64, error) {
	var num uint64
	err := sd.pool.QueryRow(context.Background(), `SELECT COUNT(subscribed_to_id)
		FROM subscriptions WHERE subscriber_id = $1`, id).Scan(&num)
	if err != nil {
		sd.logger.Warn(err)
		return 0, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return num, nil
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

func (sd SubscriptionDatabase) GetFollowers(id uint64, page int) ([]models.UserCardSQL, error) {
	var users []models.UserCardSQL
	err := pgxscan.Select(context.Background(), sd.pool, &users,
		`SELECT u.id, u.name, u.avatar, u.birthday, u.city
		FROM users u
		JOIN subscriptions s ON subscribed_to_id = $1 AND u.id = s.subscriber_id
		LIMIT 10 OFFSET $2`, id, (page-1)*10)
	if errors.As(err, &sql.ErrNoRows) || len(users) == 0 {
		return []models.UserCardSQL{}, nil
	}
	if err != nil {
		sd.logger.Warn(err)
		return []models.UserCardSQL{}, err
	}

	return users, nil
}

func (sd SubscriptionDatabase) GetSubscriptions(id uint64, page int) ([]models.UserCardSQL, error) {
	var users []models.UserCardSQL
	err := pgxscan.Select(context.Background(), sd.pool, &users,
		`SELECT u.id, u.name, u.avatar, u.birthday, u.city
		FROM users u
		JOIN subscriptions s ON subscriber_id = $1 AND u.id = s.subscribed_to_id
		LIMIT 10 OFFSET $2`, id, (page-1)*10)
	if errors.As(err, &sql.ErrNoRows) || len(users) == 0 {
		return []models.UserCardSQL{}, nil
	}
	if err != nil {
		sd.logger.Warn(err)
		return []models.UserCardSQL{}, err
	}

	return users, nil
}

func (sd SubscriptionDatabase) GetPlanningEvents(id uint64, page int) ([]models.EventCardWithDateSQL, error) {
	var events []models.EventCardWithDateSQL
	err := pgxscan.Select(context.Background(), sd.pool, &events,
		`SELECT e.id, e.title, e.place, e.description, e.start_date, e.end_date
		FROM events e
		JOIN user_event ON user_id = $1
		WHERE event_id = e.id AND is_planning = $2
		LIMIT 10 OFFSET $3`, id, true, (page-1)*10)
	if errors.As(err, &sql.ErrNoRows) || len(events) == 0 {
		return []models.EventCardWithDateSQL{}, nil
	}
	if err != nil {
		sd.logger.Warn(err)
		return []models.EventCardWithDateSQL{}, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return events, nil
}

func (sd SubscriptionDatabase) GetVisitedEvents(id uint64, page int) ([]models.EventCardWithDateSQL, error) {
	var events []models.EventCardWithDateSQL
	err := pgxscan.Select(context.Background(), sd.pool, &events,
		`SELECT e.id, e.title, e.place, e.description, e.start_date, e.end_date  
		FROM events e
		JOIN user_event ON user_id = $1
		WHERE event_id = e.id AND is_planning = $2
		LIMIT 10 OFFSET $3`, id, false, (page-1)*10)
	if errors.As(err, &sql.ErrNoRows) || len(events) == 0 {
		return []models.EventCardWithDateSQL{}, nil
	}
	if err != nil {
		sd.logger.Warn(err)
		return []models.EventCardWithDateSQL{}, err
	}

	return events, nil
}
