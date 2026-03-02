// Package all imports and registers all omnivoice providers.
// Import this package using a blank import to register all providers:
//
//	import _ "github.com/plexusone/omnivoice/providers/all"
//
// After importing, use omnivoice.GetTTSProvider or omnivoice.GetSTTProvider:
//
//	provider, err := omnivoice.GetTTSProvider("elevenlabs", omnivoice.WithAPIKey(apiKey))
package all

import (
	"fmt"

	"github.com/plexusone/omnivoice"

	// Provider packages
	elevenlabsstt "github.com/plexusone/elevenlabs-go/omnivoice/stt"
	elevenlabstts "github.com/plexusone/elevenlabs-go/omnivoice/tts"
	deepgramstt "github.com/plexusone/omnivoice-deepgram/omnivoice/stt"
	deepgramtts "github.com/plexusone/omnivoice-deepgram/omnivoice/tts"
	openaiomni "github.com/plexusone/omnivoice-openai/omnivoice"
	twiliostt "github.com/plexusone/omnivoice-twilio/stt"
	twiliotts "github.com/plexusone/omnivoice-twilio/tts"
)

func init() {
	// Register OpenAI providers
	omnivoice.RegisterTTSProvider("openai", func(config omnivoice.ProviderConfig) (omnivoice.TTSProvider, error) {
		if config.APIKey == "" {
			return nil, fmt.Errorf("openai: API key is required")
		}
		return openaiomni.NewTTSProvider(config.APIKey), nil
	})

	omnivoice.RegisterSTTProvider("openai", func(config omnivoice.ProviderConfig) (omnivoice.STTProvider, error) {
		if config.APIKey == "" {
			return nil, fmt.Errorf("openai: API key is required")
		}
		return openaiomni.NewSTTProvider(config.APIKey), nil
	})

	// Register ElevenLabs providers
	omnivoice.RegisterTTSProvider("elevenlabs", func(config omnivoice.ProviderConfig) (omnivoice.TTSProvider, error) {
		var opts []elevenlabstts.Option
		if config.APIKey != "" {
			opts = append(opts, elevenlabstts.WithAPIKey(config.APIKey))
		}
		if config.BaseURL != "" {
			opts = append(opts, elevenlabstts.WithBaseURL(config.BaseURL))
		}
		return elevenlabstts.New(opts...)
	})

	omnivoice.RegisterSTTProvider("elevenlabs", func(config omnivoice.ProviderConfig) (omnivoice.STTProvider, error) {
		var opts []elevenlabsstt.Option
		if config.APIKey != "" {
			opts = append(opts, elevenlabsstt.WithAPIKey(config.APIKey))
		}
		if config.BaseURL != "" {
			opts = append(opts, elevenlabsstt.WithBaseURL(config.BaseURL))
		}
		return elevenlabsstt.New(opts...)
	})

	// Register Deepgram providers
	omnivoice.RegisterTTSProvider("deepgram", func(config omnivoice.ProviderConfig) (omnivoice.TTSProvider, error) {
		var opts []deepgramtts.Option
		if config.APIKey != "" {
			opts = append(opts, deepgramtts.WithAPIKey(config.APIKey))
		}
		return deepgramtts.New(opts...)
	})

	omnivoice.RegisterSTTProvider("deepgram", func(config omnivoice.ProviderConfig) (omnivoice.STTProvider, error) {
		var opts []deepgramstt.Option
		if config.APIKey != "" {
			opts = append(opts, deepgramstt.WithAPIKey(config.APIKey))
		}
		return deepgramstt.New(opts...)
	})

	// Register Twilio providers
	// Note: Twilio doesn't require an API key for TTS/STT - it uses TwiML within calls
	omnivoice.RegisterTTSProvider("twilio", func(config omnivoice.ProviderConfig) (omnivoice.TTSProvider, error) {
		var opts []twiliotts.Option
		// Check for voice in extensions
		if voice, ok := config.Extensions["twilio.voice"].(string); ok && voice != "" {
			opts = append(opts, twiliotts.WithVoice(voice))
		}
		// Check for language in extensions
		if lang, ok := config.Extensions["twilio.language"].(string); ok && lang != "" {
			opts = append(opts, twiliotts.WithLanguage(lang))
		}
		return twiliotts.New(opts...)
	})

	omnivoice.RegisterSTTProvider("twilio", func(config omnivoice.ProviderConfig) (omnivoice.STTProvider, error) {
		var opts []twiliostt.Option
		// Check for language in extensions
		if lang, ok := config.Extensions["twilio.language"].(string); ok && lang != "" {
			opts = append(opts, twiliostt.WithLanguage(lang))
		}
		// Check for speech model in extensions
		if model, ok := config.Extensions["twilio.speech_model"].(string); ok && model != "" {
			opts = append(opts, twiliostt.WithSpeechModel(model))
		}
		// Check for profanity filter in extensions
		if filter, ok := config.Extensions["twilio.profanity_filter"].(bool); ok {
			opts = append(opts, twiliostt.WithProfanityFilter(filter))
		}
		return twiliostt.New(opts...)
	})
}
