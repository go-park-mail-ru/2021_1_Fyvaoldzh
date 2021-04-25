package user

type UseCase interface {
	Login(login string, password string) (uint64, error)
	CheckUser(login string, password string) (uint64, error)
	Logout(value string) error
}
