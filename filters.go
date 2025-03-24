package slogfiber

import (
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v3"
)

type Filter func(ctx fiber.Ctx) bool

// Basic
func Accept(filter Filter) Filter { return filter }
func Ignore(filter Filter) Filter { return filter }

// Method
func AcceptMethod(methods ...string) Filter {
	return func(c fiber.Ctx) bool {
		reqMethod := strings.ToLower(string(c.RequestCtx().Method()))

		for _, method := range methods {
			if strings.ToLower(method) == reqMethod {
				return true
			}
		}

		return false
	}
}

func IgnoreMethod(methods ...string) Filter {
	return func(c fiber.Ctx) bool {
		reqMethod := strings.ToLower(string(c.RequestCtx().Method()))

		for _, method := range methods {
			if strings.ToLower(method) == reqMethod {
				return false
			}
		}

		return true
	}
}

// Status
func AcceptStatus(statuses ...int) Filter {
	return func(c fiber.Ctx) bool {
		for _, status := range statuses {
			if status == c.Response().StatusCode() {
				return true
			}
		}

		return false
	}
}

func IgnoreStatus(statuses ...int) Filter {
	return func(c fiber.Ctx) bool {
		for _, status := range statuses {
			if status == c.Response().StatusCode() {
				return false
			}
		}

		return true
	}
}

func AcceptStatusGreaterThan(status int) Filter {
	return func(c fiber.Ctx) bool {
		return c.Response().StatusCode() > status
	}
}

func IgnoreStatusLessThan(status int) Filter {
	return func(c fiber.Ctx) bool {
		return c.Response().StatusCode() < status
	}
}

func AcceptStatusGreaterThanOrEqual(status int) Filter {
	return func(c fiber.Ctx) bool {
		return c.Response().StatusCode() >= status
	}
}

func IgnoreStatusLessThanOrEqual(status int) Filter {
	return func(c fiber.Ctx) bool {
		return c.Response().StatusCode() <= status
	}
}

// Path
func AcceptPath(urls ...string) Filter {
	return func(c fiber.Ctx) bool {
		for _, url := range urls {
			if c.Path() == url {
				return true
			}
		}

		return false
	}
}

func IgnorePath(urls ...string) Filter {
	return func(c fiber.Ctx) bool {
		for _, url := range urls {
			if c.Path() == url {
				return false
			}
		}

		return true
	}
}

func AcceptPathContains(parts ...string) Filter {
	return func(c fiber.Ctx) bool {
		for _, part := range parts {
			if strings.Contains(c.Path(), part) {
				return true
			}
		}

		return false
	}
}

func IgnorePathContains(parts ...string) Filter {
	return func(c fiber.Ctx) bool {
		for _, part := range parts {
			if strings.Contains(c.Path(), part) {
				return false
			}
		}

		return true
	}
}

func AcceptPathPrefix(prefixs ...string) Filter {
	return func(c fiber.Ctx) bool {
		for _, prefix := range prefixs {
			if strings.HasPrefix(c.Path(), prefix) {
				return true
			}
		}

		return false
	}
}

func IgnorePathPrefix(prefixs ...string) Filter {
	return func(c fiber.Ctx) bool {
		for _, prefix := range prefixs {
			if strings.HasPrefix(c.Path(), prefix) {
				return false
			}
		}

		return true
	}
}

func AcceptPathSuffix(prefixs ...string) Filter {
	return func(c fiber.Ctx) bool {
		for _, prefix := range prefixs {
			if strings.HasPrefix(c.Path(), prefix) {
				return true
			}
		}

		return false
	}
}

func IgnorePathSuffix(suffixs ...string) Filter {
	return func(c fiber.Ctx) bool {
		for _, suffix := range suffixs {
			if strings.HasSuffix(c.Path(), suffix) {
				return false
			}
		}

		return true
	}
}

func AcceptPathMatch(regs ...regexp.Regexp) Filter {
	return func(c fiber.Ctx) bool {
		for _, reg := range regs {
			if reg.Match([]byte(c.Path())) {
				return true
			}
		}

		return false
	}
}

func IgnorePathMatch(regs ...regexp.Regexp) Filter {
	return func(c fiber.Ctx) bool {
		for _, reg := range regs {
			if reg.Match([]byte(c.Path())) {
				return false
			}
		}

		return true
	}
}

// Host
func AcceptHost(hosts ...string) Filter {
	return func(c fiber.Ctx) bool {
		for _, host := range hosts {
			if c.Hostname() == host {
				return true
			}
		}

		return false
	}
}

func IgnoreHost(hosts ...string) Filter {
	return func(c fiber.Ctx) bool {
		for _, host := range hosts {
			if c.Hostname() == host {
				return false
			}
		}

		return true
	}
}

func AcceptHostContains(parts ...string) Filter {
	return func(c fiber.Ctx) bool {
		for _, part := range parts {
			if strings.Contains(c.Hostname(), part) {
				return true
			}
		}

		return false
	}
}

func IgnoreHostContains(parts ...string) Filter {
	return func(c fiber.Ctx) bool {
		for _, part := range parts {
			if strings.Contains(c.Hostname(), part) {
				return false
			}
		}

		return true
	}
}

func AcceptHostPrefix(prefixs ...string) Filter {
	return func(c fiber.Ctx) bool {
		for _, prefix := range prefixs {
			if strings.HasPrefix(c.Hostname(), prefix) {
				return true
			}
		}

		return false
	}
}

func IgnoreHostPrefix(prefixs ...string) Filter {
	return func(c fiber.Ctx) bool {
		for _, prefix := range prefixs {
			if strings.HasPrefix(c.Hostname(), prefix) {
				return false
			}
		}

		return true
	}
}

func AcceptHostSuffix(prefixs ...string) Filter {
	return func(c fiber.Ctx) bool {
		for _, prefix := range prefixs {
			if strings.HasPrefix(c.Hostname(), prefix) {
				return true
			}
		}

		return false
	}
}

func IgnoreHostSuffix(suffixs ...string) Filter {
	return func(c fiber.Ctx) bool {
		for _, suffix := range suffixs {
			if strings.HasSuffix(c.Hostname(), suffix) {
				return false
			}
		}

		return true
	}
}

func AcceptHostMatch(regs ...regexp.Regexp) Filter {
	return func(c fiber.Ctx) bool {
		for _, reg := range regs {
			if reg.Match([]byte(c.Hostname())) {
				return true
			}
		}

		return false
	}
}

func IgnoreHostMatch(regs ...regexp.Regexp) Filter {
	return func(c fiber.Ctx) bool {
		for _, reg := range regs {
			if reg.Match([]byte(c.Hostname())) {
				return false
			}
		}

		return true
	}
}
