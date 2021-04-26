package server

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	srepo "kudago/application/microservices/subscription/subscription/repository"
	subscriptions "kudago/application/microservices/subscription/subscription/usecase"

	"kudago/application/microservices/subscription/proto"
	"kudago/pkg/constants"
	"kudago/pkg/logger"
	"log"
	"net"
)

type Server struct {
	port string
	ss *SubscriptionServer
}

func NewServer(port string, logger *logger.Logger) *Server {
	pool, err := pgxpool.Connect(context.Background(), constants.DBConnect)
	if err != nil {
		logger.Fatal(err)
	}
	err = pool.Ping(context.Background())
	if err != nil {
		logger.Fatal(err)
	}

	sRepo := srepo.NewSubscriptionDatabase(pool, *logger)
	sUseCase := subscriptions.NewSubscription(sRepo, *logger)

	return &Server{
		port: port,
		ss: NewSubscriptionServer(sUseCase),
	}
}

func (s *Server) ListenAndServe() error {
	gServer := grpc.NewServer()

	listener, err := net.Listen("tcp", s.port)
	defer listener.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	proto.RegisterSubscriptionServer(gServer, s.ss)
	log.Println("starting server at " + s.port)
	err = gServer.Serve(listener)

	if err != nil {
		return nil
	}

	return nil
}