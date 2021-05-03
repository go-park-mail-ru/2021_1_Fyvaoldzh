package subscription

type UseCase interface {
	SubscribeUser(subscriberId uint64, subscribedToId uint64) (bool, string, error)
	UnsubscribeUser(subscriberId uint64, subscribedToId uint64) (bool, string, error)
	AddPlanning(userId uint64, eid uint64) (bool, string, error)
	AddVisited(userId uint64, eid uint64) (bool, string, error)
	RemoveEvent(userId uint64, eid uint64) (bool, string, error)
}