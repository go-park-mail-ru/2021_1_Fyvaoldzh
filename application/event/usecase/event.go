package usecase

import (
	"fmt"
	"io"
	"io/ioutil"
	"kudago/application/event"
	"kudago/application/models"
	"kudago/application/subscription"
	"kudago/pkg/constants"
	"kudago/pkg/generator"
	"kudago/pkg/logger"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

type Event struct {
	repo    event.Repository
	repoSub subscription.Repository
	logger  logger.Logger
}

func NewEvent(e event.Repository, repoSubscription subscription.Repository, logger logger.Logger) event.UseCase {
	return &Event{repo: e, repoSub: repoSubscription, logger: logger}
}

func (e Event) GetNear(coord models.Coordinates, page int) (models.EventCardsWithCoords, error) {
	sqlEvents, err := e.repo.GetNearEvents(time.Now(), coord, page)
	if err != nil {
		e.logger.Warn(err)
		return models.EventCardsWithCoords{}, err
	}

	var pageEvents models.EventCardsWithCoords

	for i := range sqlEvents {
		pageEvents = append(pageEvents, models.ConvertCoordsCard(sqlEvents[i], coord))
	}
	if len(pageEvents) == 0 {
		e.logger.Debug("page" + fmt.Sprint(page) + "is empty")
		return models.EventCardsWithCoords{}, nil
	}

	return pageEvents, nil
}

func (e Event) GetAllEvents(page int) (models.EventCards, error) {
	sqlEvents, err := e.repo.GetAllEvents(time.Now(), page)
	if err != nil {
		e.logger.Warn(err)
		return models.EventCards{}, err
	}

	var pageEvents models.EventCards

	for i := range sqlEvents {
		pageEvents = append(pageEvents, models.ConvertDateCard(sqlEvents[i]))
	}
	if len(pageEvents) == 0 {
		e.logger.Debug("page" + fmt.Sprint(page) + "is empty")
		return models.EventCards{}, nil
	}

	return pageEvents, nil
}

func (e Event) GetOneEvent(eventId uint64) (models.Event, error) {
	ev, err := e.repo.GetOneEventByID(eventId)
	if err != nil {
		e.logger.Warn(err)
		return models.Event{}, err
	}

	jsonEvent := models.ConvertEvent(ev)

	tags, err := e.repo.GetTags(eventId)
	if err != nil {
		e.logger.Warn(err)
		return jsonEvent, err
	}

	jsonEvent.Tags = tags

	followers, err := e.repoSub.GetEventFollowers(eventId)
	if err != nil {
		e.logger.Warn(err)
		return jsonEvent, err
	}

	jsonEvent.Followers = followers

	return jsonEvent, nil
}

func (e Event) GetOneEventName(eventId uint64) (string, error) {
	name, err := e.repo.GetOneEventNameByID(eventId)
	if err != nil {
		e.logger.Warn(err)
		return "", err
	}

	return name, nil
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
		e.logger.Warn(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer src.Close()

	fileName := constants.EventsPicDir + fmt.Sprint(eventId) + generator.RandStringRunes(6) + img.Filename

	dst, err := os.Create(fileName)
	if err != nil {
		e.logger.Warn(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		e.logger.Warn(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return e.repo.UpdateEventAvatar(eventId, fileName)
}

func (e Event) GetEventsByCategory(typeEvent string, page int) (models.EventCards, error) {
	var sqlEvents []models.EventCardWithDateSQL
	var err error
	if typeEvent == "" {
		sqlEvents, err = e.repo.GetAllEvents(time.Now(), page)
	} else {
		sqlEvents, err = e.repo.GetEventsByCategory(typeEvent, time.Now(), page)
	}
	if err != nil {
		e.logger.Warn(err)
		return models.EventCards{}, err
	}

	var pageEvents models.EventCards

	for i := range sqlEvents {
		pageEvents = append(pageEvents, models.ConvertDateCard(sqlEvents[i]))
	}
	if len(pageEvents) == 0 {
		e.logger.Debug("page" + fmt.Sprint(page) + "in category" + typeEvent + "empty")
		return models.EventCards{}, nil
	}

	return pageEvents, nil
}

func (e Event) GetImage(eventId uint64) ([]byte, error) {
	ev, err := e.repo.GetOneEventByID(eventId)
	if err != nil {
		e.logger.Warn(err)
		return []byte{}, err
	}

	file, err := ioutil.ReadFile(ev.Image.String)
	if err != nil {
		e.logger.Warn(errors.New("Cannot open file: " + ev.Image.String))
		return []byte{}, err
	}

	return file, nil
}

func (e Event) FindEvents(str string, category string, page int) (models.EventCards, error) {
	str = strings.ToLower(str)

	var sqlEvents []models.EventCardWithDateSQL
	var err error
	if category == "" {
		sqlEvents, err = e.repo.FindEvents(str, time.Now(), page)
	} else {
		sqlEvents, err = e.repo.CategorySearch(str, category, time.Now(), page)
	}
	if err != nil {
		e.logger.Warn(err)
		return models.EventCards{}, err
	}

	var pageEvents models.EventCards

	for i := range sqlEvents {
		pageEvents = append(pageEvents, models.ConvertDateCard(sqlEvents[i]))
	}

	if len(pageEvents) == 0 {
		e.logger.Debug("empty result for method FindEvents")
		return models.EventCards{}, nil
	}

	return pageEvents, nil
}

func (e Event) RecomendSystem(uid uint64, category string) error {
	if err := e.repo.RecomendSystem(uid, category); err != nil {
		time.Sleep(1 * time.Second)
		if err := e.repo.RecomendSystem(uid, category); err != nil {
			e.logger.Warn(err)
			return errors.New("cannot add record in user_prefer")
		}
	}
	return nil
}

func (e Event) GetRecommended(uid uint64, page int) (models.EventCards, error) {
	sqlEvents, err := e.repo.GetRecommended(uid, time.Now(), page)
	if err != nil {
		e.logger.Warn(err)
		return models.EventCards{}, err
	}

	var pageEvents models.EventCards

	for i := range sqlEvents {
		pageEvents = append(pageEvents, models.ConvertDateCard(sqlEvents[i]))
	}
	if len(pageEvents) == 0 {
		e.logger.Debug("empty result for method GetRecomended")
		return models.EventCards{}, nil
	}

	return pageEvents, nil
}
