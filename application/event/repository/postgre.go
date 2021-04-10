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
		`SELECT id, title, description, image, start_date, end_date FROM events
		ORDER BY id DESC`)

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
		`SELECT id, title, description, image, start_date, end_date FROM events
		WHERE category = $1 ORDER BY id DESC`, typeEvent)

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
		`INSERT INTO events (title, place, subway, street, description, category, start_date, end_date, image) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		newEvent.Title, newEvent.Place, newEvent.Subway, newEvent.Street, newEvent.Description,
		newEvent.Category, newEvent.StartDate, newEvent.EndDate, newEvent.Image)
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
		`SELECT DISTINCT ON(e.id) e.id, e.title, e.description,
		e.image, e.start_date, e.end_date FROM
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

func (ed EventDatabase) RecomendSystem(uid uint64, category string) error {
	if category == "Музей" {
		_, err := ed.pool.Exec(context.Background(),
			`UPDATE user_preference SET concert = concert + 1 WHERE user_id = $1`, uid)

		if err != nil {
			return err
		}
	}

	if category == "Выставка" {
		_, err := ed.pool.Exec(context.Background(),
			`UPDATE user_preference SET show = show + 1 WHERE user_id = $1`, uid)

		if err != nil {
			return err
		}
	}

	if category == "Кино" {
		_, err := ed.pool.Exec(context.Background(),
			`UPDATE user_preference SET movie = movie + 1 WHERE user_id = $1`, uid)

		if err != nil {
			return err
		}
	}

	return nil
}

func (ed EventDatabase) GetPreference(uid uint64) (models.Recomend, error) {
	var recomend models.Recomend
	err := pgxscan.Select(context.Background(), ed.pool, &recomend,
		`SELECT show, movie, concert
		FROM user_preference
		WHERE user_id = $1`, uid)

	if err != nil {
		return models.Recomend{}, err
	}
	return recomend, nil
}

//TODO сделать нормальный обсчет
func (ed EventDatabase) GetRecomended(uid uint64) ([]models.EventCardWithDateSQL, error) {
	recomend, err := ed.GetPreference(uid)
	if err == nil {
		var events []models.EventCardWithDateSQL
		err := pgxscan.Select(context.Background(), ed.pool, &events,
			`SELECT id, title, description, image, start_date, end_date FROM events
		ORDER BY id DESC`)

		if errors.As(err, &pgx.ErrNoRows) || len(events) == 0 {
			return []models.EventCardWithDateSQL{}, nil
		}

		if err != nil {
			return nil, err
		}

		return events, nil
	}
	if recomend.Concert >= recomend.Movie && recomend.Concert >= recomend.Show {
		var events []models.EventCardWithDateSQL
		err := pgxscan.Select(context.Background(), ed.pool, &events,
			`SELECT id, title, description, image, start_date, end_date FROM events
			WHERE category = 'Музей'
			ORDER BY id DESC`)
		if errors.As(err, &pgx.ErrNoRows) || len(events) == 0 {
			return []models.EventCardWithDateSQL{}, nil
		}

		if err != nil {
			return nil, err
		}
		err = pgxscan.Select(context.Background(), ed.pool, &events,
			`SELECT id, title, description, image, start_date, end_date FROM events
			WHERE category != 'Музей'
			ORDER BY id DESC`)
		if errors.As(err, &pgx.ErrNoRows) || len(events) == 0 {
			return []models.EventCardWithDateSQL{}, nil
		}

		if err != nil {
			return nil, err
		}
		return events, nil
	}
	if recomend.Movie >= recomend.Concert && recomend.Movie >= recomend.Show {
		var events []models.EventCardWithDateSQL
		err := pgxscan.Select(context.Background(), ed.pool, &events,
			`SELECT id, title, description, image, start_date, end_date FROM events
			WHERE category = 'Кино'
			ORDER BY id DESC`)
		if errors.As(err, &pgx.ErrNoRows) || len(events) == 0 {
			return []models.EventCardWithDateSQL{}, nil
		}

		if err != nil {
			return nil, err
		}
		err = pgxscan.Select(context.Background(), ed.pool, &events,
			`SELECT id, title, description, image, start_date, end_date FROM events
			WHERE category != 'Кино'
			ORDER BY id DESC`)
		if errors.As(err, &pgx.ErrNoRows) || len(events) == 0 {
			return []models.EventCardWithDateSQL{}, nil
		}

		if err != nil {
			return nil, err
		}
		return events, nil
	} else {
		var events []models.EventCardWithDateSQL
		err := pgxscan.Select(context.Background(), ed.pool, &events,
			`SELECT id, title, description, image, start_date, end_date FROM events
			WHERE category = 'Выставка'
			ORDER BY id DESC`)
		if errors.As(err, &pgx.ErrNoRows) || len(events) == 0 {
			return []models.EventCardWithDateSQL{}, nil
		}

		if err != nil {
			return nil, err
		}
		err = pgxscan.Select(context.Background(), ed.pool, &events,
			`SELECT id, title, description, image, start_date, end_date FROM events
			WHERE category != 'Выставка'
			ORDER BY id DESC`)
		if errors.As(err, &pgx.ErrNoRows) || len(events) == 0 {
			return []models.EventCardWithDateSQL{}, nil
		}

		if err != nil {
			return nil, err
		}
		return events, nil
	}
}
