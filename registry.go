package omnivoice

import (
	"fmt"
	"sync"

	"github.com/plexusone/omnivoice-core/stt"
	"github.com/plexusone/omnivoice-core/tts"
)

// ProviderConfig holds common configuration options for creating providers.
type ProviderConfig struct {
	// APIKey is the authentication key for the provider.
	APIKey string //nolint:gosec // G117: This is a config struct, not storing secrets

	// BaseURL is an optional custom API endpoint.
	BaseURL string

	// Extensions holds provider-specific configuration.
	Extensions map[string]any
}

// ProviderOption configures a ProviderConfig.
type ProviderOption func(*ProviderConfig)

// WithAPIKey sets the API key for the provider.
func WithAPIKey(apiKey string) ProviderOption {
	return func(c *ProviderConfig) {
		c.APIKey = apiKey
	}
}

// WithBaseURL sets a custom base URL for the provider.
func WithBaseURL(baseURL string) ProviderOption {
	return func(c *ProviderConfig) {
		c.BaseURL = baseURL
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

// TTSProviderFactory creates a TTS provider with the given configuration.
type TTSProviderFactory func(config ProviderConfig) (tts.Provider, error)

// STTProviderFactory creates an STT provider with the given configuration.
type STTProviderFactory func(config ProviderConfig) (stt.Provider, error)

// registry holds registered provider factories.
var registry = &providerRegistry{
	ttsFactories: make(map[string]TTSProviderFactory),
	sttFactories: make(map[string]STTProviderFactory),
}

type providerRegistry struct {
	mu           sync.RWMutex
	ttsFactories map[string]TTSProviderFactory
	sttFactories map[string]STTProviderFactory
}

// RegisterTTSProvider registers a TTS provider factory by name.
// This is typically called from provider init() functions.
func RegisterTTSProvider(name string, factory TTSProviderFactory) {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	registry.ttsFactories[name] = factory
}

// RegisterSTTProvider registers an STT provider factory by name.
// This is typically called from provider init() functions.
func RegisterSTTProvider(name string, factory STTProviderFactory) {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	registry.sttFactories[name] = factory
}

// GetTTSProvider creates a TTS provider by name with the given options.
//
// Example:
//
//	provider, err := omnivoice.GetTTSProvider("elevenlabs", omnivoice.WithAPIKey(apiKey))
func GetTTSProvider(name string, opts ...ProviderOption) (tts.Provider, error) {
	registry.mu.RLock()
	factory, ok := registry.ttsFactories[name]
	registry.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("unknown TTS provider: %s (available: %v)", name, ListTTSProviders())
	}

	config := ProviderConfig{}
	for _, opt := range opts {
		opt(&config)
	}

	return factory(config)
}

// GetSTTProvider creates an STT provider by name with the given options.
//
// Example:
//
//	provider, err := omnivoice.GetSTTProvider("deepgram", omnivoice.WithAPIKey(apiKey))
func GetSTTProvider(name string, opts ...ProviderOption) (stt.Provider, error) {
	registry.mu.RLock()
	factory, ok := registry.sttFactories[name]
	registry.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("unknown STT provider: %s (available: %v)", name, ListSTTProviders())
	}

	config := ProviderConfig{}
	for _, opt := range opts {
		opt(&config)
	}

	return factory(config)
}

// ListTTSProviders returns the names of all registered TTS providers.
func ListTTSProviders() []string {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	names := make([]string, 0, len(registry.ttsFactories))
	for name := range registry.ttsFactories {
		names = append(names, name)
	}
	return names
}

// ListSTTProviders returns the names of all registered STT providers.
func ListSTTProviders() []string {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	names := make([]string, 0, len(registry.sttFactories))
	for name := range registry.sttFactories {
		names = append(names, name)
	}
	return names
}

// HasTTSProvider checks if a TTS provider is registered.
func HasTTSProvider(name string) bool {
	registry.mu.RLock()
	defer registry.mu.RUnlock()
	_, ok := registry.ttsFactories[name]
	return ok
}

// HasSTTProvider checks if an STT provider is registered.
func HasSTTProvider(name string) bool {
	registry.mu.RLock()
	defer registry.mu.RUnlock()
	_, ok := registry.sttFactories[name]
	return ok
}
