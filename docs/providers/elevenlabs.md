# ElevenLabs

ElevenLabs provides premium voice synthesis with natural, expressive voices.

## Features

- **TTS**: Industry-leading voice quality
- **STT**: Real-time transcription (Scribe)
- **Voice Cloning**: Create custom voices
- **Streaming**: Ultra-low latency streaming
- **Languages**: 29+ languages

## Configuration

```go
import (
    "github.com/plexusone/omnivoice"
    _ "github.com/plexusone/omnivoice/providers/elevenlabs"
)

// TTS Provider
tts, err := omnivoice.GetTTSProvider("elevenlabs",
    omnivoice.WithAPIKey(os.Getenv("ELEVENLABS_API_KEY")),
)

// STT Provider
stt, err := omnivoice.GetSTTProvider("elevenlabs",
    omnivoice.WithAPIKey(os.Getenv("ELEVENLABS_API_KEY")),
)
```

## Text-to-Speech

### Popular Voices

| Voice ID | Name | Style |
|----------|------|-------|
| `pNInz6obpgDQGcFmaJgB` | Adam | Deep, narrative |
| `EXAVITQu4vr4xnSDxMaL` | Sarah | Soft, young |
| `onwK4e9ZLuTAKqWW03F9` | Daniel | British, authoritative |
| `XB0fDUnXU5powFXDhCwa` | Charlotte | Swedish, calm |

### Basic Usage

```go
result, err := tts.Synthesize(ctx, "Hello, world!", omnivoice.SynthesisConfig{
    VoiceID: "pNInz6obpgDQGcFmaJgB",
})
if err != nil {
    log.Fatal(err)
}

os.WriteFile("output.mp3", result.Audio, 0600)
```

### Models

| Model | Latency | Quality | Use Case |
|-------|---------|---------|----------|
| `eleven_turbo_v2_5` | Lowest | Good | Real-time voice agents |
| `eleven_multilingual_v2` | Low | Excellent | Multi-language apps |
| `eleven_monolingual_v1` | Medium | Excellent | English-only |

```go
result, err := tts.Synthesize(ctx, text, omnivoice.SynthesisConfig{
    VoiceID: "pNInz6obpgDQGcFmaJgB",
    Extensions: map[string]any{
        "model_id": "eleven_turbo_v2_5",
    },
})
```

### Voice Settings

```go
result, err := tts.Synthesize(ctx, text, omnivoice.SynthesisConfig{
    VoiceID: "pNInz6obpgDQGcFmaJgB",
    Extensions: map[string]any{
        "stability":        0.5,  // 0-1, higher = more consistent
        "similarity_boost": 0.75, // 0-1, higher = closer to original voice
        "style":            0.3,  // 0-1, style exaggeration
        "use_speaker_boost": true,
    },
})
```

### Output Formats

| Format | Use Case |
|--------|----------|
| `mp3_44100_128` | General purpose |
| `mp3_44100_192` | Higher quality |
| `pcm_16000` | Real-time, low latency |
| `pcm_44100` | Real-time, high quality |
| `ulaw_8000` | Telephony (Twilio) |

```go
config := omnivoice.SynthesisConfig{
    VoiceID:      "pNInz6obpgDQGcFmaJgB",
    OutputFormat: "pcm_16000",  // For real-time streaming
}
```

### Streaming

```go
stream, err := tts.SynthesizeStream(ctx, text, omnivoice.SynthesisConfig{
    VoiceID: "pNInz6obpgDQGcFmaJgB",
    Extensions: map[string]any{
        "model_id":      "eleven_turbo_v2_5",
        "optimize_streaming_latency": 3,  // 0-4, higher = lower latency
    },
})
if err != nil {
    log.Fatal(err)
}

for chunk := range stream {
    if chunk.Error != nil {
        break
    }
    // Play immediately for low latency
    playAudio(chunk.Audio)
}
```

## Speech-to-Text

### Basic Transcription

```go
result, err := stt.TranscribeFile(ctx, "audio.mp3", omnivoice.TranscriptionConfig{
    Language: "en",
})
if err != nil {
    log.Fatal(err)
}

fmt.Println(result.Text)
```

### Streaming Transcription

```go
stream, err := stt.TranscribeStream(ctx, omnivoice.TranscriptionConfig{
    Language: "en",
    Extensions: map[string]any{
        "interim_results": true,
    },
})
if err != nil {
    log.Fatal(err)
}

// Send audio
go func() {
    for audio := range audioSource {
        stream.Write(audio)
    }
    stream.Close()
}()

// Receive transcriptions
for result := range stream.Results() {
    if result.IsFinal {
        fmt.Printf("Final: %s\n", result.Text)
    } else {
        fmt.Printf("Interim: %s\r", result.Text)
    }
}
```

### Speaker Diarization

```go
result, err := stt.TranscribeFile(ctx, "meeting.mp3", omnivoice.TranscriptionConfig{
    EnableSpeakerDiarization: true,
    Extensions: map[string]any{
        "num_speakers": 2,
    },
})

for _, segment := range result.Segments {
    fmt.Printf("[Speaker %d] %s\n", segment.Speaker, segment.Text)
}
```

## Voice Cloning

Create custom voices (requires ElevenLabs account):

```go
// Use a cloned voice
result, err := tts.Synthesize(ctx, text, omnivoice.SynthesisConfig{
    VoiceID: "your-cloned-voice-id",
})
```

## Latency Optimization

For voice agents requiring minimal latency:

```go
tts, _ := omnivoice.GetTTSProvider("elevenlabs",
    omnivoice.WithAPIKey(apiKey),
)

stream, _ := tts.SynthesizeStream(ctx, response, omnivoice.SynthesisConfig{
    VoiceID:      "pNInz6obpgDQGcFmaJgB",
    OutputFormat: "pcm_16000",
    Extensions: map[string]any{
        "model_id":                    "eleven_turbo_v2_5",
        "optimize_streaming_latency": 4,  // Maximum optimization
    },
})
```

## Error Handling

```go
result, err := tts.Synthesize(ctx, text, config)
if err != nil {
    switch {
    case strings.Contains(err.Error(), "quota_exceeded"):
        log.Println("Monthly quota exceeded")
    case strings.Contains(err.Error(), "voice_not_found"):
        log.Println("Invalid voice ID")
    case strings.Contains(err.Error(), "invalid_api_key"):
        log.Println("Check ELEVENLABS_API_KEY")
    default:
        log.Printf("Error: %v", err)
    }
}
```

## Best Practices

1. **Use Turbo for real-time** - `eleven_turbo_v2_5` has lowest latency
2. **Cache common phrases** - Pre-generate greetings, confirmations
3. **Use PCM for streaming** - No encoding overhead
4. **Set optimize_streaming_latency** - Higher values reduce time-to-first-byte
5. **Monitor quota** - ElevenLabs has character quotas per plan

## Pricing

| Plan | Characters/Month | Price |
|------|------------------|-------|
| Free | 10,000 | $0 |
| Starter | 30,000 | $5 |
| Creator | 100,000 | $22 |
| Pro | 500,000 | $99 |

*Check [ElevenLabs Pricing](https://elevenlabs.io/pricing) for current rates.*

## Next Steps

- [OpenAI](openai.md) - Alternative TTS/STT
- [Deepgram](deepgram.md) - Real-time STT
- [Voice Agents](../guides/voice-agents.md) - Build conversational agents
