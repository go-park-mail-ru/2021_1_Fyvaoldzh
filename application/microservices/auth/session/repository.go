package session

type Repository interface {
	InsertSession(userId uint64, value string) error
	CheckSession(value string) (bool, uint64, error)
	DeleteSession(value string) error
}
