package slogfiber

import (
	"context"
	"net/http"
	"time"

	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Config struct {
	DefaultLevel     slog.Level
	ClientErrorLevel slog.Level
	ServerErrorLevel slog.Level

	WithRequestID bool
}

// New returns a fiber.Handler (middleware) that logs requests using slog.
//
// Requests with errors are logged using slog.Error().
// Requests without errors are logged using slog.Info().
func New(logger *slog.Logger) fiber.Handler {
	return NewWithConfig(logger, Config{
		DefaultLevel:     slog.LevelInfo,
		ClientErrorLevel: slog.LevelWarn,
		ServerErrorLevel: slog.LevelError,

		WithRequestID: true,
	})
}

// NewWithConfig returns a fiber.Handler (middleware) that logs requests using slog.
func NewWithConfig(logger *slog.Logger, config Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Path()
		start := time.Now()
		path := c.Path()

		requestID := uuid.New().String()
		if config.WithRequestID {
			c.Context().SetUserValue("request-id", requestID)
			c.Set("X-Request-ID", requestID)
		}

		err := c.Next()

		end := time.Now()
		latency := end.Sub(start)

		ip := c.Context().RemoteIP().String()
		if len(c.IPs()) > 0 {
			ip = c.IPs()[0]
		}

		attributes := []slog.Attr{
			slog.Time("time", end),
			slog.Duration("latency", latency),
			slog.String("method", string(c.Context().Method())),
			slog.String("host", c.Hostname()),
			slog.String("path", path),
			slog.String("route", c.Route().Path),
			slog.Int("status", c.Response().StatusCode()),
			slog.String("ip", ip),
			slog.String("user-agent", string(c.Context().UserAgent())),
			slog.String("referer", c.Get(fiber.HeaderReferer)),
		}

		if len(c.IPs()) > 0 {
			attributes = append(attributes, slog.Any("x-forwarded-for", c.IPs()))
		}

		if config.WithRequestID {
			attributes = append(attributes, slog.String("request-id", requestID))
		}

		// if err == nil && c.Response().StatusCode() >= http.StatusBadRequest {
		if err == nil {
			err = fiber.NewError(c.Response().StatusCode())
		}

		switch {
		case c.Response().StatusCode() >= http.StatusBadRequest && c.Response().StatusCode() < http.StatusInternalServerError:
			logger.LogAttrs(context.Background(), config.ClientErrorLevel, err.Error(), attributes...)
		case c.Response().StatusCode() >= http.StatusInternalServerError:
			logger.LogAttrs(context.Background(), config.ServerErrorLevel, err.Error(), attributes...)
		default:
			logger.LogAttrs(context.Background(), config.DefaultLevel, "Incoming request", attributes...)
		}

		return err
	}
}

// GetRequestID returns the request identifier
func GetRequestID(c *fiber.Ctx) string {
	requestID, ok := c.Context().UserValue("request-id").(string)
	if !ok {
		return ""
	}

	return requestID
}
