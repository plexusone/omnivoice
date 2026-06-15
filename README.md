# OmniVoice

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Docs][docs-mkdoc-svg]][docs-mkdoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

 [go-ci-svg]: https://github.com/plexusone/omnivoice/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/plexusone/omnivoice/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/plexusone/omnivoice/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/plexusone/omnivoice/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/plexusone/omnivoice/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/plexusone/omnivoice/actions/workflows/go-sast-codeql.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/plexusone/omnivoice
 [goreport-url]: https://goreportcard.com/report/github.com/plexusone/omnivoice
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/plexusone/omnivoice
 [docs-godoc-url]: https://pkg.go.dev/github.com/plexusone/omnivoice
 [docs-mkdoc-svg]: https://img.shields.io/badge/Go-dev%20guide-blue.svg
 [docs-mkdoc-url]: https://plexusone.dev/omnivoice
 [viz-svg]: https://img.shields.io/badge/Go-visualizaton-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=plexusone%2Fomnivoice
 [loc-svg]: https://tokei.rs/b1/github/plexusone/omnivoice
 [repo-url]: https://github.com/plexusone/omnivoice
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/plexusone/omnivoice/blob/main/LICENSE

Batteries-included voice pipeline framework for Go. This package provides a unified interface for speech-to-text (STT) and text-to-speech (TTS) with all providers included.

For a minimal dependency footprint, use [omnivoice-core](https://github.com/plexusone/omnivoice-core) instead.

## Voice Architecture

OmniVoice supports two approaches for real-time voice:

| Approach | Latency | Use Case |
|----------|---------|----------|
| **Traditional Pipeline** (STT→LLM→TTS) | 500-1500ms | Custom voices, domain-specific STT |
| **Native Voice-to-Voice** | 100-200ms | Low latency, natural conversation |

This package provides the **Traditional Pipeline** components. For native voice-to-voice, see:

- [omni-openai/omnivoice/realtime](https://github.com/plexusone/omni-openai) - OpenAI Realtime API (~100ms)
- [omni-google/omnivoice](https://github.com/plexusone/omni-google) - Gemini Live API (~200ms)

See the [Voice Architecture Guide](https://plexusone.dev/omnivoice-core/voice-architecture) for detailed comparison.

## Features

- 🎯 **Unified Interface**: Single API for all STT and TTS providers
- 🗂️ **Provider Registry**: Get providers by name - no need to import individual provider packages
- 🔌 **Multiple Providers**: OpenAI, Deepgram, ElevenLabs, Twilio, Telnyx
- ⚡ **Streaming Support**: Real-time transcription and synthesis
- 🚀 **Easy Integration**: Import and use with minimal configuration

## Installation

```bash
go get github.com/plexusone/omnivoice
```

## CLI

OmniVoice includes a command-line tool for transcription.

### Install CLI

```bash
go install github.com/plexusone/omnivoice/cmd/omnivoice@latest
```

### Usage

```bash
# Set your API key
export DEEPGRAM_API_KEY="your-api-key"

# Basic transcription (stdout)
omnivoice transcribe podcast.mp3

# Save to file
omnivoice transcribe -p deepgram -o transcript.txt podcast.mp3

# JSON output with full metadata (OmniVoice Transcript format)
omnivoice transcribe -p deepgram --diarize --timestamps -f json -o transcript.json podcast.mp3

# Generate SRT subtitles
omnivoice transcribe -p deepgram -f srt -o subtitles.srt podcast.mp3

# Generate WebVTT subtitles
omnivoice transcribe -p deepgram -f vtt -o subtitles.vtt podcast.mp3

# List available providers
omnivoice providers list
```

### Output Formats

| Format | Description |
|--------|-------------|
| `text` | Plain transcript text (default) |
| `json` | OmniVoice Transcript format with full metadata |
| `srt`  | SubRip subtitles |
| `vtt`  | WebVTT subtitles |

### Environment Variables

| Variable | Provider |
|----------|----------|
| `DEEPGRAM_API_KEY` | Deepgram |
| `OPENAI_API_KEY` | OpenAI |
| `ELEVENLABS_API_KEY` | ElevenLabs |

## Quick Start (Library)

```go
import (
    "github.com/plexusone/omnivoice"
    _ "github.com/plexusone/omnivoice/providers/all" // Register all providers
)
```

## Usage

```go
package main

import (
    "context"
    "log"
    "os"

    "github.com/plexusone/omnivoice"
    _ "github.com/plexusone/omnivoice/providers/all"
)

func main() {
    ctx := context.Background()

    // Get providers by name using the registry
    sttProvider, err := omnivoice.GetSTTProvider("deepgram",
        omnivoice.WithAPIKey(os.Getenv("DEEPGRAM_API_KEY")))
    if err != nil {
        log.Fatal(err)
    }

    ttsProvider, err := omnivoice.GetTTSProvider("elevenlabs",
        omnivoice.WithAPIKey(os.Getenv("ELEVENLABS_API_KEY")))
    if err != nil {
        log.Fatal(err)
    }

    // Transcribe audio
    result, err := sttProvider.TranscribeFile(ctx, "audio.mp3", omnivoice.TranscriptionConfig{
        Language:             "en",
        EnableWordTimestamps: true,
    })
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Transcription: %s", result.Text)

    // Synthesize speech
    audio, err := ttsProvider.Synthesize(ctx, "Hello, world!", omnivoice.SynthesisConfig{
        VoiceID: "pNInz6obpgDQGcFmaJgB", // Adam
    })
    if err != nil {
        log.Fatal(err)
    }
    // audio.Audio contains the audio bytes
}
```

## Provider Registry

Get providers by name at runtime - no need to import individual provider packages:

```go
// STT/TTS providers: "openai", "elevenlabs", "deepgram", "twilio"
ttsProvider, _ := omnivoice.GetTTSProvider("elevenlabs", omnivoice.WithAPIKey(key))
sttProvider, _ := omnivoice.GetSTTProvider("deepgram", omnivoice.WithAPIKey(key))

// Realtime providers: "openai", "gemini"
rtProvider, _ := omnivoice.GetRealtimeProvider("openai", omnivoice.WithAPIKey(key))

// Gateway providers: "twilio", "telnyx"
gateway, _ := omnivoice.GetGatewayProvider("twilio", omnivoice.WithAccountSID(sid), omnivoice.WithAuthToken(token))

// List registered providers
fmt.Println(omnivoice.ListTTSProviders())      // [openai elevenlabs deepgram twilio]
fmt.Println(omnivoice.ListSTTProviders())      // [openai elevenlabs deepgram twilio]
fmt.Println(omnivoice.ListRealtimeProviders()) // [openai gemini]
fmt.Println(omnivoice.ListGatewayProviders())  // [twilio telnyx]

// Realtime factory API (for gateway integration)
factory, _ := omnivoice.GetRealtimeFactory("openai")
fmt.Println(omnivoice.ListRealtimeFactories()) // [openai gemini]
```

## Language Codes

OmniVoice accepts language codes in [BCP-47](https://www.rfc-editor.org/info/bcp47) format, which includes ISO 639-1 two-letter codes and regional variants.

**Common codes:**

| Code | Language |
|------|----------|
| `en` | English |
| `en-US` | English (US) |
| `en-GB` | English (UK) |
| `es` | Spanish |
| `es-MX` | Spanish (Mexico) |
| `fr` | French |
| `de` | German |
| `it` | Italian |
| `pt` | Portuguese |
| `pt-BR` | Portuguese (Brazil) |
| `ja` | Japanese |
| `ko` | Korean |
| `zh` | Chinese |
| `zh-CN` | Chinese (Simplified) |
| `zh-TW` | Chinese (Traditional) |
| `ar` | Arabic |
| `hi` | Hindi |
| `ru` | Russian |

**Notes:**

- Use simple codes (`en`) for broad compatibility across providers
- Use regional variants (`en-US`) when accent/dialect matters for TTS
- Provider support varies; see provider documentation for full language lists
- STT providers generally support automatic language detection when no code is specified

## Included Providers

### STT/TTS Providers

| Provider | STT | TTS | Registry Name |
|----------|-----|-----|---------------|
| OpenAI | Whisper | TTS-1/TTS-1-HD | `"openai"` |
| ElevenLabs | Scribe | Multilingual v2 | `"elevenlabs"` |
| Deepgram | Nova-2 | Aura | `"deepgram"` |
| Twilio | Media Streams | Media Streams | `"twilio"` |

### Native Voice-to-Voice (Realtime)

| Provider | Latency | Registry Name |
|----------|---------|---------------|
| OpenAI Realtime | ~100ms | `"openai"` |
| Gemini Live | ~200ms | `"gemini"` |

### Voice Gateway

| Provider | Registry Name |
|----------|---------------|
| Twilio | `"twilio"` |
| Telnyx | `"telnyx"` |

## Related Packages

### Core

- [omnivoice-core](https://github.com/plexusone/omnivoice-core) - Core interfaces (minimal dependencies)

### Native Voice-to-Voice (Recommended for Low Latency)

- [omni-openai/omnivoice/realtime](https://github.com/plexusone/omni-openai) - OpenAI Realtime API (~100ms latency)
- [omni-google/omnivoice](https://github.com/plexusone/omni-google) - Gemini Live API (~200ms latency)

### STT/TTS Providers

- [omni-openai](https://github.com/plexusone/omni-openai) - OpenAI provider (Whisper, TTS-1)
- [omni-deepgram](https://github.com/plexusone/omni-deepgram) - Deepgram provider (Nova-2, Aura)
- [elevenlabs-go](https://github.com/plexusone/elevenlabs-go) - ElevenLabs SDK

### Voice Gateway Providers

- [omni-twilio](https://github.com/plexusone/omni-twilio) - Twilio Media Streams
- [omni-telnyx](https://github.com/plexusone/omni-telnyx) - Telnyx Media Streaming
- [omni-vonage](https://github.com/plexusone/omni-vonage) - Vonage Voice WebSocket
- [omni-plivo](https://github.com/plexusone/omni-plivo) - Plivo Stream API
- [omni-livekit](https://github.com/plexusone/omni-livekit) - LiveKit WebRTC (web/mobile)

## License

MIT License - see [LICENSE](LICENSE) for details.
