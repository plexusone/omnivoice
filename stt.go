// Package omnivoice provides a unified interface for speech-to-text and text-to-speech.
// This is the batteries-included package that imports all providers.
// For a minimal dependency footprint, use github.com/plexusone/omnivoice-core instead.
package omnivoice

import (
	"github.com/plexusone/omnivoice-core/stt"
)

// Re-export STT types from omnivoice-core
type (
	// STTProvider defines the interface for STT providers.
	STTProvider = stt.Provider

	// STTStreamingProvider extends Provider with streaming support.
	STTStreamingProvider = stt.StreamingProvider

	// TranscriptionConfig configures a STT transcription request.
	TranscriptionConfig = stt.TranscriptionConfig

	// TranscriptionResult contains the result of a STT transcription.
	TranscriptionResult = stt.TranscriptionResult

	// Segment represents a transcription segment.
	Segment = stt.Segment

	// Word represents a word with timing information.
	Word = stt.Word

	// StreamEvent represents a streaming transcription event.
	StreamEvent = stt.StreamEvent

	// STTClient is the multi-provider STT client.
	STTClient = stt.Client
)

// Re-export STT errors
var (
	ErrNoAvailableProvider = stt.ErrNoAvailableProvider
	ErrStreamingNotSupported = stt.ErrStreamingNotSupported
	ErrInvalidAudio = stt.ErrInvalidAudio
	ErrInvalidConfig = stt.ErrInvalidConfig
	ErrAudioTooLong = stt.ErrAudioTooLong
	ErrAudioTooShort = stt.ErrAudioTooShort
	ErrRateLimited = stt.ErrRateLimited
	ErrQuotaExceeded = stt.ErrQuotaExceeded
	ErrUnsupportedLanguage = stt.ErrUnsupportedLanguage
	ErrUnsupportedFormat = stt.ErrUnsupportedFormat
	ErrStreamClosed = stt.ErrStreamClosed
)

// Re-export STT functions
var NewSTTClient = stt.NewClient
