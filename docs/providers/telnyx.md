# Telnyx

Telnyx provides programmable voice calls and SMS messaging with competitive pricing.

## Features

- **Voice Calls**: Inbound and outbound PSTN calls
- **SMS**: Send and receive text messages
- **Media Streaming**: Real-time audio via WebSocket
- **Call Control API**: Programmable call manipulation
- **Global Coverage**: 150+ countries

## Configuration

```go
import (
    "github.com/plexusone/omnivoice"
    _ "github.com/plexusone/omnivoice/providers/all"
)

provider, err := omnivoice.GetCallSystemProvider("telnyx",
    omnivoice.WithAPIKey(os.Getenv("TELNYX_API_KEY")),
    omnivoice.WithPhoneNumber("+15551234567"),  // Your Telnyx number
    omnivoice.WithWebhookURL("https://your-server.com/webhook"),
)
```

## Voice Calls

### Making Outbound Calls

```go
call, err := provider.MakeCall(ctx, "+15559876543")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Call ID: %s\n", call.ID())
fmt.Printf("Status: %s\n", call.Status())
```

### With Options

```go
call, err := provider.MakeCall(ctx, "+15559876543",
    omnivoice.WithFrom("+15551234567"),
    omnivoice.WithTimeout(30 * time.Second),
)
```

### Handling Incoming Calls

```go
provider.OnIncomingCall(func(call omnivoice.Call) error {
    log.Printf("Incoming call from %s", call.From())

    // Answer the call
    if err := call.Answer(ctx); err != nil {
        return err
    }

    // Handle call logic...
    return nil
})

// HTTP handler for webhook
http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
    var event telnyxEvent
    json.NewDecoder(r.Body).Decode(&event)

    // Process event based on type
    switch event.EventType {
    case "call.initiated":
        handleCallInitiated(event)
    case "call.answered":
        handleCallAnswered(event)
    case "call.hangup":
        handleCallHangup(event)
    }

    w.WriteHeader(http.StatusOK)
})
```

### Call Control Commands

Telnyx provides granular call control:

```go
// Answer call
err := call.Answer(ctx)

// Hang up
err = call.Hangup(ctx)

// Speak using built-in TTS
call.Speak(ctx, "Hello, welcome to our service.", "en-US", "female")

// Play audio file
call.PlayAudio(ctx, "https://example.com/audio.mp3")

// Start transcription
call.StartTranscription(ctx, "en")

// Start media streaming
call.StartMediaStreaming(ctx, "wss://your-server.com/stream")

// Transfer call
call.Transfer(ctx, "+15559876543")

// Put on hold
call.Hold(ctx)

// Resume from hold
call.Unhold(ctx)
```

### Call Information

```go
fmt.Printf("Call ID: %s\n", call.ID())
fmt.Printf("Direction: %s\n", call.Direction())  // inbound, outbound
fmt.Printf("From: %s\n", call.From())
fmt.Printf("To: %s\n", call.To())
fmt.Printf("Status: %s\n", call.Status())
fmt.Printf("Duration: %s\n", call.Duration())
```

## Media Streaming

Connect real-time audio for voice agents:

### Start Streaming

```go
// Via Call Control API
call.StartMediaStreaming(ctx, "wss://your-server.com/stream")
```

### WebSocket Handler

```go
transport := provider.Transport()

// Listen for connections
connCh, _ := transport.Listen(ctx, "/stream")

http.HandleFunc("/stream", func(w http.ResponseWriter, r *http.Request) {
    transport.HandleWebSocket(w, r, "/stream")
})

// Handle audio streams
for conn := range connCh {
    go handleMediaStream(conn)
}

func handleMediaStream(conn transport.Connection) {
    defer conn.Close()

    buffer := make([]byte, 1024)
    for {
        n, err := conn.AudioOut().Read(buffer)
        if err != nil {
            break
        }

        // Process incoming audio
        processAudio(buffer[:n])

        // Send response audio
        conn.AudioIn().Write(responseAudio)
    }
}
```

### Audio Format

Telnyx media streaming uses:

- **Encoding**: PCMU (μ-law) or PCMA (A-law)
- **Sample Rate**: 8000 Hz
- **Channels**: Mono

## SMS Messaging

### Sending SMS

```go
// Cast to SMSProvider
sms := provider.(omnivoice.SMSProvider)

// Send SMS
msg, err := sms.SendSMS(ctx, "+15559876543", "Hello from OmniVoice!")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Message ID: %s\n", msg.ID)
fmt.Printf("Status: %s\n", msg.Status)
```

### From Specific Number

```go
msg, err := sms.SendSMSFrom(ctx,
    "+15559876543",  // To
    "+15551234567",  // From (your Telnyx number)
    "Hello!",
)
```

### Receiving SMS

```go
http.HandleFunc("/sms-webhook", func(w http.ResponseWriter, r *http.Request) {
    var event telnyxSMSEvent
    json.NewDecoder(r.Body).Decode(&event)

    if event.EventType == "message.received" {
        from := event.Payload.From.PhoneNumber
        body := event.Payload.Text

        log.Printf("SMS from %s: %s", from, body)

        // Send reply
        sms.SendSMSFrom(ctx, from, event.Payload.To[0].PhoneNumber, "Thanks!")
    }

    w.WriteHeader(http.StatusOK)
})
```

## Webhook Events

### Call Events

| Event Type | Description |
|------------|-------------|
| `call.initiated` | Call started |
| `call.answered` | Call answered |
| `call.bridged` | Call connected |
| `call.hangup` | Call ended |
| `call.machine.detection.ended` | Answering machine detected |
| `call.speak.started` | TTS started |
| `call.speak.ended` | TTS ended |
| `call.recording.saved` | Recording available |

### SMS Events

| Event Type | Description |
|------------|-------------|
| `message.received` | Incoming SMS |
| `message.sent` | SMS sent |
| `message.finalized` | SMS delivered/failed |

## Error Handling

```go
call, err := provider.MakeCall(ctx, to)
if err != nil {
    switch {
    case strings.Contains(err.Error(), "invalid_destination"):
        log.Println("Invalid phone number")
    case strings.Contains(err.Error(), "number_not_provisioned"):
        log.Println("From number not configured")
    case strings.Contains(err.Error(), "insufficient_balance"):
        log.Println("Account balance too low")
    case strings.Contains(err.Error(), "rate_limit"):
        log.Println("Rate limited, retry later")
    default:
        log.Printf("Call error: %v", err)
    }
}
```

## Telnyx vs Twilio

| Feature | Telnyx | Twilio |
|---------|--------|--------|
| Pricing | Lower | Higher |
| Call Control | More granular | TwiML-based |
| Built-in TTS | Yes | Via TwiML |
| Built-in STT | Yes | Via Gather |
| Media Streaming | Yes | Yes |
| Global Coverage | 150+ | 180+ |

## Best Practices

1. **Use HTTPS webhooks** - Required for production
2. **Validate signatures** - Verify webhook authenticity
3. **Handle all event types** - Process status updates
4. **Use E.164 format** - `+15551234567` for phone numbers
5. **Leverage built-in TTS** - Lower latency than external TTS
6. **Monitor connection profiles** - Ensure proper routing

## Pricing

| Service | Price |
|---------|-------|
| Outbound Call | $0.0100/min (US) |
| Inbound Call | $0.0080/min (US) |
| SMS Outbound | $0.0040/segment (US) |
| SMS Inbound | $0.0010/segment (US) |
| Phone Number | $1.00/month (US) |

*Prices vary by region. Check [Telnyx Pricing](https://telnyx.com/pricing) for current rates.*

## Next Steps

- [Voice Calls Guide](../guides/voice-calls.md) - Call handling patterns
- [SMS Guide](../guides/sms.md) - SMS messaging patterns
- [Voice Agents](../guides/voice-agents.md) - Build conversational agents
- [Twilio](twilio.md) - Alternative provider
