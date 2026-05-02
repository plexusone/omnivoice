# Speech-to-Text (STT)

Convert audio to text using multiple providers with support for batch and streaming transcription.

## Quick Start

```go
provider, _ := omnivoice.GetSTTProvider("deepgram",
    omnivoice.WithAPIKey(apiKey))

result, _ := provider.TranscribeFile(ctx, "audio.mp3", omnivoice.TranscriptionConfig{
    Language: "en",
})

fmt.Println(result.Text)
```

## Available Providers

| Provider | Registry Name | Accuracy | Latency | Best For |
|----------|---------------|----------|---------|----------|
| Deepgram | `"deepgram"` | Excellent | Very Low | Real-time, high volume |
| OpenAI | `"openai"` | Excellent | Medium | General purpose, multilingual |
| ElevenLabs | `"elevenlabs"` | Very Good | Medium | Integration with TTS |
| Twilio | `"twilio"` | Good | Low | Phone calls |

## Transcription Methods

### TranscribeFile

Transcribe a local audio file:

```go
result, err := provider.TranscribeFile(ctx, "recording.mp3", omnivoice.TranscriptionConfig{
    Language:             "en",
    EnableWordTimestamps: true,
})
```

### TranscribeURL

Transcribe audio from a URL:

```go
result, err := provider.TranscribeURL(ctx, "https://example.com/audio.mp3", omnivoice.TranscriptionConfig{
    Language: "en",
})
```

### Transcribe

Transcribe from an `io.Reader`:

```go
file, _ := os.Open("audio.mp3")
defer file.Close()

result, err := provider.Transcribe(ctx, file, omnivoice.TranscriptionConfig{
    Language: "en",
})
```

### TranscribeStream

Real-time streaming transcription:

```go
stream, err := provider.TranscribeStream(ctx, omnivoice.TranscriptionConfig{
    Language: "en",
})
if err != nil {
    panic(err)
}

// Send audio chunks
go func() {
    for chunk := range audioSource {
        stream.Write(chunk)
    }
    stream.Close()
}()

// Receive transcription results
for result := range stream.Results() {
    if result.IsFinal {
        fmt.Printf("Final: %s\n", result.Text)
    } else {
        fmt.Printf("Interim: %s\n", result.Text)
    }
}
```

## Configuration Options

```go
config := omnivoice.TranscriptionConfig{
    // Language
    Language: "en-US",  // BCP-47 language code

    // Features
    EnableWordTimestamps:  true,  // Word-level timing
    EnableSpeakerDiarization: true, // Speaker identification

    // Model selection (provider-specific)
    Model: "nova-2",

    // Provider-specific extensions
    Extensions: map[string]any{
        "smart_format": true,
        "punctuate":    true,
    },
}
```

## Provider-Specific Examples

### Deepgram

```go
provider, _ := omnivoice.GetSTTProvider("deepgram",
    omnivoice.WithAPIKey(apiKey))

result, _ := provider.TranscribeFile(ctx, "audio.mp3", omnivoice.TranscriptionConfig{
    Language: "en",
    Model:    "nova-2",
    EnableWordTimestamps: true,
    Extensions: map[string]any{
        "smart_format": true,
        "punctuate":    true,
        "diarize":      true,
        "utterances":   true,
    },
})
```

### OpenAI Whisper

```go
provider, _ := omnivoice.GetSTTProvider("openai",
    omnivoice.WithAPIKey(apiKey))

result, _ := provider.TranscribeFile(ctx, "audio.mp3", omnivoice.TranscriptionConfig{
    Language: "en",
    Extensions: map[string]any{
        "model":            "whisper-1",
        "response_format":  "verbose_json",
        "temperature":      0,
    },
})
```

### ElevenLabs Scribe

```go
provider, _ := omnivoice.GetSTTProvider("elevenlabs",
    omnivoice.WithAPIKey(apiKey))

result, _ := provider.TranscribeFile(ctx, "audio.mp3", omnivoice.TranscriptionConfig{
    Language: "en",
    EnableWordTimestamps: true,
})
```

## Working with Results

### Basic Text

```go
fmt.Println(result.Text)
```

### Word Timestamps

```go
for _, word := range result.Words {
    fmt.Printf("[%.2fs - %.2fs] %s\n",
        word.Start.Seconds(),
        word.End.Seconds(),
        word.Text)
}
```

### Speaker Diarization

```go
for _, segment := range result.Segments {
    fmt.Printf("Speaker %d: %s\n", segment.Speaker, segment.Text)
}
```

### Confidence Scores

```go
fmt.Printf("Confidence: %.2f%%\n", result.Confidence * 100)

for _, word := range result.Words {
    if word.Confidence < 0.8 {
        fmt.Printf("Low confidence word: %s (%.2f)\n", word.Text, word.Confidence)
    }
}
```

## Language Codes

OmniVoice accepts BCP-47 language codes:

| Code | Language |
|------|----------|
| `en` | English |
| `en-US` | English (US) |
| `en-GB` | English (UK) |
| `es` | Spanish |
| `fr` | French |
| `de` | German |
| `ja` | Japanese |
| `zh` | Chinese |

Most providers support automatic language detection when no code is specified.

## Error Handling

```go
result, err := provider.TranscribeFile(ctx, path, config)
if err != nil {
    switch {
    case errors.Is(err, context.DeadlineExceeded):
        log.Println("Request timed out")
    case os.IsNotExist(err):
        log.Println("Audio file not found")
    case strings.Contains(err.Error(), "unsupported_format"):
        log.Println("Audio format not supported")
    default:
        log.Printf("STT error: %v", err)
    }
    return
}
```

## Audio Format Support

| Format | Extension | Providers |
|--------|-----------|-----------|
| MP3 | `.mp3` | All |
| WAV | `.wav` | All |
| FLAC | `.flac` | Deepgram, OpenAI |
| OGG | `.ogg` | Deepgram, OpenAI |
| WebM | `.webm` | Deepgram, OpenAI |
| M4A | `.m4a` | Deepgram, OpenAI |

## Transcript Format

For persisting or sharing transcription results, use the canonical Transcript format:

```go
import "github.com/plexusone/omnivoice"

// Convert result to canonical Transcript
transcript := omnivoice.NewTranscript(result, "deepgram", "nova-2", "audio.mp3", config)

// Save to JSON (durations serialize as milliseconds)
err := transcript.SaveJSON("output.transcript.json")

// Load from JSON
loaded, err := omnivoice.LoadTranscript("output.transcript.json")

// Access timing information
fmt.Printf("Duration: %v\n", transcript.TotalDuration())
for _, seg := range transcript.Segments {
    fmt.Printf("[%v - %v] %s\n",
        seg.Start.Duration(),
        seg.End.Duration(),
        seg.Text)
}
```

The JSON format includes:

- Full text and segments with timing
- Word-level timestamps (when enabled)
- Speaker identification (when enabled)
- Provider metadata and options used
- JSON Schema URL for validation

See the [CLI Guide](cli.md) for generating transcripts from the command line with `omnivoice transcribe -f json`.

## Best Practices

1. **Use appropriate models** - Nova-2 for accuracy, Base for speed
2. **Enable word timestamps** - Essential for subtitles and alignment
3. **Handle streaming errors** - Reconnect on connection drops
4. **Choose the right provider** - Deepgram for real-time, OpenAI for accuracy
5. **Preprocess audio** - Normalize volume, remove silence
6. **Save as Transcript JSON** - Use canonical format for interoperability

## Next Steps

- [CLI Guide](cli.md) - Transcribe from the command line
- [Streaming Guide](streaming.md) - Real-time transcription
- [Subtitles Guide](subtitles.md) - Generate SRT/VTT captions
- [Voice Agents](voice-agents.md) - Build conversational agents
