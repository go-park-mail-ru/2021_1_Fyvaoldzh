package usecase

import (
	"github.com/labstack/echo"
	"kudago/application/models"
	"kudago/application/subscription"
	"kudago/pkg/logger"
	"net/http"
	"time"
)

type Subscription struct {
	repo   subscription.Repository
	Logger logger.Logger
}

func NewSubscription(subRepo subscription.Repository, logger logger.Logger) subscription.UseCase {
	return &Subscription{repo: subRepo, Logger: logger}
}

func (s Subscription) UpdateEventStatus(userId uint64, eventId uint64) error {
	return s.repo.UpdateEventStatus(userId, eventId)
}

func (s Subscription) IsAddedEvent(userId uint64, eventId uint64) (bool, error) {
	return s.repo.IsAddedEvent(userId, eventId)
}

func (s Subscription) GetFollowers(id uint64, page int) (models.UserCards, error) {
	users, err := s.repo.GetFollowers(id, page)
	if err != nil {
		s.Logger.Warn(err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var userCards models.UserCards
	for _, elem := range users {
		userCards = append(userCards, *models.ConvertUserCard(elem))
	}

	return userCards, nil
}

func (s Subscription) GetSubscriptions(id uint64, page int) (models.UserCards, error) {
	users, err := s.repo.GetFollowers(id, page)
	if err != nil {
		s.Logger.Warn(err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var userCards models.UserCards
	for _, elem := range users {
		userCards = append(userCards, *models.ConvertUserCard(elem))
	}

	return userCards, nil
}

func (s Subscription) GetPlanningEvents(id uint64, page int) (models.EventCards, error) {
	sqlEvents, err := s.repo.GetPlanningEvents(id, page)
	if err != nil {
		s.Logger.Warn(err)
		return models.EventCards{}, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var newEvents []models.EventCard
	for _, elem := range sqlEvents {
		if elem.EndDate.Before(time.Now()) {
			err := s.repo.UpdateEventStatus(id, elem.ID)
			if err != nil {
				s.Logger.Warn(err)
				return models.EventCards{}, err
			}
		} else {
			newEvents = append(newEvents, models.ConvertDateCard(elem))
		}
	}

	return newEvents, nil
}

func (s Subscription) GetVisitedEvents(id uint64, page int) (models.EventCards, error) {
	var events models.EventCards
	sqlEvents, err := s.repo.GetVisitedEvents(id, page)
	if err != nil {
		s.Logger.Warn(err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	for _, elem := range sqlEvents {
		events = append(events, models.ConvertDateCard(elem))
	}

	return events, nil
}
