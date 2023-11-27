package slogfiber

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

// New returns a fiber.Handler (middleware) that logs requests using slog.
//
// Requests with errors are logged using slog.Error().
// Requests without errors are logged using slog.Info().
//
// Deprecated: Use NewMiddleware(...) instead.
func New(logger *slog.Logger) fiber.Handler {
	return NewMiddleware(logger)
}

// NewWithFilters returns a fiber.Handler (middleware) that logs requests using slog.
//
// Requests with errors are logged using slog.Error().
// Requests without errors are logged using slog.Info().
//
// Deprecated: Use NewMiddlewareWithFilters(...) instead.
func NewWithFilters(logger *slog.Logger, filters ...Filter) fiber.Handler {
	return NewMiddlewareWithFilters(logger, filters...)
}

// NewWithConfig returns a fiber.Handler (middleware) that logs requests using slog.
//
// Deprecated: Use NewMiddlewareWithConfig(...) instead.
func NewWithConfig(logger *slog.Logger, config Config) fiber.Handler {
	return NewMiddlewareWithConfig(logger, config)
}
