package subscription

type Repository interface {
	SubscribeUser(subscriberId uint64, subscribedToId uint64) error
	UnsubscribeUser(subscriberId uint64, subscribedToId uint64) error
	AddPlanning(userId uint64, eventId uint64) error
	AddVisited(userId uint64, eventId uint64) error
	RemoveEvent(userId uint64, eventId uint64) error
}
