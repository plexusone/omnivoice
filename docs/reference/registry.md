# Provider Registry

OmniVoice uses a registry pattern for runtime provider discovery and instantiation.

## Overview

The registry allows you to:

- Register providers at init time via blank imports
- Get providers by name with configuration
- List available providers
- Check if a provider exists

## Registration

### Import All Providers

```go
import _ "github.com/plexusone/omnivoice/providers/all"
```

### Import Specific Providers

For selective imports, use the provider modules directly:

```go
import (
    _ "github.com/plexusone/omni-openai/omnivoice"
    _ "github.com/plexusone/elevenlabs-go/omnivoice/tts"
    _ "github.com/plexusone/omni-deepgram/omnivoice/tts"
    _ "github.com/plexusone/omni-twilio/omnivoice/callsystem"
    _ "github.com/plexusone/omni-telnyx/omnivoice/callsystem"
)
```

## TTS Registry

### Get Provider

```go
tts, err := omnivoice.GetTTSProvider("elevenlabs",
    omnivoice.WithAPIKey(apiKey),
)
if err != nil {
    log.Fatal(err)
}
```

### List Providers

```go
providers := omnivoice.ListTTSProviders()
fmt.Println(providers)  // ["openai", "elevenlabs", "deepgram"]
```

### Check Provider

```go
if omnivoice.HasTTSProvider("elevenlabs") {
    tts, _ := omnivoice.GetTTSProvider("elevenlabs", ...)
}
```

### Register Custom Provider

```go
omnivoice.RegisterTTSProvider("custom", func(config omnivoice.ProviderConfig) (omnivoice.TTSProvider, error) {
    return &myCustomTTS{apiKey: config.APIKey}, nil
}, omnivoice.PriorityThick)
```

## STT Registry

### Get Provider

```go
stt, err := omnivoice.GetSTTProvider("deepgram",
    omnivoice.WithAPIKey(apiKey),
)
if err != nil {
    log.Fatal(err)
}
```

### List Providers

```go
providers := omnivoice.ListSTTProviders()
fmt.Println(providers)  // ["openai", "elevenlabs", "deepgram"]
```

### Check Provider

```go
if omnivoice.HasSTTProvider("deepgram") {
    stt, _ := omnivoice.GetSTTProvider("deepgram", ...)
}
```

### Register Custom Provider

```go
omnivoice.RegisterSTTProvider("custom", func(config omnivoice.ProviderConfig) (omnivoice.STTProvider, error) {
    return &myCustomSTT{apiKey: config.APIKey}, nil
}, omnivoice.PriorityThick)
```

## CallSystem Registry

### Get Provider

```go
cs, err := omnivoice.GetCallSystemProvider("twilio",
    omnivoice.WithAccountSID(accountSID),
    omnivoice.WithAuthToken(authToken),
    omnivoice.WithPhoneNumber("+15551234567"),
    omnivoice.WithWebhookURL("https://example.com/webhook"),
)
if err != nil {
    log.Fatal(err)
}
```

### List Providers

```go
providers := omnivoice.ListCallSystemProviders()
fmt.Println(providers)  // ["twilio", "telnyx"]
```

### Check Provider

```go
if omnivoice.HasCallSystemProvider("twilio") {
    cs, _ := omnivoice.GetCallSystemProvider("twilio", ...)
}
```

### Register Custom Provider

```go
omnivoice.RegisterCallSystemProvider("custom", func(config omnivoice.ProviderConfig) (omnivoice.CallSystem, error) {
    return &myCustomCallSystem{
        accountSID: config.Extensions["accountSID"].(string),
        authToken:  config.APIKey,
    }, nil
}, omnivoice.PriorityThick)
```

## Gateway Registry

Voice gateways handle full-duplex voice calls via telephony providers.

### Get Provider

```go
gw, err := omnivoice.GetGatewayProvider("twilio",
    omnivoice.WithAccountSID(accountSID),
    omnivoice.WithAuthToken(authToken),
    omnivoice.WithPublicURL("https://your-server.com"),
    omnivoice.WithListenAddr(":8080"),
)
if err != nil {
    log.Fatal(err)
}
```

### List Providers

```go
providers := omnivoice.ListGatewayProviders()
fmt.Println(providers)  // ["twilio", "telnyx"]
```

### Check Provider

```go
if omnivoice.HasGatewayProvider("twilio") {
    gw, _ := omnivoice.GetGatewayProvider("twilio", ...)
}
```

## Realtime Registry

Realtime providers enable native voice-to-voice with ~100-200ms latency.

### Get Provider

```go
rt, err := omnivoice.GetRealtimeProvider("openai",
    omnivoice.WithAPIKey(apiKey),
    omnivoice.WithVoice("alloy"),
    omnivoice.WithInstructions("You are a helpful assistant."),
)
if err != nil {
    log.Fatal(err)
}
```

### List Providers

```go
providers := omnivoice.ListRealtimeProviders()
fmt.Println(providers)  // ["openai", "gemini"]
```

### Check Provider

```go
if omnivoice.HasRealtimeProvider("openai") {
    rt, _ := omnivoice.GetRealtimeProvider("openai", ...)
}
```

## Provider Options

### Common Options

```go
// API key (OpenAI, ElevenLabs, Deepgram, Telnyx)
omnivoice.WithAPIKey(key)

// Account credentials (Twilio)
omnivoice.WithAccountSID(sid)
omnivoice.WithAuthToken(token)

// Phone number (Twilio, Telnyx)
omnivoice.WithPhoneNumber(number)

// Webhook URL (Twilio, Telnyx)
omnivoice.WithWebhookURL(url)
```

### Extensions

For provider-specific configuration:

```go
tts, err := omnivoice.GetTTSProvider("elevenlabs",
    omnivoice.WithAPIKey(apiKey),
    omnivoice.WithExtensions(map[string]any{
        "model_id": "eleven_turbo_v2_5",
    }),
)
```

## ProviderConfig

The configuration structure (from `omnivoice-core/registry`):

```go
type ProviderConfig struct {
    APIKey     string         // Authentication key
    BaseURL    string         // Custom API endpoint
    Extensions map[string]any // Provider-specific configuration
}
```

Provider-specific options set Extensions values:

```go
omnivoice.WithAccountSID(sid)    // Extensions["accountSID"]
omnivoice.WithAuthToken(token)   // Extensions["authToken"]
omnivoice.WithPhoneNumber(num)   // Extensions["phoneNumber"]
omnivoice.WithWebhookURL(url)    // Extensions["webhookURL"]
omnivoice.WithVoice(voice)       // Extensions["voice"]
omnivoice.WithInstructions(text) // Extensions["instructions"]
```

## Priority System

Providers register with a priority level:

| Priority | Constant | Description |
|----------|----------|-------------|
| 0 | `PriorityThin` | Stdlib-only implementations |
| 10 | `PriorityThick` | Official SDK implementations |

Higher priority providers override lower priority registrations for the same name.

## Dynamic Provider Selection

Select providers at runtime:

```go
func getProvider(name string, apiKey string) (omnivoice.TTSProvider, error) {
    if !omnivoice.HasTTSProvider(name) {
        return nil, fmt.Errorf("unknown provider: %s", name)
    }
    return omnivoice.GetTTSProvider(name, omnivoice.WithAPIKey(apiKey))
}

// Use environment variable to select provider
providerName := os.Getenv("TTS_PROVIDER")
if providerName == "" {
    providerName = "openai"  // Default
}

tts, err := getProvider(providerName, os.Getenv("TTS_API_KEY"))
```

## Configuration from File

Load provider configuration from YAML/JSON:

```go
type Config struct {
    TTS struct {
        Provider string `yaml:"provider"`
        APIKey   string `yaml:"api_key"`
    } `yaml:"tts"`
    STT struct {
        Provider string `yaml:"provider"`
        APIKey   string `yaml:"api_key"`
    } `yaml:"stt"`
}

func loadProviders(cfg *Config) (omnivoice.TTSProvider, omnivoice.STTProvider, error) {
    tts, err := omnivoice.GetTTSProvider(cfg.TTS.Provider,
        omnivoice.WithAPIKey(cfg.TTS.APIKey))
    if err != nil {
        return nil, nil, err
    }

    stt, err := omnivoice.GetSTTProvider(cfg.STT.Provider,
        omnivoice.WithAPIKey(cfg.STT.APIKey))
    if err != nil {
        return nil, nil, err
    }

    return tts, stt, nil
}
```

## Thread Safety

The registry is thread-safe:

- Registration typically happens at init time
- Get/List/Has operations are safe for concurrent use
- Providers returned are independent instances

## Error Handling

```go
tts, err := omnivoice.GetTTSProvider("unknown", ...)
if err != nil {
    // err: "unknown TTS provider: unknown"
}

// Check available providers
if !omnivoice.HasTTSProvider("custom") {
    log.Fatal("Custom provider not registered. Did you import it?")
}
```

## Next Steps

- [Multi-Provider Failover](failover.md) - Reliability patterns
- [Provider Comparison](../providers/index.md) - Choose the right provider
