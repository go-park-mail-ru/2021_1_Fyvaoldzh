package subscriptions

type UseCase interface {
	SubscribeUser(uid1 uint64, uid2 uint64) error
	UnsubscribeUser(uid1 uint64, uid2 uint64) error
	AddPlanning(uid uint64, eid uint64) error
	RemovePlanning(uid uint64, eid uint64) error
	AddVisited(uid uint64, eid uint64) error
	RemoveVisited(uid uint64, eid uint64) error
}