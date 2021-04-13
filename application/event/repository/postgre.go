package repository

import (
	"context"
	"errors"
	"fmt"
	"kudago/application/event"
	"kudago/application/models"
	"kudago/pkg/constants"
	"kudago/pkg/logger"
	"net/http"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo"
)

type EventDatabase struct {
	pool   *pgxpool.Pool
	logger logger.Logger
}

func NewEventDatabase(conn *pgxpool.Pool, logger logger.Logger) event.Repository {
	return &EventDatabase{pool: conn, logger: logger}
}

func (ed EventDatabase) GetAllEvents(now time.Time, page int) ([]models.EventCardWithDateSQL, error) {
	var events []models.EventCardWithDateSQL
	err := pgxscan.Select(context.Background(), ed.pool, &events,
		`SELECT id, title, description, image, start_date, end_date FROM events
		WHERE end_date > $1
		ORDER BY id DESC
		LIMIT 6 OFFSET $2`, now, (page-1)*6)

	if errors.As(err, &pgx.ErrNoRows) || len(events) == 0 {
		ed.logger.Debug("no rows in method GetAllEvents")
		return []models.EventCardWithDateSQL{}, nil
	}

	if err != nil {
		ed.logger.Warn(err)
		return nil, err
	}

	return events, nil
}

func (ed EventDatabase) GetEventsByCategory(typeEvent string, now time.Time, page int) ([]models.EventCardWithDateSQL, error) {
	var events []models.EventCardWithDateSQL
	err := pgxscan.Select(context.Background(), ed.pool, &events,
		`SELECT id, title, description, image, start_date, end_date FROM events
		WHERE category = $1 AND end_date > $2
		ORDER BY id DESC
		LIMIT 6 OFFSET $3`, typeEvent, now, (page-1)*6)

	if errors.As(err, &pgx.ErrNoRows) || len(events) == 0 {
		ed.logger.Debug("no rows in method GetEventsByCategory")
		return []models.EventCardWithDateSQL{}, nil
	}

	if err != nil {
		ed.logger.Warn(err)
		return []models.EventCardWithDateSQL{}, err
	}

	return events, nil
}

func (ed EventDatabase) GetOneEventByID(eventId uint64) (models.EventSQL, error) {
	var ev []models.EventSQL
	err := pgxscan.Select(context.Background(), ed.pool, &ev,
		`SELECT * FROM events WHERE id = $1`, eventId)

	if errors.As(err, &pgx.ErrNoRows) || len(ev) == 0 {
		ed.logger.Debug("no event with id " + fmt.Sprint(eventId))
		return models.EventSQL{}, echo.NewHTTPError(http.StatusNotFound, errors.New("Event with id "+fmt.Sprint(eventId)+" not found"))
	}

	if err != nil {
		ed.logger.Warn(err)
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
		ed.logger.Debug("no rows in method GetTags")
		return models.Tags{}, nil
	}

	if err != nil {
		ed.logger.Warn(err)
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
		ed.logger.Warn(err)
		return err
	}

	return nil
}

func (ed EventDatabase) DeleteById(eventId uint64) error {
	resp, err := ed.pool.Exec(context.Background(),
		`DELETE FROM events WHERE id = $1`, eventId)

	if err != nil {
		ed.logger.Warn(err)
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
		ed.logger.Warn(err)
		return err
	}

	return nil
}

func (ed EventDatabase) FindEvents(str string, now time.Time, page int) ([]models.EventCardWithDateSQL, error) {
	var events []models.EventCardWithDateSQL
	err := pgxscan.Select(context.Background(), ed.pool, &events,
		`SELECT DISTINCT ON(e.id) e.id, e.title, e.description,
		e.image, e.start_date, e.end_date FROM
        events e JOIN event_tag et on e.id = et.event_id
        JOIN tags t on et.tag_id = t.id
		WHERE (LOWER(title) LIKE '%' || $1 || '%' OR LOWER(description) LIKE '%' || $1 || '%'
		OR LOWER(category) LIKE '%' || $1 || '%' OR LOWER(t.name) LIKE '%' || $1 || '%')
		AND end_date > $2
		ORDER BY e.id DESC
		LIMIT 6 OFFSET $3`, str, now, (page-1)*6)

	if errors.As(err, &pgx.ErrNoRows) || len(events) == 0 {
		ed.logger.Debug("no rows in method FindEvents with string " + str)
		return []models.EventCardWithDateSQL{}, nil
	}

	if err != nil {
		ed.logger.Warn(err)
		return []models.EventCardWithDateSQL{}, err
	}

	return events, nil
}

func (ed EventDatabase) RecomendSystem(uid uint64, category string) error {

	_, err := ed.pool.Exec(context.Background(),
		`UPDATE user_preference SET `+constants.Category[category]+`=`+constants.Category[category]+`+1 `+`WHERE user_id = $1`, uid)

	if errors.As(err, &pgx.ErrNoRows) {
		ed.logger.Debug("no rows in method RecomendSystem with id " + fmt.Sprint(uid))
		return err
	}

	if err != nil {
		ed.logger.Warn(err)
		return err
	}

	return nil
}

func (ed EventDatabase) GetSixPreference(recomend models.Recomend) models.Recomend {
	var sixPreference models.Recomend
	recomendSumm := recomend.Concert + recomend.Movie + recomend.Show
	concertProcent := float64(recomend.Concert) / float64(recomendSumm)
	showProcent := float64(recomend.Show) / float64(recomendSumm)
	sixPreference.Concert = uint64(concertProcent * 6)
	if sixPreference.Concert == 0 {
		sixPreference.Concert = 1
	}
	if sixPreference.Concert > 4 {
		sixPreference.Concert = 4
	}
	sixPreference.Show = uint64(showProcent * 6)
	if sixPreference.Show == 0 {
		sixPreference.Concert = 1
	}
	if sixPreference.Show > 4 {
		sixPreference.Concert = 4
	}
	sixPreference.Movie = 6 - sixPreference.Concert - sixPreference.Show
	return sixPreference
}

func (ed EventDatabase) GetPreference(uid uint64) (models.Recomend, error) {
	var recomend []models.Recomend
	err := pgxscan.Select(context.Background(), ed.pool, &recomend,
		`SELECT show, movie, concert
		FROM user_preference
		WHERE user_id = $1`, uid)

	if errors.As(err, &pgx.ErrNoRows) {
		ed.logger.Debug("no rows in method GetPreference with id " + fmt.Sprint(uid))
		return models.Recomend{}, err
	}

	if err != nil {
		ed.logger.Warn(err)
		return models.Recomend{}, err
	}
	//return recomend[0], nil
	return ed.GetSixPreference(recomend[0]), nil
}

func (ed EventDatabase) CategorySearch(str string, category string, now time.Time, page int) ([]models.EventCardWithDateSQL, error) {
	var events []models.EventCardWithDateSQL
	err := pgxscan.Select(context.Background(), ed.pool, &events,
		`SELECT DISTINCT ON(e.id) e.id, e.title, e.description,
		e.image, e.start_date, e.end_date FROM
        events e JOIN event_tag et on e.id = et.event_id
        JOIN tags t on et.tag_id = t.id
		WHERE (LOWER(title) LIKE '%' || $1 || '%' OR LOWER(description) LIKE '%' || $1 || '%'
		OR LOWER(t.name) LIKE '%' || $1 || '%') AND e.category = $2
		AND end_date > $3
		ORDER BY e.id DESC
		LIMIT 6 OFFSET $3`, str, category, now, (page-1)*6)

	if errors.As(err, &pgx.ErrNoRows) || len(events) == 0 {
		ed.logger.Debug("no rows in method CategorySearch with searchstring " + str)
		return []models.EventCardWithDateSQL{}, nil
	}

	if err != nil {
		ed.logger.Warn(err)
		return []models.EventCardWithDateSQL{}, err
	}

	return events, nil
}

//TODO сделать нормальный обсчет
func (ed EventDatabase) GetRecommended(uid uint64, now time.Time, page int) ([]models.EventCardWithDateSQL, error) {
	recomend, err := ed.GetPreference(uid)
	if err != nil || (recomend.Concert == recomend.Movie && recomend.Movie == recomend.Show && recomend.Show == 0) {
		ed.logger.Debug(string(err.Error()))
		return ed.GetAllEvents(now, 1)
	}
	var eventsPrefer, otherEvents []models.EventCardWithDateSQL
	err = pgxscan.Select(context.Background(), ed.pool, &eventsPrefer,
		`SELECT id, title, description, image, start_date, end_date FROM events
			WHERE category = 'Кино'
			ORDER BY id DESC`)
	/*if recomend.Concert >= recomend.Movie && recomend.Concert >= recomend.Show {
		err := pgxscan.Select(context.Background(), ed.pool, &eventsPrefer,
			`SELECT id, title, description, image, start_date, end_date FROM events
			WHERE category = 'Музей'
			ORDER BY id DESC`)

		if err != nil {
			ed.logger.Warn(err)
			return nil, err
		}
		err = pgxscan.Select(context.Background(), ed.pool, &otherEvents,
			`SELECT id, title, description, image, start_date, end_date FROM events
			WHERE category != 'Музей'
			ORDER BY id DESC`)

		if err != nil {
			ed.logger.Warn(err)
			return nil, err
		}
	}*/
	if recomend.Movie >= recomend.Concert && recomend.Movie >= recomend.Show {
		err := pgxscan.Select(context.Background(), ed.pool, &eventsPrefer,
			`SELECT id, title, description, image, start_date, end_date FROM events
			WHERE category = 'Кино'
			ORDER BY id DESC`)

		if err != nil {
			ed.logger.Warn(err)
			return nil, err
		}
		err = pgxscan.Select(context.Background(), ed.pool, &otherEvents,
			`SELECT id, title, description, image, start_date, end_date FROM events
			WHERE category != 'Кино'
			ORDER BY id DESC`)

		if err != nil {
			ed.logger.Warn(err)
			return nil, err
		}
	}
	if recomend.Show >= recomend.Concert && recomend.Show >= recomend.Movie {
		err := pgxscan.Select(context.Background(), ed.pool, &eventsPrefer,
			`SELECT id, title, description, image, start_date, end_date FROM events
			WHERE category = 'Выставка'
			ORDER BY id DESC`)

		if err != nil {
			ed.logger.Warn(err)
			return nil, err
		}
		err = pgxscan.Select(context.Background(), ed.pool, &otherEvents,
			`SELECT id, title, description, image, start_date, end_date FROM events
			WHERE category != 'Выставка'
			ORDER BY id DESC`)

		if err != nil {
			ed.logger.Warn(err)
			return nil, err
		}
	}

	eventsPrefer = append(eventsPrefer, otherEvents...)
	if len(eventsPrefer) == 0 {
		ed.logger.Debug("no rows in method GetRecomended")
		return []models.EventCardWithDateSQL{}, nil
	}
	var validEvents []models.EventCardWithDateSQL

	for _, elem := range eventsPrefer {
		if elem.EndDate.After(now) {
			validEvents = append(validEvents, elem)
		}
	}

	return validEvents, nil
}
