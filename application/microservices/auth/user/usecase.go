package user

type UseCase interface {
	CheckUser(login string, password string) (uint64, bool, error)
}
