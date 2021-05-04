package server

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	erepo "kudago/application/event/repository"
	"kudago/application/microservices/chat/chat/repository"
	"kudago/application/microservices/chat/chat/usecase"
	"kudago/application/microservices/chat/proto"
	srepo "kudago/application/subscription/repository"
	urepo "kudago/application/user/repository"
	"kudago/pkg/constants"
	"net"

	"kudago/pkg/logger"
	"log"
)

type Server struct {
	port string
	ss   *ChatServer
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

	cRepo := repository.NewChatDatabase(pool, *logger)
	sRepo := srepo.NewSubscriptionDatabase(pool, *logger)
	eRepo := erepo.NewEventDatabase(pool, *logger)
	uRepo := urepo.NewUserDatabase(pool, *logger)
	cUseCase := usecase.NewChat(cRepo, sRepo, uRepo, eRepo, *logger)

	return &Server{
		port: port,
		ss:   NewChatServer(cUseCase),
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
	proto.RegisterChatServer(gServer, s.ss)
	log.Println("starting server at " + s.port)
	err = gServer.Serve(listener)

	if err != nil {
		return nil
	}

	return nil
}
