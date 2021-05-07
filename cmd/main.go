package main

import (
	"kudago/application/server"
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

	/*conf := zap.Config{
		Encoding:         "console",
		OutputPaths:      ["console"],
		ErrorOutputPaths: ["stderr"],
	}
	if err := conf.Level.
		UnmarshalText([]byte("debug")); err != nil {
		return nil, err
	}*/

	s := server.NewServer(sugar)
	s.ListenAndServe()
}
