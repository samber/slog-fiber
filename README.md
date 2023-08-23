
# slog: Fiber middleware

[![tag](https://img.shields.io/github/tag/samber/slog-fiber.svg)](https://github.com/samber/slog-fiber/releases)
![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.21-%23007d9c)
[![GoDoc](https://godoc.org/github.com/samber/slog-fiber?status.svg)](https://pkg.go.dev/github.com/samber/slog-fiber)
![Build Status](https://github.com/samber/slog-fiber/actions/workflows/test.yml/badge.svg)
[![Go report](https://goreportcard.com/badge/github.com/samber/slog-fiber)](https://goreportcard.com/report/github.com/samber/slog-fiber)
[![Coverage](https://img.shields.io/codecov/c/github/samber/slog-fiber)](https://codecov.io/gh/samber/slog-fiber)
[![Contributors](https://img.shields.io/github/contributors/samber/slog-fiber)](https://github.com/samber/slog-fiber/graphs/contributors)
[![License](https://img.shields.io/github/license/samber/slog-fiber)](./LICENSE)

[Fiber](https://github.com/gofiber/fiber) middleware to log http requests using [slog](https://pkg.go.dev/log/slog).

**See also:**

- [slog-multi](https://github.com/samber/slog-multi): `slog.Handler` chaining, fanout, routing, failover, load balancing...
- [slog-formatter](https://github.com/samber/slog-formatter): `slog` attribute formatting
- [slog-sampling](https://github.com/samber/slog-sampling): `slog` sampling policy
- [slog-gin](https://github.com/samber/slog-gin): Gin middleware for `slog` logger
- [slog-echo](https://github.com/samber/slog-echo): Echo middleware for `slog` logger
- [slog-fiber](https://github.com/samber/slog-fiber): Fiber middleware for `slog` logger
- [slog-datadog](https://github.com/samber/slog-datadog): A `slog` handler for `Datadog`
- [slog-rollbar](https://github.com/samber/slog-rollbar): A `slog` handler for `Rollbar`
- [slog-sentry](https://github.com/samber/slog-sentry): A `slog` handler for `Sentry`
- [slog-syslog](https://github.com/samber/slog-syslog): A `slog` handler for `Syslog`
- [slog-logstash](https://github.com/samber/slog-logstash): A `slog` handler for `Logstash`
- [slog-fluentd](https://github.com/samber/slog-fluentd): A `slog` handler for `Fluentd`
- [slog-graylog](https://github.com/samber/slog-graylog): A `slog` handler for `Graylog`
- [slog-loki](https://github.com/samber/slog-loki): A `slog` handler for `Loki`
- [slog-slack](https://github.com/samber/slog-slack): A `slog` handler for `Slack`
- [slog-telegram](https://github.com/samber/slog-telegram): A `slog` handler for `Telegram`
- [slog-mattermost](https://github.com/samber/slog-mattermost): A `slog` handler for `Mattermost`
- [slog-microsoft-teams](https://github.com/samber/slog-microsoft-teams): A `slog` handler for `Microsoft Teams`
- [slog-webhook](https://github.com/samber/slog-webhook): A `slog` handler for `Webhook`
- [slog-kafka](https://github.com/samber/slog-kafka): A `slog` handler for `Kafka`
- [slog-parquet](https://github.com/samber/slog-parquet): A `slog` handler for `Parquet` + `Object Storage`

## 🚀 Install

```sh
go get github.com/samber/slog-fiber
```

**Compatibility**: go >= 1.21

No breaking changes will be made to exported APIs before v2.0.0.

## 💡 Usage

### Minimal

```go
import (
	"github.com/gofiber/fiber/v2"
	slogfiber "github.com/samber/slog-fiber"
	"log/slog"
)

// Create a slog logger, which:
//   - Logs to stdout.
logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

app := fiber.New()

app.Use(slogfiber.New(logger))

app.Get("/", func(c *fiber.Ctx) error {
	return c.SendString("Hello, World 👋!")
})

app.Listen(":4242")

// output:
// time=2023-04-10T14:00:00.000+00:00 level=INFO msg="Incoming request" status=200 method=GET path=/ ip=::1 latency=25.958µs user-agent=curl/7.77.0 time=2023-04-10T14:00:00.000+00:00 request-id=229c7fc8-64f5-4467-bc4a-940700503b0d
```

### Using custom time formatters

```go
import (
	"github.com/gofiber/fiber/v2"
	slogfiber "github.com/samber/slog-fiber"
	slogformatter "github.com/samber/slog-formatter"
	"log/slog"
)

// Create a slog logger, which:
//   - Logs to stdout.
//   - RFC3339 with UTC time format.
logger := slog.New(
	slogformatter.NewFormatterHandler(
		slogformatter.TimezoneConverter(time.UTC),
		slogformatter.TimeFormatter(time.RFC3339, nil),
	)(
		slog.NewTextHandler(os.Stdout, nil),
	),
)

app := fiber.New()

app.Use(slogfiber.New(logger))

app.Get("/", func(c *fiber.Ctx) error {
	return c.SendString("Hello, World 👋!")
})

app.Listen(":4242")

// output:
// time=2023-04-10T14:00:00.000+00:00 level=INFO msg="Incoming request" status=200 method=GET path=/ ip=::1 latency=25.958µs user-agent=curl/7.77.0 time=2023-04-10T14:00:00Z request-id=229c7fc8-64f5-4467-bc4a-940700503b0d
```

### Using custom logger sub-group

```go
import (
	"github.com/gofiber/fiber/v2"
	slogfiber "github.com/samber/slog-fiber"
	"log/slog"
)

// Create a slog logger, which:
//   - Logs to stdout.
logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

app := fiber.New()

app.Use(slogfiber.New(logger.WithGroup("http")))

app.Get("/", func(c *fiber.Ctx) error {
	return c.SendString("Hello, World 👋!")
})

app.Listen(":4242")

// output:
// time=2023-04-10T14:00:00.000+00:00 level=INFO msg="Incoming request" http.status=200 http.method=GET http.path=/ http.ip=::1 http.latency=25.958µs http.user-agent=curl/7.77.0 http.time=2023-04-10T14:00:00Z http.request-id=229c7fc8-64f5-4467-bc4a-940700503b0d
```

### Add logger to a single route

```go
import (
	"github.com/gofiber/fiber/v2"
	slogfiber "github.com/samber/slog-fiber"
	"log/slog"
)

// Create a slog logger, which:
//   - Logs to stdout.
logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

app := fiber.New()

app.Use(slogfiber.New(logger))

app.Get("/", func(c *fiber.Ctx) error {
	return c.SendString("Hello, World 👋!")
})

app.Listen(":4242")

// output:
// time=2023-04-10T14:00:00.000+00:00 level=INFO msg="Incoming request" status=200 method=GET path=/ ip=::1 latency=25.958µs user-agent=curl/7.77.0 time=2023-04-10T14:00:00Z request-id=229c7fc8-64f5-4467-bc4a-940700503b0d
```

### Adding custom attributes

```go
import (
	"github.com/gofiber/fiber/v2"
	slogfiber "github.com/samber/slog-fiber"
	"log/slog"
)

// Create a slog logger, which:
//   - Logs to stdout.
logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

// Add an attribute to all log entries made through this logger.
logger = logger.With("env", "production")

app := fiber.New()

app.Use(slogfiber.New(logger.WithGroup("http")))

app.Get("/", func(c *fiber.Ctx) error {
	return c.SendString("Hello, World 👋!")
})

app.Listen(":4242")

// output:
// time=2023-04-10T14:00:00.000+00:00 level=INFO msg="Incoming request" env=production status=200 method=GET path=/ ip=::1 latency=25.958µs user-agent=curl/7.77.0 time=2023-04-10T14:00:00Z request-id=229c7fc8-64f5-4467-bc4a-940700503b0d
```

### JSON output

```go
import (
	"github.com/gofiber/fiber/v2"
	slogfiber "github.com/samber/slog-fiber"
	"log/slog"
)

// Create a slog logger, which:
//   - Logs to stdout.
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)),

app := fiber.New()

app.Use(slogfiber.New(logger))

app.Get("/", func(c *fiber.Ctx) error {
	return c.SendString("Hello, World 👋!")
})

app.Listen(":4242")

// output:
// {"time":"2023-04-10T14:00:00Z","level":"INFO","msg":"Incoming request","status":200,"method":"GET","path":"/","ip":"::1","latency":26750,"user-agent":"curl/7.77.0","time":"2023-04-10T14:00:00Z","request-id":"04201917-d7ba-4b20-a3bb-2fffba5f2bd9"}
```

## 🤝 Contributing

- Ping me on twitter [@samuelberthe](https://twitter.com/samuelberthe) (DMs, mentions, whatever :))
- Fork the [project](https://github.com/samber/slog-fiber)
- Fix [open issues](https://github.com/samber/slog-fiber/issues) or request new features

Don't hesitate ;)

```bash
# Install some dev dependencies
make tools

# Run tests
make test
# or
make watch-test
```

## 👤 Contributors

![Contributors](https://contrib.rocks/image?repo=samber/slog-fiber)

## 💫 Show your support

Give a ⭐️ if this project helped you!

[![GitHub Sponsors](https://img.shields.io/github/sponsors/samber?style=for-the-badge)](https://github.com/sponsors/samber)

## 📝 License

Copyright © 2023 [Samuel Berthe](https://github.com/samber).

This project is [MIT](./LICENSE) licensed.
