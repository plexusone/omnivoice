# Voice Agents

Build conversational voice agents by combining STT, LLM, and TTS in a real-time pipeline.

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        Voice Agent Pipeline                     │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│   User Audio ──► STT ──► Text ──► LLM ──► Text ──► TTS ──► Audio │
│       ▲                                                    │     │
│       │                                                    │     │
│       └────────────────────────────────────────────────────┘     │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "os"

    "github.com/plexusone/omnivoice"
    _ "github.com/plexusone/omnivoice/providers/all"
)

func main() {
    ctx := context.Background()

    // Initialize providers
    stt, _ := omnivoice.GetSTTProvider("deepgram",
        omnivoice.WithAPIKey(os.Getenv("DEEPGRAM_API_KEY")))

    tts, _ := omnivoice.GetTTSProvider("elevenlabs",
        omnivoice.WithAPIKey(os.Getenv("ELEVENLABS_API_KEY")))

    // Simple conversation loop
    for {
        // 1. Get user speech
        userText := transcribeUserSpeech(ctx, stt)
        fmt.Printf("User: %s\n", userText)

        // 2. Generate response (via LLM)
        response := generateResponse(userText)
        fmt.Printf("Agent: %s\n", response)

        // 3. Speak response
        speakResponse(ctx, tts, response)
    }
}
```

## Phone-Based Voice Agent

Combine CallSystem with STT/TTS for phone agents:

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

    // Initialize providers
    callSystem, _ := omnivoice.GetCallSystemProvider("twilio",
        omnivoice.WithAccountSID(os.Getenv("TWILIO_ACCOUNT_SID")),
        omnivoice.WithAuthToken(os.Getenv("TWILIO_AUTH_TOKEN")),
        omnivoice.WithPhoneNumber("+15551234567"),
        omnivoice.WithWebhookURL("https://your-server.com/webhook"),
    )

    stt, _ := omnivoice.GetSTTProvider("deepgram",
        omnivoice.WithAPIKey(os.Getenv("DEEPGRAM_API_KEY")))

    tts, _ := omnivoice.GetTTSProvider("elevenlabs",
        omnivoice.WithAPIKey(os.Getenv("ELEVENLABS_API_KEY")))

    // Handle incoming calls
    callSystem.OnIncomingCall(func(call omnivoice.Call) error {
        return handleCall(ctx, call, stt, tts)
    })

    // Start server
    startServer(callSystem)
}

func handleCall(ctx context.Context, call omnivoice.Call, stt omnivoice.STTProvider, tts omnivoice.TTSProvider) error {
    // Answer the call
    if err := call.Answer(ctx); err != nil {
        return err
    }

    // Get transport connection for audio
    transport := call.Transport()

    // Start STT stream
    sttStream, _ := stt.TranscribeStream(ctx, omnivoice.TranscriptionConfig{
        Language: "en",
    })

    // Forward audio to STT
    go func() {
        buffer := make([]byte, 1024)
        for {
            n, err := transport.AudioOut().Read(buffer)
            if err != nil {
                break
            }
            sttStream.Write(buffer[:n])
        }
    }()

    // Process transcriptions
    for result := range sttStream.Results() {
        if !result.IsFinal {
            continue
        }

        // Generate response
        response := generateResponse(result.Text)

        // Stream TTS to call
        ttsStream, _ := tts.SynthesizeStream(ctx, response, omnivoice.SynthesisConfig{
            VoiceID: "pNInz6obpgDQGcFmaJgB",
        })

        for chunk := range ttsStream {
            transport.AudioIn().Write(chunk.Audio)
        }
    }

    return call.Hangup(ctx)
}
```

## Streaming Pipeline

For lowest latency, use streaming throughout:

```go
type VoiceAgent struct {
    stt omnivoice.STTProvider
    tts omnivoice.TTSProvider
    llm LLMClient // Your LLM client
}

func (a *VoiceAgent) HandleConversation(ctx context.Context, audioIn <-chan []byte, audioOut chan<- []byte) error {
    // Start STT stream
    sttStream, err := a.stt.TranscribeStream(ctx, omnivoice.TranscriptionConfig{
        Language: "en",
        Extensions: map[string]any{
            "interim_results": true,
        },
    })
    if err != nil {
        return err
    }
    defer sttStream.Close()

    // Forward audio to STT
    go func() {
        for audio := range audioIn {
            sttStream.Write(audio)
        }
    }()

    // Accumulate user speech
    var userText strings.Builder

    for result := range sttStream.Results() {
        if result.IsFinal {
            userText.WriteString(result.Text + " ")

            // Detect end of utterance (silence)
            if result.IsSpeechFinal {
                // Process complete utterance
                response := a.llm.Generate(ctx, userText.String())
                userText.Reset()

                // Stream response
                a.streamResponse(ctx, response, audioOut)
            }
        }
    }

    return nil
}

func (a *VoiceAgent) streamResponse(ctx context.Context, text string, audioOut chan<- []byte) {
    stream, _ := a.tts.SynthesizeStream(ctx, text, omnivoice.SynthesisConfig{
        VoiceID: "pNInz6obpgDQGcFmaJgB",
    })

    for chunk := range stream {
        audioOut <- chunk.Audio
    }
}
```

## Interruption Handling

Allow users to interrupt the agent:

```go
type InterruptibleAgent struct {
    speaking    atomic.Bool
    cancelSpeak context.CancelFunc
}

func (a *InterruptibleAgent) HandleConversation(ctx context.Context, audioIn <-chan []byte, audioOut chan<- []byte) {
    sttStream, _ := stt.TranscribeStream(ctx, config)

    go func() {
        for audio := range audioIn {
            sttStream.Write(audio)

            // If user starts speaking while agent is talking, interrupt
            if a.speaking.Load() {
                // Check for voice activity
                if hasVoiceActivity(audio) {
                    a.interrupt()
                }
            }
        }
    }()

    for result := range sttStream.Results() {
        if result.IsFinal {
            response := generateResponse(result.Text)
            a.speak(ctx, response, audioOut)
        }
    }
}

func (a *InterruptibleAgent) speak(ctx context.Context, text string, audioOut chan<- []byte) {
    speakCtx, cancel := context.WithCancel(ctx)
    a.cancelSpeak = cancel
    a.speaking.Store(true)
    defer a.speaking.Store(false)

    stream, _ := tts.SynthesizeStream(speakCtx, text, config)
    for chunk := range stream {
        select {
        case <-speakCtx.Done():
            return // Interrupted
        case audioOut <- chunk.Audio:
        }
    }
}

func (a *InterruptibleAgent) interrupt() {
    if a.cancelSpeak != nil {
        a.cancelSpeak()
    }
}
```

## Context Management

Maintain conversation context:

```go
type ConversationContext struct {
    Messages []Message
    mu       sync.Mutex
}

type Message struct {
    Role    string // "user" or "assistant"
    Content string
}

func (c *ConversationContext) AddUserMessage(text string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.Messages = append(c.Messages, Message{Role: "user", Content: text})
}

func (c *ConversationContext) AddAssistantMessage(text string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.Messages = append(c.Messages, Message{Role: "assistant", Content: text})
}

func (c *ConversationContext) GetPrompt() string {
    c.mu.Lock()
    defer c.mu.Unlock()

    var prompt strings.Builder
    prompt.WriteString("You are a helpful voice assistant. Previous conversation:\n\n")

    for _, msg := range c.Messages {
        prompt.WriteString(fmt.Sprintf("%s: %s\n", msg.Role, msg.Content))
    }

    return prompt.String()
}
```

## Latency Optimization

### 1. Use Fast Providers

```go
// STT: Deepgram Nova-2 for lowest latency
stt, _ := omnivoice.GetSTTProvider("deepgram", ...)

// TTS: Deepgram Aura or ElevenLabs Turbo
tts, _ := omnivoice.GetTTSProvider("deepgram", ...)
```

### 2. Stream Everything

```go
// Don't wait for complete transcription
for result := range sttStream.Results() {
    if result.IsFinal {
        // Process immediately
    }
}
```

### 3. Prefetch Common Responses

```go
var greetingAudio []byte

func init() {
    // Pre-generate common responses
    result, _ := tts.Synthesize(ctx, "Hello! How can I help you today?", config)
    greetingAudio = result.Audio
}
```

### 4. Use PCM Audio

```go
// Avoid encoding/decoding overhead
config := omnivoice.SynthesisConfig{
    OutputFormat: "pcm_16000", // Raw PCM
}
```

## Error Recovery

```go
func (a *VoiceAgent) HandleWithRecovery(ctx context.Context, call omnivoice.Call) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Agent panic: %v", r)
            // Apologize and hang up
            a.speak(ctx, "I'm sorry, I encountered an error. Goodbye.", call)
            call.Hangup(ctx)
        }
    }()

    for {
        err := a.handleConversation(ctx, call)
        if err != nil {
            if errors.Is(err, context.Canceled) {
                return
            }

            log.Printf("Conversation error: %v", err)

            // Try to recover
            a.speak(ctx, "I'm sorry, could you repeat that?", call)
            continue
        }
    }
}
```

## Best Practices

1. **Stream everything** - Minimize time-to-first-byte
2. **Handle interruptions** - Let users cut off the agent
3. **Maintain context** - Remember conversation history
4. **Graceful degradation** - Fall back on errors
5. **Monitor latency** - Track STT→LLM→TTS round-trip time
6. **Test with real calls** - Phone audio quality differs from local

## Latency Targets

| Stage | Target | Acceptable |
|-------|--------|------------|
| STT | < 200ms | < 300ms |
| LLM | < 300ms | < 500ms |
| TTS | < 150ms | < 250ms |
| **Total** | **< 650ms** | **< 1000ms** |

## Next Steps

- [Voice Calls](voice-calls.md) - Phone integration
- [Streaming](streaming.md) - Real-time audio
- [Provider Comparison](../providers/index.md) - Choose providers for latency
