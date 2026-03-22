# Providers

OmniVoice supports multiple providers for each capability. Choose based on your requirements.

## Provider Matrix

| Provider | TTS | STT | Streaming | CallSystem | SMS |
|----------|-----|-----|-----------|------------|-----|
| [OpenAI](openai.md) | ✓ | ✓ | ✓ | - | - |
| [ElevenLabs](elevenlabs.md) | ✓ | ✓ | ✓ | - | - |
| [Deepgram](deepgram.md) | ✓ | ✓ | ✓ | - | - |
| [Twilio](twilio.md) | - | - | - | ✓ | ✓ |
| [Telnyx](telnyx.md) | - | - | - | ✓ | ✓ |

## Comparison by Use Case

### Text-to-Speech (TTS)

| Provider | Latency | Voice Quality | Languages | Best For |
|----------|---------|---------------|-----------|----------|
| **ElevenLabs** | Low | Excellent | 29+ | Premium voice apps |
| **OpenAI** | Medium | Very Good | 50+ | Multi-language |
| **Deepgram** | Very Low | Good | 10+ | Real-time apps |

### Speech-to-Text (STT)

| Provider | Latency | Accuracy | Streaming | Best For |
|----------|---------|----------|-----------|----------|
| **Deepgram** | Very Low | Excellent | ✓ | Real-time transcription |
| **OpenAI** | High | Excellent | - | Batch transcription |
| **ElevenLabs** | Low | Very Good | ✓ | Voice apps with ElevenLabs TTS |

### Phone/SMS (CallSystem)

| Provider | Call Quality | Geographic Coverage | Pricing |
|----------|-------------|---------------------|---------|
| **Twilio** | Excellent | Global | Pay-per-use |
| **Telnyx** | Excellent | Global | Lower cost |

## Quick Comparison

### Lowest Latency (Real-Time)

```go
// For voice agents requiring <500ms response time
stt, _ := omnivoice.GetSTTProvider("deepgram", ...)  // ~200ms
tts, _ := omnivoice.GetTTSProvider("deepgram", ...)  // ~150ms
```

### Highest Quality

```go
// For premium voice applications
stt, _ := omnivoice.GetSTTProvider("deepgram", ...)     // Best accuracy
tts, _ := omnivoice.GetTTSProvider("elevenlabs", ...)   // Best voice quality
```

### Best Value

```go
// For cost-sensitive applications
stt, _ := omnivoice.GetSTTProvider("openai", ...)   // Good balance
tts, _ := omnivoice.GetTTSProvider("openai", ...)   // Good balance
callSystem, _ := omnivoice.GetCallSystemProvider("telnyx", ...)  // Lower rates
```

## Registering Providers

### Import All Providers

```go
import _ "github.com/plexusone/omnivoice/providers/all"
```

### Import Specific Providers

```go
// Only TTS/STT providers
import (
    _ "github.com/plexusone/omnivoice/providers/openai"
    _ "github.com/plexusone/omnivoice/providers/elevenlabs"
    _ "github.com/plexusone/omnivoice/providers/deepgram"
)

// Only telephony providers
import (
    _ "github.com/plexusone/omnivoice/providers/twilio"
    _ "github.com/plexusone/omnivoice/providers/telnyx"
)
```

## Environment Variables

| Provider | Required Variables |
|----------|-------------------|
| OpenAI | `OPENAI_API_KEY` |
| ElevenLabs | `ELEVENLABS_API_KEY` |
| Deepgram | `DEEPGRAM_API_KEY` |
| Twilio | `TWILIO_ACCOUNT_SID`, `TWILIO_AUTH_TOKEN` |
| Telnyx | `TELNYX_API_KEY` |

## Multi-Provider Patterns

### Failover

```go
// Try primary, fall back on failure
tts1, _ := omnivoice.GetTTSProvider("elevenlabs", ...)
tts2, _ := omnivoice.GetTTSProvider("openai", ...)

result, err := tts1.Synthesize(ctx, text, config)
if err != nil {
    result, err = tts2.Synthesize(ctx, text, config)
}
```

### Load Balancing

```go
// Round-robin across providers
providers := []omnivoice.TTSProvider{tts1, tts2, tts3}
provider := providers[requestCount % len(providers)]
```

### Cost Optimization

```go
// Use cheaper provider for non-critical requests
func selectProvider(priority string) omnivoice.TTSProvider {
    if priority == "high" {
        return elevenlabsProvider  // Premium quality
    }
    return openaiProvider  // Lower cost
}
```

## Next Steps

- [OpenAI](openai.md) - OpenAI TTS and Whisper STT
- [ElevenLabs](elevenlabs.md) - Premium voice synthesis
- [Deepgram](deepgram.md) - Real-time transcription
- [Twilio](twilio.md) - Phone calls and SMS
- [Telnyx](telnyx.md) - Phone calls and SMS
