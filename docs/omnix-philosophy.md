# OmniX Philosophy and Repository Patterns

This document describes the naming conventions and architectural patterns used across PlexusOne's OmniX family of Go libraries.

## Overview

The OmniX libraries provide unified interfaces for interacting with multiple providers in specific domains:

| Library | Domain | Providers |
|---------|--------|-----------|
| OmniLLM | Large Language Models | OpenAI, Anthropic, Gemini, Ollama, etc. |
| OmniVoice | Voice (TTS, STT, Agents) | ElevenLabs, OpenAI, Deepgram, Twilio, etc. |
| OmniVault | Secret Management | AWS, 1Password, Keyring, etc. |
| OmniSerp | Search | Google, Bing, etc. |
| OmniAgent | AI Agents | Claude Code, Custom agents |

## Repository Naming Patterns

### Pattern 1: `omni<type>-core` — Interface Definitions

**Purpose:** Contains only interface definitions and shared types with zero external dependencies.

**Examples:**
- `omnivoice-core` — Voice interfaces (tts.Provider, stt.Provider, agent.Provider)
- `omnillm-core` — LLM interfaces (provider.Provider, Capabilities, types)

**Characteristics:**
- Minimal dependencies (typically zero)
- Defines the contracts that providers implement
- Versioned independently
- Other packages depend on this, not vice versa

```go
// omnivoice-core/tts/provider.go
type Provider interface {
    Synthesize(ctx context.Context, text string, config Config) (*Result, error)
    Capabilities() Capabilities
    Close() error
}
```

### Pattern 2: `omni<type>-<provider>` — Provider Adapters (External SDK)

**Purpose:** Implements omni<type>-core interfaces using an external/official SDK.

**When to use:** When the provider has an official Go SDK that we wrap.

**Examples:**
- `omni-openai` — Uses `github.com/openai/openai-go`
- `omni-deepgram` — Uses Deepgram's SDK
- `omni-openai/omnillm` — Uses `github.com/openai/openai-go`
- `omnillm-anthropic` — Uses `github.com/anthropics/anthropic-sdk-go`

**Dependencies:**
```
omni-openai
├── github.com/openai/openai-go (official SDK)
└── github.com/plexusone/omnivoice-core (interfaces)
```

### Pattern 3: `<provider>-go` — Full Provider SDK

**Purpose:** A complete SDK for a provider's API, written by PlexusOne.

**When to use:** When no official SDK exists, or we need deeper integration.

**Examples:**
- `elevenlabs-go` — Full ElevenLabs API SDK
- `opik-go` — Full Comet Opik SDK
- `phoenix-go` — Full Arize Phoenix SDK

**Characteristics:**
- Covers all/most provider API endpoints
- May include CLI tools
- Generated from OpenAPI specs (using ogen) when available

### Pattern 4: `<provider>-go/omni<type>/` — Embedded Adapter

**Purpose:** OmniX adapter embedded within a full SDK we own.

**When to use:** When we write the full SDK (`<provider>-go`), the adapter lives inside as a subdirectory.

**Examples:**
- `elevenlabs-go/omnivoice/` — OmniVoice adapter using elevenlabs-go
- `elevenlabs-go/ax/` — Additional integrations

**Structure:**
```
elevenlabs-go/
├── client.go          # Core SDK
├── audio.go           # API endpoints
├── omnivoice/         # OmniVoice adapter
│   ├── tts/
│   └── stt/
└── cmd/               # CLI tools
```

### Pattern 5: `omni<type>` — Batteries-Included

**Purpose:** Convenience package that pulls in multiple providers for easy setup.

**Examples:**
- `omnivoice` — Imports multiple voice providers
- `omnivault` — Imports multiple vault providers
- `omnillm` — Multi-provider LLM client with fallback, caching, etc.

**Characteristics:**
- Many dependencies (pulls in providers)
- Provides registry/factory for dynamic provider selection
- Higher-level features (fallback, load balancing, caching)
- Like [rclone](https://github.com/rclone/rclone) pattern

```go
// omnivoice/registry.go
import (
    _ "github.com/plexusone/omni-openai/omnivoice"
    _ "github.com/plexusone/omni-deepgram/omnivoice/tts"
    _ "github.com/plexusone/elevenlabs-go/omnivoice/tts"
)
```

## Decision Tree

```
Do you need to interact with provider X?
│
├─► Does provider X have an official Go SDK?
│   │
│   ├─► YES: Create `omni<type>-<provider>` adapter
│   │        Example: omni-openai wraps openai/openai-go
│   │
│   └─► NO: Do you need full API coverage?
│       │
│       ├─► YES: Create `<provider>-go` SDK
│       │        with `<provider>-go/omni<type>/` adapter inside
│       │        Example: elevenlabs-go with elevenlabs-go/omnivoice/
│       │
│       └─► NO: Create `omni<type>-<provider>` with direct HTTP
│                (Consider upgrading to full SDK later)
```

## Versioning Strategy

Each package is versioned independently:

| Package | Version | Notes |
|---------|---------|-------|
| `omnivoice-core` | v0.8.0 | Stable interfaces |
| `omni-openai` | v0.5.0 | Tracks openai-go updates |
| `elevenlabs-go` | v0.12.0 | Full SDK releases |
| `omnivoice` | v0.9.0 | Batteries-included |

**Compatibility rules:**
- `omni<type>-core` changes require major version bump for breaking changes
- Provider adapters can update independently
- Batteries-included packages pin to compatible provider versions

## Provider Capabilities

All providers implement a `Capabilities()` method to advertise supported features:

```go
type Capabilities struct {
    Tools             bool  // Tool/function calling
    Streaming         bool  // Streaming responses
    Vision            bool  // Image inputs
    JSON              bool  // JSON response mode
    SystemRole        bool  // System message support
    MaxContextWindow  int   // Max tokens
}
```

This allows runtime feature detection:

```go
if provider.Capabilities().Tools {
    // Use tool calling
} else {
    // Fallback to prompt-based approach
}
```

## Migration Guide: omnillm

Current state:
```
omnillm/
├── provider/          # Interfaces (→ omnillm-core)
├── providers/         # Embedded implementations (→ omnillm-{openai,anthropic,...})
├── client.go          # Stays in omnillm
└── fallback.go        # Stays in omnillm
```

Target state:
```
omnillm-core/          # NEW: Interfaces only
├── provider.go
├── types.go
└── capabilities.go

omni-openai/omnillm/        # NEW: Uses official SDK
└── adapter.go

omnillm-anthropic/     # NEW: Uses official SDK
└── adapter.go

omnillm/               # Batteries-included (depends on above)
├── client.go
├── fallback.go
└── registry.go
```

## References

- [rclone](https://github.com/rclone/rclone) — Inspiration for batteries-included pattern
- [openai-go](https://github.com/openai/openai-go) — Official OpenAI SDK
- [anthropic-sdk-go](https://github.com/anthropics/anthropic-sdk-go) — Official Anthropic SDK
