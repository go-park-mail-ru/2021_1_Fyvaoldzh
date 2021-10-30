package constants

const (
	Localhost         = "127.0.0.1"
	CookieLength      = uint8(32)
	SessionCookieName = "SID"
	UserIdKey         = "user_id"
	PageKey           = "page"
	IdKey             = "id"
	CSRFHeader        = "header:X-XSRF-TOKEN"
	DBConnect         = " dbname=qdago host=localhost port=5432 sslmode=disable pool_max_conns=10"
	DateFormat        = "2006-01-02"
	DateTimeFormat    = "2006-01-02T15:04:05"
	TarantoolAddress        = "127.0.0.1:3301"
	AuthServicePort         = ":3001"
	SubscriptionServicePort = ":3002"
	ChatServicePort         = ":3003"
	KudagoServicePort       = ":3004"
	TarantoolSpaceName      = "qdago"
	TarantoolSpaceName2     = "user_count"
	TarantoolNotifications  = "notifications"
	TarantoolMessages       = "chat"
	UserPicDir              = "public/users/"
	EventsPicDir            = "public/events/"
	SaltLength              = 8
	//CookiePath              = "/"
	EventsPerPage = 6
	//UsersPerPage            = 10
	ChatPerPage     = 100
	MailingText     = " приглашает Вас на мероприятие "
	MailingAddress  = "https://qdaqda.ru/event"
	MailNotif       = "Mail"
	MailNotifText   = " приглашает Вас посетить новое мероприятие"
	EventNotif      = "Event"
	EventNotifText1 = "Не забудьте посетить мероприятие "
	EventNotifText2 = " которое пройдет через 5 часов"
)

var Category = map[string]string{
	"Развлечения": "entertainment",
	"Образование": "education",
	"Кино":        "cinema",
	"Выставка":    "exhibition",
	"Фестиваль":   "festival",
	"Экскурсия":   "tour",
	"Концерт":     "concert",
}

var LetterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
