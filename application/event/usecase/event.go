package usecase

import (
	"kudago/application/event"
	"kudago/application/models"
)

type Event struct {
	repo event.Repository
}

func NewEvent(e event.Repository) event.UseCase {
	return &Event{repo: e}
}

func (e Event) GetAllEvents() (models.EventCards, error) {
	sqlEvents, err := e.repo.GetAllEvents()
	if err != nil {
		return models.EventCards{}, err
	}

	var events models.EventCards
	for _, elem := range sqlEvents {
		events = append(events, models.ConvertCard(elem))
	}

	return events, nil
}

func (e Event) GetOneEvent(eventId uint64) (models.Event, error) {
	ev, err := e.repo.GetOneEventByID(eventId)
	if err != nil {
		return models.Event{}, err
	}

	jsonEvent := models.ConvertEvent(ev)

	massiv, err := e.repo.GetCategoryTags(eventId)
	if err != nil {
		return models.Event{}, err
	}

	var desc models.CategoryTag
	desc.ID = massiv[0].CategoryID
	desc.Name = massiv[0].CategoryName

	jsonEvent.TypeEvent = append(jsonEvent.TypeEvent, desc)

	for _, elem := range massiv {
		desc.ID = elem.TagID
		desc.Name = elem.TagName
		jsonEvent.TypeEvent = append(jsonEvent.TypeEvent, desc)
	}

	return jsonEvent, nil
}

