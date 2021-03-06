package subscription

type Repository interface {
	SubscribeUser(subscriberId uint64, subscribedToId uint64) error
	UnsubscribeUser(subscriberId uint64, subscribedToId uint64) error
	CheckSubscription(subscriberId uint64, subscribedToId uint64) (bool, error)
	CheckEventInList(eventId uint64) (bool, error)
	CheckEventAdded(userId uint64, eventId uint64) (bool, error)
	AddPlanning(userId uint64, eventId uint64) error
	AddVisited(userId uint64, eventId uint64) error
	RemoveEvent(userId uint64, eventId uint64) error
	AddUserEventAction(userId uint64, eventId uint64) error
	AddSubscriptionAction(subscriberId uint64, subscribedToId uint64) error
	RemoveSubscriptionAction(userId uint64, eventId uint64) error
	RemoveUserEventAction(userId uint64, eventId uint64) error
}
