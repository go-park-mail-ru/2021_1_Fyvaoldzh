package usecase

import (
	"github.com/labstack/echo"
	"kudago/application/subscription"
	"net/http"
)

// юзер1 подписывается на юзера2 !1!!

type Subscription struct {
	repo subscription.Repository
}

func NewSubscription(s subscription.Repository) subscription.UseCase {
	return &Subscription{repo: s}
}

// TODO: заменить invalid data на че-нить более говорящее

func (s Subscription) SubscribeUser(uid1 uint64, uid2 uint64) error {
	if uid1 == uid2 || uid1 <= 0 || uid2 <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid data")
	}

	err := s.repo.SubscribeUser(uid1, uid2)

	if err != nil {
		return err
	}

	return nil
}

func (s Subscription) UnsubscribeUser(uid1 uint64, uid2 uint64) error {
	if uid1 == uid2 || uid1 <= 0 || uid2 <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid data")
	}

	return s.repo.UnsubscribeUser(uid1, uid2)
}

func (s Subscription) AddPlanning(uid uint64, eid uint64) error {
	if uid <= 0 || eid <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid data")
	}

	return s.repo.AddPlanning(uid, eid)
}

func (s Subscription) RemovePlanning(uid uint64, eid uint64) error {
	if uid <= 0 || eid <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid data")
	}

	return s.repo.RemovePlanning(uid, eid)
}

func (s Subscription) AddVisited(uid uint64, eid uint64) error {
	if uid <= 0 || eid <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid data")
	}

	return s.repo.AddVisited(uid, eid)
}

func (s Subscription) RemoveVisited(uid uint64, eid uint64) error {
	if uid <= 0 || eid <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid data")
	}

	return s.repo.RemoveVisited(uid, eid)
}

func (s Subscription) GetPlanning(uid uint64) error {
	panic("implement me")
	return nil
}

func (s Subscription) GetVisited(uid uint64) error {
	panic("implement me")
	return nil
}
