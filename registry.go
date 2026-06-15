// Package omnivoice provides a unified interface for speech-to-text and text-to-speech.
// This is the batteries-included package that imports all providers.
// For a minimal dependency footprint, use github.com/plexusone/omnivoice-core instead.
package omnivoice

import (
	core "github.com/plexusone/omnivoice-core"
	"github.com/plexusone/omnivoice-core/callsystem"
	"github.com/plexusone/omnivoice-core/registry"
	"github.com/plexusone/omnivoice-core/stt"
	"github.com/plexusone/omnivoice-core/tts"
)

// Re-export registry types from omnivoice-core.
type (
	// ProviderConfig holds common configuration options for creating providers.
	ProviderConfig = registry.ProviderConfig

	// ProviderOption configures a ProviderConfig.
	ProviderOption = registry.ProviderOption
)

// Re-export priority constants from omnivoice-core.
const (
	// PriorityThin is the priority for thin (stdlib-only) provider implementations.
	PriorityThin = core.PriorityThin

	// PriorityThick is the priority for thick (official SDK) provider implementations.
	PriorityThick = core.PriorityThick
)

// Re-export registry option functions from omnivoice-core.
var (
	// Core options
	WithAPIKey    = registry.WithAPIKey
	WithBaseURL   = registry.WithBaseURL
	WithExtension = registry.WithExtension

	// CallSystem options
	WithAccountSID  = registry.WithAccountSID
	WithAuthToken   = registry.WithAuthToken
	WithPhoneNumber = registry.WithPhoneNumber
	WithWebhookURL  = registry.WithWebhookURL
	WithRegion      = registry.WithRegion

	// Gateway options
	WithListener     = registry.WithListener
	WithPublicURL    = registry.WithPublicURL
	WithListenAddr   = registry.WithListenAddr
	WithConnectionID = registry.WithConnectionID

	// Realtime options
	WithVoice        = registry.WithVoice
	WithModel        = registry.WithModel
	WithInstructions = registry.WithInstructions

	// Pipeline options
	WithSTTProvider     = registry.WithSTTProvider
	WithSTTAPIKey       = registry.WithSTTAPIKey
	WithSTTModel        = registry.WithSTTModel
	WithSTTLanguage     = registry.WithSTTLanguage
	WithTTSProvider     = registry.WithTTSProvider
	WithTTSAPIKey       = registry.WithTTSAPIKey
	WithTTSVoiceID      = registry.WithTTSVoiceID
	WithTTSModel        = registry.WithTTSModel
	WithLLMProvider     = registry.WithLLMProvider
	WithLLMAPIKey       = registry.WithLLMAPIKey
	WithLLMModel        = registry.WithLLMModel
	WithLLMSystemPrompt = registry.WithLLMSystemPrompt

	// Session options
	WithGreeting           = registry.WithGreeting
	WithMaxSessionDuration = registry.WithMaxSessionDuration
	WithInterruptionMode   = registry.WithInterruptionMode
	WithLogger             = registry.WithLogger
	WithPipelineMode       = registry.WithPipelineMode
)

// STT Provider Registry - delegates to omnivoice-core

// RegisterSTTProvider registers an STT provider factory with the given name and priority.
// Higher priority values override lower priority registrations.
func RegisterSTTProvider(name string, factory registry.STTProviderFactory, priority int) {
	core.RegisterSTTProvider(name, factory, priority)
}

// GetSTTProvider creates an STT provider instance from the registry.
func GetSTTProvider(name string, opts ...ProviderOption) (stt.Provider, error) {
	return core.GetSTTProvider(name, opts...)
}

// ListSTTProviders returns a list of all registered STT provider names.
func ListSTTProviders() []string {
	return core.ListSTTProviders()
}

// HasSTTProvider returns true if an STT provider with the given name is registered.
func HasSTTProvider(name string) bool {
	return core.HasSTTProvider(name)
}

// GetSTTProviderPriority returns the priority of the registered STT provider.
func GetSTTProviderPriority(name string) int {
	return core.GetSTTProviderPriority(name)
}

// TTS Provider Registry - delegates to omnivoice-core

// RegisterTTSProvider registers a TTS provider factory with the given name and priority.
// Higher priority values override lower priority registrations.
func RegisterTTSProvider(name string, factory registry.TTSProviderFactory, priority int) {
	core.RegisterTTSProvider(name, factory, priority)
}

// GetTTSProvider creates a TTS provider instance from the registry.
func GetTTSProvider(name string, opts ...ProviderOption) (tts.Provider, error) {
	return core.GetTTSProvider(name, opts...)
}

// ListTTSProviders returns a list of all registered TTS provider names.
func ListTTSProviders() []string {
	return core.ListTTSProviders()
}

// HasTTSProvider returns true if a TTS provider with the given name is registered.
func HasTTSProvider(name string) bool {
	return core.HasTTSProvider(name)
}

// GetTTSProviderPriority returns the priority of the registered TTS provider.
func GetTTSProviderPriority(name string) int {
	return core.GetTTSProviderPriority(name)
}

// CallSystem Provider Registry - delegates to omnivoice-core

// RegisterCallSystemProvider registers a CallSystem provider factory with the given name and priority.
// Higher priority values override lower priority registrations.
func RegisterCallSystemProvider(name string, factory registry.CallSystemProviderFactory, priority int) {
	core.RegisterCallSystemProvider(name, factory, priority)
}

// GetCallSystemProvider creates a CallSystem provider instance from the registry.
func GetCallSystemProvider(name string, opts ...ProviderOption) (callsystem.CallSystem, error) {
	return core.GetCallSystemProvider(name, opts...)
}

// ListCallSystemProviders returns a list of all registered CallSystem provider names.
func ListCallSystemProviders() []string {
	return core.ListCallSystemProviders()
}

// HasCallSystemProvider returns true if a CallSystem provider with the given name is registered.
func HasCallSystemProvider(name string) bool {
	return core.HasCallSystemProvider(name)
}

// GetCallSystemProviderPriority returns the priority of the registered CallSystem provider.
func GetCallSystemProviderPriority(name string) int {
	return core.GetCallSystemProviderPriority(name)
}

// Gateway Provider Registry - delegates to omnivoice-core

// RegisterGatewayProvider registers a Gateway provider factory with the given name and priority.
// Higher priority values override lower priority registrations.
func RegisterGatewayProvider(name string, factory registry.GatewayProviderFactory, priority int) {
	core.RegisterGatewayProvider(name, factory, priority)
}

// GetGatewayProvider creates a Gateway provider instance from the registry.
func GetGatewayProvider(name string, opts ...ProviderOption) (registry.Gateway, error) {
	return core.GetGatewayProvider(name, opts...)
}

// ListGatewayProviders returns a list of all registered Gateway provider names.
func ListGatewayProviders() []string {
	return core.ListGatewayProviders()
}

// HasGatewayProvider returns true if a Gateway provider with the given name is registered.
func HasGatewayProvider(name string) bool {
	return core.HasGatewayProvider(name)
}

// GetGatewayProviderPriority returns the priority of the registered Gateway provider.
func GetGatewayProviderPriority(name string) int {
	return core.GetGatewayProviderPriority(name)
}

// Realtime Provider Registry - delegates to omnivoice-core

// RegisterRealtimeProvider registers a Realtime provider factory with the given name and priority.
// Higher priority values override lower priority registrations.
func RegisterRealtimeProvider(name string, factory registry.RealtimeProviderFactory, priority int) {
	core.RegisterRealtimeProvider(name, factory, priority)
}

// GetRealtimeProvider creates a Realtime provider instance from the registry.
//
// Realtime providers enable native voice-to-voice conversations with low latency
// (~100-300ms). Available providers include:
//   - "openai": OpenAI Realtime API (~100ms latency)
//   - "gemini": Google Gemini Live API (~200ms latency)
//
// Example:
//
//	provider, err := omnivoice.GetRealtimeProvider("openai",
//	    omnivoice.WithAPIKey(os.Getenv("OPENAI_API_KEY")))
func GetRealtimeProvider(name string, opts ...ProviderOption) (registry.RealtimeProvider, error) {
	return core.GetRealtimeProvider(name, opts...)
}

// ListRealtimeProviders returns a list of all registered Realtime provider names.
func ListRealtimeProviders() []string {
	return core.ListRealtimeProviders()
}

// HasRealtimeProvider returns true if a Realtime provider with the given name is registered.
func HasRealtimeProvider(name string) bool {
	return core.HasRealtimeProvider(name)
}

// GetRealtimeProviderPriority returns the priority of the registered Realtime provider.
func GetRealtimeProviderPriority(name string) int {
	return core.GetRealtimeProviderPriority(name)
}
