package web

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/http2"
)

var (
	ErrInvalidListenAddress = errors.New("invalid listen address")
	ErrInvalidHost          = errors.New("invalid host")
	ErrServerNotInitialized = errors.New("server not initialized")
	ErrInvalidSSLConfig     = errors.New("invalid SSL configuration")
)

// Engine wraps Gin engine and implements HttpServer interface
type Engine struct {
	router    *gin.Engine
	server    *http.Server
	addr      string
	config    *Config
	tlsConfig *tls.Config
}

// Config holds the configuration for creating a new Gin engine
type Config struct {
	// Basic server config
	Host      string
	Port      uint16
	DebugMode bool

	// Timeouts
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration

	// SSL/TLS config
	UseSSL        bool
	SSLCert       string
	SSLKey        string
	MinTLSVersion uint16

	// HTTP/2 config
	EnableHTTP2 bool

	// Route handling
	NoRouteTo   string
	NoRouteJSON bool

	// Middleware config
	UseRecovery    bool
	TrustedProxies []string

	// Request limits
	MaxPayloadSize int64
}

// DefaultConfig returns a Config with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		Host:           "",
		Port:           8080,
		DebugMode:      false,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    60 * time.Second,
		MinTLSVersion:  tls.VersionTLS12,
		UseSSL:         false,
		EnableHTTP2:    true,
		UseRecovery:    true,
		NoRouteJSON:    true,
		MaxPayloadSize: 10 * 1024 * 1024, // 10MB default
	}
}

// RouteSetup is a function type for setting up routes
type RouteSetup func(*gin.Engine)

// MiddlewareSetup is a function type for setting up middleware
type MiddlewareSetup func(*gin.Engine)

// NewGin creates a new Gin engine with the specified configuration
func NewGin(host string, port uint16, debugMode bool) (*Engine, error) {
	config := DefaultConfig()
	config.Host = host
	config.Port = port
	config.DebugMode = debugMode

	return NewGinWithConfig(config)
}

// NewGinWithConfig creates a new Gin engine with the provided configuration
func NewGinWithConfig(config *Config) (*Engine, error) {
	if config == nil {
		config = DefaultConfig()
	}

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	// Set Gin mode
	setGinMode(config.DebugMode)

	// Create router
	router := createRouter(config)

	// Build address
	addr := buildAddress(config.Host, config.Port)

	// Create engine
	engine := &Engine{
		router: router,
		addr:   addr,
		config: config,
	}

	// Configure TLS if needed
	if config.UseSSL {
		engine.tlsConfig = createTLSConfig(config)
	}

	// Create HTTP server
	engine.server = engine.createHTTPServer()

	// Setup default middleware
	engine.setupDefaultMiddleware()

	// Setup 404 handler
	engine.setupNoRouteHandler()

	return engine, nil
}

// validateConfig validates the engine configuration
func validateConfig(config *Config) error {
	if config.Port == 0 {
		return ErrInvalidListenAddress
	}

	// Validate host if provided
	if config.Host != "" && config.Host != ":" {
		if _, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", config.Host, config.Port)); err != nil {
			return fmt.Errorf("%w: %v", ErrInvalidHost, err)
		}
	}

	// Validate SSL config
	if config.UseSSL {
		if config.SSLCert == "" || config.SSLKey == "" {
			return ErrInvalidSSLConfig
		}
	}

	return nil
}

// setGinMode configures the Gin mode
func setGinMode(debugMode bool) {
	if debugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
}

// createRouter creates and configures a new Gin router
func createRouter(config *Config) *gin.Engine {
	// Create router without default middleware
	router := gin.New()

	// Set trusted proxies if configured
	if len(config.TrustedProxies) > 0 {
		router.SetTrustedProxies(config.TrustedProxies)
	}

	return router
}

// buildAddress constructs the server address
func buildAddress(host string, port uint16) string {
	if host == "" {
		return fmt.Sprintf(":%d", port)
	}
	return fmt.Sprintf("%s:%d", host, port)
}

// createTLSConfig creates TLS configuration
func createTLSConfig(config *Config) *tls.Config {
	minVersion := config.MinTLSVersion
	if minVersion == 0 {
		minVersion = tls.VersionTLS12
	}

	return &tls.Config{
		MinVersion: minVersion,
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519,
		},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}
}

// createHTTPServer creates the underlying HTTP server
func (e *Engine) createHTTPServer() *http.Server {
	server := &http.Server{
		Addr:         e.addr,
		Handler:      e.router,
		ReadTimeout:  e.config.ReadTimeout,
		WriteTimeout: e.config.WriteTimeout,
		IdleTimeout:  e.config.IdleTimeout,
		TLSConfig:    e.tlsConfig,
	}

	// Configure HTTP/2 if enabled and using SSL
	if e.config.EnableHTTP2 && e.config.UseSSL {
		http2.ConfigureServer(server, &http2.Server{})
	}

	return server
}

// setupDefaultMiddleware sets up default middleware based on config
func (e *Engine) setupDefaultMiddleware() {
	// Recovery middleware
	if e.config.UseRecovery {
		e.router.Use(gin.Recovery())
	}

	// Payload size limit
	if e.config.MaxPayloadSize > 0 {
		e.Use(createPayloadLimitMiddleware(e.config.MaxPayloadSize))
	}
}

// createPayloadLimitMiddleware creates middleware to limit request payload size
func createPayloadLimitMiddleware(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
		c.Next()
	}
}

// setupNoRouteHandler sets up the 404 handler based on config
func (e *Engine) setupNoRouteHandler() {
	e.router.NoRoute(func(c *gin.Context) {
		if e.config.NoRouteTo != "" && !e.config.NoRouteJSON {
			c.Redirect(http.StatusTemporaryRedirect, e.config.NoRouteTo)
			return
		}

		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "The requested resource was not found",
			"path":    c.Request.URL.Path,
		})
	})
}

// SetupRoutes applies route setup function to the router
func (e *Engine) SetupRoutes(setup RouteSetup) {
	if setup != nil && e.router != nil {
		setup(e.router)
	}
}

// SetupMiddleware applies middleware setup function to the router
func (e *Engine) SetupMiddleware(setup MiddlewareSetup) {
	if setup != nil && e.router != nil {
		setup(e.router)
	}
}

// Use adds middleware to the router
func (e *Engine) Use(middleware ...gin.HandlerFunc) {
	if e.router != nil {
		e.router.Use(middleware...)
	}
}

// Listen starts the HTTP server
func (e *Engine) Listen() error {
	if e.server == nil {
		return ErrServerNotInitialized
	}

	if e.config.UseSSL {
		return e.server.ListenAndServeTLS(e.config.SSLCert, e.config.SSLKey)
	}

	return e.server.ListenAndServe()
}

// ListenAndServe is an alias for Listen
func (e *Engine) ListenAndServe() error {
	return e.Listen()
}

// ListenAndServeTLS starts the server with TLS
func (e *Engine) ListenAndServeTLS(certFile, keyFile string) error {
	if e.server == nil {
		return ErrServerNotInitialized
	}
	return e.server.ListenAndServeTLS(certFile, keyFile)
}

// Shutdown gracefully shuts down the server
func (e *Engine) Shutdown(ctx context.Context) error {
	if e.server == nil {
		return ErrServerNotInitialized
	}
	return e.server.Shutdown(ctx)
}

// Router provides access to the underlying Gin router
func (e *Engine) Router() *gin.Engine {
	return e.router
}

// Server provides access to the underlying HTTP server
func (e *Engine) Server() *http.Server {
	return e.server
}

// Addr returns the server address
func (e *Engine) Addr() string {
	return e.addr
}

// Config returns the engine configuration
func (e *Engine) Config() *Config {
	return e.config
}

// IsSSLEnabled returns whether SSL is enabled
func (e *Engine) IsSSLEnabled() bool {
	return e.config != nil && e.config.UseSSL
}

// Run is a convenience method to start the server
func (e *Engine) Run() error {
	return e.Listen()
}

// RunTLS is a convenience method to start the server with TLS
func (e *Engine) RunTLS(certFile, keyFile string) error {
	return e.ListenAndServeTLS(certFile, keyFile)
}
