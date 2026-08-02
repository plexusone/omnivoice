# OmniVoice

Batteries-included voice pipeline framework for Go. Single import, all providers included.

## Why OmniVoice?

- 🎯 **Single Import** - One package for all STT, TTS, and telephony needs
- 🗂️ **Provider Registry** - Get providers by name at runtime
- 🔌 **Multiple Providers** - OpenAI, ElevenLabs, Deepgram, Twilio, Telnyx
- ⚡ **Streaming Support** - Real-time transcription and synthesis
- 📞 **Voice Calls** - Make and receive phone calls with CallSystem
- 💬 **SMS** - Send text messages via SMSProvider
- 🖥️ **CLI Tool** - Transcribe audio from the command line
- 📄 **Transcript Format** - Canonical JSON format with timestamps and metadata
- ⚡ **Native Voice-to-Voice** - OpenAI Realtime (~100ms) and Gemini Live (~200ms)

## Quick Example

```go
package main

import (
    "context"
    "log"
    "os"

    "github.com/plexusone/omnivoice"
    _ "github.com/plexusone/omnivoice/providers/all" // Register all providers
)

func main() {
    ctx := context.Background()

    // Get a TTS provider by name
    tts, err := omnivoice.GetTTSProvider("elevenlabs",
        omnivoice.WithAPIKey(os.Getenv("ELEVENLABS_API_KEY")))
    if err != nil {
        log.Fatal(err)
    }

    // Synthesize speech
    result, err := tts.Synthesize(ctx, "Hello, world!", omnivoice.SynthesisConfig{
        VoiceID: "pNInz6obpgDQGcFmaJgB",
    })
    if err != nil {
        log.Fatal(err)
    }

    // result.Audio contains the audio bytes
    log.Printf("Generated %d bytes of audio", len(result.Audio))
}
```

## Installation

```bash
go get github.com/plexusone/omnivoice
```

## Use Cases

| Use Case | Guide | Providers |
|----------|-------|-----------|
| CLI Transcription | [CLI Guide](guides/cli.md) | All STT providers |
| Text-to-Speech | [TTS Guide](guides/tts.md) | OpenAI, ElevenLabs, Deepgram, Twilio |
| Speech-to-Text | [STT Guide](guides/stt.md) | OpenAI, ElevenLabs, Deepgram, Twilio |
| Voice Calls | [Voice Calls](guides/voice-calls.md) | Twilio, Telnyx |
| SMS Messaging | [SMS Guide](guides/sms.md) | Twilio, Telnyx |
| Real-time Streaming | [Streaming](guides/streaming.md) | All providers |
| Subtitles | [Subtitles](guides/subtitles.md) | All STT providers |
| Voice Agents | [Voice Agents](guides/voice-agents.md) | Combined stack |
| Native Voice-to-Voice | [v0.9.0 Release](releases/v0.9.0.md) | OpenAI Realtime, Gemini Live |

## Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                                  OmniVoice                                  │
│                            (batteries-included)                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   import _ "github.com/plexusone/omnivoice/providers/all"                   │
│                                                                             │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                           Provider Registry                         │   │
│   ├────────────┬────────────┬─────────────┬─────────────┬───────────────┤   │
│   │    TTS     │     STT    │  CallSystem │   Realtime  │  SMSProvider  │   │
│   ├────────────┼────────────┼─────────────┼─────────────┼───────────────┤   │
│   │ elevenlabs │ elevenlabs │   twilio    │ openai-     │    twilio     │   │
│   │ openai     │ openai     │   telnyx    │  realtime   │    telnyx     │   │
│   │ deepgram   │ deepgram   │             │ gemini-live │               │   │
│   │ twilio     │ twilio     │             │             │               │   │
│   └────────────┴────────────┴─────────────┴─────────────┴───────────────┘   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Minimal Dependencies Alternative

For a minimal dependency footprint, use [omnivoice-core](https://github.com/plexusone/omnivoice-core) with only the providers you need:

```go
import (
    "github.com/plexusone/omnivoice-core/tts"
    elevenlabs "github.com/plexusone/elevenlabs-go/omnivoice/tts"
)
```

## Next Steps

- [Getting Started](getting-started.md) - Detailed setup guide
- [Provider Comparison](providers/index.md) - Choose the right providers
- [API Reference](https://pkg.go.dev/github.com/plexusone/omnivoice) - GoDoc
