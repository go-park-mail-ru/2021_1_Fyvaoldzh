package session_manager

import (
	"github.com/labstack/echo"
	"github.com/tarantool/go-tarantool"
	"kudago/pkg/constants"
	"kudago/pkg/infrastructure"
	"net/http"
)

type SessionManager struct {
	Conn *tarantool.Connection
}

func NewSessionManager(c *tarantool.Connection) infrastructure.SessionTarantool {
	return &SessionManager{Conn: c}
}

func (sm SessionManager) CheckSession(value string) (bool, uint64, error) {
	resp, err := sm.Conn.Select(constants.TarantoolSpaceName,
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

func (sm SessionManager) InsertSession(uid uint64, value string) error {
	_, err := sm.Conn.Insert(constants.TarantoolSpaceName, []interface{}{value, uid})

	if err != nil {
		return err
	}

	return nil
}

func (sm SessionManager) DeleteSession(value string) error {
	_, err := sm.Conn.Delete(constants.TarantoolSpaceName, "primary", []interface{}{value})

	if err != nil {
		return err
	}

	return nil
}
