package models

type Event struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	TypeEvent   string `json:"typeEvent"`
	Image       string `json:"image"`
}

//easyjson:json
type Events []Event
