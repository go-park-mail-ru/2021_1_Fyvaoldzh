package kudago_usecase

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/microcosm-cc/bluemonday"
	"kudago/application/microservices/api_kudago/kudago"
	"kudago/application/models"
	"kudago/pkg/constants"
	"kudago/pkg/generator"
	"kudago/pkg/logger"
	"strconv"
	"strings"
	"time"
)

type KudagoUsecase struct {
	repo   kudago.Repository
	logger logger.Logger
}

func NewKudagoUsecase(r kudago.Repository, logger logger.Logger) kudago.Usecase {
	return &KudagoUsecase{repo: r, logger: logger}
}

func (k KudagoUsecase) AddEvent(elem models.Elem, place models.Place) (bool, error) {
	flag, eventId, err := k.repo.IsExistingEvent(elem.Id)
	if err != nil {
		return false, err
	}
	if flag {
		return false, nil
	}


	newEvent, err := k.ConvertToNewEvent(elem, place)
	if err != nil {
		return false, err
	}
	if newEvent.ID == uint64(0) {
		return false, nil
	}

	eventId, err = k.repo.AddEvent(newEvent)
	if err != nil {
		return false, err
	}
	for _, tag := range elem.Tags {
		if len(tag) <= 60 {
			if len(tag) != len([]rune(tag)){
				tag = strings.ToUpper(tag[:2]) + tag[2:]
			} else {
				tag = strings.ToUpper(tag[:1]) + tag[1:]
			}
			flag, tagId, err := k.repo.IsExistingTag(tag)
			if err != nil {
				return false, err
			}
			if !flag {
				tagId, err = k.AddTag(tag)
				if err != nil {
					return false, err
				}
			}
			err = k.repo.AddEventTag(eventId, tagId)
			if err != nil {
				return false, err
			}
		}

	}

	return true, nil
}

func (k KudagoUsecase) AddTag(name string) (uint32, error) {
	return k.repo.AddTag(name)
}

func (k KudagoUsecase) ConvertToNewEvent(elem models.Elem, place models.Place) (models.Event, error) {
	var newEvent models.Event

	sanitizer := bluemonday.UGCPolicy()
	newEvent.ID = elem.Id
	newEvent.Place = place.Title
	firstWord := strings.Split(elem.Title, " ")[0]
	if len(firstWord) != len([]rune(firstWord)){
		firstWord = strings.ToUpper(firstWord[:2])
		elem.Title = firstWord + elem.Title[2:]
	} else {
		firstWord = strings.ToUpper(firstWord[:1]) + firstWord[1:]
		elem.Title = firstWord + elem.Title[1:]
	}
	newEvent.Title = elem.Title
	newEvent.Title = strings.ToUpper(newEvent.Title[:2]) + newEvent.Title[2:]
	newEvent.Description = sanitizer.Sanitize(elem.BodyText)
	p := strings.NewReader(newEvent.Description)
	doc, _ := goquery.NewDocumentFromReader(p)
	doc.Find("script").Each(func(i int, el *goquery.Selection) {
		el.Remove()
	})
	newEvent.Description = strings.Trim(doc.Text(), " ")
	newEvent.StartDate = time.Unix(elem.Date[0].Start, 0).Format(constants.DateTimeFormat)
	newEvent.EndDate = time.Unix(elem.Date[0].End, 0).Format(constants.DateTimeFormat)
	newEvent.Subway = place.Subway
	newEvent.Street = place.Address
	newEvent.Latitude = place.Map.Latitude
	newEvent.Longitude = place.Map.Longitude
	switch elem.Categories[0] {
	case "cinema":
		newEvent.Category = "Кино"
	case "education":
		newEvent.Category = "Образование"
	case "entertainment":
		newEvent.Category = "Развлечения"
	case "exhibition":
		newEvent.Category = "Выставка"
	case "festival":
		newEvent.Category = "Фестиваль"
	case "tour":
		newEvent.Category = "Экскурсия"
	case "concert":
		newEvent.Category = "Концерт"
	default:
		return models.Event{ID: uint64(0)}, nil
	}

	fileName := constants.EventsPicDir + strconv.FormatUint(elem.Id, 10) + generator.RandStringRunes(12) +".jpg"
	err := k.repo.AddImage(elem.Images[0].Path, fileName)
	if err != nil {
		return models.Event{}, err
	}
	newEvent.Image = fileName

	return newEvent, nil
}