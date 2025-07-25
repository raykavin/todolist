package middleware

import (
	"time"
	"todolist/internal/adapter/delivery/http"
	"todolist/pkg/logger"

	"github.com/gin-gonic/gin"
)

// RequestLog creates a middleware for request logging
func RequestLogger(log logger.ExtendedLog) gin.HandlerFunc {
	return http.WrapHandler(func(ctx http.RequestContext) {
		start := time.Now()

		// Process request
		ctx.Next()

		// Calculate duration
		duration := time.Since(start)

		// Log the request
		log.API(
			ctx.Request().Method,
			ctx.Request().URL.Path,
			ctx.StatusCode(),
			duration,
		)
	})
}
