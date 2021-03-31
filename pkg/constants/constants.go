package constants

const (
	CookieLength      = uint8(32)
	SessionCookieName = "SID"
	DBConnect         = "user=postgre dbname=qdago password=fyvaoldzh host=127.0.0.1 port=5432 sslmode=disable pool_max_conns=10"
	TimeFormat        = "2006-01-02"
	DefaultAvatar     = "public/default.png"
)

var LetterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
