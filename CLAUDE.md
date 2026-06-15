# CLAUDE.md

Project-specific guidelines for omnivoice.

## Overview

`omnivoice` is the **batteries-included** voice abstraction package. It re-exports the core interfaces from `omnivoice-core` and imports all provider packages for automatic registration.

For detailed architecture documentation, see [`omnivoice-core/CLAUDE.md`](https://github.com/plexusone/omnivoice-core/blob/main/CLAUDE.md).

## Package Structure

```
omnivoice/
├── providers/all/all.go   # Imports and registers all providers
├── *.go                   # Re-exports from omnivoice-core
└── cmd/                   # CLI tools
```

## Usage

Import `omnivoice` for the batteries-included experience:

```go
import (
    "github.com/plexusone/omnivoice"
    _ "github.com/plexusone/omnivoice/providers/all"  // Register all providers
)

func main() {
    // All providers are available
    tts, _ := omnivoice.GetTTSProvider("elevenlabs", omnivoice.WithAPIKey(key))
    stt, _ := omnivoice.GetSTTProvider("deepgram", omnivoice.WithAPIKey(key))
    cs, _ := omnivoice.GetCallSystemProvider("twilio", omnivoice.WithAPIKey(token),
        omnivoice.WithExtension("accountSID", accountSID))

    // Native voice-to-voice (realtime)
    rt, _ := omnivoice.GetRealtimeProvider("openai-realtime", omnivoice.WithAPIKey(key))
}
```

## Registered Providers

| Type | Provider | Package | Latency |
|------|----------|---------|---------|
| STT | `openai` | `omni-openai` | - |
| STT | `deepgram` | `omni-deepgram` | - |
| STT | `elevenlabs` | `elevenlabs-go` | - |
| STT | `twilio` | `omni-twilio` | - |
| TTS | `openai` | `omni-openai` | - |
| TTS | `deepgram` | `omni-deepgram` | - |
| TTS | `elevenlabs` | `elevenlabs-go` | - |
| TTS | `twilio` | `omni-twilio` | - |
| CallSystem | `twilio` | `omni-twilio` | - |
| CallSystem | `telnyx` | `omni-telnyx` | - |
| Realtime | `openai-realtime` | `omni-openai` | ~100ms |
| Realtime | `gemini-live` | `omni-google` | ~200ms |

## Dependency Architecture

```
omnivoice-core           ← Core interfaces + registry
     ↑
provider packages        ← Implement interfaces, register via init()
     ↑
omnivoice               ← Imports all providers (this package)
```

**Key principle:** Provider packages (omni-twilio, omni-deepgram, etc.) depend on `omnivoice-core`, NOT on this package. This package imports providers, not the other way around.

## Adding New Providers

1. Create provider package implementing interfaces from `omnivoice-core`
2. Register provider using `omnivoice.RegisterXXXProvider()`
3. Add import to `providers/all/all.go`

See `omnivoice-core/CLAUDE.md` for detailed provider implementation guidelines.

## Testing

```bash
go test ./...
```

## Linting

```bash
golangci-lint run
```
