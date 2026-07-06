// Package omnivoice provides a unified interface for speech-to-text and text-to-speech.
// This is the batteries-included package that imports all providers.
// For a minimal dependency footprint, use github.com/plexusone/omnivoice-core instead.
package omnivoice

import (
	"fmt"

	"github.com/plexusone/omnivoice-core/gateway"

	deepgramRealtime "github.com/plexusone/omni-deepgram/omnivoice/realtime"
	googleRealtime "github.com/plexusone/omni-google/omnivoice/realtime"
	openaiRealtime "github.com/plexusone/omni-openai/omnivoice/realtime"
)

// GetRealtimeFactory returns a realtime provider factory by name.
// Supported providers:
//   - "openai": OpenAI Realtime API (~100ms latency)
//   - "gemini": Google Gemini Live API (~200ms latency)
//   - "deepgram": Deepgram Voice Agent API (~100-300ms latency)
//
// Example:
//
//	factory, err := omnivoice.GetRealtimeFactory("openai")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// Use factory with gateway configuration
func GetRealtimeFactory(name string) (gateway.RealtimeProviderFactory, error) {
	switch name {
	case "openai":
		return openaiRealtime.NewFactory(), nil
	case "gemini":
		return googleRealtime.NewFactory(), nil
	case "deepgram":
		return deepgramRealtime.NewFactory(), nil
	default:
		return nil, fmt.Errorf("unknown realtime factory: %s (supported: openai, gemini, deepgram)", name)
	}
}

// MustGetRealtimeFactory returns a realtime provider factory by name.
// It panics if the factory is not found.
func MustGetRealtimeFactory(name string) gateway.RealtimeProviderFactory {
	factory, err := GetRealtimeFactory(name)
	if err != nil {
		panic(err)
	}
	return factory
}

// ListRealtimeFactories returns the names of all available realtime factories.
func ListRealtimeFactories() []string {
	return []string{"openai", "gemini", "deepgram"}
}

// HasRealtimeFactory returns true if a realtime factory with the given name exists.
func HasRealtimeFactory(name string) bool {
	switch name {
	case "openai", "gemini", "deepgram":
		return true
	default:
		return false
	}
}
