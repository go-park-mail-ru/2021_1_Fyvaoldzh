package repository

import (
	"context"
	"fmt"
	"kudago/pkg/logger"
	"log"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/jackc/pgmock"
	"github.com/jackc/pgproto3/v2"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setUp(t *testing.T) *pgxpool.Pool {
	script := &pgmock.Script{
		Steps: pgmock.AcceptUnauthenticatedConnRequestSteps(),
	}
	script.Steps = append(script.Steps, pgmock.ExpectMessage(&pgproto3.Query{String: "SELECT id, title, description, image, start_date, end_date FROM events WHERE end_date > $1 ORDER BY id DESC LIMIT 6 OFFSET $2"}))
	script.Steps = append(script.Steps, pgmock.SendMessage(&pgproto3.RowDescription{
		Fields: []pgproto3.FieldDescription{
			{
				Name:                 []byte("id"),
				TableOID:             0,
				TableAttributeNumber: 0,
				DataTypeOID:          23,
				DataTypeSize:         8,
				TypeModifier:         -1,
				Format:               0,
			},
			{
				Name:                 []byte("title"),
				TableOID:             0,
				TableAttributeNumber: 0,
				DataTypeOID:          23,
				DataTypeSize:         60,
				TypeModifier:         -1,
				Format:               0,
			},
			{
				Name:                 []byte("description"),
				TableOID:             0,
				TableAttributeNumber: 0,
				DataTypeOID:          23,
				DataTypeSize:         60,
				TypeModifier:         -1,
				Format:               0,
			},
			{
				Name:                 []byte("image"),
				TableOID:             0,
				TableAttributeNumber: 0,
				DataTypeOID:          23,
				DataTypeSize:         60,
				TypeModifier:         -1,
				Format:               0,
			},
			{
				Name:                 []byte("start_date"),
				TableOID:             0,
				TableAttributeNumber: 0,
				DataTypeOID:          23,
				DataTypeSize:         60,
				TypeModifier:         -1,
				Format:               0,
			},
			{
				Name:                 []byte("end_date"),
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

func TestUserDatabase_GetEventsOk(t *testing.T) {
	pool := setUp(t)
	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	sugar := l.Sugar()
	zap.NewAtomicLevelAt(zapcore.DebugLevel)

	h := NewEventDatabase(pool, logger.NewLogger(sugar))
	_, err = h.GetAllEvents(time.Now(), 1)

	assert.Nil(t, err)
}

func TestUserDatabase_GetEventsByCategoryOk(t *testing.T) {
	pool := setUp(t)
	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	sugar := l.Sugar()
	zap.NewAtomicLevelAt(zapcore.DebugLevel)

	h := NewEventDatabase(pool, logger.NewLogger(sugar))
	_, err = h.GetEventsByCategory("", time.Now(), 1)

	assert.Nil(t, err)
}

func TestUserDatabase_GetOneEventByIDError(t *testing.T) {
	pool := setUp(t)
	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	sugar := l.Sugar()
	zap.NewAtomicLevelAt(zapcore.DebugLevel)

	h := NewEventDatabase(pool, logger.NewLogger(sugar))
	_, err = h.GetOneEventByID(0)

	assert.NotNil(t, err)
}

func TestUserDatabase_GetTagsOk(t *testing.T) {
	pool := setUp(t)
	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	sugar := l.Sugar()
	zap.NewAtomicLevelAt(zapcore.DebugLevel)

	h := NewEventDatabase(pool, logger.NewLogger(sugar))
	_, err = h.GetTags(1)

	assert.Nil(t, err)
}

func TestUserDatabase_DeleteByIdOk(t *testing.T) {
	pool := setUp(t)
	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	sugar := l.Sugar()
	zap.NewAtomicLevelAt(zapcore.DebugLevel)

	h := NewEventDatabase(pool, logger.NewLogger(sugar))
	err = h.DeleteById(1)

	assert.NotNil(t, err)
}

func TestUserDatabase_CategorySearchOk(t *testing.T) {
	pool := setUp(t)
	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	sugar := l.Sugar()
	zap.NewAtomicLevelAt(zapcore.DebugLevel)

	h := NewEventDatabase(pool, logger.NewLogger(sugar))
	_, err = h.CategorySearch("", "", time.Now(), 1)

	assert.Nil(t, err)
}

func TestUserDatabase_GetRecommendOk(t *testing.T) {
	pool := setUp(t)
	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	sugar := l.Sugar()
	zap.NewAtomicLevelAt(zapcore.DebugLevel)

	h := NewEventDatabase(pool, logger.NewLogger(sugar))
	_, err = h.GetRecommended(1, time.Now(), 1)

	assert.Nil(t, err)
}

func TestUserDatabase_RecommendSystemOk(t *testing.T) {
	pool := setUp(t)
	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	sugar := l.Sugar()
	zap.NewAtomicLevelAt(zapcore.DebugLevel)

	h := NewEventDatabase(pool, logger.NewLogger(sugar))
	err = h.RecomendSystem(1, "")

	assert.NotNil(t, err)
}

func TestUserDatabase_FindEventsOk(t *testing.T) {
	pool := setUp(t)
	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	sugar := l.Sugar()
	zap.NewAtomicLevelAt(zapcore.DebugLevel)

	h := NewEventDatabase(pool, logger.NewLogger(sugar))
	_, err = h.FindEvents("", time.Now(), 1)

	assert.Nil(t, err)
}

func TestUserDatabase_UpdateEventAvatarError(t *testing.T) {
	pool := setUp(t)
	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	sugar := l.Sugar()
	zap.NewAtomicLevelAt(zapcore.DebugLevel)

	h := NewEventDatabase(pool, logger.NewLogger(sugar))
	err = h.UpdateEventAvatar(1, "")

	assert.NotNil(t, err)
}
