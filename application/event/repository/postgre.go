package repository

import (
	"context"
	"errors"
	"fmt"
	"kudago/application/event"
	"kudago/application/models"
	"net/http"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo"
)

type EventDatabase struct {
	pool *pgxpool.Pool
}

func NewEventDatabase(conn *pgxpool.Pool) event.Repository {
	return &EventDatabase{pool: conn}
}

func (ed EventDatabase) GetAllEvents() ([]models.EventCardWithDateSQL, error) {
	var events []models.EventCardWithDateSQL
	err := pgxscan.Select(context.Background(), ed.pool, &events,
		`SELECT id, title, description, image, start_date FROM events`)

	if errors.As(err, &pgx.ErrNoRows) || len(events) == 0 {
		return []models.EventCardWithDateSQL{}, nil
	}

	if err != nil {
		return nil, err
	}

	return events, nil
}

func (ed EventDatabase) GetEventsByCategory(typeEvent string) ([]models.EventCardWithDateSQL, error) {
	var events []models.EventCardWithDateSQL
	err := pgxscan.Select(context.Background(), ed.pool, &events,
		`SELECT id, title, description, image, start_date FROM events
		WHERE category = $1`, typeEvent)

	if errors.As(err, &pgx.ErrNoRows) || len(events) == 0 {
		return []models.EventCardWithDateSQL{}, nil
	}

	if err != nil {
		return []models.EventCardWithDateSQL{}, err
	}

	return events, nil
}

func (ed EventDatabase) GetOneEventByID(eventId uint64) (models.EventSQL, error) {
	var ev []models.EventSQL
	err := pgxscan.Select(context.Background(), ed.pool, &ev,
		`SELECT * FROM events WHERE id = $1`, eventId)

	if errors.As(err, &pgx.ErrNoRows) || len(ev) == 0 {
		return models.EventSQL{}, echo.NewHTTPError(http.StatusNotFound, errors.New("Event with id "+fmt.Sprint(eventId)+" not found"))
	}

	if err != nil {
		return models.EventSQL{}, err
	}

	return ev[0], nil
}

func (ed EventDatabase) GetTags(eventId uint64) (models.Tags, error) {
	var parameters models.Tags
	err := pgxscan.Select(context.Background(), ed.pool, &parameters,
		`SELECT t.id AS id, t.name AS name
		FROM tags t
		JOIN event_tag e on e.tag_id = t.id
        WHERE e.event_id = $1`, eventId)

	if errors.As(err, &pgx.ErrNoRows) || len(parameters) == 0 {
		return models.Tags{}, nil
	}

	if err != nil {
		return models.Tags{}, err
	}

	return parameters, err
}

func (ed EventDatabase) AddEvent(newEvent *models.Event) error {
	// TODO: добавить промежуточный sql, который будет в базу null пихать
	_, err := ed.pool.Exec(context.Background(),
		`INSERT INTO events (title, place, subway, street, description, date, image) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		newEvent.Title, newEvent.Place, newEvent.Subway, newEvent.Street, newEvent.Description, newEvent.StartDate, newEvent.Image)
	if err != nil {
		return err
	}

	return nil
}

func (ed EventDatabase) DeleteById(eventId uint64) error {
	resp, err := ed.pool.Exec(context.Background(),
		`DELETE FROM events WHERE id = $1`, eventId)

	if err != nil {
		return err
	}

	if resp.RowsAffected() == 0 {
		return echo.NewHTTPError(http.StatusNotFound, errors.New("Event with id "+fmt.Sprint(eventId)+" not found"))
	}

	return nil
}

func (ed EventDatabase) UpdateEventAvatar(eventId uint64, path string) error {
	_, err := ed.pool.Exec(context.Background(),
		`UPDATE events SET image = $1 WHERE id = $2`, path, eventId)

	if err != nil {
		return err
	}

	return nil
}

func (ed EventDatabase) FindEvents(str string) ([]models.EventCardWithDateSQL, error) {
	var events []models.EventCardWithDateSQL
	err := pgxscan.Select(context.Background(), ed.pool, &events,
		`SELECT DISTINCT e.id as id, e.title as title, e.description as description, e.image as image, e.start_date as start_date FROM
        events e JOIN event_tag et on e.id = et.event_id
        JOIN tags t on et.tag_id = t.id
		WHERE LOWER(title) LIKE '%' || $1 || '%' OR LOWER(description) LIKE '%' || $1 || '%'
		OR LOWER(category) LIKE '%' || $1 || '%' OR LOWER(t.name) LIKE '%' || $1 || '%'`, str)

	if errors.As(err, &pgx.ErrNoRows) || len(events) == 0 {
		return []models.EventCardWithDateSQL{}, nil
	}

	if err != nil {
		return []models.EventCardWithDateSQL{}, err
	}

	return events, nil
}