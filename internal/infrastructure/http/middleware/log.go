package middleware

import (
	"time"
	"todolist/internal/adapter/http"
	gHttp "todolist/internal/adapter/http/gin"
	"todolist/pkg/log"

	"github.com/gin-gonic/gin"
)

// RequestLog creates a middleware for request logging
func RequestLog(log log.ExtendedLog) gin.HandlerFunc {
	return gHttp.WrapMiddleware(func(ctx http.RequestContext) {
		start := time.Now()

		// Process request
		ctx.Next()

		// Calculate duration
		duration := time.Since(start)

		// Log the request
		log.API(
			ctx.Request().Method,
			ctx.Request().URL.Path,
			ctx.Request().Response.StatusCode,
			duration,
		)
	})
}
