package server

import (
	"context"
	"kudago/application/microservices/subscription/proto"
	"kudago/application/microservices/subscription/subscription"
)

type SubscriptionServer struct {
	usecase subscription.UseCase
}

func NewSubscriptionServer(usecase subscription.UseCase) *SubscriptionServer {
	return &SubscriptionServer{usecase: usecase}
}

func (s *SubscriptionServer) Subscribe(_ context.Context, users *proto.Users) (*proto.SubscriptionAnswer, error) {
	subscriberId := users.SubscriberId
	subscribedToId := users.SubscribedToId

	flag, msg, err := s.usecase.SubscribeUser(subscriberId, subscribedToId)
	if err != nil {
		return nil, err
	}

	return &proto.SubscriptionAnswer{Flag: flag, Msg: msg}, nil
}

func (s *SubscriptionServer) Unsubscribe(_ context.Context, users *proto.Users) (*proto.SubscriptionAnswer, error) {
	subscriberId := users.SubscriberId
	subscribedToId := users.SubscribedToId

	flag, msg, err := s.usecase.UnsubscribeUser(subscriberId, subscribedToId)
	if err != nil {
		return nil, err
	}

	return &proto.SubscriptionAnswer{Flag: flag, Msg: msg}, nil
}

func (s *SubscriptionServer) AddPlanningEvent(_ context.Context, userEvent *proto.UserEvent) (*proto.SubscriptionAnswer, error) {
	userId := userEvent.UserId
	eventId := userEvent.EventId

	flag, msg, err := s.usecase.AddPlanning(userId, eventId)
	if err != nil {
		return nil, err
	}

	return &proto.SubscriptionAnswer{Flag: flag, Msg: msg}, nil
}

func (s *SubscriptionServer) AddVisitedEvent(_ context.Context, userEvent *proto.UserEvent) (*proto.SubscriptionAnswer, error) {
	userId := userEvent.UserId
	eventId := userEvent.EventId

	flag, msg, err := s.usecase.AddVisited(userId, eventId)
	if err != nil {
		return nil, err
	}

	return &proto.SubscriptionAnswer{Flag: flag, Msg: msg}, nil
}

func (s *SubscriptionServer) RemoveEvent(_ context.Context, userEvent *proto.UserEvent) (*proto.SubscriptionAnswer, error) {
	userId := userEvent.UserId
	eventId := userEvent.EventId

	flag, msg, err := s.usecase.RemoveEvent(userId, eventId)
	if err != nil {
		return nil, err
	}

	return &proto.SubscriptionAnswer{Flag: flag, Msg: msg}, nil
}
