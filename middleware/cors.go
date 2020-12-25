package middleware

import (
	"net/http"

	"github.com/rs/cors"
)

type CorsHandler struct {
	http.Handler
	allowedOrigins []string
	next           http.Handler
}

func (handler CorsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := cors.New(cors.Options{
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
		},
		AllowedMethods: []string{http.MethodPost, http.MethodOptions},
		AllowedOrigins: handler.allowedOrigins,
	})
	c.Handler(handler.next).ServeHTTP(w, r)
}

func WithCors(
	allowedOrigins []string,
	next http.Handler,
) CorsHandler {
	return CorsHandler{
		allowedOrigins: allowedOrigins,
		next:           next,
	}
}
