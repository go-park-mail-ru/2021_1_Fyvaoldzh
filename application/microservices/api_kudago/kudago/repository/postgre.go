package kudago_repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	"io"
	"kudago/application/microservices/api_kudago/kudago"
	"kudago/application/models"
	"kudago/pkg/logger"
	"net/http"
	"os"
)

type KudagoDatabase struct {
	pool   *pgxpool.Pool
	logger logger.Logger
}

func NewKudagoDatabase(conn *pgxpool.Pool, logger logger.Logger) kudago.Repository {
	return &KudagoDatabase{pool: conn, logger: logger}
}

func (kd KudagoDatabase) AddEvent(newEvent models.Event) (uint64, error) {
	var id uint64
	err := kd.pool.QueryRow(context.Background(),
		`INSERT INTO events (id, title, place, subway, street, description, category, start_date, end_date, image, latitude, longitude) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)  RETURNING id`,
		newEvent.ID, newEvent.Title, newEvent.Place, newEvent.Subway, newEvent.Street, newEvent.Description,
		newEvent.Category, newEvent.StartDate, newEvent.EndDate, newEvent.Image, newEvent.Latitude, newEvent.Longitude).Scan(&id)

	if err != nil {
		kd.logger.Warn(err)
		return 0, err
	}

	return id, nil
}

func (kd KudagoDatabase) AddEventTag(eventId uint64, tagId uint32) error {
	_, err := kd.pool.Exec(context.Background(),
		`INSERT INTO event_tag (event_id, tag_id) VALUES ($1, $2)`, eventId, tagId)

	if err != nil {
		kd.logger.Warn(err)
		return err
	}

	return nil
}

func (kd KudagoDatabase) IsExistingTag(tagName string) (bool, uint32, error) {
	var id uint32
	err := kd.pool.
		QueryRow(context.Background(),
			`SELECT id FROM tags WHERE name = $1`, tagName).Scan(&id)

	if errors.As(err, &pgx.ErrNoRows) {
		return false, 0, nil
	}
	if err != nil {
		kd.logger.Warn(err)
		return false, 0, err
	}

	return true, id, nil
}

func (kd KudagoDatabase) IsExistingEvent(id uint64) (bool, uint64, error) {
	var eventId uint64
	err := kd.pool.
		QueryRow(context.Background(),
			`SELECT id FROM events WHERE id = $1`, id).Scan(&eventId)

	if errors.As(err, &pgx.ErrNoRows) {
		return false, 0, nil
	}
	if err != nil {
		kd.logger.Warn(err)
		return false, 0, err
	}

	return true, eventId, nil
}

func (kd KudagoDatabase) AddTag(name string) (uint32, error) {
	var id uint32
	err := kd.pool.QueryRow(context.Background(),
		`INSERT INTO tags (name) VALUES ($1)  RETURNING id`, name).Scan(&id)

	if err != nil {
		kd.logger.Warn(err)
		return 0, err
	}

	return id, nil
}

func (kd KudagoDatabase) AddImage(URL string, fileName string) error {
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("received non 200 response code")
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
