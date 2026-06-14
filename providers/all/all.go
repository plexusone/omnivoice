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
	deepgramstt "github.com/plexusone/omni-deepgram/omnivoice/stt"
	deepgramtts "github.com/plexusone/omni-deepgram/omnivoice/tts"
	googlerealtime "github.com/plexusone/omni-google/omnivoice/realtime"
	openaiomni "github.com/plexusone/omni-openai/omnivoice"
	openairealtime "github.com/plexusone/omni-openai/omnivoice/realtime"
	telnyxcallsystem "github.com/plexusone/omni-telnyx/omnivoice/callsystem"
	twiliocallsystem "github.com/plexusone/omni-twilio/omnivoice/callsystem"
	twiliostt "github.com/plexusone/omni-twilio/omnivoice/stt"
	twiliotts "github.com/plexusone/omni-twilio/omnivoice/tts"
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

	// Register Twilio CallSystem
	omnivoice.RegisterCallSystemProvider("twilio", func(config omnivoice.ProviderConfig) (omnivoice.CallSystem, error) {
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
	})

	// Register Telnyx CallSystem
	omnivoice.RegisterCallSystemProvider("telnyx", func(config omnivoice.ProviderConfig) (omnivoice.CallSystem, error) {
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
	})

	// Register OpenAI Realtime provider
	omnivoice.RegisterRealtimeProvider("openai-realtime", func(config omnivoice.ProviderConfig) (omnivoice.RealtimeProvider, error) {
		if config.APIKey == "" {
			return nil, fmt.Errorf("openai-realtime: API key is required")
		}
		var opts []openairealtime.Option
		// Check for voice in extensions
		if voice, ok := config.Extensions["voice"].(string); ok && voice != "" {
			opts = append(opts, openairealtime.WithVoice(voice))
		}
		// Check for instructions in extensions
		if instructions, ok := config.Extensions["instructions"].(string); ok && instructions != "" {
			opts = append(opts, openairealtime.WithInstructions(instructions))
		}
		return openairealtime.NewProvider(config.APIKey, opts...), nil
	})

	// Register Gemini Live provider
	omnivoice.RegisterRealtimeProvider("gemini-live", func(config omnivoice.ProviderConfig) (omnivoice.RealtimeProvider, error) {
		if config.APIKey == "" {
			return nil, fmt.Errorf("gemini-live: API key is required")
		}
		var opts []googlerealtime.Option
		// Check for voice in extensions
		if voice, ok := config.Extensions["voice"].(string); ok && voice != "" {
			opts = append(opts, googlerealtime.WithVoice(voice))
		}
		// Check for instructions in extensions
		if instructions, ok := config.Extensions["instructions"].(string); ok && instructions != "" {
			opts = append(opts, googlerealtime.WithInstructions(instructions))
		}
		return googlerealtime.NewRealtimeProvider(config.APIKey, opts...), nil
	})
}
