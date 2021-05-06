package client

type ISubscriptionClient interface {
	Subscribe(subscriberId uint64, subscribedToId uint64) (error, int)
	Unsubscribe(subscriberId uint64, subscribedToId uint64) (error, int)
	AddPlanningEvent(userId uint64, eventId uint64) (error, int)
	RemoveEvent(userId uint64, eventId uint64) (error, int)
	AddVisitedEvent(userId uint64, eventId uint64) (error, int)
	Close()
}
