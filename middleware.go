package slogfiber

import (
	"context"
	"net/http"
	"strings"
	"time"

	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	customAttributesCtxKey = "slog-fiber.custom-attributes"
)

var (
	HiddenRequestHeaders = map[string]struct{}{
		"authorization": {},
		"cookie":        {},
		"set-cookie":    {},
		"x-auth-token":  {},
		"x-csrf-token":  {},
		"x-xsrf-token":  {},
	}
	HiddenResponseHeaders = map[string]struct{}{
		"set-cookie": {},
	}
)

type Config struct {
	DefaultLevel     slog.Level
	ClientErrorLevel slog.Level
	ServerErrorLevel slog.Level

	WithRequestID      bool
	WithRequestBody    bool
	WithRequestHeader  bool
	WithResponseBody   bool
	WithResponseHeader bool

	Filters []Filter
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

		WithRequestID:      true,
		WithRequestBody:    false,
		WithRequestHeader:  false,
		WithResponseBody:   false,
		WithResponseHeader: false,

		Filters: []Filter{},
	})
}

// NewWithFilters returns a fiber.Handler (middleware) that logs requests using slog.
//
// Requests with errors are logged using slog.Error().
// Requests without errors are logged using slog.Info().
func NewWithFilters(logger *slog.Logger, filters ...Filter) fiber.Handler {
	return NewWithConfig(logger, Config{
		DefaultLevel:     slog.LevelInfo,
		ClientErrorLevel: slog.LevelWarn,
		ServerErrorLevel: slog.LevelError,

		WithRequestID:      true,
		WithRequestBody:    false,
		WithRequestHeader:  false,
		WithResponseBody:   false,
		WithResponseHeader: false,

		Filters: filters,
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

		// request
		if config.WithRequestBody {
			attributes = append(attributes, slog.Group("request", slog.String("body", string(c.Body()))))
		}
		if config.WithRequestHeader {
			for k, v := range c.GetReqHeaders() {
				if _, found := HiddenRequestHeaders[strings.ToLower(k)]; found {
					continue
				}
				attributes = append(attributes, slog.Group("request", slog.Group("header", slog.Any(k, v))))
			}
		}

		// response
		if config.WithResponseBody {
			attributes = append(attributes, slog.Group("response", slog.String("body", string(c.Response().Body()))))
		}
		if config.WithResponseHeader {
			for k, v := range c.GetRespHeaders() {
				if _, found := HiddenResponseHeaders[strings.ToLower(k)]; found {
					continue
				}
				attributes = append(attributes, slog.Group("response", slog.Group("header", slog.Any(k, v))))
			}
		}

		// custom context values
		if v := c.Context().UserValue(customAttributesCtxKey); v != nil {
			switch attrs := v.(type) {
			case []slog.Attr:
				attributes = append(attributes, attrs...)
			}
		}

		logErr := err
		if logErr == nil {
			logErr = fiber.NewError(c.Response().StatusCode())
		}

		for _, filter := range config.Filters {
			if !filter(c) {
				return err
			}
		}

		switch {
		case c.Response().StatusCode() >= http.StatusBadRequest && c.Response().StatusCode() < http.StatusInternalServerError:
			logger.LogAttrs(context.Background(), config.ClientErrorLevel, logErr.Error(), attributes...)
		case c.Response().StatusCode() >= http.StatusInternalServerError:
			logger.LogAttrs(context.Background(), config.ServerErrorLevel, logErr.Error(), attributes...)
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

func AddCustomAttributes(c *fiber.Ctx, attr slog.Attr) {
	v := c.Context().UserValue(customAttributesCtxKey)
	if v == nil {
		c.Context().SetUserValue(customAttributesCtxKey, []slog.Attr{attr})
		return
	}

	switch attrs := v.(type) {
	case []slog.Attr:
		c.Context().SetUserValue(customAttributesCtxKey, append(attrs, attr))
	}
}
