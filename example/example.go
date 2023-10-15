package main

import (
	"fmt"
	"os"
	"time"

	"log/slog"

	"github.com/gofiber/fiber/v2"
	slogfiber "github.com/samber/slog-fiber"
	slogformatter "github.com/samber/slog-formatter"
)

func main() {
	// Create a slog logger, which:
	//   - Logs to stdout.
	//   - RFC3339 with UTC time format.
	logger := slog.New(
		slogformatter.NewFormatterHandler(
			slogformatter.TimezoneConverter(time.UTC),
			slogformatter.TimeFormatter(time.RFC3339, nil),
		)(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}),
		),
	)

	// Add an attribute to all log entries made through this logger.
	logger = logger.With("env", "production")

	app := fiber.New()

	app.Use(slogfiber.New(logger.WithGroup("http")))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})
	app.Get("/foobar/:id", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	err := app.Listen(":4242")
	if err != nil {
		fmt.Println(err.Error())
	}

	// output:
	// time=2023-04-10T14:00:00.000+00:00 level=INFO msg="Incoming request" env=production http.status=200 http.method=GET http.path=/ http.route=/ http.ip=::1 http.latency=25.958µs http.user-agent=curl/7.77.0 http.time=2023-04-10T14:00:00Z http.request-id=229c7fc8-64f5-4467-bc4a-940700503b0d
}
