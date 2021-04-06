package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo"
	"kudago/application/event"
	"kudago/application/models"
	"log"
	"net/http"
)

type EventDatabase struct {
	pool *pgxpool.Pool
}

func NewEventDatabase(conn *pgxpool.Pool) event.Repository {
	return &EventDatabase{pool: conn}
}

func (ed EventDatabase) GetAllEvents() ([]models.EventCardSQL, error) {
	var events []models.EventCardSQL
	err := pgxscan.Select(context.Background(), ed.pool, &events,
		`SELECT id, title, description, image FROM events`)

	if errors.As(err, &pgx.ErrNoRows) || len(events) == 0 {
		return []models.EventCardSQL{}, nil
	}

	if err != nil {
		return nil, err
	}

	return events, nil
}

func (ed EventDatabase) GetEventsByType(typeEvent string) ([]models.EventCardSQL, error) {
	var events []models.EventCardSQL
	err := pgxscan.Select(context.Background(), ed.pool, &events,
		`SELECT DISTINCT e.id, e.title, e.description, e.image FROM events e
		JOIN ev_cat_tag ect ON ect.eid = e.id
		JOIN categories c ON c.id = ect.cid
		JOIN tags t on t.id = ect.tid
		WHERE t.name = $1 OR c.name = $1;`, typeEvent)

	if errors.As(err, &pgx.ErrNoRows) || len(events) == 0 {
		return []models.EventCardSQL{}, nil
	}

	if err != nil {
		return []models.EventCardSQL{}, err
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

func (ed EventDatabase) GetCategoryTags(eventId uint64) ([]models.CategoryTagDescription, error) {
	var parameters []models.CategoryTagDescription
	err := pgxscan.Select(context.Background(), ed.pool, &parameters,
		`SELECT c.id AS category_id, c.name AS category_name, 
		t.id AS tag_id, t.name AS tag_name
		FROM categories c, tags t
		JOIN ev_cat_tag on eid = $1
		WHERE c.id = cid AND t.id = tid`, eventId)

	if err != nil {
		log.Println(err)
		return []models.CategoryTagDescription{}, err
	}

	return parameters, err
}

func (ed EventDatabase) AddEvent(newEvent *models.Event) error {
	// TODO: добавить промежуточный sql, который будет в базу null пихать
	_, err := ed.pool.Exec(context.Background(),
		`INSERT INTO events (title, place, subway, street, description, date, image) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		newEvent.Title, newEvent.Place, newEvent.Subway, newEvent.Street, newEvent.Description, newEvent.Date, newEvent.Image)
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

func (ed EventDatabase) FindEvents(str string) ([]models.EventCardSQL, error) {
	var events []models.EventCardSQL
	err := pgxscan.Select(context.Background(), ed.pool, &events,
		`SELECT DISTINCT id, title, description, image FROM events 
		WHERE LOWER(title) LIKE '%' || $1 || '%' OR LOWER(description) LIKE '%' || $1 || '%';`, str)


	if errors.As(err, &pgx.ErrNoRows) || len(events) == 0 {
		return []models.EventCardSQL{}, nil
	}

	if err != nil {
		return []models.EventCardSQL{}, err
	}

	return events, nil
}