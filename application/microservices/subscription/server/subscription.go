package server

import (
	"context"
	"kudago/application/microservices/subscription/proto"
	"kudago/application/microservices/subscription/subscription"
)

type SubscriptionServer struct {
	usecase subscription.UseCase
}

func NewSubscriptionServer(usecase subscription.UseCase) *SubscriptionServer{
	return &SubscriptionServer{usecase: usecase}
}

func (s *SubscriptionServer) Subscribe(ctx context.Context, users *proto.Users) (*proto.Nothing, error) {
	subscriberId := users.SubscriberId
	subscribedToId := users.SubscribedToId

	err := s.usecase.SubscribeUser(subscriberId, subscribedToId)
	if err != nil {
		return nil, err
	}

	return &proto.Nothing{}, nil
}

func (s *SubscriptionServer) Unsubscribe(ctx context.Context, users *proto.Users) (*proto.Nothing, error) {
	subscriberId := users.SubscriberId
	subscribedToId := users.SubscribedToId

	err := s.usecase.UnsubscribeUser(subscriberId, subscribedToId)
	if err != nil {
		return nil, err
	}

	return &proto.Nothing{}, nil
}

func (s *SubscriptionServer) AddPlanningEvent(ctx context.Context, userEvent *proto.UserEvent) (*proto.Nothing, error) {
	userId := userEvent.UserId
	eventId := userEvent.EventId

	err := s.usecase.AddPlanning(userId, eventId)
	if err != nil {
		return nil, err
	}

	return &proto.Nothing{}, nil
}

func (s *SubscriptionServer) AddVisitedEvent(ctx context.Context, userEvent *proto.UserEvent) (*proto.Nothing, error) {
	userId := userEvent.UserId
	eventId := userEvent.EventId

	err := s.usecase.AddVisited(userId, eventId)
	if err != nil {
		return nil, err
	}

	return &proto.Nothing{}, nil
}

func (s *SubscriptionServer) RemoveEvent(ctx context.Context, userEvent *proto.UserEvent) (*proto.Nothing, error) {
	userId := userEvent.UserId
	eventId := userEvent.EventId

	err := s.usecase.RemoveEvent(userId, eventId)
	if err != nil {
		return nil, err
	}

	return &proto.Nothing{}, nil
}
