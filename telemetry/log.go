package telemetry

import (
	"io"
	"os"
	"time"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/trace"
)

var globalLogger zerolog.Logger

func GetLogger() *zerolog.Logger {
	return &globalLogger
}

func InitLogger() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	globalLogger = log.Output(
		zerolog.ConsoleWriter{
			Out:     os.Stdout,
			NoColor: false,
		},
	)
}

func LoggerMiddleware() gin.HandlerFunc {
	return logger.SetLogger(
		logger.WithLogger(func(c *gin.Context, out io.Writer, latency time.Duration) zerolog.Logger {
			return zerolog.New(out).With().
				Str("path", c.Request.URL.Path).
				Str("method", c.Request.Method).
				Str("agent", c.Request.UserAgent()).
				// Int("status", c.Request.Response.StatusCode).
				Dur("duration", latency).
				Str("traceID", trace.SpanContextFromContext(c.Request.Context()).TraceID().String()).
				Logger()
		}),
	)
}
