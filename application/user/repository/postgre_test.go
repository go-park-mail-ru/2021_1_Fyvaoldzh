package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"kudago/application/models"
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

/*
import (
	"context"
	"database/sql"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"kudago/application/models"
	"kudago/pkg/constants"
	"kudago/pkg/logger"
	mock_pool "kudago/pkg/pool/mocks"
	"log"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/jackc/pgx/v4/stdlib"

	"testing"
)

var (
	userId          uint64 = 1
	pageNum         = 1
	login           = "userlogin"
	name            = "username"
	frontPassword   = "123456"
	backPassword    = "IvJrQEdIeoTzLsMX_839spM7MzaXS7aJ_b3xTzmYqbotq3HRKAs="
	badBackPassword = "1111IvJrQEdIeoTzLsMX_839spM7MzaXS7aJ_b3xTzmYqbotq3HRKAs="
	email           = "email@mail.ru"
	birthdayStr     = "1999-01-01"
	birthday, err   = time.Parse(constants.TimeFormat, "1999-01-01")
	city            = "City"
	about           = "some personal information"
	avatar          = "public/users/default.png"
	imageName       = "image.png"
	evPlanningSQL   = models.EventCardWithDateSQL{
		ID:        1,
		StartDate: time.Now(),
		EndDate:   time.Now().Add(10 * time.Hour),
	}
	evPlanning = models.EventCard{
		ID:        1,
		StartDate: evPlanningSQL.StartDate.String(),
		EndDate:   evPlanningSQL.EndDate.String(),
	}
	evVisitedSQL = models.EventCardWithDateSQL{
		ID:        2,
		StartDate: time.Now(),
		EndDate:   time.Now(),
	}
	evVisited = models.EventCard{
		ID:        2,
		StartDate: evVisitedSQL.StartDate.String(),
		EndDate:   evVisitedSQL.EndDate.String(),
	}
	eventsPlanningSQL = []models.EventCardWithDateSQL{
		evPlanningSQL, evVisitedSQL,
	}
	eventsVisitedSQL = []models.EventCardWithDateSQL{
		evVisitedSQL,
	}
	eventsPlanning = []models.EventCard{
		evPlanning,
	}
	eventsVisited = []models.EventCard{
		evVisited,
	}

	followers = []uint64{2, 2, 3}
)

var testUserFront = &models.User{
	Login:    login,
	Password: frontPassword,
}

var testUserBack = &models.User{
	Id:       userId,
	Login:    login,
	Password: backPassword,
}


func setUp(t *testing.T) (*pgxpool.Pool, logger.Logger) {
	ctrl := gomock.NewController(t)
	bd := mock_pool.NewMockCustomPool(ctrl)

	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	sugar := l.Sugar()
	zap.NewAtomicLevelAt(zapcore.DebugLevel)

	rep := NewUserDatabase(bd, logger.NewLogger(sugar))

	return newBd, db, mock, logger.NewLogger(sugar)
}

func TestUserDatabase_IsCorrect(t *testing.T) {
	conn, _, mock, l := setUp(t)
	defer conn.Close(context.Background())
	// good query
	rows := sqlmock.
		NewRows([]string{"id", "password"})

	expect := testUserBack
	rows = rows.AddRow(expect.Id, expect.Password)
	mock.ExpectQuery(`SELECT id, password FROM users WHERE`).
		WithArgs(expect.Login).
		WillReturnRows(rows)

	ud := &UserDatabase{
		conn: conn,
		logger: l,
	}
	gotUser, err := ud.GetUser(testUserBack)
	if err != nil {
		t.Error(err.Error())
	}
	require.Equal(t, testUserBack, gotUser)
}


func TestUserDatabase_IsExisting(t *testing.T) {
	conn, _, mock, l := setUp(t)
	defer conn.Close(context.Background())

	expect := testUserBack
	mock.ExpectQuery(`id FROM users WHERE`).
		WithArgs(expect.Login).
		WillReturnError(sql.ErrNoRows)

	ud := &UserDatabase{
		conn: conn,
		logger: l,
	}
	flag, err := ud.IsExisting(testUserBack.Login)

	require.Error(t, err, sql.ErrNoRows)
	require.Equal(t, flag, false)
}


*/

var (
	userId   uint64 = 1
	str             = "str"
	login           = "login"
	password        = "123456"
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

func setUp(t *testing.T) *pgxpool.Pool {
	script := &pgmock.Script{
		Steps: pgmock.AcceptUnauthenticatedConnRequestSteps(),
	}
	script.Steps = append(script.Steps, pgmock.ExpectMessage(&pgproto3.Query{String: "SELECT id FROM users WHERE login = $1"}))
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

func TestUserDatabase_AddError(t *testing.T) {
	pool := setUp(t)
	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	sugar := l.Sugar()
	zap.NewAtomicLevelAt(zapcore.DebugLevel)

	h := NewUserDatabase(pool, logger.NewLogger(sugar))
	_, err = h.Add(testRegData)

	assert.Error(t, err)
}

func TestUserDatabase_ChangeAvatarError(t *testing.T) {
	pool := setUp(t)
	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	sugar := l.Sugar()
	zap.NewAtomicLevelAt(zapcore.DebugLevel)

	h := NewUserDatabase(pool, logger.NewLogger(sugar))
	err = h.ChangeAvatar(userId, str)

	assert.Error(t, err)
}

func TestUserDatabase_AddToPreferencesError(t *testing.T) {
	pool := setUp(t)
	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	sugar := l.Sugar()
	zap.NewAtomicLevelAt(zapcore.DebugLevel)

	h := NewUserDatabase(pool, logger.NewLogger(sugar))
	err = h.AddToPreferences(userId)

	assert.Error(t, err)
}

func TestUserDatabase_IsExistingNoRows(t *testing.T) {
	pool := setUp(t)
	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	sugar := l.Sugar()
	zap.NewAtomicLevelAt(zapcore.DebugLevel)

	h := NewUserDatabase(pool, logger.NewLogger(sugar))
	_, err = h.IsExisting(str)
	// returns err.NoRows
	// wtf?!

	assert.Nil(t, err)
}

func TestUserDatabase_IsCorrectError(t *testing.T) {
	pool := setUp(t)
	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	sugar := l.Sugar()
	zap.NewAtomicLevelAt(zapcore.DebugLevel)

	h := NewUserDatabase(pool, logger.NewLogger(sugar))
	_, err = h.IsCorrect(testUserBack)

	assert.Error(t, err)
}

func TestUserDatabase_UpdateError(t *testing.T) {
	pool := setUp(t)
	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	sugar := l.Sugar()
	zap.NewAtomicLevelAt(zapcore.DebugLevel)

	h := NewUserDatabase(pool, logger.NewLogger(sugar))
	err = h.Update(userId, testUserData)

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
	pool := setUp(t)
	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	sugar := l.Sugar()
	zap.NewAtomicLevelAt(zapcore.DebugLevel)

	h := NewUserDatabase(pool, logger.NewLogger(sugar))
	_, err = h.IsExistingEmail(str)
	// тож false с nil

	assert.Nil(t, err)
}

func TestUserDatabase_IsExistingUserIdError(t *testing.T) {
	pool := setUp(t)
	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	sugar := l.Sugar()
	zap.NewAtomicLevelAt(zapcore.DebugLevel)

	h := NewUserDatabase(pool, logger.NewLogger(sugar))
	err = h.IsExistingUserId(userId)
	// тож false с nil

	assert.Error(t, err)
}

func TestUserDatabase_GetUsersError(t *testing.T) {
	pool := setUp(t)
	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	sugar := l.Sugar()
	zap.NewAtomicLevelAt(zapcore.DebugLevel)

	h := NewUserDatabase(pool, logger.NewLogger(sugar))
	_, err = h.GetUsers(1)
	// ErrNoRows ???

	assert.Nil(t, err)
}