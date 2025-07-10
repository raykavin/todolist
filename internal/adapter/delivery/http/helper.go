package http

import "github.com/gin-gonic/gin"

// WrapHandler converts a RequestContext-based handler to a Gin handler
func WrapHandler(handler func(ctx RequestContext)) gin.HandlerFunc {
	return func(c *gin.Context) {
		adapter := &GinAdapter{Ctx: c}
		handler(adapter)
	}
}

// WrapMiddleware converts a RequestContext-based middleware to a Gin middleware
func WrapMiddleware(middleware func(ctx RequestContext)) gin.HandlerFunc {
	return func(c *gin.Context) {
		adapter := &GinAdapter{Ctx: c}
		middleware(adapter)
	}
}

// GetAdapter extracts the GinAdapter from a Gin context
// Useful if you need to access Gin-specific features
func GetAdapter(c *gin.Context) *GinAdapter {
	return &GinAdapter{Ctx: c}
}
