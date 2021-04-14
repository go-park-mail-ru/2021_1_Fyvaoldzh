package repository

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"kudago/application/models"
	"kudago/pkg/constants"
	"kudago/pkg/logger"
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


func setUp(t *testing.T) (*pgx.Conn, *sql.DB, sqlmock.Sqlmock, logger.Logger) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(err.Error())
	}

	newBd, _ := stdlib.AcquireConn(db)

	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	sugar := l.Sugar()
	zap.NewAtomicLevelAt(zapcore.DebugLevel)

	return newBd, db, mock, logger.NewLogger(sugar)
}

/*
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
	gotUser, err := ud.IsCorrect(testUserBack)
	if err != nil {
		t.Error(err.Error())
	}
	require.Equal(t, testUserBack, gotUser)
}

 */

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


