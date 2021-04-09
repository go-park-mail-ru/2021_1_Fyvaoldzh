package models

import (
	"database/sql"
	"time"
)

type Event struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Place       string `json:"place"`
	Description string `json:"description"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	Subway      string `json:"subway"`
	Street      string `json:"street"`
	Tags        Tags   `json:"tags"`
	Category    string `json:"category"`
	Image       string `json:"image"`
}

type EventSQL struct {
	ID          uint64
	Title       string
	Place       string
	Description string
	StartDate   sql.NullTime
	EndDate     sql.NullTime
	Subway      sql.NullString
	Street      sql.NullString
	Category    string
	Image       sql.NullString
}

type EventCard struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	StartDate   string `json:"startDate"`
}

type EventCardWithDateSQL struct {
	ID          uint64
	Title       string
	Description string
	Image       sql.NullString
	StartDate   time.Time
}

type Tag struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

func ConvertDateCard(old EventCardWithDateSQL) EventCard {
	var newCard EventCard
	newCard.ID = old.ID
	newCard.Title = old.Title
	newCard.Description = old.Description
	newCard.Image = old.Image.String
	newCard.StartDate = old.StartDate.String()
	return newCard
}

func ConvertEvent(old EventSQL) Event {
	var newEvent Event
	newEvent.ID = old.ID
	newEvent.Title = old.Title
	newEvent.Place = old.Place
	newEvent.Description = old.Description
	newEvent.StartDate = old.StartDate.Time.String()
	newEvent.EndDate = old.EndDate.Time.String()
	newEvent.Subway = old.Subway.String
	newEvent.Street = old.Street.String
	newEvent.Category = old.Category
	newEvent.Image = old.Image.String
	return newEvent
}

//easyjson:json
type Events []Event

//easyjson:json
type EventCards []EventCard

//easyjson:json
type Tags []Tag
