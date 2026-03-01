# OmniVoice

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![License][license-svg]][license-url]

Batteries-included voice pipeline framework for Go. This package provides a unified interface for speech-to-text (STT) and text-to-speech (TTS) with all providers included.

For a minimal dependency footprint, use [omnivoice-core](https://github.com/plexusone/omnivoice-core) instead.

## Features

- **Unified Interface**: Single API for all STT and TTS providers
- **Multiple Providers**: OpenAI, Deepgram, ElevenLabs, Twilio
- **Streaming Support**: Real-time transcription and synthesis
- **Easy Integration**: Import and use with minimal configuration

## Installation

```bash
go get github.com/plexusone/omnivoice
```

## Quick Start

Import with all providers:

```go
import (
    "github.com/plexusone/omnivoice"
    _ "github.com/plexusone/omnivoice/providers/all"
)
```

Or import specific providers:

```go
import (
    "github.com/plexusone/omnivoice"
    openai "github.com/plexusone/omnivoice-openai/omnivoice"
)
```

## Usage

```go
package main

import (
    "context"
    "log"

    "github.com/plexusone/omnivoice"
    openaistt "github.com/plexusone/omnivoice-openai/omnivoice/stt"
    openaitts "github.com/plexusone/omnivoice-openai/omnivoice/tts"
)

func main() {
    ctx := context.Background()

    // Create STT provider and client
    sttProvider := openaistt.NewProvider()
    sttClient := omnivoice.NewSTTClient(sttProvider)

    // Transcribe audio
    result, err := sttClient.Transcribe(ctx, audioData, omnivoice.TranscriptionConfig{
        Language: "en",
    })
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Transcription: %s", result.Text)

    // Create TTS provider and client
    ttsProvider := openaitts.NewProvider()
    ttsClient := omnivoice.NewTTSClient(ttsProvider)

    // Synthesize speech
    audio, err := ttsClient.Synthesize(ctx, "Hello, world!", omnivoice.SynthesisConfig{})
    if err != nil {
        log.Fatal(err)
    }
    // audio.Data contains the audio bytes
}
```

## Included Providers

| Provider | STT | TTS | Package |
|----------|-----|-----|---------|
| OpenAI | Whisper | TTS-1/TTS-1-HD | `omnivoice-openai` |
| Deepgram | Nova-2 | Aura | `omnivoice-deepgram` |
| ElevenLabs | - | Multilingual v2 | `elevenlabs-go` |
| Twilio | Media Streams | Media Streams | `omnivoice-twilio` |

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
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/plexusone/omnivoice/blob/main/LICENSE
