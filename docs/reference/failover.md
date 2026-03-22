# Multi-Provider Failover

Build reliable voice applications by using multiple providers with automatic failover.

## Overview

Provider failover ensures your application stays operational when:

- A provider has an outage
- API rate limits are hit
- Network connectivity issues occur
- Provider-specific errors happen

## CallSystemClient

For phone calls, use `CallSystemClient` with automatic failover:

```go
import "github.com/plexusone/omnivoice"

// Create providers
twilioProvider, _ := omnivoice.GetCallSystemProvider("twilio",
    omnivoice.WithAccountSID(os.Getenv("TWILIO_ACCOUNT_SID")),
    omnivoice.WithAuthToken(os.Getenv("TWILIO_AUTH_TOKEN")),
    omnivoice.WithPhoneNumber(os.Getenv("TWILIO_PHONE_NUMBER")),
)

telnyxProvider, _ := omnivoice.GetCallSystemProvider("telnyx",
    omnivoice.WithAPIKey(os.Getenv("TELNYX_API_KEY")),
    omnivoice.WithPhoneNumber(os.Getenv("TELNYX_PHONE_NUMBER")),
)

// Create client with failover
client := omnivoice.NewCallSystemClient(twilioProvider, telnyxProvider)
client.SetPrimary("twilio")
client.SetFallbacks("telnyx")

// MakeCall automatically fails over on error
call, err := client.MakeCall(ctx, "+15559876543")
```

## TTS Failover

Implement failover for text-to-speech:

```go
type TTSWithFallback struct {
    primary   omnivoice.TTSProvider
    fallbacks []omnivoice.TTSProvider
}

func NewTTSWithFallback(primary omnivoice.TTSProvider, fallbacks ...omnivoice.TTSProvider) *TTSWithFallback {
    return &TTSWithFallback{
        primary:   primary,
        fallbacks: fallbacks,
    }
}

func (t *TTSWithFallback) Synthesize(ctx context.Context, text string, config omnivoice.SynthesisConfig) (*omnivoice.SynthesisResult, error) {
    // Try primary
    result, err := t.primary.Synthesize(ctx, text, config)
    if err == nil {
        return result, nil
    }
    log.Printf("Primary TTS failed: %v, trying fallbacks", err)

    // Try fallbacks in order
    for i, fb := range t.fallbacks {
        result, err = fb.Synthesize(ctx, text, config)
        if err == nil {
            return result, nil
        }
        log.Printf("Fallback %d failed: %v", i, err)
    }

    return nil, fmt.Errorf("all TTS providers failed")
}

// Usage
elevenlabs, _ := omnivoice.GetTTSProvider("elevenlabs", ...)
openai, _ := omnivoice.GetTTSProvider("openai", ...)
deepgram, _ := omnivoice.GetTTSProvider("deepgram", ...)

tts := NewTTSWithFallback(elevenlabs, openai, deepgram)
result, err := tts.Synthesize(ctx, "Hello!", config)
```

## STT Failover

Implement failover for speech-to-text:

```go
type STTWithFallback struct {
    primary   omnivoice.STTProvider
    fallbacks []omnivoice.STTProvider
}

func (s *STTWithFallback) TranscribeFile(ctx context.Context, path string, config omnivoice.TranscriptionConfig) (*omnivoice.TranscriptionResult, error) {
    // Try primary
    result, err := s.primary.TranscribeFile(ctx, path, config)
    if err == nil {
        return result, nil
    }
    log.Printf("Primary STT failed: %v, trying fallbacks", err)

    // Try fallbacks
    for _, fb := range s.fallbacks {
        result, err = fb.TranscribeFile(ctx, path, config)
        if err == nil {
            return result, nil
        }
    }

    return nil, fmt.Errorf("all STT providers failed")
}

// Usage
deepgram, _ := omnivoice.GetSTTProvider("deepgram", ...)
openai, _ := omnivoice.GetSTTProvider("openai", ...)

stt := &STTWithFallback{primary: deepgram, fallbacks: []omnivoice.STTProvider{openai}}
```

## SMS Failover

```go
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
        log.Printf("SMS provider failed: %v, trying next", err)
    }

    return nil, fmt.Errorf("all SMS providers failed: %w", lastErr)
}
```

## Circuit Breaker Pattern

Avoid hammering a failing provider:

```go
type CircuitBreaker struct {
    failures    int
    threshold   int
    lastFailure time.Time
    cooldown    time.Duration
    mu          sync.Mutex
}

func (cb *CircuitBreaker) Allow() bool {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    if cb.failures >= cb.threshold {
        // Check if cooldown has passed
        if time.Since(cb.lastFailure) < cb.cooldown {
            return false  // Circuit is open
        }
        cb.failures = 0  // Reset after cooldown
    }
    return true
}

func (cb *CircuitBreaker) RecordFailure() {
    cb.mu.Lock()
    defer cb.mu.Unlock()
    cb.failures++
    cb.lastFailure = time.Now()
}

func (cb *CircuitBreaker) RecordSuccess() {
    cb.mu.Lock()
    defer cb.mu.Unlock()
    cb.failures = 0
}

// Usage
type ProviderWithCircuitBreaker struct {
    provider omnivoice.TTSProvider
    cb       *CircuitBreaker
}

func (p *ProviderWithCircuitBreaker) Synthesize(ctx context.Context, text string, config omnivoice.SynthesisConfig) (*omnivoice.SynthesisResult, error) {
    if !p.cb.Allow() {
        return nil, fmt.Errorf("circuit breaker open")
    }

    result, err := p.provider.Synthesize(ctx, text, config)
    if err != nil {
        p.cb.RecordFailure()
        return nil, err
    }

    p.cb.RecordSuccess()
    return result, nil
}
```

## Load Balancing

Distribute load across providers:

```go
type LoadBalancedTTS struct {
    providers []omnivoice.TTSProvider
    current   uint64
}

func (lb *LoadBalancedTTS) Synthesize(ctx context.Context, text string, config omnivoice.SynthesisConfig) (*omnivoice.SynthesisResult, error) {
    // Round-robin selection
    idx := atomic.AddUint64(&lb.current, 1) % uint64(len(lb.providers))
    provider := lb.providers[idx]

    return provider.Synthesize(ctx, text, config)
}
```

## Weighted Selection

Prefer certain providers based on cost or quality:

```go
type WeightedProvider struct {
    Provider omnivoice.TTSProvider
    Weight   int  // Higher = more likely to be selected
}

type WeightedTTS struct {
    providers   []WeightedProvider
    totalWeight int
}

func NewWeightedTTS(providers ...WeightedProvider) *WeightedTTS {
    total := 0
    for _, p := range providers {
        total += p.Weight
    }
    return &WeightedTTS{providers: providers, totalWeight: total}
}

func (w *WeightedTTS) Select() omnivoice.TTSProvider {
    r := rand.Intn(w.totalWeight)
    for _, p := range w.providers {
        r -= p.Weight
        if r < 0 {
            return p.Provider
        }
    }
    return w.providers[0].Provider
}

// Usage: 70% ElevenLabs, 20% OpenAI, 10% Deepgram
tts := NewWeightedTTS(
    WeightedProvider{elevenlabs, 70},
    WeightedProvider{openai, 20},
    WeightedProvider{deepgram, 10},
)
```

## Health Checks

Monitor provider health:

```go
type HealthChecker struct {
    providers map[string]omnivoice.TTSProvider
    healthy   map[string]bool
    mu        sync.RWMutex
}

func (h *HealthChecker) StartHealthChecks(ctx context.Context, interval time.Duration) {
    ticker := time.NewTicker(interval)
    go func() {
        for {
            select {
            case <-ctx.Done():
                return
            case <-ticker.C:
                h.checkAll(ctx)
            }
        }
    }()
}

func (h *HealthChecker) checkAll(ctx context.Context) {
    for name, provider := range h.providers {
        // Quick health check with minimal text
        _, err := provider.Synthesize(ctx, "test", omnivoice.SynthesisConfig{
            VoiceID: "default",
        })

        h.mu.Lock()
        h.healthy[name] = (err == nil)
        h.mu.Unlock()

        if err != nil {
            log.Printf("Provider %s unhealthy: %v", name, err)
        }
    }
}

func (h *HealthChecker) GetHealthy() []omnivoice.TTSProvider {
    h.mu.RLock()
    defer h.mu.RUnlock()

    var healthy []omnivoice.TTSProvider
    for name, provider := range h.providers {
        if h.healthy[name] {
            healthy = append(healthy, provider)
        }
    }
    return healthy
}
```

## Retry with Backoff

Retry failed requests with exponential backoff:

```go
func synthesizeWithRetry(ctx context.Context, provider omnivoice.TTSProvider, text string, config omnivoice.SynthesisConfig) (*omnivoice.SynthesisResult, error) {
    maxRetries := 3
    baseDelay := 100 * time.Millisecond

    var lastErr error
    for i := 0; i < maxRetries; i++ {
        result, err := provider.Synthesize(ctx, text, config)
        if err == nil {
            return result, nil
        }

        lastErr = err

        // Don't retry non-transient errors
        if strings.Contains(err.Error(), "invalid_api_key") {
            return nil, err
        }

        // Exponential backoff
        delay := baseDelay * time.Duration(1<<i)
        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        case <-time.After(delay):
        }
    }

    return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}
```

## Best Practices

1. **Order fallbacks by priority** - Fastest/cheapest first
2. **Use circuit breakers** - Prevent cascading failures
3. **Log provider switches** - Track failover events
4. **Monitor all providers** - Health checks even for fallbacks
5. **Test failover regularly** - Ensure backup paths work
6. **Consider cost implications** - Some fallbacks may be expensive

## Next Steps

- [Provider Registry](registry.md) - Dynamic provider selection
- [Provider Comparison](../providers/index.md) - Choose the right providers
