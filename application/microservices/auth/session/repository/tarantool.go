package repository

import (
	"github.com/labstack/echo"
	"github.com/tarantool/go-tarantool"
	"kudago/application/microservices/auth/session"
	"kudago/pkg/constants"
	"net/http"
)

type SessionRepository struct {
	Conn *tarantool.Connection
}

func NewSessionRepository(c *tarantool.Connection) session.Repository {
	return &SessionRepository{Conn: c}
}

func (s SessionRepository) InsertSession(userId uint64, value string) error {
	_, err := s.Conn.Insert(constants.TarantoolSpaceName, []interface{}{value, userId})
	if err != nil {
		return err
	}

	return nil

}

func (s SessionRepository) CheckSession(value string) (bool, uint64, error) {
	resp, err := s.Conn.Select(constants.TarantoolSpaceName,
		"primary", 0, 1, tarantool.IterEq, []interface{}{value})
	if err != nil {
		return false, 0, err
	}

	if len(resp.Data) != 0 {
		data := resp.Data[0]
		d, ok := data.([]interface{})
		if !ok {
			return false, 0, echo.NewHTTPError(http.StatusBadRequest, "cast error")
		}

		sid, ok := d[1].(uint64)
		if !ok {
			return false, 0, echo.NewHTTPError(http.StatusBadRequest, "cast error")
		}

		return true, sid, nil
	}

	return false, 0, nil
}

func (s SessionRepository) DeleteSession(value string) error {
	_, err := s.Conn.Delete(constants.TarantoolSpaceName, "primary", []interface{}{value})
	if err != nil {
		return err
	}

	return nil
}