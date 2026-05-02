# Twilio

Twilio provides programmable voice calls and SMS messaging.

## Features

- **Voice Calls**: Inbound and outbound PSTN calls
- **SMS**: Send and receive text messages
- **Media Streams**: Real-time audio via WebSocket
- **TwiML**: Declarative call control
- **Global Coverage**: 180+ countries

## Configuration

```go
import (
    "github.com/plexusone/omnivoice"
    _ "github.com/plexusone/omnivoice/providers/all"
)

provider, err := omnivoice.GetCallSystemProvider("twilio",
    omnivoice.WithAccountSID(os.Getenv("TWILIO_ACCOUNT_SID")),
    omnivoice.WithAuthToken(os.Getenv("TWILIO_AUTH_TOKEN")),
    omnivoice.WithPhoneNumber("+15551234567"),  // Your Twilio number
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

fmt.Printf("Call SID: %s\n", call.ID())
fmt.Printf("Status: %s\n", call.Status())
```

### With Custom From Number

```go
call, err := provider.MakeCall(ctx, "+15559876543",
    omnivoice.WithFrom("+15551234567"),
    omnivoice.WithTimeout(30 * time.Second),
)
```

### Handling Incoming Calls

Set up webhook handler:

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
    r.ParseForm()
    callSid := r.Form.Get("CallSid")
    from := r.Form.Get("From")
    to := r.Form.Get("To")

    // Return TwiML response
    w.Header().Set("Content-Type", "application/xml")
    w.Write([]byte(`<?xml version="1.0"?>
        <Response>
            <Say>Welcome to our service.</Say>
            <Connect>
                <Stream url="wss://your-server.com/media-stream" />
            </Connect>
        </Response>`))
})
```

### Call Control

```go
// Answer
err := call.Answer(ctx)

// Hang up
err = call.Hangup(ctx)

// Get status
status := call.Status()  // ringing, answered, ended

// Get call info
fmt.Printf("Direction: %s\n", call.Direction())  // inbound, outbound
fmt.Printf("From: %s\n", call.From())
fmt.Printf("To: %s\n", call.To())
fmt.Printf("Duration: %s\n", call.Duration())
```

## Media Streams

Connect real-time audio for voice agents:

### TwiML Setup

```xml
<?xml version="1.0"?>
<Response>
    <Connect>
        <Stream url="wss://your-server.com/media-stream">
            <Parameter name="customParam" value="value" />
        </Stream>
    </Connect>
</Response>
```

### WebSocket Handler

```go
transport := provider.Transport()

// Listen for connections
connCh, _ := transport.Listen(ctx, "/media-stream")

http.HandleFunc("/media-stream", func(w http.ResponseWriter, r *http.Request) {
    transport.HandleWebSocket(w, r, "/media-stream")
})

// Handle audio streams
for conn := range connCh {
    go handleMediaStream(conn)
}

func handleMediaStream(conn transport.Connection) {
    defer conn.Close()

    buffer := make([]byte, 1024)
    for {
        n, err := conn.AudioOut().Read(buffer)  // Audio from caller
        if err != nil {
            break
        }

        // Process audio (send to STT)
        processAudio(buffer[:n])

        // Send response audio (from TTS)
        conn.AudioIn().Write(responseAudio)
    }
}
```

### Audio Format

Twilio Media Streams use:

- **Encoding**: mulaw (μ-law)
- **Sample Rate**: 8000 Hz
- **Channels**: Mono

```go
// Configure TTS for Twilio-compatible output
ttsConfig := omnivoice.SynthesisConfig{
    VoiceID:      "aura-asteria-en",
    OutputFormat: "mulaw_8000",  // Twilio-compatible
}
```

## SMS Messaging

### Sending SMS

```go
// Cast to SMSProvider
sms := provider.(omnivoice.SMSProvider)

// Send using default number
msg, err := sms.SendSMS(ctx, "+15559876543", "Hello from OmniVoice!")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Message SID: %s\n", msg.ID)
fmt.Printf("Status: %s\n", msg.Status)
```

### From Specific Number

```go
msg, err := sms.SendSMSFrom(ctx,
    "+15559876543",  // To
    "+15551234567",  // From
    "Hello!",
)
```

### Receiving SMS

```go
http.HandleFunc("/sms-webhook", func(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    from := r.Form.Get("From")
    body := r.Form.Get("Body")

    log.Printf("SMS from %s: %s", from, body)

    // Reply with TwiML
    w.Header().Set("Content-Type", "application/xml")
    w.Write([]byte(`<?xml version="1.0"?>
        <Response>
            <Message>Thanks for your message!</Message>
        </Response>`))
})
```

## TwiML Examples

### Greeting with Speech Input

```xml
<?xml version="1.0"?>
<Response>
    <Say voice="Polly.Joanna">Welcome. How can I help you?</Say>
    <Gather input="speech" action="/handle-speech" timeout="3">
        <Say>Please speak after the beep.</Say>
    </Gather>
</Response>
```

### Transfer Call

```xml
<?xml version="1.0"?>
<Response>
    <Say>Transferring you now.</Say>
    <Dial>+15559876543</Dial>
</Response>
```

### Play Audio

```xml
<?xml version="1.0"?>
<Response>
    <Play>https://example.com/audio.mp3</Play>
</Response>
```

## Webhook Security

Validate Twilio webhook signatures:

```go
import "github.com/twilio/twilio-go/client"

func validateWebhook(r *http.Request, authToken string) bool {
    validator := client.NewRequestValidator(authToken)

    signature := r.Header.Get("X-Twilio-Signature")
    url := "https://your-server.com" + r.URL.Path

    return validator.Validate(url, r.Form, signature)
}
```

## Error Handling

```go
call, err := provider.MakeCall(ctx, to)
if err != nil {
    switch {
    case strings.Contains(err.Error(), "invalid_phone_number"):
        log.Println("Phone number format invalid")
    case strings.Contains(err.Error(), "number_not_verified"):
        log.Println("Trial account: verify number first")
    case strings.Contains(err.Error(), "insufficient_funds"):
        log.Println("Account balance too low")
    case strings.Contains(err.Error(), "unverified_caller"):
        log.Println("From number not verified")
    default:
        log.Printf("Call error: %v", err)
    }
}
```

## Best Practices

1. **Use HTTPS webhooks** - Twilio requires secure endpoints
2. **Validate signatures** - Verify webhook authenticity
3. **Handle status callbacks** - Track call progress
4. **Use E.164 format** - `+15551234567` for phone numbers
5. **Set timeout on calls** - Avoid billing for unanswered calls
6. **Monitor usage** - Track minutes and SMS counts

## Pricing

| Service | Price |
|---------|-------|
| Outbound Call | $0.0140/min (US) |
| Inbound Call | $0.0085/min (US) |
| SMS Outbound | $0.0079/segment (US) |
| SMS Inbound | $0.0079/segment (US) |
| Phone Number | $1.15/month (US) |

*Prices vary by region. Check [Twilio Pricing](https://www.twilio.com/pricing) for current rates.*

## Next Steps

- [Voice Calls Guide](../guides/voice-calls.md) - Call handling patterns
- [SMS Guide](../guides/sms.md) - SMS messaging patterns
- [Voice Agents](../guides/voice-agents.md) - Build conversational agents
- [Telnyx](telnyx.md) - Alternative provider
