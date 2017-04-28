package etag

import (
	"net/http"

	"github.com/acoshift/middleware"
)

// Config is the etag config
type Config struct {
	Skipper middleware.Skipper
}

// DefaultConfig is the default config for middleware
var DefaultConfig = Config{
	Skipper: middleware.DefaultSkipper,
}

// New creates new etag middleware
func New(c Config) middleware.Middleware {
	if c.Skipper == nil {
		c.Skipper = DefaultConfig.Skipper
	}

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if c.Skipper(r) {
				h.ServeHTTP(w, r)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
