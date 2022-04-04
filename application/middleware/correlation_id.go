package middleware

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"service_admin_contractor/infrastructure/logging"
)

const (
	CorrelationIdCtxKey    = "CorrelationId"
	correlationIdHeaderKey = "X-Correlation-ID"
)

func CorrelationHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(correlationIdHeaderKey)
		if id == "" {
			id = uuid.NewV4().String()
		}

		logEntry := logging.GetLogEntry(r).WithField("correlation_id", id)
		ctx := logging.ContextWithLogEntry(r, logEntry)
		ctx = context.WithValue(ctx, CorrelationIdCtxKey, id)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
