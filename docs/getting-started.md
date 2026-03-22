# Getting Started

This guide walks you through setting up OmniVoice and making your first API calls.

## Installation

```bash
go get github.com/plexusone/omnivoice
```

## Basic Setup

Import OmniVoice and register all providers:

```go
import (
    "github.com/plexusone/omnivoice"
    _ "github.com/plexusone/omnivoice/providers/all" // Register all providers
)
```

The blank import `_ "github.com/plexusone/omnivoice/providers/all"` registers all available providers with the registry during `init()`.

## API Keys

Each provider requires its own API key. Set them as environment variables:

```bash
# TTS/STT Providers
export OPENAI_API_KEY="sk-..."
export ELEVENLABS_API_KEY="..."
export DEEPGRAM_API_KEY="..."

# Telephony Providers
export TWILIO_ACCOUNT_SID="AC..."
export TWILIO_AUTH_TOKEN="..."
export TELNYX_API_KEY="KEY..."
```

## Your First TTS Request

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

    // Get ElevenLabs TTS provider
    provider, err := omnivoice.GetTTSProvider("elevenlabs",
        omnivoice.WithAPIKey(os.Getenv("ELEVENLABS_API_KEY")))
    if err != nil {
        log.Fatal(err)
    }

    // Synthesize speech
    result, err := provider.Synthesize(ctx, "Hello from OmniVoice!", omnivoice.SynthesisConfig{
        VoiceID:      "pNInz6obpgDQGcFmaJgB", // Adam
        OutputFormat: "mp3_44100_128",
    })
    if err != nil {
        log.Fatal(err)
    }

    // Save to file
    if err := os.WriteFile("hello.mp3", result.Audio, 0644); err != nil {
        log.Fatal(err)
    }

    log.Println("Saved hello.mp3")
}
```

## Your First STT Request

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

    // Get Deepgram STT provider
    provider, err := omnivoice.GetSTTProvider("deepgram",
        omnivoice.WithAPIKey(os.Getenv("DEEPGRAM_API_KEY")))
    if err != nil {
        log.Fatal(err)
    }

    // Transcribe audio file
    result, err := provider.TranscribeFile(ctx, "audio.mp3", omnivoice.TranscriptionConfig{
        Language:             "en",
        EnableWordTimestamps: true,
    })
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Transcription: %s", result.Text)
}
```

## Making a Phone Call

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

    // Get Twilio CallSystem provider
    provider, err := omnivoice.GetCallSystemProvider("twilio",
        omnivoice.WithAccountSID(os.Getenv("TWILIO_ACCOUNT_SID")),
        omnivoice.WithAuthToken(os.Getenv("TWILIO_AUTH_TOKEN")),
        omnivoice.WithPhoneNumber("+15551234567"),
        omnivoice.WithWebhookURL("https://your-server.com/webhook"),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Make outbound call
    call, err := provider.MakeCall(ctx, "+15559876543")
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Call initiated: %s", call.ID())
}
```

## Sending SMS

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

    // Get Twilio CallSystem provider (also implements SMSProvider)
    provider, err := omnivoice.GetCallSystemProvider("twilio",
        omnivoice.WithAccountSID(os.Getenv("TWILIO_ACCOUNT_SID")),
        omnivoice.WithAuthToken(os.Getenv("TWILIO_AUTH_TOKEN")),
        omnivoice.WithPhoneNumber("+15551234567"),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Cast to SMSProvider
    smsProvider := provider.(omnivoice.SMSProvider)

    // Send SMS
    msg, err := smsProvider.SendSMS(ctx, "+15559876543", "Hello from OmniVoice!")
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Message sent: %s", msg.ID)
}
```

## Listing Available Providers

```go
package main

import (
    "fmt"

    "github.com/plexusone/omnivoice"
    _ "github.com/plexusone/omnivoice/providers/all"
)

func main() {
    fmt.Println("TTS Providers:", omnivoice.ListTTSProviders())
    fmt.Println("STT Providers:", omnivoice.ListSTTProviders())
    fmt.Println("CallSystem Providers:", omnivoice.ListCallSystemProviders())
}
```

Output:

```
TTS Providers: [openai elevenlabs deepgram twilio]
STT Providers: [openai elevenlabs deepgram twilio]
CallSystem Providers: [twilio telnyx]
```

## Next Steps

- [TTS Guide](guides/tts.md) - Deep dive into text-to-speech
- [STT Guide](guides/stt.md) - Deep dive into speech-to-text
- [Voice Calls](guides/voice-calls.md) - Phone call handling
- [Provider Comparison](providers/index.md) - Choose the right providers
