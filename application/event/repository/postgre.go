package repository

import (
	"context"
	"errors"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"kudago/application/event"
	"kudago/application/models"
	"log"
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

func (ed EventDatabase) GetOneEventByID(eventId uint64) (models.EventSQL, error) {
	var ev []models.EventSQL
	err := pgxscan.Select(context.Background(), ed.pool, &ev,
		`SELECT * FROM events WHERE id = $1`, eventId)

	if errors.As(err, &pgx.ErrNoRows){
		return models.EventSQL{}, nil
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
