# Streaming

Real-time streaming for TTS and STT enables low-latency voice applications.

## Why Streaming?

| Approach | Latency | Use Case |
|----------|---------|----------|
| Batch | High (seconds) | Offline processing, file transcription |
| Streaming | Low (milliseconds) | Real-time apps, voice agents, live captions |

## Streaming TTS

Generate audio in real-time as chunks arrive:

```go
provider, _ := omnivoice.GetTTSProvider("elevenlabs",
    omnivoice.WithAPIKey(apiKey))

stream, err := provider.SynthesizeStream(ctx, "This is a longer text that will be streamed in chunks...", omnivoice.SynthesisConfig{
    VoiceID: "pNInz6obpgDQGcFmaJgB",
})
if err != nil {
    log.Fatal(err)
}

// Process audio chunks as they arrive
for chunk := range stream {
    if chunk.Error != nil {
        log.Printf("Stream error: %v", chunk.Error)
        break
    }

    // Play or forward audio immediately
    playAudio(chunk.Audio)
}
```

### Streaming to HTTP Response

```go
http.HandleFunc("/speak", func(w http.ResponseWriter, r *http.Request) {
    text := r.URL.Query().Get("text")

    stream, err := provider.SynthesizeStream(ctx, text, omnivoice.SynthesisConfig{
        VoiceID:      "pNInz6obpgDQGcFmaJgB",
        OutputFormat: "mp3_44100_128",
    })
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    w.Header().Set("Content-Type", "audio/mpeg")
    w.Header().Set("Transfer-Encoding", "chunked")

    flusher, _ := w.(http.Flusher)
    for chunk := range stream {
        w.Write(chunk.Audio)
        flusher.Flush()
    }
})
```

### Streaming to WebSocket

```go
func handleWebSocket(conn *websocket.Conn, text string) {
    stream, _ := provider.SynthesizeStream(ctx, text, config)

    for chunk := range stream {
        // Send audio chunk over WebSocket
        conn.WriteMessage(websocket.BinaryMessage, chunk.Audio)
    }
}
```

## Streaming STT

Transcribe audio in real-time:

```go
provider, _ := omnivoice.GetSTTProvider("deepgram",
    omnivoice.WithAPIKey(apiKey))

stream, err := provider.TranscribeStream(ctx, omnivoice.TranscriptionConfig{
    Language: "en",
    Extensions: map[string]any{
        "interim_results": true,
        "punctuate":       true,
    },
})
if err != nil {
    log.Fatal(err)
}

// Send audio in a goroutine
go func() {
    defer stream.Close()

    // Read from microphone or audio source
    for {
        audioChunk := readAudioChunk()
        if audioChunk == nil {
            break
        }
        stream.Write(audioChunk)
    }
}()

// Receive transcription results
for result := range stream.Results() {
    if result.IsFinal {
        fmt.Printf("Final: %s\n", result.Text)
    } else {
        fmt.Printf("Interim: %s\r", result.Text) // Overwrite line
    }
}
```

### From Microphone

```go
import "github.com/gordonklaus/portaudio"

func streamFromMicrophone(stream omnivoice.STTStream) error {
    portaudio.Initialize()
    defer portaudio.Terminate()

    // 16kHz mono audio
    buffer := make([]int16, 1600) // 100ms chunks

    audioStream, err := portaudio.OpenDefaultStream(1, 0, 16000, len(buffer), buffer)
    if err != nil {
        return err
    }
    defer audioStream.Close()

    audioStream.Start()
    defer audioStream.Stop()

    for {
        audioStream.Read()

        // Convert int16 to bytes
        bytes := make([]byte, len(buffer)*2)
        for i, sample := range buffer {
            bytes[i*2] = byte(sample)
            bytes[i*2+1] = byte(sample >> 8)
        }

        stream.Write(bytes)
    }
}
```

### From WebSocket

```go
func handleAudioWebSocket(conn *websocket.Conn) {
    stream, _ := provider.TranscribeStream(ctx, config)
    defer stream.Close()

    // Forward results back to client
    go func() {
        for result := range stream.Results() {
            response, _ := json.Marshal(result)
            conn.WriteMessage(websocket.TextMessage, response)
        }
    }()

    // Receive audio from client
    for {
        _, audio, err := conn.ReadMessage()
        if err != nil {
            break
        }
        stream.Write(audio)
    }
}
```

## Bidirectional Streaming

For voice agents, combine STT and TTS streaming:

```go
func voiceAgent(audioIn <-chan []byte, audioOut chan<- []byte) {
    // Start STT stream
    sttStream, _ := sttProvider.TranscribeStream(ctx, sttConfig)

    // Forward incoming audio to STT
    go func() {
        for audio := range audioIn {
            sttStream.Write(audio)
        }
        sttStream.Close()
    }()

    // Process transcriptions and generate responses
    for result := range sttStream.Results() {
        if !result.IsFinal {
            continue
        }

        // Generate response (e.g., from LLM)
        response := generateResponse(result.Text)

        // Stream TTS response
        ttsStream, _ := ttsProvider.SynthesizeStream(ctx, response, ttsConfig)
        for chunk := range ttsStream {
            audioOut <- chunk.Audio
        }
    }
}
```

## Audio Formats for Streaming

### STT Input Formats

| Provider | Recommended Format |
|----------|-------------------|
| Deepgram | Linear16 PCM, 16kHz, mono |
| OpenAI | Various (file-based only) |
| ElevenLabs | Linear16 PCM, 16kHz, mono |

### TTS Output Formats

| Provider | Streaming Formats |
|----------|------------------|
| ElevenLabs | `mp3_44100_128`, `pcm_16000`, `pcm_44100` |
| OpenAI | `mp3`, `opus`, `aac`, `flac` |
| Deepgram | `mp3`, `linear16` |

## Latency Optimization

### Reduce Chunk Size

```go
// Smaller chunks = lower latency, more overhead
config := omnivoice.SynthesisConfig{
    Extensions: map[string]any{
        "chunk_size": 256, // Smaller chunks
    },
}
```

### Use Faster Models

```go
// Deepgram: Use base model for lowest latency
sttConfig := omnivoice.TranscriptionConfig{
    Model: "nova-2", // Fast and accurate
}

// OpenAI: Use tts-1 (faster) vs tts-1-hd (quality)
ttsConfig := omnivoice.SynthesisConfig{
    Extensions: map[string]any{
        "model": "tts-1",
    },
}
```

### Connection Pooling

```go
// Reuse providers instead of creating new ones
var (
    ttsProvider omnivoice.TTSProvider
    sttProvider omnivoice.STTProvider
    once        sync.Once
)

func getProviders() (omnivoice.TTSProvider, omnivoice.STTProvider) {
    once.Do(func() {
        ttsProvider, _ = omnivoice.GetTTSProvider("elevenlabs", ...)
        sttProvider, _ = omnivoice.GetSTTProvider("deepgram", ...)
    })
    return ttsProvider, sttProvider
}
```

## Error Handling

```go
stream, err := provider.SynthesizeStream(ctx, text, config)
if err != nil {
    log.Fatal(err)
}

for chunk := range stream {
    if chunk.Error != nil {
        if errors.Is(chunk.Error, context.Canceled) {
            log.Println("Stream canceled")
        } else if errors.Is(chunk.Error, io.EOF) {
            log.Println("Stream completed")
        } else {
            log.Printf("Stream error: %v", chunk.Error)
        }
        break
    }

    // Process chunk
}
```

## Best Practices

1. **Buffer appropriately** - Small buffers for latency, larger for stability
2. **Handle backpressure** - Don't overwhelm consumers
3. **Graceful shutdown** - Always close streams properly
4. **Monitor latency** - Track time-to-first-byte
5. **Use appropriate formats** - PCM for low latency, MP3 for bandwidth

## Next Steps

- [Voice Agents](voice-agents.md) - Build real-time conversational agents
- [Voice Calls](voice-calls.md) - Integrate with phone systems
- [Subtitles](subtitles.md) - Generate captions from streaming STT
