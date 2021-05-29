package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"kudago/application/microservices/subscription/server"
	"kudago/pkg/constants"
	"kudago/pkg/logger"
	"log"
)

func main() {
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

	s := server.NewServer(constants.SubscriptionServicePort, &l)
	_ = s.ListenAndServe()
}
