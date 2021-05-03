package session

type UseCase interface {
	Login(login string, password string) (string, bool, error)
	Check(value string) (bool, uint64, error)
	Logout(value string) error
}