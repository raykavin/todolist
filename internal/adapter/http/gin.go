package http

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GinAdapter implements the RequestContext interface using Gin
type GinAdapter struct {
	Ctx *gin.Context
}

// Abort implements http.RequestContext.
func (g *GinAdapter) Abort() {
	g.Ctx.Abort()
}

// BindJSON implements http.RequestContext.
func (g *GinAdapter) BindJSON(dest any) error {
	return g.Ctx.BindJSON(dest)
}

// BindQuery implements http.RequestContext.
func (g *GinAdapter) BindQuery(dest any) error {
	return g.Ctx.BindQuery(dest)
}

// Context implements http.RequestContext.
func (g *GinAdapter) Context() context.Context {
	return g.Ctx.Request.Context()
}

// Get implements http.RequestContext.
func (g *GinAdapter) Get(key string) (any, bool) {
	return g.Ctx.Get(key)
}

// GetParam implements http.RequestContext.
func (g *GinAdapter) GetParam(key string) string {
	return g.Ctx.Param(key)
}

// GetQuery implements http.RequestContext.
func (g *GinAdapter) GetQuery(key string) string {
	return g.Ctx.Query(key)
}

// JSON implements http.RequestContext.
func (g *GinAdapter) JSON(statusCode int, data any) {
	g.Ctx.JSON(statusCode, data)
}

// Next implements http.RequestContext.
func (g *GinAdapter) Next() {
	g.Ctx.Next()
}

// Redirect implements http.RequestContext.
func (g *GinAdapter) Redirect(statusCode int, to string) {
	g.Ctx.Redirect(statusCode, to)
}

// Request implements http.RequestContext.
func (g *GinAdapter) Request() *http.Request {
	return g.Ctx.Request
}

// Set implements http.RequestContext.
func (g *GinAdapter) Set(key string, value any) {
	g.Ctx.Set(key, value)
}

// SetCookie implements http.RequestContext.
func (g *GinAdapter) SetCookie(c *http.Cookie) {
	g.Ctx.SetCookie(c.Name, c.Value, c.MaxAge, c.Path, c.Domain, c.Secure, c.HttpOnly)
}

// Writer implements http.RequestContext.
func (g *GinAdapter) Writer() http.ResponseWriter {
	return g.Ctx.Writer
}
