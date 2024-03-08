package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/labstack/echo/v4"
)

func generateRandomString(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}

type CustomContext struct {
	echo.Context
	htmxNonce           string
	responseTargetNonce string
	twNonce             string
}

func (cc *CustomContext) Get(key string) interface{} {
	val := cc.Context.Get(key)
	if val != nil {
		return val
	}
	return cc.Request().Context().Value(key)
}

func (cc *CustomContext) Set(key string, val interface{}) {
	cc.SetRequest(
		cc.Request().WithContext(
			context.WithValue(cc.Request().Context(), key, val),
		),
	)
}

func CSP() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		htmxNonce := generateRandomString(16)
		responseTargetNonce := generateRandomString(16)
		twNonce := generateRandomString(16)

		return func(c echo.Context) error {
			cc := &CustomContext{
				Context:             c,
				htmxNonce:           htmxNonce,
				responseTargetNonce: responseTargetNonce,
				twNonce:             twNonce,
			}

			htmxCSSHash := "sha256-4Su6mBWzEIFnH4pAGMOuaeBrstwJN4Z3pq/s1Kn4/KQ="

			cspHeader := fmt.Sprintf(
				"default-src 'self'; script-src 'nonce-%s' 'nonce-%s'; style-src 'nonce-%s' '%s';",
				htmxNonce,
				responseTargetNonce,
				twNonce,
				htmxCSSHash,
			)

			cc.Set("htmxNonce", htmxNonce)
			cc.Set("responseTargetNonce", responseTargetNonce)
			cc.Set("twNonce", twNonce)
			cc.Response().Header().Set("Content-Security-Policy", cspHeader)
			cc.Context.Set("htmxNonce", htmxNonce)
			cc.Context.Set("responseTargetNonce", responseTargetNonce)
			cc.Context.Set("twNonce", twNonce)
			return next(cc)
		}
	}
}
