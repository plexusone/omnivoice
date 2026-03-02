# Release Notes - v0.5.0

**Release Date:** 2026-03-01

## Overview

OmniVoice v0.5.0 is the first release of the batteries-included voice pipeline framework for Go. This package provides a unified interface for speech-to-text (STT) and text-to-speech (TTS) with all providers included out of the box.

## Highlights

- **Unified Voice Interface**: Single API for all STT and TTS providers, eliminating the need to learn provider-specific APIs
- **Provider Registry Pattern**: Consumers only need to import `github.com/plexusone/omnivoice` - no direct provider package imports required
- **Multi-Provider Support**: OpenAI (Whisper, TTS-1), Deepgram (Nova-2, Aura), ElevenLabs (Scribe, Multilingual v2), and Twilio (Media Streams)
- **Streaming Support**: Real-time transcription and synthesis for low-latency voice applications

## What's New

### Provider Registry

The key architectural feature of omnivoice is the provider registry. Instead of importing individual provider packages:

```go
// OLD WAY - requires knowing provider-specific packages
import (
    elevenlabstts "github.com/plexusone/elevenlabs-go/omnivoice/tts"
    deepgramtts "github.com/plexusone/omnivoice-deepgram/omnivoice/tts"
)
provider, _ := elevenlabstts.New(elevenlabstts.WithAPIKey(key))
```

You now use the unified registry:

```go
// NEW WAY - only import omnivoice
import (
    "github.com/plexusone/omnivoice"
    _ "github.com/plexusone/omnivoice/providers/all"
)
provider, _ := omnivoice.GetTTSProvider("elevenlabs", omnivoice.WithAPIKey(key))
```

### Core Features

- **STTClient**: Unified client for speech-to-text transcription across providers
- **TTSClient**: Unified client for text-to-speech synthesis across providers
- **TranscriptionConfig**: Provider-agnostic configuration for transcription requests
- **SynthesisConfig**: Provider-agnostic configuration for synthesis requests
- **providers/all**: Single import to register all available providers

### Registry Functions

| Function | Description |
|----------|-------------|
| `GetTTSProvider(name, opts...)` | Create a TTS provider by name |
| `GetSTTProvider(name, opts...)` | Create an STT provider by name |
| `ListTTSProviders()` | List registered TTS providers |
| `ListSTTProviders()` | List registered STT providers |
| `WithAPIKey(key)` | Option to set API key |
| `WithBaseURL(url)` | Option to set custom endpoint |
| `WithExtension(key, value)` | Option for provider-specific config |

### Included Providers

| Provider | STT | TTS | Registry Name |
|----------|-----|-----|---------------|
| OpenAI | Whisper | TTS-1/TTS-1-HD | `"openai"` |
| ElevenLabs | Scribe | Multilingual v2 | `"elevenlabs"` |
| Deepgram | Nova-2 | Aura | `"deepgram"` |
| Twilio | Media Streams | Media Streams | `"twilio"` |

## Installation

```bash
go get github.com/plexusone/omnivoice
```

## Quick Start

```go
import (
    "github.com/plexusone/omnivoice"
    _ "github.com/plexusone/omnivoice/providers/all"
)

// Get provider by name - runtime selection
ttsProvider, err := omnivoice.GetTTSProvider("elevenlabs", omnivoice.WithAPIKey(apiKey))
sttProvider, err := omnivoice.GetSTTProvider("deepgram", omnivoice.WithAPIKey(apiKey))
```

## Related Packages

For a minimal dependency footprint, use [omnivoice-core](https://github.com/plexusone/omnivoice-core) which contains only the interface definitions.

## License

MIT License - see [LICENSE](LICENSE) for details.
