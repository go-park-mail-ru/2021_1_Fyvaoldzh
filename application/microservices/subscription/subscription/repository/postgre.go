package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/tarantool/go-tarantool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"kudago/application/microservices/subscription/subscription"
	"kudago/pkg/constants"
	"kudago/pkg/logger"
	"time"
)

type SubscriptionDatabase struct {
	pool   *pgxpool.Pool
	ttool  *tarantool.Connection
	logger logger.Logger
}

func NewSubscriptionDatabase(conn *pgxpool.Pool, ttool *tarantool.Connection, logger logger.Logger) subscription.Repository {
	return &SubscriptionDatabase{pool: conn, ttool: ttool, logger: logger}
}

func (sd SubscriptionDatabase) SubscribeUser(subscriberId uint64, subscribedToId uint64) error {
	_, err := sd.pool.Exec(context.Background(),
		`INSERT INTO subscriptions (subscriber_id, subscribed_to_id) VALUES ($1, $2)`,
		subscriberId, subscribedToId)
	if err != nil {
		sd.logger.Warn(err)
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

func (sd SubscriptionDatabase) UnsubscribeUser(subscriberId uint64, subscribedToId uint64) error {
	_, err := sd.pool.Exec(context.Background(),
		`DELETE FROM subscriptions WHERE subscriber_id = $1 AND subscribed_to_id = $2`,
		subscriberId, subscribedToId)
	if err != nil {
		sd.logger.Warn(err)
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

func (sd SubscriptionDatabase) CheckSubscription(subscriberId uint64, subscribedToId uint64) (bool, error) {
	var id1, id2 uint64
	err := sd.pool.
		QueryRow(context.Background(),
			`SELECT subscriber_id, subscribed_to_id 
			FROM subscriptions WHERE subscriber_id = $1 AND subscribed_to_id = $2`,
			subscriberId, subscribedToId).Scan(&id1, &id2)
	if errors.As(err, &sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		sd.logger.Warn(err)
		return false, status.Error(codes.Internal, err.Error())
	}

	return true, nil
}

func (sd SubscriptionDatabase) AddPlanning(userId uint64, eventId uint64) error {
	_, err := sd.pool.Exec(context.Background(),
		`INSERT INTO user_event (user_id, event_id, is_planning) VALUES ($1, $2, $3)`,
		userId, eventId, true)
	if err != nil {
		sd.logger.Warn(err)
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

func (sd SubscriptionDatabase) AddVisited(userId uint64, eventId uint64) error {
	_, err := sd.pool.Exec(context.Background(),
		`INSERT INTO user_event (user_id, event_id, is_planning) VALUES ($1, $2, $3)`,
		userId, eventId, false)
	if err != nil {
		sd.logger.Warn(err)
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

func (sd SubscriptionDatabase) RemoveEvent(userId uint64, eventId uint64) error {
	_, err := sd.pool.Exec(context.Background(),
		`DELETE FROM user_event WHERE user_id = $1 AND event_id = $2`,
		userId, eventId)
	if err != nil {
		sd.logger.Warn(err)
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

func (sd SubscriptionDatabase) CheckEventAdded(userId uint64, eventId uint64) (bool, error) {
	var id1, id2 uint64
	err := sd.pool.
		QueryRow(context.Background(),
			`SELECT user_id, event_id
			FROM user_event WHERE user_id = $1 AND event_id = $2`,
			userId, eventId).Scan(&id1, &id2)
	if errors.As(err, &sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		sd.logger.Warn(err)
		return false, status.Error(codes.Internal, err.Error())
	}

	return true, nil
}

func (sd SubscriptionDatabase) AddUserEventAction(userId uint64, eventId uint64) error {
	t := time.Now()
	_, err := sd.pool.Exec(context.Background(),
		`INSERT INTO actions_user_event (user_id, event_id, time) VALUES ($1, $2, $3)`,
		userId, eventId, t)
	if err != nil {
		sd.logger.Warn(err)
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

func (sd SubscriptionDatabase) RemoveUserEventAction(userId uint64, eventId uint64) error {
	_, err := sd.pool.Exec(context.Background(),
		`DELETE FROM actions_user_event WHERE user_id = $1 AND event_id = $2`,
		userId, eventId)
	if err != nil {
		sd.logger.Warn(err)
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

func (sd SubscriptionDatabase) AddSubscriptionAction(subscriberId uint64, subscribedToId uint64) error {
	t := time.Now()
	_, err := sd.pool.Exec(context.Background(),
		`INSERT INTO actions_subscription (subscriber_id, subscribed_to_id, time) VALUES ($1, $2, $3)`,
		subscriberId, subscribedToId, t)
	if err != nil {
		sd.logger.Warn(err)
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

func (sd SubscriptionDatabase) RemoveSubscriptionAction(userId uint64, eventId uint64) error {
	_, err := sd.pool.Exec(context.Background(),
		`DELETE FROM actions_subscription WHERE subscriber_id = $1 AND subscribed_to_id = $2`,
		userId, eventId)
	if err != nil {
		sd.logger.Warn(err)
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

func (sd SubscriptionDatabase) CheckEventInList(eventId uint64) (bool, error) {
	var id uint64
	err := sd.pool.
		QueryRow(context.Background(),
			`SELECT id
			FROM events WHERE id = $1`,
			eventId).Scan(&id)
	if errors.As(err, &sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		sd.logger.Warn(err)
		return false, status.Error(codes.Internal, err.Error())
	}

	return true, nil
}

func (sd SubscriptionDatabase) GetTimeEvent(eventId uint64) (time.Time, error) {
	var date sql.NullTime
	err := sd.pool.
		QueryRow(context.Background(),
			`SELECT start_date
			FROM events WHERE id = $1`,
			eventId).Scan(&date)
	if err != nil {
		sd.logger.Warn(err)
		return time.Time{}, status.Error(codes.Internal, err.Error())
	}

	return date.Time, nil
}

func (sd SubscriptionDatabase) AddPlanningNotification(eventId uint64, userId uint64, eventDate time.Time) error {
	_, err := sd.pool.Exec(context.Background(),
		`INSERT INTO notifications 
		VALUES ($1, $2, $3, $4, default)`,
		eventId, constants.EventNotif, userId, eventDate)
	if err != nil {
		sd.logger.Warn(err)
		return status.Error(codes.Internal, err.Error())
	}
	return nil
}

func (sd SubscriptionDatabase) RemovePlanningNotification(eventId uint64, userId uint64) error {
	_, err := sd.pool.Exec(context.Background(),
		`DELETE FROM notifications WHERE id = $1 AND id_to = $2 AND type = $3`,
		eventId, userId, constants.EventNotif)
	if err != nil {
		sd.logger.Warn(err)
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

func (sd SubscriptionDatabase) AddCountNotification(id uint64) error {
	_, err := sd.ttool.Update(constants.TarantoolSpaceName2, "primary",
		[]interface{}{id}, []interface{}{[]interface{}{"+", constants.TarantoolNotifications, 1}})
	if err != nil {
		sd.logger.Warn(err)
		return err
	}

	return nil
}
