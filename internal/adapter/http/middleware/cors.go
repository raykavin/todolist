package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	defaultAllowOrigin      = "*"
	defaultAllowCredentials = "true"
)

var (
	defaultAllowedHeaders = []string{
		"Content-Type",
		"Content-Length",
		"Accept-Encoding",
		"X-CSRF-Token",
		"Authorization",
		"Accept",
		"Origin",
		"Cache-Control",
		"X-Requested-With",
	}

	defaultAllowedMethods = []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodOptions,
	}
)

// CORS returns a Gin middleware that sets CORS headers.
// If customCors is empty, it uses sensible defaults.
func CORS(customCors ...map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Apply custom headers if provided
		if len(customCors) > 0 {
			for key, value := range customCors[0] {
				c.Writer.Header().Set(key, value)
			}
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Origin", defaultAllowOrigin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", defaultAllowCredentials)
			c.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join(defaultAllowedHeaders, ", "))
			c.Writer.Header().Set("Access-Control-Allow-Methods", strings.Join(defaultAllowedMethods, ", "))
		}

		// handle preflight request
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		// Continue to next handler
		c.Next()
	}
}
