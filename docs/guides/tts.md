# Text-to-Speech (TTS)

Convert text to natural-sounding speech using multiple providers.

## Quick Start

```go
provider, _ := omnivoice.GetTTSProvider("elevenlabs",
    omnivoice.WithAPIKey(apiKey))

result, _ := provider.Synthesize(ctx, "Hello, world!", omnivoice.SynthesisConfig{
    VoiceID: "pNInz6obpgDQGcFmaJgB",
})

// result.Audio contains the audio bytes
```

## Available Providers

| Provider | Registry Name | Quality | Latency | Best For |
|----------|---------------|---------|---------|----------|
| ElevenLabs | `"elevenlabs"` | Excellent | Medium | Natural voices, cloning |
| OpenAI | `"openai"` | Very Good | Low | General purpose |
| Deepgram | `"deepgram"` | Good | Very Low | Real-time, low latency |
| Twilio | `"twilio"` | Good | Low | Phone calls (TwiML) |

## Basic Synthesis

```go
package main

import (
    "context"
    "os"

    "github.com/plexusone/omnivoice"
    _ "github.com/plexusone/omnivoice/providers/all"
)

func main() {
    ctx := context.Background()

    provider, _ := omnivoice.GetTTSProvider("elevenlabs",
        omnivoice.WithAPIKey(os.Getenv("ELEVENLABS_API_KEY")))

    result, err := provider.Synthesize(ctx, "Welcome to OmniVoice!", omnivoice.SynthesisConfig{
        VoiceID:      "pNInz6obpgDQGcFmaJgB", // Adam
        OutputFormat: "mp3_44100_128",
        SampleRate:   44100,
    })
    if err != nil {
        panic(err)
    }

    os.WriteFile("output.mp3", result.Audio, 0644)
}
```

## Streaming TTS

For real-time applications, use streaming synthesis:

```go
stream, err := provider.SynthesizeStream(ctx, "This is a longer text that will be streamed...", omnivoice.SynthesisConfig{
    VoiceID: "pNInz6obpgDQGcFmaJgB",
})
if err != nil {
    panic(err)
}

// Read audio chunks as they arrive
for chunk := range stream {
    // Process chunk.Audio in real-time
    playAudio(chunk.Audio)
}
```

## Configuration Options

```go
config := omnivoice.SynthesisConfig{
    // Voice selection
    VoiceID: "voice-id",         // Provider-specific voice ID

    // Audio format
    OutputFormat: "mp3_44100_128", // Format string (provider-specific)
    SampleRate:   44100,           // Sample rate in Hz

    // Voice settings (ElevenLabs)
    Stability:       0.5,  // 0.0-1.0, lower = more expressive
    SimilarityBoost: 0.75, // 0.0-1.0, higher = closer to original

    // Provider-specific extensions
    Extensions: map[string]any{
        "model_id": "eleven_multilingual_v2",
    },
}
```

## Provider-Specific Examples

### ElevenLabs

```go
provider, _ := omnivoice.GetTTSProvider("elevenlabs",
    omnivoice.WithAPIKey(apiKey))

result, _ := provider.Synthesize(ctx, text, omnivoice.SynthesisConfig{
    VoiceID:      "pNInz6obpgDQGcFmaJgB",
    OutputFormat: "mp3_44100_128",
    Extensions: map[string]any{
        "model_id": "eleven_multilingual_v2",
    },
})
```

### OpenAI

```go
provider, _ := omnivoice.GetTTSProvider("openai",
    omnivoice.WithAPIKey(apiKey))

result, _ := provider.Synthesize(ctx, text, omnivoice.SynthesisConfig{
    VoiceID: "alloy", // alloy, echo, fable, onyx, nova, shimmer
    Extensions: map[string]any{
        "model": "tts-1-hd", // tts-1 or tts-1-hd
        "speed": 1.0,        // 0.25 to 4.0
    },
})
```

### Deepgram

```go
provider, _ := omnivoice.GetTTSProvider("deepgram",
    omnivoice.WithAPIKey(apiKey))

result, _ := provider.Synthesize(ctx, text, omnivoice.SynthesisConfig{
    VoiceID: "aura-asteria-en", // Aura voices
    Extensions: map[string]any{
        "encoding": "mp3",
    },
})
```

## SSML Support

Some providers support SSML for fine-grained control:

```go
ssml := `<speak>
    Hello! <break time="500ms"/>
    This is <emphasis level="strong">important</emphasis>.
    <prosody rate="slow">Speaking slowly now.</prosody>
</speak>`

result, _ := provider.Synthesize(ctx, ssml, omnivoice.SynthesisConfig{
    VoiceID: voiceID,
    Extensions: map[string]any{
        "use_ssml": true,
    },
})
```

## Voice Selection

### ElevenLabs Voices

| Voice ID | Name | Description |
|----------|------|-------------|
| `pNInz6obpgDQGcFmaJgB` | Adam | Deep, narrative |
| `EXAVITQu4vr4xnSDxMaL` | Sarah | Soft, conversational |
| `21m00Tcm4TlvDq8ikWAM` | Rachel | Calm, professional |

### OpenAI Voices

| Voice ID | Description |
|----------|-------------|
| `alloy` | Neutral, balanced |
| `echo` | Warm, conversational |
| `fable` | Expressive, British |
| `onyx` | Deep, authoritative |
| `nova` | Friendly, upbeat |
| `shimmer` | Clear, professional |

## Error Handling

```go
result, err := provider.Synthesize(ctx, text, config)
if err != nil {
    switch {
    case errors.Is(err, context.DeadlineExceeded):
        log.Println("Request timed out")
    case strings.Contains(err.Error(), "invalid_api_key"):
        log.Println("Invalid API key")
    case strings.Contains(err.Error(), "quota_exceeded"):
        log.Println("Rate limit or quota exceeded")
    default:
        log.Printf("TTS error: %v", err)
    }
    return
}
```

## Best Practices

1. **Cache audio** - Store generated audio to avoid repeated API calls
2. **Use streaming** - For long texts or real-time applications
3. **Choose appropriate quality** - Lower quality = faster, smaller files
4. **Handle errors gracefully** - Implement retries for transient failures
5. **Consider latency** - Deepgram for real-time, ElevenLabs for quality

## Next Steps

- [Streaming Guide](streaming.md) - Real-time TTS streaming
- [Voice Agents](voice-agents.md) - Combine TTS with STT and LLMs
- [Provider Details](../providers/index.md) - Provider-specific features
