package subscription

type UseCase interface {
	UpdateEventStatus(userId uint64, eventId uint64) error
	IsAddedEvent(userId uint64, eventId uint64) (bool, error)
}
