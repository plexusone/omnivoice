// Package omnivoice provides a unified interface for speech-to-text and text-to-speech.
// This is the batteries-included package that imports all providers.
// For a minimal dependency footprint, use github.com/plexusone/omnivoice-core instead.
package omnivoice

import (
	"github.com/plexusone/omnivoice-core/gateway"
	"github.com/plexusone/omnivoice-core/registry"
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
