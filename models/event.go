package models

type Event struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Place       string `json:"place"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Subway      string `json:"subway"`
	Street      string `json:"street"`
	TypeEvent   string `json:"typeEvent"`
	Image       string `json:"image"`
}

//easyjson:json
type Events []Event

var BaseEvents = Events{
	{125, "Идущий к реке", "У реки", "Я в своем познаии настолько преисполнился", "11.11.1111 11:11", " ", "Пушкина", "cognition", "kudago/125"},
	{126, "Димон заминированный тапок", "На дороге", "Мне абсолютно все равно", "12.12.1212 12:12", " ", "Колотушкина", "cognition", "kudago/126"},
	{127, "exampleTitle", "examplePlace", "exampleText", "01.01.0001 00:00", "exampleSubway", "exampleStreet", "1", "kudago/127"},
	{128, "Пример", "Место", "Пример без картинки", "12:00", "Примерное метро", "Примерная улица", "noImg", ""},
}
