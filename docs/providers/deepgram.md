# Deepgram

Deepgram provides ultra-low latency STT and TTS optimized for real-time applications.

## Features

- **STT**: Industry-leading real-time transcription (Nova-2)
- **TTS**: Aura voices with streaming support
- **Streaming**: Sub-200ms latency
- **Languages**: 36+ languages
- **Features**: Diarization, punctuation, smart formatting

## Configuration

```go
import (
    "github.com/plexusone/omnivoice"
    _ "github.com/plexusone/omnivoice/providers/all"
)

// STT Provider
stt, err := omnivoice.GetSTTProvider("deepgram",
    omnivoice.WithAPIKey(os.Getenv("DEEPGRAM_API_KEY")),
)

// TTS Provider
tts, err := omnivoice.GetTTSProvider("deepgram",
    omnivoice.WithAPIKey(os.Getenv("DEEPGRAM_API_KEY")),
)
```

## Speech-to-Text

### Models

| Model | Accuracy | Speed | Use Case |
|-------|----------|-------|----------|
| `nova-2` | Highest | Fast | General purpose (recommended) |
| `nova` | High | Fast | Previous generation |
| `enhanced` | Good | Fast | Budget option |
| `base` | Basic | Fastest | High-volume, low-cost |

### Basic Transcription

```go
result, err := stt.TranscribeFile(ctx, "audio.mp3", omnivoice.TranscriptionConfig{
    Language: "en",
    Model:    "nova-2",
})
if err != nil {
    log.Fatal(err)
}

fmt.Println(result.Text)
```

### Real-Time Streaming

```go
stream, err := stt.TranscribeStream(ctx, omnivoice.TranscriptionConfig{
    Language: "en",
    Model:    "nova-2",
    Extensions: map[string]any{
        "interim_results": true,  // Get partial transcripts
        "punctuate":       true,
        "smart_format":    true,
    },
})
if err != nil {
    log.Fatal(err)
}

// Send audio chunks
go func() {
    defer stream.Close()
    for {
        audio := readAudioChunk()  // 100ms chunks recommended
        if audio == nil {
            break
        }
        stream.Write(audio)
    }
}()

// Receive transcriptions in real-time
for result := range stream.Results() {
    if result.IsFinal {
        fmt.Printf("Final: %s\n", result.Text)
    } else {
        fmt.Printf("Interim: %s\r", result.Text)
    }
}
```

### Word Timestamps

```go
result, err := stt.TranscribeFile(ctx, "audio.mp3", omnivoice.TranscriptionConfig{
    EnableWordTimestamps: true,
    Extensions: map[string]any{
        "punctuate":    true,
        "smart_format": true,
    },
})

for _, word := range result.Words {
    fmt.Printf("[%.2fs] %s\n", word.Start, word.Word)
}
```

### Speaker Diarization

```go
result, err := stt.TranscribeFile(ctx, "meeting.mp3", omnivoice.TranscriptionConfig{
    EnableSpeakerDiarization: true,
    Extensions: map[string]any{
        "diarize": true,
    },
})

for _, segment := range result.Segments {
    fmt.Printf("[Speaker %d] %s\n", segment.Speaker, segment.Text)
}
```

### Audio Formats

| Format | Sample Rate | Channels | Notes |
|--------|-------------|----------|-------|
| `linear16` | 16000 Hz | Mono | Recommended for streaming |
| `linear16` | 8000 Hz | Mono | Telephony |
| `mp3` | Any | Any | File transcription |
| `flac` | Any | Any | High quality files |

```go
stream, _ := stt.TranscribeStream(ctx, omnivoice.TranscriptionConfig{
    Extensions: map[string]any{
        "encoding":    "linear16",
        "sample_rate": 16000,
        "channels":    1,
    },
})
```

### Endpointing

Detect end of speech for voice agents:

```go
stream, _ := stt.TranscribeStream(ctx, omnivoice.TranscriptionConfig{
    Extensions: map[string]any{
        "endpointing":    300,  // 300ms silence = end of utterance
        "interim_results": true,
        "utterance_end_ms": 1000,  // Max wait for utterance end
    },
})

for result := range stream.Results() {
    if result.IsSpeechFinal {
        // User finished speaking, generate response
        response := generateResponse(result.Text)
        speak(response)
    }
}
```

## Text-to-Speech

### Aura Voices

| Voice | Gender | Style |
|-------|--------|-------|
| `aura-asteria-en` | Female | Warm, friendly |
| `aura-luna-en` | Female | Soft, calm |
| `aura-stella-en` | Female | Professional |
| `aura-athena-en` | Female | Authoritative |
| `aura-hera-en` | Female | Mature |
| `aura-orion-en` | Male | Deep, resonant |
| `aura-arcas-en` | Male | Conversational |
| `aura-perseus-en` | Male | Young, energetic |
| `aura-angus-en` | Male | Irish accent |
| `aura-orpheus-en` | Male | Storytelling |
| `aura-helios-en` | Male | British |
| `aura-zeus-en` | Male | Authoritative |

### Basic Usage

```go
result, err := tts.Synthesize(ctx, "Hello, world!", omnivoice.SynthesisConfig{
    VoiceID: "aura-asteria-en",
})
if err != nil {
    log.Fatal(err)
}

os.WriteFile("output.mp3", result.Audio, 0600)
```

### Streaming

```go
stream, err := tts.SynthesizeStream(ctx, text, omnivoice.SynthesisConfig{
    VoiceID:      "aura-asteria-en",
    OutputFormat: "linear16",  // PCM for lowest latency
})
if err != nil {
    log.Fatal(err)
}

for chunk := range stream {
    if chunk.Error != nil {
        break
    }
    playAudio(chunk.Audio)
}
```

### Output Formats

| Format | Use Case |
|--------|----------|
| `linear16` | Real-time streaming (recommended) |
| `mp3` | File storage |
| `opus` | WebRTC |
| `flac` | Archival |
| `mulaw` | Telephony |
| `alaw` | Telephony (EU) |

## Latency Optimization

For voice agents:

```go
// STT: Use streaming with fast endpointing
sttStream, _ := stt.TranscribeStream(ctx, omnivoice.TranscriptionConfig{
    Model: "nova-2",
    Extensions: map[string]any{
        "interim_results": true,
        "endpointing":     200,  // Fast response
        "encoding":        "linear16",
        "sample_rate":     16000,
    },
})

// TTS: Use PCM output for zero encoding overhead
ttsStream, _ := tts.SynthesizeStream(ctx, response, omnivoice.SynthesisConfig{
    VoiceID:      "aura-asteria-en",
    OutputFormat: "linear16",
})
```

## Error Handling

```go
result, err := stt.TranscribeFile(ctx, "audio.mp3", config)
if err != nil {
    switch {
    case strings.Contains(err.Error(), "insufficient_funds"):
        log.Println("Account balance depleted")
    case strings.Contains(err.Error(), "invalid_credentials"):
        log.Println("Check DEEPGRAM_API_KEY")
    case strings.Contains(err.Error(), "unsupported_format"):
        log.Println("Audio format not supported")
    default:
        log.Printf("Error: %v", err)
    }
}
```

## Best Practices

1. **Use Nova-2** - Best accuracy with low latency
2. **Stream in 100ms chunks** - Optimal balance for real-time
3. **Enable endpointing** - Detect when user stops speaking
4. **Use linear16** - No encoding/decoding overhead
5. **Enable smart_format** - Better punctuation and formatting

## Pricing

| Service | Model | Price |
|---------|-------|-------|
| STT | Nova-2 | $0.0043/min |
| STT | Nova | $0.0036/min |
| STT | Enhanced | $0.0145/min |
| STT | Base | $0.0125/min |
| TTS | Aura | $0.015/1K chars |

*Check [Deepgram Pricing](https://deepgram.com/pricing) for current rates.*

## Next Steps

- [Streaming Guide](../guides/streaming.md) - Real-time audio patterns
- [Voice Agents](../guides/voice-agents.md) - Build conversational agents
- [Subtitles](../guides/subtitles.md) - Generate captions
