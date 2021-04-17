package logger

import (
	"time"

	"github.com/labstack/echo"
	"go.uber.org/zap"
)

type Logger struct {
	logger *zap.SugaredLogger
}

func NewLogger(logger *zap.SugaredLogger) Logger {
	return Logger{logger: logger}
}

func (log Logger) LogInfo(c echo.Context, start time.Time, request_id string) {
	log.logger.Info(c.Request().URL.Path,
		zap.String("method", c.Request().Method),
		zap.String("remote_addr", c.Request().RemoteAddr),
		zap.String("url", c.Request().URL.Path),
		zap.Duration("work_time", time.Since(start)),
		zap.String("request_id", request_id),
	)
}

func (log Logger) LogWarn(c echo.Context, start time.Time, request_id string, err error) {
	log.logger.Warn(c.Request().URL.Path,
		zap.String("method", c.Request().Method),
		zap.String("remote_addr", c.Request().RemoteAddr),
		zap.String("url", c.Request().URL.Path),
		zap.Duration("work_time", time.Since(start)),
		zap.String("request_id", request_id),
		zap.Errors("error", []error{err}),
	)
}

func (log Logger) LogError(c echo.Context, start time.Time, request_id string, err error) {
	log.logger.Error(c.Request().URL.Path,
		zap.String("method", c.Request().Method),
		zap.String("remote_addr", c.Request().RemoteAddr),
		zap.String("url", c.Request().URL.Path),
		zap.Duration("work_time", time.Since(start)),
		zap.String("request_id", request_id),
		zap.Errors("error", []error{err}),
	)
}

//For init server
func (log Logger) Fatal(err error) {
	log.logger.Fatal(err)
}

func (log Logger) Warn(err error) {
	log.logger.Warn(err)
}

func (log Logger) Debug(str string) {
	log.logger.Debug(str)
}
