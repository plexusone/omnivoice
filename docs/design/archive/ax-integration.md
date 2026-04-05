# AX Integration Design

Adding Agent Experience (AX) awareness to OmniVoice for intelligent error handling, retry logic, and self-healing capabilities.

## Status

**Status:** Implemented
**Author:** @johnwang
**Created:** 2024-04-03
**Last Updated:** 2024-04-03

### Implementation Summary

All phases completed:

| Phase | Status | Version |
|-------|--------|---------|
| Phase 1: Core Resilience | ✅ Complete | omnivoice-core v0.8.0 |
| Phase 2: AX Bridge | ✅ Complete | elevenlabs-go v0.10.0 |
| Phase 3: Smart Fallback | ✅ Complete | omnivoice-core v0.8.0 |
| Phase 4: elevenlabs-go Enhancements | ✅ Complete | elevenlabs-go v0.10.0 |
| Phase 5: Integration Testing | ✅ Complete | elevenlabs-go v0.10.0 |
| Phase 6: Documentation | ✅ Complete | Both repos |

See [FEAT_AX_PLAN.md](./FEAT_AX_PLAN.md) for detailed implementation tracking.

## Goals

### Primary Goals

| ID | Goal | Metric | Target | Measurement |
|----|------|--------|--------|-------------|
| G1 | **Auto-recover from transient errors** | Error recovery rate | >70% of retryable errors recovered | Count recovered vs total retryable errors |
| G2 | **Reduce wasted API calls** | Validation savings | >80% of validation errors caught pre-flight | Count pre-flight catches vs API 400 errors |
| G3 | **Intelligent retry with backoff** | Retry success rate | >60% of retried operations succeed | Count successful retries vs total retries |
| G4 | **Smart fallback decisions** | Fallback appropriateness | Only fallback on permanent errors | Count unnecessary fallbacks (retryable errors that triggered fallback) |
| G5 | **Preserve error context** | Error diagnostics | 100% of errors have category + suggestion | Audit error returns for metadata |

### Secondary Goals

| ID | Goal | Metric | Target |
|----|------|--------|--------|
| G6 | **Minimize latency overhead** | Added latency | <5ms per operation (excluding retries) |
| G7 | **Backward compatible** | Breaking changes | Zero breaking changes to public API |
| G8 | **Provider agnostic patterns** | Code reuse | Resilience logic works for all providers |

## Non-Goals

- Implementing AX for non-ElevenLabs providers (future work)
- Modifying elevenlabs-go's core API client
- Adding circuit breaker patterns (future work)
- Quota management across providers (future work)

## Background

### Current Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Application                          │
│              (omniagent, videoascode)                   │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│                      OmniVoice                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────┐ │
│  │ TTS Client  │  │ STT Client  │  │ CallSystem      │ │
│  │ (fallback)  │  │ (fallback)  │  │ Client          │ │
│  └──────┬──────┘  └──────┬──────┘  └────────┬────────┘ │
│         │                │                   │          │
│  ┌──────┴──────────────────┴─────────────────┴────────┐ │
│  │              Provider Registry                     │ │
│  └────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────┘
                          │
          ┌───────────────┼───────────────┐
          ▼               ▼               ▼
   ┌────────────┐  ┌────────────┐  ┌────────────┐
   │ ElevenLabs │  │  Deepgram  │  │   OpenAI   │
   │  Provider  │  │  Provider  │  │  Provider  │
   └─────┬──────┘  └────────────┘  └────────────┘
         │
         ▼
   ┌────────────┐
   │elevenlabs- │
   │    go      │
   │  ┌──────┐  │
   │  │  ax/ │  │  ← AX metadata exists here
   │  └──────┘  │
   └────────────┘
```

### Current Behavior

1. **Error handling:** Providers return errors, clients attempt fallback on ANY error
2. **No retry logic:** Single attempt per provider before fallback
3. **No error categorization:** All errors treated equally
4. **No pre-flight validation:** Requests sent without checking required fields

### Problems

| Problem | Impact | Example |
|---------|--------|---------|
| Rate limits trigger fallback | Unnecessary provider switch when retry would work | 429 → fallback to Deepgram instead of waiting 1s |
| No backoff on transient errors | Repeated failures, potential ban | Immediate retry hammers rate-limited API |
| Generic error messages | Poor debugging, no recovery hints | "tts failed" vs "rate limited, retry in 60s" |
| Wasted API calls | Cost and latency | Missing voice_id caught at API, not pre-flight |

## Proposed Design

### Architecture Changes

```
┌─────────────────────────────────────────────────────────┐
│                    Application                          │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│                      OmniVoice                          │
│  ┌─────────────────────────────────────────────────┐   │
│  │              Resilience Layer (NEW)              │   │
│  │  ┌──────────┐ ┌──────────┐ ┌────────────────┐   │   │
│  │  │  Retry   │ │ Backoff  │ │  Error         │   │   │
│  │  │  Policy  │ │ Strategy │ │  Categorizer   │   │   │
│  │  └──────────┘ └──────────┘ └────────────────┘   │   │
│  └─────────────────────────────────────────────────┘   │
│                          │                              │
│  ┌─────────────┐  ┌──────┴──────┐  ┌─────────────────┐ │
│  │ TTS Client  │  │ STT Client  │  │ CallSystem      │ │
│  │ (smart      │  │ (smart      │  │ Client          │ │
│  │  fallback)  │  │  fallback)  │  │                 │ │
│  └──────┬──────┘  └──────┬──────┘  └────────┬────────┘ │
└─────────┼────────────────┼──────────────────┼──────────┘
          │                │                   │
          ▼                ▼                   ▼
   ┌────────────┐  ┌────────────┐  ┌────────────┐
   │ ElevenLabs │  │  Deepgram  │  │   OpenAI   │
   │  Provider  │  │  Provider  │  │  Provider  │
   │ (AX-aware) │  │            │  │            │
   └─────┬──────┘  └────────────┘  └────────────┘
         │
         ▼
   ┌────────────┐
   │elevenlabs- │
   │    go/ax   │ ← Error codes, retry policies, validation
   └────────────┘
```

### Component Design

#### 1. Resilience Package (omnivoice-core)

New package: `github.com/plexusone/omnivoice-core/resilience`

```go
// resilience/error.go

// ErrorCategory classifies errors for handling decisions
type ErrorCategory string

const (
    CategoryTransient   ErrorCategory = "transient"    // Retry with backoff
    CategoryRateLimit   ErrorCategory = "rate_limit"   // Retry with longer backoff
    CategoryValidation  ErrorCategory = "validation"   // Don't retry, fix input
    CategoryAuth        ErrorCategory = "auth"         // Don't retry, re-authenticate
    CategoryNotFound    ErrorCategory = "not_found"    // Don't retry, resource missing
    CategoryServer      ErrorCategory = "server"       // Retry with backoff
    CategoryQuota       ErrorCategory = "quota"        // Don't retry, limit exceeded
    CategoryUnknown     ErrorCategory = "unknown"      // Use default behavior
)

// ErrorInfo provides actionable metadata about an error
type ErrorInfo struct {
    Category   ErrorCategory
    Retryable  bool
    Code       string  // Machine-readable code (e.g., "RATE_LIMITED")
    Message    string  // Human-readable message
    Suggestion string  // Recovery suggestion
    RetryAfter time.Duration // Hint for backoff (0 = use default)
}

// ErrorClassifier categorizes errors from any source
type ErrorClassifier interface {
    Classify(err error) ErrorInfo
}

// ProviderError wraps provider errors with AX metadata
type ProviderError struct {
    Provider string
    Op       string  // Operation name
    Err      error   // Original error
    Info     ErrorInfo
}

func (e *ProviderError) Error() string
func (e *ProviderError) Unwrap() error
func (e *ProviderError) Is(target error) bool
```

```go
// resilience/retry.go

// RetryConfig controls retry behavior
type RetryConfig struct {
    MaxAttempts   int           // Max retry attempts (default: 3)
    InitialDelay  time.Duration // First retry delay (default: 1s)
    MaxDelay      time.Duration // Cap on delay (default: 30s)
    Multiplier    float64       // Backoff multiplier (default: 2.0)
    Jitter        float64       // Random jitter factor 0-1 (default: 0.1)
    RetryIf       func(err error) bool // Custom retry predicate
}

// DefaultRetryConfig returns sensible defaults
func DefaultRetryConfig() RetryConfig

// Retry executes fn with retries according to config
func Retry(ctx context.Context, config RetryConfig, fn func() error) error

// RetryWithResult executes fn and returns result on success
func RetryWithResult[T any](ctx context.Context, config RetryConfig, fn func() (T, error)) (T, error)
```

```go
// resilience/backoff.go

// BackoffStrategy computes delay between retries
type BackoffStrategy interface {
    NextDelay(attempt int, err error) time.Duration
    Reset()
}

// ExponentialBackoff implements exponential backoff with jitter
type ExponentialBackoff struct {
    Initial    time.Duration
    Max        time.Duration
    Multiplier float64
    Jitter     float64
}

func (b *ExponentialBackoff) NextDelay(attempt int, err error) time.Duration

// RateLimitBackoff respects Retry-After hints
type RateLimitBackoff struct {
    Fallback BackoffStrategy
}

func (b *RateLimitBackoff) NextDelay(attempt int, err error) time.Duration
```

#### 2. AX Bridge (omnivoice)

New file: `providers/elevenlabs/ax_bridge.go`

```go
// ax_bridge.go - Bridges elevenlabs-go/ax to omnivoice resilience

import (
    "github.com/plexusone/elevenlabs-go/ax"
    "github.com/plexusone/omnivoice-core/resilience"
)

// ElevenLabsClassifier implements ErrorClassifier using AX metadata
type ElevenLabsClassifier struct{}

func (c *ElevenLabsClassifier) Classify(err error) resilience.ErrorInfo {
    code, ok := elevenlabs.GetAXErrorCode(err)
    if !ok {
        return resilience.ErrorInfo{Category: resilience.CategoryUnknown}
    }

    axInfo := ax.GetErrorInfo(code)
    return resilience.ErrorInfo{
        Category:   mapCategory(axInfo.Category),
        Retryable:  axInfo.Retryable,
        Code:       code,
        Message:    axInfo.Description,
        Suggestion: axInfo.Suggestion,
    }
}

// mapCategory maps AX categories to resilience categories
func mapCategory(axCat string) resilience.ErrorCategory {
    switch axCat {
    case "rate_limit":
        return resilience.CategoryRateLimit
    case "auth", "authentication":
        return resilience.CategoryAuth
    case "validation":
        return resilience.CategoryValidation
    case "not_found":
        return resilience.CategoryNotFound
    case "server", "internal":
        return resilience.CategoryServer
    case "quota":
        return resilience.CategoryQuota
    default:
        return resilience.CategoryUnknown
    }
}

// ValidateRequest checks required fields before API call
func ValidateRequest(operationID string, fields map[string]bool) error {
    if msg := ax.ValidateFields(operationID, fields); msg != "" {
        return &resilience.ProviderError{
            Provider: "elevenlabs",
            Op:       operationID,
            Info: resilience.ErrorInfo{
                Category:   resilience.CategoryValidation,
                Retryable:  false,
                Code:       "VALIDATION_FAILED",
                Message:    msg,
                Suggestion: "Provide all required fields before calling",
            },
        }
    }
    return nil
}
```

#### 3. AX-Aware ElevenLabs Provider (omnivoice)

Update: `providers/elevenlabs/tts.go`

```go
// tts.go - ElevenLabs TTS provider with AX awareness

type Provider struct {
    client     *elevenlabs.Client
    classifier *ElevenLabsClassifier
    retry      resilience.RetryConfig
    metrics    *Metrics  // For goal measurement
}

type Metrics struct {
    TotalCalls          int64
    SuccessfulCalls     int64
    RetriedCalls        int64
    RetrySuccesses      int64
    ValidationCatches   int64
    FallbacksAvoided    int64  // Retryable errors that succeeded on retry
}

func (p *Provider) Synthesize(ctx context.Context, text string, config tts.Config) (*tts.Result, error) {
    // G2: Pre-flight validation
    fields := map[string]bool{
        "text": text != "",
        // Add other required fields from config
    }
    if err := ValidateRequest("text_to_speech", fields); err != nil {
        atomic.AddInt64(&p.metrics.ValidationCatches, 1)
        return nil, err
    }

    atomic.AddInt64(&p.metrics.TotalCalls, 1)

    // G1, G3: Retry with backoff for transient errors
    var result *tts.Result
    var lastErr error

    attempt := 0
    backoff := &resilience.ExponentialBackoff{
        Initial:    p.retry.InitialDelay,
        Max:        p.retry.MaxDelay,
        Multiplier: p.retry.Multiplier,
        Jitter:     p.retry.Jitter,
    }

    for attempt < p.retry.MaxAttempts {
        result, lastErr = p.doSynthesize(ctx, text, config)
        if lastErr == nil {
            atomic.AddInt64(&p.metrics.SuccessfulCalls, 1)
            if attempt > 0 {
                atomic.AddInt64(&p.metrics.RetrySuccesses, 1)
                atomic.AddInt64(&p.metrics.FallbacksAvoided, 1)
            }
            return result, nil
        }

        // G5: Classify error with full context
        info := p.classifier.Classify(lastErr)

        // G4: Only retry if error is retryable
        if !info.Retryable {
            return nil, p.wrapError("Synthesize", lastErr, info)
        }

        attempt++
        if attempt >= p.retry.MaxAttempts {
            break
        }

        atomic.AddInt64(&p.metrics.RetriedCalls, 1)

        // Respect Retry-After if present, otherwise use backoff
        delay := backoff.NextDelay(attempt, lastErr)
        if info.RetryAfter > 0 {
            delay = info.RetryAfter
        }

        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        case <-time.After(delay):
            // Continue to next attempt
        }
    }

    // All retries exhausted
    info := p.classifier.Classify(lastErr)
    return nil, p.wrapError("Synthesize", lastErr, info)
}

func (p *Provider) wrapError(op string, err error, info resilience.ErrorInfo) error {
    return &resilience.ProviderError{
        Provider: "elevenlabs",
        Op:       op,
        Err:      err,
        Info:     info,
    }
}

func (p *Provider) doSynthesize(ctx context.Context, text string, config tts.Config) (*tts.Result, error) {
    // Actual API call to elevenlabs-go
    // ...
}

// GetMetrics returns metrics for goal measurement
func (p *Provider) GetMetrics() Metrics {
    return Metrics{
        TotalCalls:        atomic.LoadInt64(&p.metrics.TotalCalls),
        SuccessfulCalls:   atomic.LoadInt64(&p.metrics.SuccessfulCalls),
        RetriedCalls:      atomic.LoadInt64(&p.metrics.RetriedCalls),
        RetrySuccesses:    atomic.LoadInt64(&p.metrics.RetrySuccesses),
        ValidationCatches: atomic.LoadInt64(&p.metrics.ValidationCatches),
        FallbacksAvoided:  atomic.LoadInt64(&p.metrics.FallbacksAvoided),
    }
}
```

#### 4. Smart Fallback in TTS Client (omnivoice-core)

Update: `tts/tts.go`

```go
// tts.go - TTS client with smart fallback

func (c *Client) Synthesize(ctx context.Context, text string, config Config) (*Result, error) {
    // Try primary provider
    if p, ok := c.providers[c.primary]; ok {
        result, err := p.Synthesize(ctx, text, config)
        if err == nil {
            return result, nil
        }

        // NEW: Check if error is retryable - provider already handled retries
        // Only fallback on permanent errors
        if pErr, ok := err.(*resilience.ProviderError); ok {
            if pErr.Info.Retryable {
                // Provider exhausted retries but error is still retryable
                // This shouldn't happen if provider implements retry correctly
                // Log for debugging
                c.logRetryableExhausted(pErr)
            }
            // Permanent error - proceed to fallback
        }

        c.lastError = err
    }

    // Try fallbacks (only reached for permanent errors)
    for _, name := range c.fallbacks {
        if p, ok := c.providers[name]; ok {
            result, err := p.Synthesize(ctx, text, config)
            if err == nil {
                return result, nil
            }
            c.lastError = err
        }
    }

    if c.lastError != nil {
        return nil, c.lastError
    }
    return nil, ErrNoAvailableProvider
}
```

#### 5. elevenlabs-go/ax Additions

New file: `ax/omnivoice.go`

```go
// omnivoice.go - Helpers for OmniVoice integration

// CategoryForCode returns the error category for an AX error code
func CategoryForCode(code string) string {
    if info, ok := errorMetadata[code]; ok {
        return info.Category
    }
    return "unknown"
}

// IsRetryableCode returns whether the error code is retryable
func IsRetryableCode(code string) bool {
    if info, ok := errorMetadata[code]; ok {
        return info.Retryable
    }
    return false
}

// SuggestionForCode returns the recovery suggestion for an error code
func SuggestionForCode(code string) string {
    if info, ok := errorMetadata[code]; ok {
        return info.Suggestion
    }
    return ""
}

// OperationRequiredFields returns required fields for pre-flight validation
func OperationRequiredFields(operationID string) []string {
    return RequiredFieldsForOperation[operationID]
}
```

## Implementation Plan

### Phase 1: Core Resilience (omnivoice-core)

**Duration:** 1 week
**Owner:** TBD

| Task | Files | Description |
|------|-------|-------------|
| 1.1 | `resilience/error.go` | Error types and categories |
| 1.2 | `resilience/retry.go` | Retry logic with generics |
| 1.3 | `resilience/backoff.go` | Backoff strategies |
| 1.4 | `resilience/retry_test.go` | Unit tests for retry |
| 1.5 | `resilience/doc.go` | Package documentation |

**Deliverables:**
- [ ] `resilience` package with error categorization
- [ ] Retry logic with configurable backoff
- [ ] 90%+ test coverage
- [ ] Documentation

### Phase 2: AX Bridge (omnivoice)

**Duration:** 1 week
**Owner:** TBD
**Depends on:** Phase 1

| Task | Files | Description |
|------|-------|-------------|
| 2.1 | `providers/elevenlabs/ax_bridge.go` | Error classifier |
| 2.2 | `providers/elevenlabs/ax_bridge_test.go` | Bridge tests |
| 2.3 | `providers/elevenlabs/tts.go` | Update TTS provider |
| 2.4 | `providers/elevenlabs/stt.go` | Update STT provider |
| 2.5 | `providers/elevenlabs/metrics.go` | Metrics collection |

**Deliverables:**
- [ ] ElevenLabs error classifier using AX metadata
- [ ] Pre-flight validation integration
- [ ] Metrics for goal measurement
- [ ] Integration tests

### Phase 3: Smart Fallback (omnivoice-core)

**Duration:** 0.5 week
**Owner:** TBD
**Depends on:** Phase 1

| Task | Files | Description |
|------|-------|-------------|
| 3.1 | `tts/tts.go` | Smart fallback logic |
| 3.2 | `stt/stt.go` | Smart fallback logic |
| 3.3 | `tts/tts_test.go` | Fallback tests |

**Deliverables:**
- [ ] Fallback only on permanent errors
- [ ] Preserve error context through fallback chain
- [ ] Tests for fallback scenarios

### Phase 4: elevenlabs-go Enhancements

**Duration:** 0.5 week
**Owner:** TBD

| Task | Files | Description |
|------|-------|-------------|
| 4.1 | `ax/omnivoice.go` | OmniVoice helper functions |
| 4.2 | `ax/omnivoice_test.go` | Helper tests |
| 4.3 | `errors.go` | Ensure AX extraction works |

**Deliverables:**
- [ ] Helper functions for OmniVoice integration
- [ ] Verified error code extraction

### Phase 5: Integration Testing

**Duration:** 1 week
**Owner:** TBD
**Depends on:** Phases 1-4

| Task | Description |
|------|-------------|
| 5.1 | End-to-end tests with mock ElevenLabs API |
| 5.2 | Rate limit simulation tests |
| 5.3 | Fallback scenario tests |
| 5.4 | Metrics validation tests |
| 5.5 | Performance benchmarks (G6) |

**Deliverables:**
- [ ] Integration test suite
- [ ] Performance benchmarks
- [ ] Goal measurement validation

### Phase 6: Documentation & Release

**Duration:** 0.5 week
**Owner:** TBD

| Task | Description |
|------|-------------|
| 6.1 | Update omnivoice docs with resilience guide |
| 6.2 | Update CHANGELOG |
| 6.3 | Release notes |

## Success Criteria

### Goal Measurement

After implementation, run the experiments defined in `ax-spec/TASKS.md`:

| Goal | Measurement Method | Success Threshold |
|------|-------------------|-------------------|
| G1: Auto-recovery rate | `RetrySuccesses / (RetriedCalls)` | >70% |
| G2: Validation savings | `ValidationCatches / (ValidationCatches + API400Errors)` | >80% |
| G3: Retry success rate | `RetrySuccesses / RetriedCalls` | >60% |
| G4: Smart fallback | `FallbacksAvoided / TotalRetryableErrors` | >90% |
| G5: Error context | Audit: all errors have category + code | 100% |
| G6: Latency overhead | Benchmark: operation without retry | <5ms |
| G7: Backward compat | API compatibility check | 0 breaks |

### Test Scenarios

| Scenario | Expected Behavior |
|----------|-------------------|
| Rate limit (429) | Retry with backoff, succeed on retry |
| Auth error (401) | No retry, return immediately with context |
| Validation error (400) | No retry, clear error message |
| Server error (500) | Retry with backoff |
| Network timeout | Retry with backoff |
| Missing required field | Caught pre-flight, no API call |

## Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| Breaking API changes | High | Use new types, don't modify existing signatures |
| Performance regression | Medium | Benchmark critical paths, optimize hot paths |
| Complex error unwrapping | Medium | Implement proper `Is()` and `Unwrap()` |
| Provider-specific quirks | Low | Isolate provider logic in classifier |

## Future Work

- Circuit breaker pattern for repeated failures
- Quota tracking across providers
- AX integration for Deepgram, OpenAI providers
- Distributed rate limit coordination
- Adaptive retry based on historical success rates

## References

- [AX Spec](https://github.com/plexusone/ax-spec)
- [DIRECT Principles](https://github.com/grokify/direct-principles)
- [elevenlabs-go/ax](https://github.com/plexusone/elevenlabs-go/tree/main/ax)
- [OmniVoice Architecture](../index.md)
