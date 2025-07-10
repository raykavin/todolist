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

// webConfig defines the configuration for the web server.
type webConfig struct {
	Listen       uint16            `mapstructure:"listen"`
	ReadTimeout  time.Duration     `mapstructure:"read_timeout"`
	WriteTimeout time.Duration     `mapstructure:"write_timeout"`
	SSLCert      string            `mapstructure:"ssl_cert"`
	SSLKey       string            `mapstructure:"ssl_key"`
	NoRouteTo    string            `mapstructure:"no_route_to"`
	CORS         map[string]string `mapstructure:"cors"`
}

// GetListen returns the port to listen on.
func (w webConfig) GetListen() uint16 { return w.Listen }

// GetReadTimeout returns the read timeout duration.
func (w webConfig) GetReadTimeout() time.Duration { return w.ReadTimeout }

// GetWriteTimeout returns the write timeout duration.
func (w webConfig) GetWriteTimeout() time.Duration { return w.WriteTimeout }

// GetSSLCert returns the path to the SSL certificate file.
func (w webConfig) GetSSLCert() string { return w.SSLCert }

// GetSSLKey returns the path to the SSL key file.
func (w webConfig) GetSSLKey() string { return w.SSLKey }

// GetNoRouteTo returns the custom 404 page path.
func (w webConfig) GetNoRouteTo() string { return w.NoRouteTo }

// GetCORS returns the CORS headers.
func (w webConfig) GetCORS() map[string]string { return w.CORS }
