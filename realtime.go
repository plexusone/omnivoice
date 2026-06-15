// Package omnivoice provides a unified interface for speech-to-text and text-to-speech.
// This is the batteries-included package that imports all providers.
// For a minimal dependency footprint, use github.com/plexusone/omnivoice-core instead.
package omnivoice

import (
	"github.com/plexusone/omnivoice-core/realtime"
	"github.com/plexusone/omnivoice-core/registry"
)

// Re-export Realtime types from omnivoice-core
type (
	// RealtimeProvider defines the interface for real-time voice-to-voice providers.
	// This is the full interface from realtime package with ProcessAudioStream.
	RealtimeProvider = realtime.Provider

	// RealtimeProviderMinimal is the minimal interface used by the registry.
	// Use type assertion to access the full RealtimeProvider interface.
	RealtimeProviderMinimal = registry.RealtimeProvider

	// ProcessConfig configures a real-time audio processing session.
	ProcessConfig = realtime.ProcessConfig

	// FunctionDeclaration describes a function the model can call.
	FunctionDeclaration = realtime.FunctionDeclaration

	// RealtimeAudioChunk represents a chunk of audio data from the model.
	RealtimeAudioChunk = realtime.AudioChunk

	// RealtimeTranscript represents a transcript update during a realtime conversation.
	RealtimeTranscript = realtime.Transcript

	// RealtimeClient is the multi-provider Realtime client.
	RealtimeClient = realtime.Client
)

// Re-export Realtime functions
var NewRealtimeClient = realtime.NewClient
