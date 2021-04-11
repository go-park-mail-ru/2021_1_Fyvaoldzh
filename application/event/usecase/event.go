package usecase

import (
	"fmt"
	"io"
	"io/ioutil"
	"kudago/application/event"
	"kudago/application/models"
	"kudago/pkg/constants"
	"kudago/pkg/generator"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

type Event struct {
	repo event.Repository
}

func NewEvent(e event.Repository) event.UseCase {
	return &Event{repo: e}
}

func (e Event) GetAllEvents(page int) (models.EventCards, error) {
	sqlEvents, err := e.repo.GetAllEvents(time.Now())
	if err != nil {
		return models.EventCards{}, err
	}

	var pageEvents models.EventCards

	for i := (page - 1) * 6; i < page*6; i++ {
		if i >= len(sqlEvents) {
			break
		}
		pageEvents = append(pageEvents, models.ConvertDateCard(sqlEvents[i]))
	}
	if len(pageEvents) == 0 {
		return models.EventCards{}, nil
	}

	return pageEvents, nil
}

func (e Event) GetOneEvent(eventId uint64) (models.Event, error) {
	ev, err := e.repo.GetOneEventByID(eventId)
	if err != nil {
		return models.Event{}, err
	}

	jsonEvent := models.ConvertEvent(ev)

	tags, err := e.repo.GetTags(eventId)

	if err != nil {
		return models.Event{}, err
	}

	jsonEvent.Tags = tags

	return jsonEvent, nil
}

func (e Event) Delete(eventId uint64) error {
	return e.repo.DeleteById(eventId)
}

func (e Event) CreateNewEvent(newEvent *models.Event) error {
	// TODO где-то здесь должна быть проверка на поля
	return e.repo.AddEvent(newEvent)
}

func (e Event) SaveImage(eventId uint64, img *multipart.FileHeader) error {
	src, err := img.Open()
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer src.Close()

	fileName := constants.EventsPicDir + fmt.Sprint(eventId) + generator.RandStringRunes(6) + img.Filename

	dst, err := os.Create(fileName)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return e.repo.UpdateEventAvatar(eventId, fileName)
}

func (e Event) GetEventsByCategory(typeEvent string, page int) (models.EventCards, error) {
	sqlEvents, err := e.repo.GetEventsByCategory(typeEvent, time.Now())
	if err != nil {
		return models.EventCards{}, err
	}

	if len(sqlEvents) == 0 {
		return models.EventCards{}, err
	}

	var pageEvents models.EventCards

	for i := (page - 1) * 6; i < page*6; i++ {
		if i >= len(sqlEvents) {
			break
		}
		pageEvents = append(pageEvents, models.ConvertDateCard(sqlEvents[i]))
	}
	if len(pageEvents) == 0 {
		return models.EventCards{}, nil
	}

	return pageEvents, nil
}

func (e Event) GetImage(eventId uint64) ([]byte, error) {
	ev, err := e.repo.GetOneEventByID(eventId)
	if err != nil {
		return []byte{}, err
	}

	file, err := ioutil.ReadFile(ev.Image.String)
	if err != nil {
		log.Println("Cannot open file: " + ev.Image.String)
		return []byte{}, err
	}

	return file, nil
}

func (e Event) FindEvents(str string, category string, page int) (models.EventCards, error) {
	str = strings.ToLower(str)

	var sqlEvents []models.EventCardWithDateSQL
	var err error
	if category == "" {
		sqlEvents, err = e.repo.FindEvents(str, time.Now())
	} else {
		sqlEvents, err = e.repo.CategorySearch(str, category, time.Now())
	}
	if err != nil {
		return models.EventCards{}, err
	}

	if len(sqlEvents) == 0 {
		return models.EventCards{}, err
	}

	var pageEvents models.EventCards

	for i := (page - 1) * 6; i < page*6; i++ {
		if i >= len(sqlEvents) {
			break
		}
		pageEvents = append(pageEvents, models.ConvertDateCard(sqlEvents[i]))
	}
	if len(pageEvents) == 0 {
		return models.EventCards{}, nil
	}

	return pageEvents, nil
}

func (e Event) RecomendSystem(uid uint64, category string) error {
	if err := e.repo.RecomendSystem(uid, category); err != nil {
		time.Sleep(1 * time.Second)
		if err := e.repo.RecomendSystem(uid, category); err != nil {
			return errors.New("cannot add record in user_prefer")
		}
	}
	return nil
}

func (e Event) GetRecomended(uid uint64, page int) (models.EventCards, error) {
	sqlEvents, err := e.repo.GetRecomended(uid, time.Now())
	if err != nil {
		return models.EventCards{}, err
	}

	if len(sqlEvents) == 0 {
		return models.EventCards{}, err
	}

	var pageEvents models.EventCards

	for i := (page - 1) * 6; i < page*6; i++ {
		if i >= len(sqlEvents) {
			break
		}
		pageEvents = append(pageEvents, models.ConvertDateCard(sqlEvents[i]))
	}
	if len(pageEvents) == 0 {
		return models.EventCards{}, nil
	}

	return pageEvents, nil
}
