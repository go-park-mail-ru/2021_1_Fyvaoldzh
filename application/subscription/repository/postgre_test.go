package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgmock"
	"github.com/jackc/pgproto3/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"kudago/application/subscription"
	"kudago/pkg/logger"
	"log"
	"net"
	"strings"
	"testing"
	"time"
)

var (
	userId  uint64 = 1
	eventId uint64 = 1
	pageNum        = 1
)

func newDb(t *testing.T) subscription.Repository {
	pool := setUp(t)
	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	sugar := l.Sugar()
	zap.NewAtomicLevelAt(zapcore.DebugLevel)

	h := NewSubscriptionDatabase(pool, logger.NewLogger(sugar))
	return h
}
func setUp(t *testing.T) *pgxpool.Pool {
	script := &pgmock.Script{
		Steps: pgmock.AcceptUnauthenticatedConnRequestSteps(),
	}
	script.Steps = append(script.Steps, pgmock.ExpectMessage(&pgproto3.Query{String: ""}))
	script.Steps = append(script.Steps, pgmock.SendMessage(&pgproto3.RowDescription{
		Fields: []pgproto3.FieldDescription{
			pgproto3.FieldDescription{
				Name:                 []byte("id"),
				TableOID:             0,
				TableAttributeNumber: 0,
				DataTypeOID:          23,
				DataTypeSize:         8,
				TypeModifier:         -1,
				Format:               0,
			},
			pgproto3.FieldDescription{
				Name:                 []byte("name"),
				TableOID:             0,
				TableAttributeNumber: 0,
				DataTypeOID:          23,
				DataTypeSize:         60,
				TypeModifier:         -1,
				Format:               0,
			},
			pgproto3.FieldDescription{
				Name:                 []byte("login"),
				TableOID:             0,
				TableAttributeNumber: 0,
				DataTypeOID:          23,
				DataTypeSize:         60,
				TypeModifier:         -1,
				Format:               0,
			},
			pgproto3.FieldDescription{
				Name:                 []byte("password"),
				TableOID:             0,
				TableAttributeNumber: 0,
				DataTypeOID:          23,
				DataTypeSize:         60,
				TypeModifier:         -1,
				Format:               0,
			},
		},
	}))
	script.Steps = append(script.Steps, pgmock.SendMessage(&pgproto3.DataRow{
		Values: [][]byte{[]byte("1")},
	}))
	script.Steps = append(script.Steps, pgmock.SendMessage(&pgproto3.CommandComplete{CommandTag: []byte("SELECT 1")}))
	script.Steps = append(script.Steps, pgmock.SendMessage(&pgproto3.ReadyForQuery{TxStatus: 'I'}))
	script.Steps = append(script.Steps, pgmock.ExpectMessage(&pgproto3.Terminate{}))

	ln, err := net.Listen("tcp", "127.0.0.1:")
	require.NoError(t, err)
	defer ln.Close()

	serverErrChan := make(chan error, 1)
	go func() {
		defer close(serverErrChan)

		conn, err := ln.Accept()
		if err != nil {
			serverErrChan <- err
			return
		}
		defer conn.Close()

		err = conn.SetDeadline(time.Now().Add(time.Second))
		if err != nil {
			serverErrChan <- err
			return
		}

		err = script.Run(pgproto3.NewBackend(pgproto3.NewChunkReader(conn), conn))
		if err != nil {
			serverErrChan <- err
			return
		}
	}()

	parts := strings.Split(ln.Addr().String(), ":")
	host := parts[0]
	port := parts[1]
	connStr := fmt.Sprintf("sslmode=disable host=%s port=%s", host, port)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	pool, err := pgxpool.Connect(ctx, connStr)
	require.NoError(t, err)
	return pool
}

func TestSubscriptionDatabase_CountUserFollowersError(t *testing.T) {
	h := newDb(t)
	_, err := h.CountUserFollowers(userId)

	assert.Error(t, err)
}

func TestSubscriptionDatabase_CountUserSubscriptionsError(t *testing.T) {
	h := newDb(t)
	_, err := h.CountUserSubscriptions(userId)

	assert.Error(t, err)
}

func TestSubscriptionDatabase_UpdateEventStatusError(t *testing.T) {
	h := newDb(t)
	err := h.UpdateEventStatus(userId, eventId)

	assert.Error(t, err)
}

func TestSubscriptionDatabase_GetEventFollowers(t *testing.T) {
	h := newDb(t)
	_, err := h.GetEventFollowers(eventId)

	assert.Nil(t, err)
}

func TestSubscriptionDatabase_IsAddedEvent(t *testing.T) {
	h := newDb(t)
	_, err := h.IsAddedEvent(userId, eventId)

	assert.Nil(t, err)
}

func TestSubscriptionDatabase_GetFollowers(t *testing.T) {
	h := newDb(t)
	_, err := h.GetFollowers(userId, pageNum)

	assert.Nil(t, err)
}

func TestSubscriptionDatabase_GetSubscriptions(t *testing.T) {
	h := newDb(t)
	_, err := h.GetSubscriptions(userId, pageNum)

	assert.Nil(t, err)
}

func TestSubscriptionDatabase_GetPlanningEvents(t *testing.T) {
	h := newDb(t)
	_, err := h.GetPlanningEvents(userId, pageNum)

	assert.Nil(t, err)
}

func TestSubscriptionDatabase_GetVisitedEvents(t *testing.T) {
	h := newDb(t)
	_, err := h.GetVisitedEvents(userId, pageNum)

	assert.Nil(t, err)
}
