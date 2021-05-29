package kudago

import "kudago/application/models"

type Repository interface {
	AddTag(name string) (uint32, error)
	AddEvent(newEvent models.Event) (uint64, error)
	AddImage(URL string, fileName string) error
	AddEventTag(eventId uint64, tagId uint32) error
	IsExistingTag(tagName string) (bool, uint32, error)
	IsExistingEvent(id uint64) (bool, uint64, error)

}