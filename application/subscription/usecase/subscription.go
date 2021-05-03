package usecase

import (
	"github.com/labstack/echo"
	"kudago/application/models"
	"kudago/application/subscription"
	"kudago/pkg/logger"
	"net/http"
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

func (s Subscription) GetFollowers(id uint64) (models.UserCards, error) {
	users, err := s.repo.GetFollowers(id)
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

func (s Subscription) GetSubscriptions(id uint64) (models.UserCards, error) {
	users, err := s.repo.GetFollowers(id)
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
