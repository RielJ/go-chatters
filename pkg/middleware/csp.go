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

			htmxCSSHash := "sha256-d7rFBVhb3n/Drrf+EpNWYdITkos3kQRFpB0oSOycXg4="

			cspHeader := fmt.Sprintf(
				"default-src 'self'; script-src 'nonce-%s' 'nonce-%s'; style-src 'self' 'nonce-%s' '%s';",
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

// func TextHTMLMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-Type", "text/html charset=utf-8")
// 		next.ServeHTTP(w, r)
// 	})
// }

func TextHTML() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
			return next(c)
		}
	}
}
