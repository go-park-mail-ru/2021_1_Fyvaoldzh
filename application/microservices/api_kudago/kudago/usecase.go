package kudago

import "kudago/application/models"

type Usecase interface {
	AddEvent(elem models.Elem, place models.Place) (bool, error)
	AddTag(name string) (uint32, error)
	ConvertToNewEvent(elem models.Elem, place models.Place) (models.Event, error)
}