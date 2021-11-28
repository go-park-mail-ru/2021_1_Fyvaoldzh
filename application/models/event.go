package models

import (
	"database/sql"
	"time"
)

type Event struct {
	ID          uint64       `json:"id"`
	Title       string       `json:"title"`
	Place       string       `json:"place"`
	Description string       `json:"description"`
	StartDate   string       `json:"startDate"`
	EndDate     string       `json:"endDate"`
	Subway      string       `json:"subway"`
	Street      string       `json:"street"`
	Tags        Tags         `json:"tags"`
	Category    string       `json:"category"`
	Coordinates []float64    `json:"coordinates"`
	Image       string       `json:"image"`
	Followers   UsersOnEvent `json:"followers"`
	Latitude    float64
	Longitude   float64
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
	Latitude    sql.NullFloat64
	Longitude   sql.NullFloat64
	Image       sql.NullString
}

type EventCard struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Place       string `json:"place"`
	Description string `json:"description"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
}

type EventCardWithCoords struct {
	ID          uint64  `json:"id"`
	Title       string  `json:"title"`
	Place       string  `json:"place"`
	Description string  `json:"description"`
	StartDate   string  `json:"startDate"`
	EndDate     string  `json:"endDate"`
	Distance    float64 `json:"distance"`
}

type EventCardWithDateSQL struct {
	ID          uint64
	Title       string
	Place       string
	Description string
	StartDate   time.Time
	EndDate     time.Time
}

type EventCardWithCoordsSQL struct {
	ID          uint64
	Title       string
	Place       string
	Description string
	StartDate   time.Time
	EndDate     time.Time
	Distance    float64 `db:"distance"`
}

type Tag struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type IsAddedEvent struct {
	UserId  uint64 `json:"userId"`
	EventId uint64 `json:"eventId"`
	IsAdded bool   `json:"isAdded"`
}

type Recomend struct {
	Entertainment uint64
	Education     uint64
	Cinema        uint64
	Exhibition    uint64
	Festival      uint64
	Tour          uint64
	Concert       uint64
}

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func ConvertDateCard(old EventCardWithDateSQL) EventCard {
	var newCard EventCard
	newCard.ID = old.ID
	newCard.Title = old.Title
	newCard.Description = old.Description
	newCard.Place = old.Place
	newCard.StartDate = old.StartDate.String()
	newCard.EndDate = old.EndDate.String()
	return newCard
}

func ConvertCoordsCard(old EventCardWithCoordsSQL, point Coordinates) EventCardWithCoords {
	var newCard EventCardWithCoords
	newCard.ID = old.ID
	newCard.Title = old.Title
	newCard.Description = old.Description
	newCard.Place = old.Place
	newCard.StartDate = old.StartDate.String()
	newCard.EndDate = old.EndDate.String()
	newCard.Distance = old.Distance
	return newCard
}

func ConvertEvent(old EventSQL) Event {
	var newEvent Event
	newEvent.ID = old.ID
	newEvent.Title = old.Title
	newEvent.Place = old.Place
	newEvent.Description = old.Description
	if old.StartDate.Valid {
		newEvent.StartDate = old.StartDate.Time.String()
	}
	if old.EndDate.Valid {
		newEvent.EndDate = old.EndDate.Time.String()
	}
	newEvent.Subway = old.Subway.String
	newEvent.Street = old.Street.String
	newEvent.Category = old.Category
	newEvent.Image = old.Image.String
	if old.Latitude.Valid && old.Longitude.Valid {
		newEvent.Coordinates = append(newEvent.Coordinates, old.Latitude.Float64, old.Longitude.Float64)
	}
	return newEvent
}

//easyjson:json
type Events []Event

//easyjson:json
type EventCards []EventCard

//easyjson:json
type Tags []Tag

//easyjson:json
type EventCardsWithCoords []EventCardWithCoords

type ViewData struct {
	Id    uint64
	Title string
}
