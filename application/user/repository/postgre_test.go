package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"kudago/application/models"
	"kudago/application/user"
	"kudago/pkg/logger"
	"log"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgmock"
	"github.com/jackc/pgproto3/v2"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	userId   uint64 = 1
	str             = "str"
	login           = "login"
	password        = "123456"
	pageNum         = 1
)

var testRegData = &models.RegData{
	Login:    login,
	Password: password,
}

var testUserBack = &models.User{
	Id:       userId,
	Login:    login,
	Password: password,
}

var testUserData = &models.UserDataSQL{
	Id:    userId,
	Login: login,
}

func newDb(t *testing.T) user.Repository {
	pool := setUp(t)
	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	sugar := l.Sugar()
	zap.NewAtomicLevelAt(zapcore.DebugLevel)

	h := NewUserDatabase(pool, logger.NewLogger(sugar))
	return h
}
func setUp(t *testing.T) *pgxpool.Pool {
	script := &pgmock.Script{
		Steps: pgmock.AcceptUnauthenticatedConnRequestSteps(),
	}
	script.Steps = append(script.Steps, pgmock.ExpectMessage(&pgproto3.Query{String: ""}))
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
				Name:                 []byte("name"),
				TableOID:             0,
				TableAttributeNumber: 0,
				DataTypeOID:          23,
				DataTypeSize:         60,
				TypeModifier:         -1,
				Format:               0,
			},
			{
				Name:                 []byte("login"),
				TableOID:             0,
				TableAttributeNumber: 0,
				DataTypeOID:          23,
				DataTypeSize:         60,
				TypeModifier:         -1,
				Format:               0,
			},
			{
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

func TestUserDatabase_AddError(t *testing.T) {
	h := newDb(t)
	_, err := h.Add(testRegData)

	assert.Error(t, err)
}

func TestUserDatabase_ChangeAvatarError(t *testing.T) {
	h := newDb(t)
	err := h.ChangeAvatar(userId, str)

	assert.Error(t, err)
}

func TestUserDatabase_AddToPreferencesError(t *testing.T) {
	h := newDb(t)
	err := h.AddToPreferences(userId)

	assert.Error(t, err)
}

func TestUserDatabase_IsExistingNoRows(t *testing.T) {
	h := newDb(t)
	_, err := h.IsExisting(str)
	// returns err.NoRows
	// wtf?!

	assert.Nil(t, err)
}

func TestUserDatabase_IsCorrectError(t *testing.T) {
	h := newDb(t)
	_, err := h.IsCorrect(testUserBack)

	assert.Error(t, err)
}

func TestUserDatabase_UpdateError(t *testing.T) {
	h := newDb(t)
	err := h.Update(userId, testUserData)

	assert.Error(t, err)
}

func TestUserDatabase_GetByIdOwnError(t *testing.T) {
	pool := setUp(t)
	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	sugar := l.Sugar()
	zap.NewAtomicLevelAt(zapcore.DebugLevel)

	h := NewUserDatabase(pool, logger.NewLogger(sugar))
	_, err = h.GetByIdOwn(userId)

	assert.Error(t, err)
}

func TestUserDatabase_IsExistingEmail(t *testing.T) {
	h := newDb(t)
	_, err := h.IsExistingEmail(str)
	// тож false с nil

	assert.Nil(t, err)
}

func TestUserDatabase_IsExistingUserIdError(t *testing.T) {
	h := newDb(t)
	err := h.IsExistingUserId(userId)
	// тож false с nil

	assert.Error(t, err)
}

func TestUserDatabase_GetUsersError(t *testing.T) {
	h := newDb(t)
	_, err := h.GetUsers(pageNum)
	// ErrNoRows ???

	assert.Nil(t, err)
}

func TestUserDatabase_FindUsersError(t *testing.T) {
	h := newDb(t)
	_, err := h.FindUsers(str, pageNum)
	// ErrNoRows ???

	assert.Nil(t, err)
}

func TestUserDatabase_GetUserByIDError(t *testing.T) {
	h := newDb(t)
	_, err := h.GetUserByID(userId)

	assert.Error(t, err)
}

func TestUserDatabase_GetActionsError(t *testing.T) {
	h := newDb(t)
	_, err := h.GetActions(userId, pageNum)

	assert.Nil(t, err)
}
