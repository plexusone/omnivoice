# OpenAI

OpenAI provides TTS via their text-to-speech API and STT via Whisper.

## Features

- **TTS**: High-quality voices with natural prosody
- **STT**: Whisper model for accurate transcription
- **Languages**: 50+ languages supported
- **Streaming**: TTS streaming supported

## Configuration

```go
import (
    "github.com/plexusone/omnivoice"
    _ "github.com/plexusone/omnivoice/providers/openai"
)

// TTS Provider
tts, err := omnivoice.GetTTSProvider("openai",
    omnivoice.WithAPIKey(os.Getenv("OPENAI_API_KEY")),
)

// STT Provider
stt, err := omnivoice.GetSTTProvider("openai",
    omnivoice.WithAPIKey(os.Getenv("OPENAI_API_KEY")),
)
```

## Text-to-Speech

### Available Voices

| Voice | Description |
|-------|-------------|
| `alloy` | Neutral, balanced |
| `echo` | Warm, conversational |
| `fable` | British, storytelling |
| `onyx` | Deep, authoritative |
| `nova` | Friendly, upbeat |
| `shimmer` | Soft, gentle |

### Basic Usage

```go
result, err := tts.Synthesize(ctx, "Hello, world!", omnivoice.SynthesisConfig{
    VoiceID: "alloy",
})
if err != nil {
    log.Fatal(err)
}

os.WriteFile("output.mp3", result.Audio, 0600)
```

### Models

| Model | Quality | Speed | Use Case |
|-------|---------|-------|----------|
| `tts-1` | Good | Fast | Real-time applications |
| `tts-1-hd` | Better | Slower | Pre-generated content |

```go
result, err := tts.Synthesize(ctx, text, omnivoice.SynthesisConfig{
    VoiceID: "alloy",
    Extensions: map[string]any{
        "model": "tts-1-hd",  // Higher quality
    },
})
```

### Output Formats

```go
config := omnivoice.SynthesisConfig{
    VoiceID:      "alloy",
    OutputFormat: "mp3",  // mp3, opus, aac, flac
}
```

### Streaming

```go
stream, err := tts.SynthesizeStream(ctx, text, omnivoice.SynthesisConfig{
    VoiceID: "alloy",
})
if err != nil {
    log.Fatal(err)
}

for chunk := range stream {
    if chunk.Error != nil {
        log.Printf("Error: %v", chunk.Error)
        break
    }
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

### With Word Timestamps

```go
result, err := stt.TranscribeFile(ctx, "audio.mp3", omnivoice.TranscriptionConfig{
    Language:             "en",
    EnableWordTimestamps: true,
})

for _, word := range result.Words {
    fmt.Printf("[%.2f-%.2f] %s\n", word.Start, word.End, word.Word)
}
```

### Whisper Models

| Model | Size | Speed | Accuracy |
|-------|------|-------|----------|
| `whisper-1` | Large | Medium | High |

```go
result, err := stt.TranscribeFile(ctx, "audio.mp3", omnivoice.TranscriptionConfig{
    Model:    "whisper-1",
    Language: "en",
})
```

### Translation

Translate audio to English:

```go
result, err := stt.TranscribeFile(ctx, "french_audio.mp3", omnivoice.TranscriptionConfig{
    Extensions: map[string]any{
        "task": "translate",  // Translate to English
    },
})
```

## Limitations

- **STT Streaming**: Not supported (file-based only)
- **File Size**: Max 25MB per request
- **Audio Length**: Max ~2 hours

## Error Handling

```go
result, err := tts.Synthesize(ctx, text, config)
if err != nil {
    switch {
    case strings.Contains(err.Error(), "rate_limit"):
        log.Println("Rate limited, retry with backoff")
    case strings.Contains(err.Error(), "invalid_api_key"):
        log.Println("Check OPENAI_API_KEY")
    default:
        log.Printf("Error: %v", err)
    }
}
```

## Best Practices

1. **Use tts-1 for real-time** - Lower latency than tts-1-hd
2. **Batch transcriptions** - OpenAI STT is file-based, not streaming
3. **Set language explicitly** - Improves accuracy for STT
4. **Use opus for streaming** - Smaller chunks, lower latency

## Pricing

| Service | Model | Price |
|---------|-------|-------|
| TTS | tts-1 | $15/1M chars |
| TTS | tts-1-hd | $30/1M chars |
| STT | whisper-1 | $0.006/min |

*Prices as of 2024. Check [OpenAI Pricing](https://openai.com/pricing) for current rates.*

## Next Steps

- [ElevenLabs](elevenlabs.md) - Premium voice quality
- [Deepgram](deepgram.md) - Real-time STT streaming
- [Provider Comparison](index.md) - Compare all providers
