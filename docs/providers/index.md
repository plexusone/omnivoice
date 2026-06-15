# Providers

OmniVoice supports multiple providers for each capability. Choose based on your requirements.

## Provider Matrix

### STT/TTS Providers

| Provider | TTS | STT | Streaming | Notes |
|----------|-----|-----|-----------|-------|
| [OpenAI](openai.md) | ✓ | ✓ | ✓ | Whisper for STT, TTS-1 for TTS |
| [ElevenLabs](elevenlabs.md) | ✓ | ✓ | ✓ | Premium voice synthesis |
| [Deepgram](deepgram.md) | ✓ | ✓ | ✓ | Real-time transcription |

### Native Voice-to-Voice (Realtime)

| Provider | Latency | Registry Name | Notes |
|----------|---------|---------------|-------|
| OpenAI Realtime | ~100ms | `openai` | GPT-4o voice model |
| Gemini Live | ~200ms | `gemini` | Gemini 2.0 Flash |

### Voice Gateway Providers (PSTN)

| Provider | Voice | SMS | Audio Format | Notes |
|----------|-------|-----|--------------|-------|
| [Twilio](twilio.md) | ✓ | ✓ | mulaw 8kHz | Most popular, Media Streams |
| [Telnyx](telnyx.md) | ✓ | ✓ | L16 16kHz | Competitive pricing |
| Vonage | ✓ | ✓ | L16 16kHz | JWT auth, NCCO call control |
| Plivo | ✓ | ✓ | L16 16kHz | Good international coverage |

### Voice Gateway Providers (WebRTC)

| Provider | Protocol | Latency | Notes |
|----------|----------|---------|-------|
| LiveKit | WebRTC | <200ms | Browser/mobile, open source |

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

### PSTN Voice Gateway

| Provider | Call Quality | Geographic Coverage | Pricing |
|----------|-------------|---------------------|---------|
| **Twilio** | Excellent | Global | Pay-per-use |
| **Telnyx** | Excellent | Global | Lower cost |
| **Vonage** | Excellent | Global | Competitive |
| **Plivo** | Excellent | Global | Simple pricing |

### WebRTC Voice Gateway

| Provider | Latency | Use Case | Pricing |
|----------|---------|----------|---------|
| **LiveKit** | <200ms | Custom web/mobile apps | Infrastructure only |

### Native Voice-to-Voice (Realtime)

| Provider | Latency | Voice Quality | Best For |
|----------|---------|---------------|----------|
| **OpenAI Realtime** | ~100ms | Excellent | Conversational agents |
| **Gemini Live** | ~200ms | Excellent | Multi-turn dialogue |

## Quick Comparison

### Lowest Latency (Native Voice-to-Voice)

```go
// For voice agents requiring <500ms response time
rt, _ := omnivoice.GetRealtimeProvider("openai", ...)  // ~100ms
rt, _ := omnivoice.GetRealtimeProvider("gemini", ...)  // ~200ms
```

### Lowest Latency (Traditional Pipeline)

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

For selective imports, use the provider modules directly:

```go
// TTS/STT providers
import (
    _ "github.com/plexusone/omni-openai/omnivoice"
    _ "github.com/plexusone/elevenlabs-go/omnivoice/tts"
    _ "github.com/plexusone/elevenlabs-go/omnivoice/stt"
    _ "github.com/plexusone/omni-deepgram/omnivoice/tts"
    _ "github.com/plexusone/omni-deepgram/omnivoice/stt"
)

// Telephony providers
import (
    _ "github.com/plexusone/omni-twilio/omnivoice/callsystem"
    _ "github.com/plexusone/omni-telnyx/omnivoice/callsystem"
)

// Realtime (native voice-to-voice) providers
import (
    _ "github.com/plexusone/omni-openai/omnivoice/realtime"
    _ "github.com/plexusone/omni-google/omnivoice/realtime"
)
```

## Environment Variables

### STT/TTS Providers

| Provider | Required Variables |
|----------|-------------------|
| OpenAI | `OPENAI_API_KEY` |
| ElevenLabs | `ELEVENLABS_API_KEY` |
| Deepgram | `DEEPGRAM_API_KEY` |

### Native Voice-to-Voice (Realtime)

| Provider | Required Variables |
|----------|-------------------|
| OpenAI Realtime | `OPENAI_API_KEY` |
| Gemini Live | `GOOGLE_API_KEY` |

### PSTN Voice Gateway

| Provider | Required Variables |
|----------|-------------------|
| Twilio | `TWILIO_ACCOUNT_SID`, `TWILIO_AUTH_TOKEN` |
| Telnyx | `TELNYX_API_KEY`, `TELNYX_CONNECTION_ID` |
| Vonage | `VONAGE_APPLICATION_ID`, `VONAGE_PRIVATE_KEY` |
| Plivo | `PLIVO_AUTH_ID`, `PLIVO_AUTH_TOKEN` |

### WebRTC Voice Gateway

| Provider | Required Variables |
|----------|-------------------|
| LiveKit | `LIVEKIT_URL`, `LIVEKIT_API_KEY`, `LIVEKIT_API_SECRET` |

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

### STT/TTS

- [OpenAI](openai.md) - OpenAI TTS and Whisper STT
- [ElevenLabs](elevenlabs.md) - Premium voice synthesis
- [Deepgram](deepgram.md) - Real-time transcription

### PSTN Voice Gateway

- [Twilio](twilio.md) - Phone calls and SMS
- [Telnyx](telnyx.md) - Phone calls and SMS
- Vonage - JWT auth, NCCO call control
- Plivo - Simple pricing, international coverage

### WebRTC Voice Gateway

- [LiveKit](https://github.com/plexusone/omni-livekit) - Browser/mobile apps
