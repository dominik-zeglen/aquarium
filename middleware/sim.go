package middleware

import (
	"context"
	"net/http"

	"github.com/dominik-zeglen/aquarium/sim"
	opentracing "github.com/opentracing/opentracing-go"
)

// SimContextKey defines key holding website config data in request context
const SimContextKey = ContextKey("config")

// WithSim provides sim by context
func WithSim(next http.Handler, s *sim.Sim) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			simSpan, _ := opentracing.StartSpanFromContext(r.Context(), "lock")
			ctx := context.WithValue(r.Context(), SimContextKey, s)
			s.Lock()
			simSpan.Finish()
			defer s.Unlock()

			next.ServeHTTP(w, r.WithContext(ctx))
			return
		},
	)
}
