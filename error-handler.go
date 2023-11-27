package slogfiber

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// NewErrorHandler returns a fiber.ErrorHandler that logs requests using slog.
//
// Requests with errors are logged using slog.Error().
// Requests without errors are logged using slog.Info().
func NewErrorHandler(logger *slog.Logger) fiber.ErrorHandler {
	return NewErrorHandlerWithConfig(logger, Config{
		DefaultLevel:     slog.LevelInfo,
		ClientErrorLevel: slog.LevelWarn,
		ServerErrorLevel: slog.LevelError,

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

// NewErrorHandlerWithFilters returns a fiber.ErrorHandler that logs requests using slog.
//
// Requests with errors are logged using slog.Error().
// Requests without errors are logged using slog.Info().
func NewErrorHandlerWithFilters(logger *slog.Logger, filters ...Filter) fiber.ErrorHandler {
	return NewErrorHandlerWithConfig(logger, Config{
		DefaultLevel:     slog.LevelInfo,
		ClientErrorLevel: slog.LevelWarn,
		ServerErrorLevel: slog.LevelError,

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

// NewErrorHandlerWithConfig returns a fiber.ErrorHandler that logs requests using slog.
func NewErrorHandlerWithConfig(logger *slog.Logger, config Config) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		requestID := uuid.New().String()
		// start := ??

		return handle(config, logger, c, err, requestID, nil)
	}
}
