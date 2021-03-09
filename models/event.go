package models

type Event struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Place       string `json:"place"`
	Description string `json:"description"`
	Date        string `json:"time"`
	Subway      string `json:"subway"`
	Street      string `json:"street"`
	TypeEvent   string `json:"typeEvent"`
	Image       string `json:"image"`
}

//easyjson:json
type Events []Event
