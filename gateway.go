// Package omnivoice provides a unified interface for speech-to-text and text-to-speech.
// This is the batteries-included package that imports all providers.
// For a minimal dependency footprint, use github.com/plexusone/omnivoice-core instead.
package omnivoice

import (
	"context"

	"github.com/plexusone/omnivoice-core/gateway"
	"github.com/plexusone/omnivoice-core/registry"

	// Provider packages for gateway-specific options
	telnyxGateway "github.com/plexusone/omni-telnyx/omnivoice/gateway"
	twilioGateway "github.com/plexusone/omni-twilio/omnivoice/gateway"
)

// Re-export Gateway types from omnivoice-core
type (
	// Gateway defines the interface for voice gateway providers.
	Gateway = gateway.Gateway

	// GatewayMinimal is the minimal interface used by the registry.
	// Use type assertion to access the full Gateway interface.
	GatewayMinimal = registry.Gateway

	// GatewayConfig provides common configuration for voice gateways.
	GatewayConfig = gateway.Config

	// GatewaySession represents an active voice conversation session.
	GatewaySession = gateway.Session

	// GatewayCallInfo contains information about a call.
	GatewayCallInfo = gateway.CallInfo

	// GatewayCallHandler is called when a new call is received.
	GatewayCallHandler = gateway.CallHandler

	// GatewayTurn represents a single conversation turn.
	GatewayTurn = gateway.Turn

	// GatewayToolCall represents a tool invocation during conversation.
	GatewayToolCall = gateway.ToolCall

	// GatewayMetrics contains session performance metrics.
	GatewayMetrics = gateway.Metrics

	// GatewayEvent represents a session event.
	GatewayEvent = gateway.Event

	// GatewayEventType identifies the type of session event.
	GatewayEventType = gateway.EventType

	// RealtimeConfig configures a realtime provider for voice-to-voice.
	RealtimeConfig = gateway.RealtimeConfig

	// RealtimeProviderFactory creates realtime providers from configuration.
	RealtimeProviderFactory = gateway.RealtimeProviderFactory

	// LLMProvider defines the interface for LLM integration with voice gateways.
	LLMProvider = gateway.LLMProvider
)

// Re-export Gateway provider name constants
const (
	GatewayProviderTwilio  = gateway.ProviderTwilio
	GatewayProviderTelnyx  = gateway.ProviderTelnyx
	GatewayProviderVonage  = gateway.ProviderVonage
	GatewayProviderPlivo   = gateway.ProviderPlivo
	GatewayProviderLiveKit = gateway.ProviderLiveKit
)

// Re-export Pipeline mode constants
const (
	PipelineModeText     = gateway.PipelineModeText
	PipelineModeRealtime = gateway.PipelineModeRealtime
)

// Re-export Gateway event type constants
const (
	EventSessionStarted   = gateway.EventSessionStarted
	EventSessionEnded     = gateway.EventSessionEnded
	EventUserSpeechStart  = gateway.EventUserSpeechStart
	EventUserSpeechEnd    = gateway.EventUserSpeechEnd
	EventUserTranscript   = gateway.EventUserTranscript
	EventAgentThinking    = gateway.EventAgentThinking
	EventAgentSpeechStart = gateway.EventAgentSpeechStart
	EventAgentSpeechEnd   = gateway.EventAgentSpeechEnd
	EventAgentTranscript  = gateway.EventAgentTranscript
	EventToolCall         = gateway.EventToolCall
	EventInterruption     = gateway.EventInterruption
	EventError            = gateway.EventError
	EventAudioReceived    = gateway.EventAudioReceived
	EventAudioSent        = gateway.EventAudioSent
)

// Re-export audio format types and constants
type AudioFormat = gateway.AudioFormat

var (
	AudioFormatTwilio       = gateway.AudioFormatTwilio
	AudioFormatTelnyx       = gateway.AudioFormatTelnyx
	AudioFormatOpenAI       = gateway.AudioFormatOpenAI
	AudioFormatGeminiInput  = gateway.AudioFormatGeminiInput
	AudioFormatGeminiOutput = gateway.AudioFormatGeminiOutput
)

// ToolDefinition defines a tool that can be called during voice conversations.
// This is a generic type that works with any gateway provider.
type ToolDefinition struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters,omitempty"`
}

// ToolHandler is a function that handles tool calls during voice conversations.
type ToolHandler func(ctx context.Context, args map[string]any) (string, error)

// WithGatewayTools returns a provider option that configures tools for a gateway.
// The provider parameter specifies which gateway provider to configure ("twilio" or "telnyx").
func WithGatewayTools(provider string, tools []ToolDefinition) registry.ProviderOption {
	switch provider {
	case "twilio":
		twilioTools := make([]twilioGateway.ToolDefinition, len(tools))
		for i, t := range tools {
			twilioTools[i] = twilioGateway.ToolDefinition{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  t.Parameters,
			}
		}
		return twilioGateway.WithTools(twilioTools)
	case "telnyx":
		telnyxTools := make([]telnyxGateway.ToolDefinition, len(tools))
		for i, t := range tools {
			telnyxTools[i] = telnyxGateway.ToolDefinition{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  t.Parameters,
			}
		}
		return telnyxGateway.WithTools(telnyxTools)
	default:
		// Return no-op option for unknown providers
		return func(c *registry.ProviderConfig) {}
	}
}

// WithGatewayToolHandlers returns a provider option that configures tool handlers for a gateway.
// The provider parameter specifies which gateway provider to configure ("twilio" or "telnyx").
func WithGatewayToolHandlers(provider string, handlers map[string]ToolHandler) registry.ProviderOption {
	switch provider {
	case "twilio":
		twilioHandlers := make(map[string]twilioGateway.ToolHandler, len(handlers))
		for name, handler := range handlers {
			h := handler // capture for closure
			twilioHandlers[name] = func(ctx context.Context, args map[string]any) (string, error) {
				return h(ctx, args)
			}
		}
		return twilioGateway.WithToolHandlers(twilioHandlers)
	case "telnyx":
		telnyxHandlers := make(map[string]telnyxGateway.ToolHandler, len(handlers))
		for name, handler := range handlers {
			h := handler // capture for closure
			telnyxHandlers[name] = func(ctx context.Context, args map[string]any) (string, error) {
				return h(ctx, args)
			}
		}
		return telnyxGateway.WithToolHandlers(telnyxHandlers)
	default:
		// Return no-op option for unknown providers
		return func(c *registry.ProviderConfig) {}
	}
}

// WithRealtimeFactory returns a provider option that configures a realtime provider factory.
// This is used for native voice-to-voice mode with gateways.
func WithRealtimeFactory(factory gateway.RealtimeProviderFactory) registry.ProviderOption {
	// Currently only Twilio supports realtime factory configuration
	return twilioGateway.WithRealtimeProviderFactory(factory)
}

// WithGatewayRealtimeConfig returns a provider option that configures realtime settings.
// This is used for native voice-to-voice mode with gateways.
func WithGatewayRealtimeConfig(config *gateway.RealtimeConfig) registry.ProviderOption {
	// Currently only Twilio supports realtime config
	return twilioGateway.WithRealtimeConfig(config)
}
