package models

import (
	"database/sql"
)

type Event struct {
	ID          uint64       `json:"id"`
	Title       string       `json:"title"`
	Place       string       `json:"place"`
	Description string       `json:"description"`
	Date        string       `json:"date"`
	Subway      string       `json:"subway"`
	Street      string       `json:"street"`
	TypeEvent   CategoryTags `json:"typeEvent"`
	Image       string       `json:"image"`
}

type EventSQL struct {
	ID          uint64
	Title       string
	Place       string
	Description string
	Date        sql.NullTime
	Subway      sql.NullString
	Street      sql.NullString
	Image       sql.NullString
}

type EventCard struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

type EventCardSQL struct {
	ID          uint64         `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Image       sql.NullString `json:"image"`
}

type CategoryTag struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type CategoryTagDescription struct {
	CategoryID   uint64
	CategoryName string
	TagID        uint64
	TagName      string
}

func ConvertCard(old EventCardSQL) EventCard {
	var newCard EventCard
	newCard.ID = old.ID
	newCard.Title = old.Title
	newCard.Description = old.Description
	newCard.Image = old.Image.String
	return newCard
}

func ConvertEvent(old EventSQL) Event {
	var newEvent Event
	newEvent.ID = old.ID
	newEvent.Title = old.Title
	newEvent.Place = old.Place
	newEvent.Description = old.Description
	newEvent.Date = old.Date.Time.String()
	newEvent.Subway = old.Subway.String
	newEvent.Street = old.Street.String
	newEvent.Image = old.Image.String
	return newEvent
}

//easyjson:json
type Events []Event

//easyjson:json
type EventCards []EventCard

//easyjson:json
type CategoryTags []CategoryTag
