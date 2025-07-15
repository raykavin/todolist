package http

import (
	"github.com/gin-gonic/gin"
)

// WrapHandler converts a RequestContext-based handler to a Gin handler
func WrapHandler(handler func(ctx RequestContext)) gin.HandlerFunc {
	return func(c *gin.Context) {
		adapter := &GinAdapter{Ctx: c}
		handler(adapter)
	}
}
