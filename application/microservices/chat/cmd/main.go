package main

import (
	"kudago/application/microservices/chat/server"
	"kudago/pkg/constants"
	"kudago/pkg/logger"
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	os.Setenv("DB_PASSWORD", "fyvaoldzh")
	os.Setenv("POSTGRE_USER", "postgre")
	os.Setenv("TARANTOOL_USER", "admin")
	lg, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := lg.Sync(); err != nil {
			log.Fatalf(`error '%s' while closing resource`, err)
		}
	}()
	sugar := lg.Sugar()
	zap.NewAtomicLevelAt(zapcore.DebugLevel)
	l := logger.NewLogger(sugar)

	s := server.NewServer(constants.ChatServicePort, &l)
	_ = s.ListenAndServe()
}
