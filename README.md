# OmniVoice

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

Batteries-included voice pipeline framework for Go. This package provides a unified interface for speech-to-text (STT) and text-to-speech (TTS) with all providers included.

For a minimal dependency footprint, use [omnivoice-core](https://github.com/plexusone/omnivoice-core) instead.

## Features

- **Unified Interface**: Single API for all STT and TTS providers
- **Provider Registry**: Get providers by name - no need to import individual provider packages
- **Multiple Providers**: OpenAI, Deepgram, ElevenLabs, Twilio
- **Streaming Support**: Real-time transcription and synthesis
- **Easy Integration**: Import and use with minimal configuration

## Installation

```bash
go get github.com/plexusone/omnivoice
```

## Quick Start

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
// Available providers: "openai", "elevenlabs", "deepgram", "twilio"
ttsProvider, _ := omnivoice.GetTTSProvider("elevenlabs", omnivoice.WithAPIKey(key))
sttProvider, _ := omnivoice.GetSTTProvider("deepgram", omnivoice.WithAPIKey(key))

// List registered providers
fmt.Println(omnivoice.ListTTSProviders()) // [openai elevenlabs deepgram twilio]
fmt.Println(omnivoice.ListSTTProviders()) // [openai elevenlabs deepgram twilio]
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

| Provider | STT | TTS | Registry Name |
|----------|-----|-----|---------------|
| OpenAI | Whisper | TTS-1/TTS-1-HD | `"openai"` |
| ElevenLabs | Scribe | Multilingual v2 | `"elevenlabs"` |
| Deepgram | Nova-2 | Aura | `"deepgram"` |
| Twilio | Media Streams | Media Streams | `"twilio"` |

## Related Packages

- [omnivoice-core](https://github.com/plexusone/omnivoice-core) - Core interfaces (minimal dependencies)
- [omnivoice-openai](https://github.com/plexusone/omnivoice-openai) - OpenAI provider
- [omnivoice-deepgram](https://github.com/plexusone/omnivoice-deepgram) - Deepgram provider
- [omnivoice-twilio](https://github.com/plexusone/omnivoice-twilio) - Twilio provider
- [elevenlabs-go](https://github.com/plexusone/elevenlabs-go) - ElevenLabs SDK

## License

MIT License - see [LICENSE](LICENSE) for details.

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
 [viz-svg]: https://img.shields.io/badge/visualizaton-Go-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=plexusone%2Fomnivoice
 [loc-svg]: https://tokei.rs/b1/github/plexusone/omnivoice
 [repo-url]: https://github.com/plexusone/omnivoice
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/plexusone/omnivoice/blob/master/LICENSE