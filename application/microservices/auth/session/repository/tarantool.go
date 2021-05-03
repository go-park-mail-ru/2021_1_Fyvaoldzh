package repository

import (
	"errors"
	"github.com/tarantool/go-tarantool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"kudago/application/microservices/auth/session"
	"kudago/pkg/constants"
	"kudago/pkg/logger"
)

type SessionRepository struct {
	Conn *tarantool.Connection
	lg logger.Logger
}

func NewSessionRepository(c *tarantool.Connection, l logger.Logger) session.Repository {
	return &SessionRepository{Conn: c, lg: l}
}

func (s SessionRepository) InsertSession(userId uint64, value string) error {
	_, err := s.Conn.Insert(constants.TarantoolSpaceName, []interface{}{value, userId})
	if err != nil {
		s.lg.Warn(err)
		return err
	}

	return nil

}

func (s SessionRepository) CheckSession(value string) (bool, uint64, error) {
	resp, err := s.Conn.Select(constants.TarantoolSpaceName,
		"primary", 0, 1, tarantool.IterEq, []interface{}{value})
	if err != nil {
		s.lg.Warn(err)
		return false, 0, err
	}

	if len(resp.Data) != 0 {
		data := resp.Data[0]
		d, ok := data.([]interface{})
		if !ok {
			return false, 0, errors.New("cast error")
		}

		sid, ok := d[1].(uint64)
		if !ok {
			return false, 0, errors.New("cast error")
		}

		return true, sid, nil
	}

	return false, 0, nil
}

func (s SessionRepository) DeleteSession(value string) error {
	_, err := s.Conn.Delete(constants.TarantoolSpaceName, "primary", []interface{}{value})
	if err != nil {
		s.lg.Warn(err)
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}