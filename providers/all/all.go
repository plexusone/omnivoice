// Package all imports all omnivoice providers.
// Import this package using a blank import to register all providers:
//
//	import _ "github.com/plexusone/omnivoice/providers/all"
//
// Alternatively, import individual providers:
//
//	import _ "github.com/plexusone/omnivoice-openai/omnivoice"
//	import _ "github.com/plexusone/omnivoice-deepgram/omnivoice"
package all

import (
	// Import all provider packages for side effects (self-registration).
	// These imports ensure providers are available when using the omnivoice client.

	// OpenAI provider (Whisper STT, TTS)
	_ "github.com/plexusone/omnivoice-openai/omnivoice"

	// Deepgram provider (STT, TTS)
	_ "github.com/plexusone/omnivoice-deepgram/omnivoice/stt"
	_ "github.com/plexusone/omnivoice-deepgram/omnivoice/tts"

	// ElevenLabs provider (STT, TTS) - from elevenlabs-go SDK
	_ "github.com/plexusone/elevenlabs-go/omnivoice"

	// Twilio provider (STT, TTS, Transport, CallSystem)
	_ "github.com/plexusone/omnivoice-twilio/stt"
	_ "github.com/plexusone/omnivoice-twilio/tts"
)
