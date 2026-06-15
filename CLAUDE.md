# CLAUDE.md

Project-specific guidelines for omnivoice.

## Overview

`omnivoice` is the **batteries-included** voice abstraction package. It re-exports the core interfaces from `omnivoice-core` and delegates to its registry for provider discovery.

For detailed architecture documentation, see [`omnivoice-core/CLAUDE.md`](https://github.com/plexusone/omnivoice-core/blob/main/CLAUDE.md).

## Package Structure

```
omnivoice/
├── registry.go            # Delegates to omnivoice-core registry
├── realtime.go            # Re-exports realtime types
├── gateway.go             # Re-exports gateway types
├── stt.go, tts.go, etc.   # Re-exports from omnivoice-core
├── providers/all/all.go   # Imports and registers all providers
└── cmd/                   # CLI tools
```

## Registry Pattern

`omnivoice` delegates all registry operations to `omnivoice-core`. This means:

1. Provider packages register with `omnivoice-core` via `init()`
2. `omnivoice.Get*Provider()` calls delegate to `omnivoice-core`
3. Both packages share the same registry state

```go
import (
    "github.com/plexusone/omnivoice"
    _ "github.com/plexusone/omnivoice/providers/all"  // Register all providers
)

func main() {
    // STT/TTS providers
    tts, _ := omnivoice.GetTTSProvider("elevenlabs", omnivoice.WithAPIKey(key))
    stt, _ := omnivoice.GetSTTProvider("deepgram", omnivoice.WithAPIKey(key))

    // CallSystem providers
    cs, _ := omnivoice.GetCallSystemProvider("twilio",
        omnivoice.WithAccountSID(accountSID),
        omnivoice.WithAuthToken(authToken))

    // Gateway providers (voice pipelines)
    gw, _ := omnivoice.GetGatewayProvider("twilio",
        omnivoice.WithAccountSID(accountSID),
        omnivoice.WithAuthToken(authToken),
        omnivoice.WithPublicURL(publicURL))

    // Realtime providers (voice-to-voice)
    rt, _ := omnivoice.GetRealtimeProvider("openai", omnivoice.WithAPIKey(key))
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
| Gateway | `twilio` | `omni-twilio` | - |
| Gateway | `telnyx` | `omni-telnyx` | - |
| Realtime | `openai` | `omni-openai` | ~100ms |
| Realtime | `gemini` | `omni-google` | ~200ms |

## Dependency Architecture

```
omnivoice-core           ← Core interfaces + global registry
     ↑
provider packages        ← Implement interfaces, register via init()
     ↑
omnivoice               ← Imports all providers, delegates to core registry
```

**Key principles:**

1. Provider packages depend on `omnivoice-core`, NOT on this package
2. This package imports providers for side-effect registration
3. Registry functions delegate to `omnivoice-core` (single source of truth)

## Adding New Providers

### Realtime/Gateway Providers (auto-registration)

1. Create provider package with `init()` that registers with `omnivoice-core`
2. Add side-effect import to `providers/all/all.go`:
   ```go
   _ "github.com/plexusone/omni-newprovider/omnivoice/realtime"
   ```

### STT/TTS/CallSystem Providers (manual registration)

1. Create provider package implementing interfaces from `omnivoice-core`
2. Add factory registration in `providers/all/all.go`:
   ```go
   omnivoice.RegisterSTTProvider("newprovider", func(config registry.ProviderConfig) (stt.Provider, error) {
       return newprovider.New(config.APIKey), nil
   }, omnivoice.PriorityThick)
   ```

See `omnivoice-core/CLAUDE.md` for detailed provider implementation guidelines.

## Testing

```bash
go test ./...
```

## Linting

```bash
golangci-lint run
```
