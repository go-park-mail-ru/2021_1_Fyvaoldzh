package client

type IAuthClient interface {
	Login(login string, password string, value string) (uint64, string, error, int)
	Logout(value string) (error, int)
	Check(value string) (bool, uint64, error, int)
	Close()
}
