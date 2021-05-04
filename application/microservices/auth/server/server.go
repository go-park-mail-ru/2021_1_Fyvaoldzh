package server

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/tarantool/go-tarantool"
	"google.golang.org/grpc"
	"kudago/application/microservices/auth/proto"
	srepo "kudago/application/microservices/auth/session/repository"
	sessions "kudago/application/microservices/auth/session/usecase"
	urepo "kudago/application/microservices/auth/user/repository"
	users "kudago/application/microservices/auth/user/usecase"
	"kudago/pkg/constants"
	"kudago/pkg/logger"
	"log"
	"net"
)

type Server struct {
	port string
	auth *AuthServer
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

	conn, err := tarantool.Connect(constants.TarantoolAddress, tarantool.Opts{
		User: constants.TarantoolUser,
		Pass: constants.TarantoolPassword,
	})
	if err != nil {
		logger.Fatal(err)
	}

	_, err = conn.Ping()
	if err != nil {
		logger.Fatal(err)
	}

	sRepo := srepo.NewSessionRepository(conn, *logger)
	uRepo := urepo.NewUserDatabase(pool, *logger)
	uUseCase := users.NewUserUseCase(uRepo, *logger)
	s := sessions.NewSessionUseCase(sRepo, uUseCase, *logger)

	return &Server{
		port: port,
		auth: NewAuthServer(s),
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
	proto.RegisterAuthServer(gServer, s.auth)
	log.Println("starting server at " + s.port)
	err = gServer.Serve(listener)

	if err != nil {
		return nil
	}

	return nil
}
