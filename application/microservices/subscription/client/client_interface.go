package client

type ISubClient interface {
	Subscribe(subscriberId uint64, subscribedToId uint64) error
	Unsubscribe(subscriberId uint64, subscribedToId uint64) error
	AddPlanningEvent(userId uint64, eventId uint64) error

}