# Voice Calls

OmniVoice supports two types of voice connections:

1. **PSTN (Phone Calls)** - Traditional phone calls via Twilio, Telnyx, Vonage, or Plivo
2. **WebRTC (Browser/Mobile)** - Direct connections via LiveKit for web and mobile apps

## PSTN Providers (Phone Calls)

| Provider | Registry Name | Audio Format | Auth Method |
|----------|---------------|--------------|-------------|
| Twilio | `"twilio"` | mulaw 8kHz | Account SID + Token |
| Telnyx | `"telnyx"` | L16 16kHz | API Key |
| Vonage | `"vonage"` | L16 16kHz | JWT (RS256) |
| Plivo | `"plivo"` | L16 16kHz | Auth ID + Token |

## WebRTC Providers (Browser/Mobile)

| Provider | Registry Name | Audio Format | Best For |
|----------|---------------|--------------|----------|
| LiveKit | `"livekit"` | Opus (configurable) | Custom apps, low latency |

## Quick Start (PSTN)

```go
provider, _ := omnivoice.GetCallSystemProvider("twilio",
    omnivoice.WithAccountSID(accountSID),
    omnivoice.WithAuthToken(authToken),
    omnivoice.WithPhoneNumber("+15551234567"),
)

call, _ := provider.MakeCall(ctx, "+15559876543")
fmt.Printf("Call ID: %s\n", call.ID())
```

## Quick Start (WebRTC)

```go
import "github.com/plexusone/omni-livekit/omnivoice/gateway"

gw, _ := gateway.New(gateway.Config{
    LiveKitURL:    os.Getenv("LIVEKIT_URL"),
    LiveKitAPIKey: os.Getenv("LIVEKIT_API_KEY"),
    LiveKitSecret: os.Getenv("LIVEKIT_API_SECRET"),
    RoomName:      "voice-agent",
    AgentIdentity: "ai-agent",
    SampleRate:    24000,
})

gw.OnParticipantJoined(func(p *gateway.ParticipantInfo) error {
    log.Printf("User joined: %s", p.Identity)
    return nil
})

gw.Start(ctx)
```

## Making Outbound Calls

### Basic Call

```go
provider, _ := omnivoice.GetCallSystemProvider("twilio",
    omnivoice.WithAccountSID(os.Getenv("TWILIO_ACCOUNT_SID")),
    omnivoice.WithAuthToken(os.Getenv("TWILIO_AUTH_TOKEN")),
    omnivoice.WithPhoneNumber("+15551234567"),
    omnivoice.WithWebhookURL("https://your-server.com/webhook"),
)

call, err := provider.MakeCall(ctx, "+15559876543")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Call initiated: %s\n", call.ID())
fmt.Printf("Status: %s\n", call.Status())
```

### With Options

```go
call, err := provider.MakeCall(ctx, "+15559876543",
    omnivoice.WithFrom("+15551234567"),      // Override default number
    omnivoice.WithTimeout(30 * time.Second), // Ring timeout
)
```

## Handling Incoming Calls

Set up a handler for incoming calls:

```go
provider.OnIncomingCall(func(call omnivoice.Call) error {
    log.Printf("Incoming call from %s to %s", call.From(), call.To())

    // Answer the call
    if err := call.Answer(ctx); err != nil {
        return err
    }

    // Handle the call...
    return nil
})
```

### Webhook Handler

You need an HTTP endpoint to receive call webhooks:

```go
http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
    // Parse webhook (provider-specific)
    // For Twilio:
    r.ParseForm()
    callSID := r.Form.Get("CallSid")
    from := r.Form.Get("From")
    to := r.Form.Get("To")

    // Handle the incoming call
    call, err := provider.HandleIncomingWebhook(callSID, from, to)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    // Return TwiML response
    w.Header().Set("Content-Type", "application/xml")
    w.Write([]byte(`<Response><Say>Hello!</Say></Response>`))
})
```

## Call Control

### Answer and Hangup

```go
// Answer an incoming call
err := call.Answer(ctx)

// Hang up
err = call.Hangup(ctx)
```

### Call Status

```go
status := call.Status()

switch status {
case omnivoice.StatusRinging:
    fmt.Println("Call is ringing")
case omnivoice.StatusAnswered:
    fmt.Println("Call is active")
case omnivoice.StatusEnded:
    fmt.Println("Call has ended")
}
```

### Call Information

```go
fmt.Printf("Call ID: %s\n", call.ID())
fmt.Printf("Direction: %s\n", call.Direction()) // inbound or outbound
fmt.Printf("From: %s\n", call.From())
fmt.Printf("To: %s\n", call.To())
fmt.Printf("Duration: %s\n", call.Duration())
fmt.Printf("Start Time: %s\n", call.StartTime())
```

## Media Streaming

Connect real-time audio for voice agents:

```go
// Get the transport provider
transport := provider.Transport()

// Listen for Media Stream connections
connCh, _ := transport.Listen(ctx, "/media-stream")

// Handle WebSocket connections
http.HandleFunc("/media-stream", func(w http.ResponseWriter, r *http.Request) {
    transport.HandleWebSocket(w, r, "/media-stream")
})

// Process audio
for conn := range connCh {
    go handleMediaStream(conn)
}

func handleMediaStream(conn transport.Connection) {
    defer conn.Close()

    // Read audio from caller
    audio := make([]byte, 1024)
    for {
        n, err := conn.AudioOut().Read(audio)
        if err != nil {
            break
        }

        // Process audio (e.g., send to STT)
        processAudio(audio[:n])

        // Send response audio
        conn.AudioIn().Write(responseAudio)
    }
}
```

## Twilio-Specific Features

### TwiML Responses

```go
// Generate TwiML for call handling
twiml := `<?xml version="1.0"?>
<Response>
    <Say voice="Polly.Joanna">Welcome to our service.</Say>
    <Gather input="speech" action="/handle-speech" timeout="3">
        <Say>Please tell us how we can help you.</Say>
    </Gather>
</Response>`
```

### Connect to Media Streams

```xml
<Response>
    <Connect>
        <Stream url="wss://your-server.com/media-stream" />
    </Connect>
</Response>
```

## Telnyx-Specific Features

### Call Control Commands

```go
// Cast to Telnyx call for advanced features
telnyxCall := call.(*telnyxcallsystem.Call)

// Speak using built-in TTS
telnyxCall.Speak(ctx, "Hello, how can I help you?", "", "")

// Play audio file
telnyxCall.PlayAudio(ctx, "https://example.com/greeting.mp3")

// Start transcription
telnyxCall.StartTranscription(ctx, "en")

// Start media streaming
telnyxCall.StartMediaStreaming(ctx, "wss://your-server.com/stream")
```

## Multi-Provider Failover

Use CallSystemClient for automatic failover:

```go
// Create providers
twilioProvider, _ := omnivoice.GetCallSystemProvider("twilio", ...)
telnyxProvider, _ := omnivoice.GetCallSystemProvider("telnyx", ...)

// Create client with failover
client := omnivoice.NewCallSystemClient(twilioProvider, telnyxProvider)
client.SetPrimary("twilio")
client.SetFallbacks("telnyx")

// MakeCall automatically falls back on failure
call, err := client.MakeCall(ctx, "+15559876543",
    omnivoice.WithFrom("+15551234567"),
)
```

## Error Handling

```go
call, err := provider.MakeCall(ctx, to)
if err != nil {
    switch {
    case strings.Contains(err.Error(), "invalid_phone_number"):
        log.Println("Invalid phone number format")
    case strings.Contains(err.Error(), "insufficient_funds"):
        log.Println("Account balance too low")
    case strings.Contains(err.Error(), "number_not_verified"):
        log.Println("Number not verified for trial account")
    default:
        log.Printf("Call error: %v", err)
    }
    return
}
```

## Best Practices

1. **Always set webhook URLs** - Required for call events and media
2. **Use HTTPS** - Webhooks must use secure connections
3. **Handle call state** - Track status changes properly
4. **Implement failover** - Use CallSystemClient for reliability
5. **Secure webhooks** - Validate webhook signatures

## Next Steps

- [SMS Messaging](sms.md) - Send text messages
- [Voice Agents](voice-agents.md) - Build conversational agents
- [Multi-Provider Failover](../reference/failover.md) - Reliability patterns
