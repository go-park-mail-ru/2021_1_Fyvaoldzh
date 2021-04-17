package infrastructure

type SessionTarantool interface {
	CheckSession(value string) (bool, uint64, error)
	InsertSession(uid uint64, value string) error
	DeleteSession(value string) error
}
