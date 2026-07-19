// Package registry provides types for avatar provider registration and discovery.
package registry

import (
	"github.com/plexusone/omniavatar-core/live"
	"github.com/plexusone/omniavatar-core/render"
)

// ProviderConfig holds common configuration options for creating providers.
type ProviderConfig struct {
	// APIKey is the authentication key for the provider.
	APIKey string //nolint:gosec // G117: This is a config struct, not storing secrets

	// BaseURL is an optional custom API endpoint.
	BaseURL string

	// Extensions holds provider-specific configuration.
	// Keys and values depend on the provider.
	Extensions map[string]any
}

// ProviderOption configures a ProviderConfig.
type ProviderOption func(*ProviderConfig)

// ApplyOptions applies provider options to a config.
func ApplyOptions(opts ...ProviderOption) ProviderConfig {
	config := ProviderConfig{
		Extensions: make(map[string]any),
	}
	for _, opt := range opts {
		opt(&config)
	}
	return config
}

// WithAPIKey sets the API key for the provider.
func WithAPIKey(key string) ProviderOption {
	return func(c *ProviderConfig) {
		c.APIKey = key
	}
}

// WithBaseURL sets a custom API base URL.
func WithBaseURL(url string) ProviderOption {
	return func(c *ProviderConfig) {
		c.BaseURL = url
	}
}

// WithExtension sets a provider-specific configuration value.
func WithExtension(key string, value any) ProviderOption {
	return func(c *ProviderConfig) {
		if c.Extensions == nil {
			c.Extensions = make(map[string]any)
		}
		c.Extensions[key] = value
	}
}

// LiveProviderFactory creates a live (real-time session) Provider from
// configuration.
type LiveProviderFactory func(config ProviderConfig) (live.Provider, error)

// RenderProviderFactory creates a render (batch generation) Provider
// from configuration.
type RenderProviderFactory func(config ProviderConfig) (render.Provider, error)

// GetString retrieves a string extension value with a default.
func (c ProviderConfig) GetString(key, defaultValue string) string {
	if v, ok := c.Extensions[key].(string); ok {
		return v
	}
	return defaultValue
}

// GetBool retrieves a bool extension value with a default.
func (c ProviderConfig) GetBool(key string, defaultValue bool) bool {
	if v, ok := c.Extensions[key].(bool); ok {
		return v
	}
	return defaultValue
}

// GetInt retrieves an int extension value with a default.
func (c ProviderConfig) GetInt(key string, defaultValue int) int {
	switch v := c.Extensions[key].(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	default:
		return defaultValue
	}
}
