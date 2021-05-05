package client

type IAuthClient interface {
	Login(login string, password string, value string) (uint64, string, error)
	Logout(value string) error
	Check(value string) (bool, uint64, error)
	Close()
}
