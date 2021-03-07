package models

type EventInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	TypeEvent   string `json:"typeEvent"`
}

type Event struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	TypeEvent   string `json:"typeEvent"`
}

//easyjson:json
type Events []Event
