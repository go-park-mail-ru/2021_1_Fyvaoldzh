package main

import (
	"kudago/server"
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Fatalf(`error '%s' while closing resource`, err)
		}
	}()
	sugar := logger.Sugar()
	zap.NewAtomicLevelAt(zapcore.DebugLevel)

	/*conf := zap.Config{
		Encoding:         "console",
		OutputPaths:      ["console"],
		ErrorOutputPaths: ["stderr"],
	}
	if err := conf.Level.
		UnmarshalText([]byte("debug")); err != nil {
		return nil, err
	}*/

	e := server.NewServer(sugar)
	server.ListenAndServe(e)
}
