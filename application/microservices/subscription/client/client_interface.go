package client

type ISubscriptionClient interface {
	Subscribe(subscriberId uint64, subscribedToId uint64) error
	Unsubscribe(subscriberId uint64, subscribedToId uint64) error
	AddPlanningEvent(userId uint64, eventId uint64) error
	RemoveEvent(userId uint64, eventId uint64) error
	AddVisitedEvent(userId uint64, eventId uint64) error
	Close()
}