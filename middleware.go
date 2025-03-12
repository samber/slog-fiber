package slogfiber

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel/trace"
)

type customAttributesCtxKeyType struct{}

var customAttributesCtxKey = customAttributesCtxKeyType{}

var (
	TraceIDKey   = "trace_id"
	SpanIDKey    = "span_id"
	RequestIDKey = "id"

	RequestBodyMaxSize  = 64 * 1024 // 64KB
	ResponseBodyMaxSize = 64 * 1024 // 64KB

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

	// Formatted with http.CanonicalHeaderKey
	RequestIDHeaderKey = "X-Request-Id"
)

type Config struct {
	DefaultLevel     slog.Level
	ClientErrorLevel slog.Level
	ServerErrorLevel slog.Level

	WithUserAgent      bool
	WithRequestID      bool
	WithRequestBody    bool
	WithRequestHeader  bool
	WithResponseBody   bool
	WithResponseHeader bool
	WithSpanID         bool
	WithTraceID        bool

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

		WithUserAgent:      false,
		WithRequestID:      true,
		WithRequestBody:    false,
		WithRequestHeader:  false,
		WithResponseBody:   false,
		WithResponseHeader: false,
		WithSpanID:         false,
		WithTraceID:        false,

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

		WithUserAgent:      false,
		WithRequestID:      true,
		WithRequestBody:    false,
		WithRequestHeader:  false,
		WithResponseBody:   false,
		WithResponseHeader: false,
		WithSpanID:         false,
		WithTraceID:        false,

		Filters: filters,
	})
}

// NewWithConfig returns a fiber.Handler (middleware) that logs requests using slog.
func NewWithConfig(logger *slog.Logger, config Config) fiber.Handler {
	var (
		once       sync.Once
		errHandler fiber.ErrorHandler
	)

	return func(c *fiber.Ctx) error {
		once.Do(func() {
			errHandler = c.App().ErrorHandler
		})

		start := time.Now()
		path := c.Path()
		query := string(c.Request().URI().QueryString())

		requestID := c.Get(RequestIDHeaderKey)
		if config.WithRequestID {
			if requestID == "" {
				requestID = uuid.New().String()
			}
			c.Context().SetUserValue("request-id", requestID)
			c.Set("X-Request-ID", requestID)
		}

		err := c.Next()
		if err != nil {
			if err = errHandler(c, err); err != nil {
				_ = c.SendStatus(fiber.StatusInternalServerError) //nolint:errcheck
			}
		}

		// Pass thru filters and skip early the code below, to prevent unnecessary processing.
		for _, filter := range config.Filters {
			if !filter(c) {
				return err
			}
		}

		status := c.Response().StatusCode()
		method := c.Context().Method()
		host := c.Hostname()
		params := c.AllParams()
		route := c.Route().Path
		end := time.Now()
		latency := end.Sub(start)
		userAgent := c.Context().UserAgent()
		referer := c.Get(fiber.HeaderReferer)

		ip := c.Context().RemoteIP().String()
		if len(c.IPs()) > 0 {
			ip = c.IPs()[0]
		}

		baseAttributes := []slog.Attr{}

		requestAttributes := []slog.Attr{
			slog.Time("time", start.UTC()),
			slog.String("method", string(method)),
			slog.String("host", host),
			slog.String("path", path),
			slog.String("query", query),
			slog.Any("params", params),
			slog.String("route", route),
			slog.String("ip", ip),
			slog.Any("x-forwarded-for", c.IPs()),
			slog.String("referer", referer),
		}

		responseAttributes := []slog.Attr{
			slog.Time("time", end.UTC()),
			slog.Duration("latency", latency),
			slog.Int("status", status),
		}

		if config.WithRequestID {
			baseAttributes = append(baseAttributes, slog.String(RequestIDKey, requestID))
		}

		// otel
		baseAttributes = append(baseAttributes, extractTraceSpanID(c.UserContext(), config.WithTraceID, config.WithSpanID)...)

		// request body
		requestAttributes = append(requestAttributes, slog.Int("length", len((c.Body()))))
		if config.WithRequestBody {
			body := c.Body()
			if len(body) > RequestBodyMaxSize {
				body = body[:RequestBodyMaxSize]
			}
			requestAttributes = append(requestAttributes, slog.String("body", string(body)))
		}

		// request headers
		if config.WithRequestHeader {
			kv := []any{}

			for k, v := range c.GetReqHeaders() {
				if _, found := HiddenRequestHeaders[strings.ToLower(k)]; found {
					continue
				}
				kv = append(kv, slog.Any(k, v))
			}

			requestAttributes = append(requestAttributes, slog.Group("header", kv...))
		}

		if config.WithUserAgent {
			requestAttributes = append(requestAttributes, slog.String("user-agent", string(userAgent)))
		}

		// response body
		responseAttributes = append(responseAttributes, slog.Int("length", len(c.Response().Body())))
		if config.WithResponseBody {
			body := c.Response().Body()
			if len(body) > ResponseBodyMaxSize {
				body = body[:ResponseBodyMaxSize]
			}
			responseAttributes = append(responseAttributes, slog.String("body", string(body)))
		}

		// response headers
		if config.WithResponseHeader {
			kv := []any{}

			for k, v := range c.GetRespHeaders() {
				if _, found := HiddenResponseHeaders[strings.ToLower(k)]; found {
					continue
				}
				kv = append(kv, slog.Any(k, v))
			}

			responseAttributes = append(responseAttributes, slog.Group("header", kv...))
		}

		attributes := append(
			[]slog.Attr{
				{
					Key:   "request",
					Value: slog.GroupValue(requestAttributes...),
				},
				{
					Key:   "response",
					Value: slog.GroupValue(responseAttributes...),
				},
			},
			baseAttributes...,
		)

		// custom context values
		if v := c.Context().UserValue(customAttributesCtxKey); v != nil {
			switch attrs := v.(type) {
			case []slog.Attr:
				attributes = append(attributes, attrs...)
			}
		}

		logErr := err
		if logErr == nil {
			logErr = fiber.NewError(status)
		}

		level := config.DefaultLevel
		msg := "Incoming request"
		if status >= http.StatusInternalServerError {
			level = config.ServerErrorLevel
			msg = logErr.Error()
			if msg == "" {
				msg = fmt.Sprintf("HTTP error: %d %s", status, strings.ToLower(http.StatusText(status)))
			}
		} else if status >= http.StatusBadRequest && status < http.StatusInternalServerError {
			level = config.ClientErrorLevel
			msg = logErr.Error()
			if msg == "" {
				msg = fmt.Sprintf("HTTP error: %d %s", status, strings.ToLower(http.StatusText(status)))
			}
		}

		logger.LogAttrs(c.UserContext(), level, msg, attributes...)

		return err
	}
}

// GetRequestID returns the request identifier.
func GetRequestID(c *fiber.Ctx) string {
	return GetRequestIDFromContext(c.Context())
}

// GetRequestIDFromContext returns the request identifier from the context.
func GetRequestIDFromContext(ctx *fasthttp.RequestCtx) string {
	requestID, ok := ctx.UserValue("request-id").(string)
	if !ok {
		return ""
	}

	return requestID
}

// AddCustomAttributes adds custom attributes to the request context.
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

func extractTraceSpanID(ctx context.Context, withTraceID bool, withSpanID bool) []slog.Attr {
	if !(withTraceID || withSpanID) {
		return []slog.Attr{}
	}

	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return []slog.Attr{}
	}

	attrs := []slog.Attr{}
	spanCtx := span.SpanContext()

	if withTraceID && spanCtx.HasTraceID() {
		traceID := trace.SpanFromContext(ctx).SpanContext().TraceID().String()
		attrs = append(attrs, slog.String(TraceIDKey, traceID))
	}

	if withSpanID && spanCtx.HasSpanID() {
		spanID := spanCtx.SpanID().String()
		attrs = append(attrs, slog.String(SpanIDKey, spanID))
	}

	return attrs
}
