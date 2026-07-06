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
	"github.com/plexusone/omnivoice-core/registry"

	// STT/TTS provider packages (manual registration)
	elevenlabsstt "github.com/plexusone/elevenlabs-go/omnivoice/stt"
	elevenlabstts "github.com/plexusone/elevenlabs-go/omnivoice/tts"
	deepgramstt "github.com/plexusone/omni-deepgram/omnivoice/stt"
	deepgramtts "github.com/plexusone/omni-deepgram/omnivoice/tts"
	openaiomni "github.com/plexusone/omni-openai/omnivoice"
	twiliostt "github.com/plexusone/omni-twilio/omnivoice/stt"
	twiliotts "github.com/plexusone/omni-twilio/omnivoice/tts"

	// CallSystem provider packages (manual registration)
	telnyxcallsystem "github.com/plexusone/omni-telnyx/omnivoice/callsystem"
	twiliocallsystem "github.com/plexusone/omni-twilio/omnivoice/callsystem"

	// Realtime providers - auto-register via init()
	_ "github.com/plexusone/omni-deepgram/omnivoice/realtime"
	_ "github.com/plexusone/omni-google/omnivoice/realtime"
	_ "github.com/plexusone/omni-openai/omnivoice/realtime"

	// Gateway providers - auto-register via init()
	_ "github.com/plexusone/omni-telnyx/omnivoice/gateway"
	_ "github.com/plexusone/omni-twilio/omnivoice/gateway"
)

func init() {
	// Register OpenAI STT/TTS providers
	omnivoice.RegisterTTSProvider("openai", func(config registry.ProviderConfig) (omnivoice.TTSProvider, error) {
		if config.APIKey == "" {
			return nil, fmt.Errorf("openai: API key is required")
		}
		return openaiomni.NewTTSProvider(config.APIKey), nil
	}, omnivoice.PriorityThick)

	omnivoice.RegisterSTTProvider("openai", func(config registry.ProviderConfig) (omnivoice.STTProvider, error) {
		if config.APIKey == "" {
			return nil, fmt.Errorf("openai: API key is required")
		}
		return openaiomni.NewSTTProvider(config.APIKey), nil
	}, omnivoice.PriorityThick)

	// Register ElevenLabs providers
	omnivoice.RegisterTTSProvider("elevenlabs", func(config registry.ProviderConfig) (omnivoice.TTSProvider, error) {
		var opts []elevenlabstts.Option
		if config.APIKey != "" {
			opts = append(opts, elevenlabstts.WithAPIKey(config.APIKey))
		}
		if config.BaseURL != "" {
			opts = append(opts, elevenlabstts.WithBaseURL(config.BaseURL))
		}
		return elevenlabstts.New(opts...)
	}, omnivoice.PriorityThick)

	omnivoice.RegisterSTTProvider("elevenlabs", func(config registry.ProviderConfig) (omnivoice.STTProvider, error) {
		var opts []elevenlabsstt.Option
		if config.APIKey != "" {
			opts = append(opts, elevenlabsstt.WithAPIKey(config.APIKey))
		}
		if config.BaseURL != "" {
			opts = append(opts, elevenlabsstt.WithBaseURL(config.BaseURL))
		}
		return elevenlabsstt.New(opts...)
	}, omnivoice.PriorityThick)

	// Register Deepgram providers
	omnivoice.RegisterTTSProvider("deepgram", func(config registry.ProviderConfig) (omnivoice.TTSProvider, error) {
		var opts []deepgramtts.Option
		if config.APIKey != "" {
			opts = append(opts, deepgramtts.WithAPIKey(config.APIKey))
		}
		return deepgramtts.New(opts...)
	}, omnivoice.PriorityThick)

	omnivoice.RegisterSTTProvider("deepgram", func(config registry.ProviderConfig) (omnivoice.STTProvider, error) {
		var opts []deepgramstt.Option
		if config.APIKey != "" {
			opts = append(opts, deepgramstt.WithAPIKey(config.APIKey))
		}
		return deepgramstt.New(opts...)
	}, omnivoice.PriorityThick)

	// Register Twilio STT/TTS providers
	// Note: Twilio doesn't require an API key for TTS/STT - it uses TwiML within calls
	omnivoice.RegisterTTSProvider("twilio", func(config registry.ProviderConfig) (omnivoice.TTSProvider, error) {
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
	}, omnivoice.PriorityThick)

	omnivoice.RegisterSTTProvider("twilio", func(config registry.ProviderConfig) (omnivoice.STTProvider, error) {
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
	}, omnivoice.PriorityThick)

	// Register Twilio CallSystem
	omnivoice.RegisterCallSystemProvider("twilio", func(config registry.ProviderConfig) (omnivoice.CallSystem, error) {
		var opts []twiliocallsystem.Option

		// Account SID is required
		accountSID, _ := config.Extensions["accountSID"].(string)
		if accountSID == "" {
			return nil, fmt.Errorf("twilio: accountSID is required")
		}
		opts = append(opts, twiliocallsystem.WithAccountSID(accountSID))

		// Auth token from APIKey or extensions
		authToken := config.APIKey
		if authToken == "" {
			authToken, _ = config.Extensions["authToken"].(string)
		}
		if authToken == "" {
			return nil, fmt.Errorf("twilio: authToken is required")
		}
		opts = append(opts, twiliocallsystem.WithAuthToken(authToken))

		// Optional phone number
		if phoneNumber, ok := config.Extensions["phoneNumber"].(string); ok && phoneNumber != "" {
			opts = append(opts, twiliocallsystem.WithPhoneNumber(phoneNumber))
		}

		// Optional webhook URL
		if webhookURL, ok := config.Extensions["webhookURL"].(string); ok && webhookURL != "" {
			opts = append(opts, twiliocallsystem.WithWebhookURL(webhookURL))
		}

		return twiliocallsystem.New(opts...)
	}, omnivoice.PriorityThick)

	// Register Telnyx CallSystem
	omnivoice.RegisterCallSystemProvider("telnyx", func(config registry.ProviderConfig) (omnivoice.CallSystem, error) {
		var opts []telnyxcallsystem.Option

		// API key is required
		if config.APIKey == "" {
			return nil, fmt.Errorf("telnyx: API key is required")
		}
		opts = append(opts, telnyxcallsystem.WithAPIKey(config.APIKey))

		// Optional phone number
		if phoneNumber, ok := config.Extensions["phoneNumber"].(string); ok && phoneNumber != "" {
			opts = append(opts, telnyxcallsystem.WithPhoneNumber(phoneNumber))
		}

		// Optional webhook URL
		if webhookURL, ok := config.Extensions["webhookURL"].(string); ok && webhookURL != "" {
			opts = append(opts, telnyxcallsystem.WithWebhookURL(webhookURL))
		}

		// Optional connection ID
		if connectionID, ok := config.Extensions["connectionID"].(string); ok && connectionID != "" {
			opts = append(opts, telnyxcallsystem.WithConnectionID(connectionID))
		}

		return telnyxcallsystem.New(opts...)
	}, omnivoice.PriorityThick)

	// Note: Realtime providers (openai, gemini, deepgram) and Gateway providers (twilio, telnyx)
	// are auto-registered via side-effect imports above. They register with omnivoice-core
	// in their init() functions.
}
