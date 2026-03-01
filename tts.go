// Package omnivoice provides a unified interface for speech-to-text and text-to-speech.
// This is the batteries-included package that imports all providers.
// For a minimal dependency footprint, use github.com/plexusone/omnivoice-core instead.
package omnivoice

import (
	"github.com/plexusone/omnivoice-core/tts"
)

// Re-export TTS types from omnivoice-core
type (
	// TTSProvider defines the interface for TTS providers.
	TTSProvider = tts.Provider

	// TTSStreamingProvider extends Provider with input streaming support.
	TTSStreamingProvider = tts.StreamingProvider

	// Voice represents a voice configuration for TTS.
	Voice = tts.Voice

	// SynthesisConfig configures a TTS synthesis request.
	SynthesisConfig = tts.SynthesisConfig

	// SynthesisResult contains the result of a TTS synthesis.
	SynthesisResult = tts.SynthesisResult

	// StreamChunk represents a chunk of streaming audio.
	TTSStreamChunk = tts.StreamChunk

	// TTSClient is the multi-provider TTS client.
	TTSClient = tts.Client
)

// Re-export TTS errors
var (
	ErrTTSNoAvailableProvider = tts.ErrNoAvailableProvider
	ErrVoiceNotFound          = tts.ErrVoiceNotFound
	ErrTTSInvalidConfig       = tts.ErrInvalidConfig
	ErrTTSRateLimited         = tts.ErrRateLimited
	ErrTTSQuotaExceeded       = tts.ErrQuotaExceeded
	ErrTTSStreamClosed        = tts.ErrStreamClosed
)

// Re-export TTS functions
var NewTTSClient = tts.NewClient
