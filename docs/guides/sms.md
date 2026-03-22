# SMS Messaging

Send SMS messages using Twilio or Telnyx via the SMSProvider interface.

## Quick Start

```go
provider, _ := omnivoice.GetCallSystemProvider("twilio",
    omnivoice.WithAccountSID(accountSID),
    omnivoice.WithAuthToken(authToken),
    omnivoice.WithPhoneNumber("+15551234567"),
)

// Cast to SMSProvider
sms := provider.(omnivoice.SMSProvider)

msg, _ := sms.SendSMS(ctx, "+15559876543", "Hello from OmniVoice!")
fmt.Printf("Message ID: %s\n", msg.ID)
```

## Available Providers

| Provider | Registry Name | Features |
|----------|---------------|----------|
| Twilio | `"twilio"` | SMS, MMS, delivery receipts |
| Telnyx | `"telnyx"` | SMS, delivery receipts |

## Sending SMS

### Using Default Number

```go
provider, _ := omnivoice.GetCallSystemProvider("twilio",
    omnivoice.WithAccountSID(os.Getenv("TWILIO_ACCOUNT_SID")),
    omnivoice.WithAuthToken(os.Getenv("TWILIO_AUTH_TOKEN")),
    omnivoice.WithPhoneNumber("+15551234567"), // Default from number
)

smsProvider := provider.(omnivoice.SMSProvider)

// Send using default number
msg, err := smsProvider.SendSMS(ctx, "+15559876543", "Hello!")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Message sent: %s\n", msg.ID)
```

### From Specific Number

```go
// Send from a specific number
msg, err := smsProvider.SendSMSFrom(ctx,
    "+15559876543",  // To
    "+15551234567",  // From
    "Hello from this specific number!",
)
```

## Message Response

The `SMSMessage` struct contains:

```go
type SMSMessage struct {
    ID     string // Provider message ID
    To     string // Recipient number
    From   string // Sender number
    Body   string // Message content
    Status string // Message status
}
```

### Checking Message Status

```go
msg, _ := smsProvider.SendSMS(ctx, to, body)

fmt.Printf("ID: %s\n", msg.ID)
fmt.Printf("Status: %s\n", msg.Status)
// Status values: queued, sent, delivered, failed, etc.
```

## Provider Configuration

### Twilio

```go
provider, _ := omnivoice.GetCallSystemProvider("twilio",
    omnivoice.WithAccountSID(os.Getenv("TWILIO_ACCOUNT_SID")),
    omnivoice.WithAuthToken(os.Getenv("TWILIO_AUTH_TOKEN")),
    omnivoice.WithPhoneNumber("+15551234567"),
)
```

### Telnyx

```go
provider, _ := omnivoice.GetCallSystemProvider("telnyx",
    omnivoice.WithAPIKey(os.Getenv("TELNYX_API_KEY")),
    omnivoice.WithPhoneNumber("+15551234567"),
)
```

## Sending to Multiple Recipients

```go
recipients := []string{
    "+15559876543",
    "+15559876544",
    "+15559876545",
}

for _, to := range recipients {
    msg, err := smsProvider.SendSMS(ctx, to, "Bulk message!")
    if err != nil {
        log.Printf("Failed to send to %s: %v", to, err)
        continue
    }
    log.Printf("Sent to %s: %s", to, msg.ID)
}
```

## Message Templates

```go
func sendWelcome(ctx context.Context, sms omnivoice.SMSProvider, to, name string) error {
    body := fmt.Sprintf("Welcome %s! Thanks for signing up.", name)
    _, err := sms.SendSMS(ctx, to, body)
    return err
}

func sendVerificationCode(ctx context.Context, sms omnivoice.SMSProvider, to, code string) error {
    body := fmt.Sprintf("Your verification code is: %s", code)
    _, err := sms.SendSMS(ctx, to, body)
    return err
}

func sendOrderConfirmation(ctx context.Context, sms omnivoice.SMSProvider, to, orderID string) error {
    body := fmt.Sprintf("Order %s confirmed! Track at: https://example.com/orders/%s", orderID, orderID)
    _, err := sms.SendSMS(ctx, to, body)
    return err
}
```

## Error Handling

```go
msg, err := smsProvider.SendSMS(ctx, to, body)
if err != nil {
    switch {
    case strings.Contains(err.Error(), "invalid_phone_number"):
        log.Println("Invalid phone number format")
    case strings.Contains(err.Error(), "unsubscribed"):
        log.Println("Recipient has opted out")
    case strings.Contains(err.Error(), "insufficient_funds"):
        log.Println("Account balance too low")
    case strings.Contains(err.Error(), "rate_limit"):
        log.Println("Rate limit exceeded, retry later")
    default:
        log.Printf("SMS error: %v", err)
    }
    return
}
```

## Phone Number Formatting

OmniVoice expects E.164 format:

```go
// Correct format
"+15551234567"  // US number
"+447911123456" // UK number
"+81312345678"  // Japan number

// Incorrect formats (will likely fail)
"5551234567"    // Missing country code
"(555) 123-4567" // Formatted
"555-123-4567"  // Formatted
```

### Normalizing Numbers

```go
import "regexp"

func normalizePhone(phone string) string {
    // Remove all non-digit characters except leading +
    re := regexp.MustCompile(`[^\d+]`)
    normalized := re.ReplaceAllString(phone, "")

    // Add + if missing
    if !strings.HasPrefix(normalized, "+") {
        // Assume US if 10 digits
        if len(normalized) == 10 {
            normalized = "+1" + normalized
        }
    }

    return normalized
}
```

## SMS Character Limits

| Encoding | Characters per Segment |
|----------|----------------------|
| GSM-7 | 160 |
| Unicode (UCS-2) | 70 |

Messages exceeding these limits are split into multiple segments.

```go
// Check message length
body := "Your message here..."
if len(body) > 160 {
    log.Printf("Warning: Message will be split into %d segments", (len(body)/153)+1)
}
```

## Best Practices

1. **Use E.164 format** - Always include country code with +
2. **Handle opt-outs** - Respect unsubscribe requests
3. **Rate limiting** - Don't send too fast (provider limits)
4. **Message length** - Keep under 160 chars when possible
5. **Include opt-out** - "Reply STOP to unsubscribe"
6. **Validate numbers** - Check format before sending

## Multi-Provider Failover

```go
// Create providers
twilioProvider, _ := omnivoice.GetCallSystemProvider("twilio", ...)
telnyxProvider, _ := omnivoice.GetCallSystemProvider("telnyx", ...)

// Try primary, fall back to secondary
func sendSMSWithFallback(ctx context.Context, to, body string) (*omnivoice.SMSMessage, error) {
    providers := []omnivoice.SMSProvider{
        twilioProvider.(omnivoice.SMSProvider),
        telnyxProvider.(omnivoice.SMSProvider),
    }

    var lastErr error
    for _, p := range providers {
        msg, err := p.SendSMS(ctx, to, body)
        if err == nil {
            return msg, nil
        }
        lastErr = err
        log.Printf("Provider failed: %v, trying next...", err)
    }

    return nil, fmt.Errorf("all providers failed: %w", lastErr)
}
```

## Next Steps

- [Voice Calls](voice-calls.md) - Make and receive phone calls
- [Voice Agents](voice-agents.md) - Build conversational agents
- [Provider Details](../providers/twilio.md) - Twilio-specific features
