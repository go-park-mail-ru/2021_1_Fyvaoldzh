package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"kudago/application/microservices/subscription/subscription"
	"kudago/pkg/logger"
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

func (sd SubscriptionDatabase) CheckEvent(userId uint64, eventId uint64) (bool, error) {
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

