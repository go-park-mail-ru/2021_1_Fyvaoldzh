package constants

const (
	Localhost               = "127.0.0.1"
	CookieLength            = uint8(32)
	SessionCookieName       = "SID"
	UserIdKey               = "user_id"
	PageKey                 = "page"
	IdKey                   = "id"
	CSRFHeader              = "header:X-XSRF-TOKEN"
	DBConnect               = "user=postgre dbname=qdago password=fyvaoldzh host=localhost port=5432 sslmode=disable pool_max_conns=10"
	DateFormat              = "2006-01-02"
	DateTimeFormat          = "2006-01-02T15:04:05"
	DefaultAvatar           = "public/default.png"
	TarantoolAddress        = "127.0.0.1:3301"
	TarantoolUser           = "admin"
	AuthServicePort         = ":3001"
	SubscriptionServicePort = ":3002"
	ChatServicePort         = ":3003"
	TarantoolPassword       = "fyvaoldzh"
	TarantoolSpaceName      = "qdago"
	UserPicDir              = "public/users/"
	EventsPicDir            = "public/events/"
	SaltLength              = 8
	CookiePath              = "/"
	EventsPerPage           = 6
	UsersPerPage            = 10
	ChatPerPage             = 100
	MailingText             = " приглашает Вас на мероприятие "
	MailingAddress          = "95.163.180.8:3000/event"
)

var Category = map[string]string{
	"Музей":    "concert",
	"Выставка": "show",
	"Кино":     "movie",
}

var LetterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
