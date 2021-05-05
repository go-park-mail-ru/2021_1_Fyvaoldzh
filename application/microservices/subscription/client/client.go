package client

/*
protoc --go_out=plugins=grpc:. *.proto
*/

import (
	"context"
	"github.com/labstack/echo"
	traceutils "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"kudago/application/microservices/subscription/proto"
	"kudago/pkg/constants"
	"kudago/pkg/logger"
	"net/http"
)

type SubscriptionClient struct {
	client proto.SubscriptionClient
	gConn  *grpc.ClientConn
	logger logger.Logger
}

func NewSubscriptionClient(port string, logger logger.Logger, tracer opentracing.Tracer) (ISubscriptionClient, error) {
	gConn, err := grpc.Dial(
		constants.Localhost+port,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(traceutils.OpenTracingClientInterceptor(tracer)),
	)
	if err != nil {
		logger.Warn(err)
		return nil, err
	}

	return &SubscriptionClient{client: proto.NewSubscriptionClient(gConn), gConn: gConn, logger: logger}, nil
}

func (s *SubscriptionClient) Subscribe(subscriberId uint64, subscribedToId uint64) error {
	users := &proto.Users{
		SubscriberId:   subscriberId,
		SubscribedToId: subscribedToId,
	}

	answer, err := s.client.Subscribe(context.Background(), users)
	if err != nil {
		s.logger.Warn(err)
		return err
	}
	if answer.Flag {
		return echo.NewHTTPError(http.StatusBadRequest, answer.Msg)
	}

	return nil
}

func (s *SubscriptionClient) Unsubscribe(subscriberId uint64, subscribedToId uint64) error {
	users := &proto.Users{
		SubscriberId:   subscriberId,
		SubscribedToId: subscribedToId,
	}

	answer, err := s.client.Unsubscribe(context.Background(), users)
	if err != nil {
		s.logger.Warn(err)
		return err
	}
	if answer.Flag {
		return echo.NewHTTPError(http.StatusBadRequest, answer.Msg)
	}

	return nil
}

func (s *SubscriptionClient) AddPlanningEvent(userId uint64, eventId uint64) error {
	userEvent := &proto.UserEvent{
		UserId:  userId,
		EventId: eventId,
	}

	answer, err := s.client.AddPlanningEvent(context.Background(), userEvent)
	if err != nil {
		s.logger.Warn(err)
		return err
	}
	if answer.Flag {
		return echo.NewHTTPError(http.StatusBadRequest, answer.Msg)
	}

	return nil
}

func (s *SubscriptionClient) RemoveEvent(userId uint64, eventId uint64) error {
	userEvent := &proto.UserEvent{
		UserId:  userId,
		EventId: eventId,
	}

	answer, err := s.client.RemoveEvent(context.Background(), userEvent)
	if err != nil {
		s.logger.Warn(err)
		return err
	}
	if answer.Flag {
		return echo.NewHTTPError(http.StatusBadRequest, answer.Msg)
	}

	return nil
}

func (s *SubscriptionClient) AddVisitedEvent(userId uint64, eventId uint64) error {
	userEvent := &proto.UserEvent{
		UserId:  userId,
		EventId: eventId,
	}

	answer, err := s.client.AddVisitedEvent(context.Background(), userEvent)
	if err != nil {
		s.logger.Warn(err)
		return err
	}
	if answer.Flag {
		return echo.NewHTTPError(http.StatusBadRequest, answer.Msg)
	}

	return nil
}

func (s *SubscriptionClient) Close() {
	if err := s.gConn.Close(); err != nil {
		s.logger.Warn(err)
	}
}
