package constants

const (
	Localhost          = "127.0.0.1"
	CookieLength       = uint8(32)
	SessionCookieName  = "SID"
	CSRFHeader         = "header:X-XSRF-TOKEN"
	DBConnect          = "user=postgre dbname=qdago password=fyvaoldzh host=localhost port=5432 sslmode=disable pool_max_conns=10"
	TimeFormat         = "2006-01-02"
	DefaultAvatar      = "public/default.png"
	TarantoolAddress   = "127.0.0.1:3301"
	TarantoolUser      = "admin"
	AuthServicePort    = ":3001"
	TarantoolPassword  = "fyvaoldzh"
	TarantoolSpaceName = "qdago"
	UserPicDir         = "public/users/"
	EventsPicDir       = "public/events/"
	SaltLength         = 8
	CookiePath         = "/"
)

var Category = map[string]string{
	"Музей":    "concert",
	"Выставка": "show",
	"Кино":     "movie",
}

var LetterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
