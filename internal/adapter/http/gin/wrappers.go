package http

import (
	"todolist/internal/adapter/http"

	"github.com/gin-gonic/gin"
)

// WrapHandler converts a RequestContext-based handler to a Gin handler
func WrapHandler(handler func(ctx http.RequestContext)) gin.HandlerFunc {
	return func(c *gin.Context) {
		adapter := &GinAdapter{Ctx: c}
		handler(adapter)
	}
}

// WrapMiddleware converts a RequestContext-based middleware to a Gin middleware
func WrapMiddleware(middleware func(ctx http.RequestContext)) gin.HandlerFunc {
	return func(c *gin.Context) {
		adapter := &GinAdapter{Ctx: c}
		middleware(adapter)
	}
}
