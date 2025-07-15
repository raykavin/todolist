package config

import "time"

/*
 * web.go
 *
 * This file defines web server configuration settings.
 *
 * Examples include server port, CORS settings, timeouts, allowed hosts,
 * or SSL/TLS certificates.
 *
 * These settings help you run and secure your HTTP server properly
 * in different environments.
 */

var _ WebConfigProvider = (*webConfig)(nil)

// webConfig defines the configuration for the web server.
type webConfig struct {
	Listen         uint16            `mapstructure:"listen"`           // Server listening port
	ReadTimeout    time.Duration     `mapstructure:"read_timeout"`     // Maximum duration for reading the entire request
	WriteTimeout   time.Duration     `mapstructure:"write_timeout"`    // Maximum duration before timing out writes of the response
	IdleTimeout    time.Duration     `mapstructure:"idle_timeout"`     // Maximum amount of time to wait for the next request when keep-alives are enabled
	UseSSL         bool              `mapstructure:"use_ssl"`          // Enable SSL
	SSLCert        string            `mapstructure:"crt"`              // Path to the SSL certificate file
	SSLKey         string            `mapstructure:"key"`              // Path to the SSL key file
	MaxPayloadSize int64             `mapstructure:"max_payload_size"` // Maximum allowed payload size in bytes
	NoRouteTo      string            `mapstructure:"no_router"`        // URL to redirect when route not found (404)
	CORS           map[string]string `mapstructure:"cors"`             // CORS headers configuration
}

// GetListen returns the port number to listen on.
func (w webConfig) GetListen() uint16 { return w.Listen }

// GetReadTimeout returns the read timeout duration.
func (w webConfig) GetReadTimeout() time.Duration { return w.ReadTimeout }

// GetWriteTimeout returns the write timeout duration.
func (w webConfig) GetWriteTimeout() time.Duration { return w.WriteTimeout }

// GetIdleTimeout returns the idle timeout duration.
func (w webConfig) GetIdleTimeout() time.Duration { return w.IdleTimeout }

// GetUseSSL indicates whether SSL is enabled.
func (w webConfig) GetUseSSL() bool { return w.UseSSL }

// GetSSLCert returns the SSL certificate file path.
func (w webConfig) GetSSLCert() string { return w.SSLCert }

// GetSSLKey returns the SSL key file path.
func (w webConfig) GetSSLKey() string { return w.SSLKey }

// GetMaxPayloadSize returns the maximum payload size allowed.
func (w webConfig) GetMaxPayloadSize() int64 { return w.MaxPayloadSize }

// GetNoRouteTo returns the URL to redirect to when no route matches.
func (w webConfig) GetNoRouteTo() string { return w.NoRouteTo }

// GetCORS returns the CORS headers configuration.
func (w webConfig) GetCORS() map[string]string { return w.CORS }
