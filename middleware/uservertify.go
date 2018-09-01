package middleware

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	// KeyAuthConfig defines the config for KeyAuth middleware.
	KeyAuthConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper middleware.Skipper

		// Validator is a function to validate key.
		// Required.
		Validator UserValidator
	}

	// UserValidator defines a function to validate KeyAuth credentials.
	UserValidator func(echo.Context) (bool, error)
)

var (
	// DefaultUserVertifyConfig is the default KeyAuth middleware config.
	DefaultUserVertifyConfig = KeyAuthConfig{
		Skipper: middleware.DefaultSkipper,
	}
)

// UserVertify returns an KeyAuth middleware.
//
// For valid key it calls the next handler.
// For invalid key, it sends "401 - Unauthorized" response.
// For missing key, it sends "400 - Bad Request" response.
func UserVertify(fn UserValidator) echo.MiddlewareFunc {
	c := DefaultUserVertifyConfig
	c.Validator = fn
	return UserVertifyWithConfig(c)
}

// UserVertifyWithConfig returns an KeyAuth middleware with config.
// See `KeyAuth()`.
func UserVertifyWithConfig(config KeyAuthConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultUserVertifyConfig.Skipper
	}

	if config.Validator == nil {
		panic("echo: user vertify middleware requires a validator function")
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			valid, err := config.Validator(c)
			if err != nil {
				return err
			} else if valid {
				return next(c)
			}

			return echo.ErrUnauthorized
		}
	}
}
