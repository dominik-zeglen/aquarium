package middleware

import (
	"net/http"

	opentracing "github.com/opentracing/opentracing-go"
)

// WithTracing provides tracer context
func WithTracing(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			span := opentracing.GlobalTracer().StartSpan("api-handler")
			defer span.Finish()
			ctx := opentracing.ContextWithSpan(r.Context(), span)

			next.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}
