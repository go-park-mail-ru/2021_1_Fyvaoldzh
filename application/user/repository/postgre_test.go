package repository

/*
import (
	"context"
	"github.com/stretchr/testify/require"
	"log"

	"github.com/DATA-DOG/go-sqlmock"


	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v4/stdlib"
	"kudago/application/models"

	"testing"
)


func SetupDB() () {


}



func TestUserDatabase_IsCorrect(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()


	pgdb, err := stdlib.AcquireConn(db)
	log.Println(pgdb)
	if err != nil {
		log.Println(err.Error())
	}

	str := pgdb.Config().ConnString()
	pgxPool, err := pgxpool.Connect(context.Background(), str)
	if err != nil {
		t.Fatalf(err.Error())

	}
	repo := &UserDatabase{
		pool: pgxPool,
	}

	// good query
	rows := sqlmock.
		NewRows([]string{"id", "password"})

	expect := testUser
	rows = rows.AddRow(expect.Id, expect.Password)
	mock.ExpectQuery(`SELECT id, password FROM users WHERE`).
		WithArgs(expect.Login).
		WillReturnRows(rows)

	gotUser, err := repo.IsCorrect(&testUser)
	if err != nil {
		t.Error(err.Error())
	}
	require.Equal(t, testUser, gotUser)
}

 */