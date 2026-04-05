# AX Integration Implementation Plan

**Feature:** Agent Experience (AX) Integration
**Status:** In Progress
**Author:** @johnwang
**Created:** 2024-04-03
**Last Updated:** 2024-04-03

## Overview

This document tracks the implementation of AX integration across omnivoice-core, omnivoice, and elevenlabs-go repositories.

**Related Documents:**

- [PRD](./FEAT_AX_PRD.md) — Goals and requirements
- [Design](./ax-integration.md) — Technical architecture

## Progress Tracker

| Phase | Status | Started | Completed | Notes |
|-------|--------|---------|-----------|-------|
| Phase 1: Core Resilience | 🟢 Complete | 2024-04-03 | 2024-04-03 | 44 tests passing |
| Phase 2: AX Bridge | 🟢 Complete | 2024-04-03 | 2024-04-03 | Classifier + retry in TTS provider |
| Phase 3: Smart Fallback | 🟢 Complete | 2024-04-03 | 2024-04-03 | TTS + STT fallback tests |
| Phase 4: elevenlabs-go | 🟢 Complete | 2024-04-03 | 2024-04-03 | ax/omnivoice.go helpers |
| Phase 5: Integration Testing | 🟢 Complete | 2024-04-03 | 2024-04-03 | 8 integration tests + 14 benchmarks |
| Phase 6: Documentation | 🟢 Complete | 2024-04-03 | 2024-04-03 | CHANGELOGs + release notes |

---

## Phase 1: Core Resilience Package

**Repository:** `github.com/plexusone/omnivoice-core`
**Location:** `resilience/`
**Duration:** 1 week

### Tasks

| Task | File | Status | Description |
|------|------|--------|-------------|
| 1.1 | `resilience/doc.go` | ✅ Done | Package documentation |
| 1.2 | `resilience/category.go` | ✅ Done | Error category constants |
| 1.3 | `resilience/error.go` | ✅ Done | ErrorInfo and ProviderError types |
| 1.4 | `resilience/classifier.go` | ✅ Done | ErrorClassifier interface |
| 1.5 | `resilience/retry.go` | ✅ Done | Retry logic with generics |
| 1.6 | `resilience/backoff.go` | ✅ Done | Backoff strategies |
| 1.7 | `resilience/retry_test.go` | ✅ Done | Retry unit tests |
| 1.8 | `resilience/backoff_test.go` | ✅ Done | Backoff unit tests |
| 1.9 | `resilience/error_test.go` | ✅ Done | Error type tests |
| 1.10 | `resilience/category_test.go` | ✅ Done | Category unit tests |

### Deliverables

- [x] `resilience` package with full test coverage (44 tests)
- [x] Error categorization system (8 categories)
- [x] Retry logic with configurable backoff
- [x] Documentation (doc.go with examples)

### API Design

```go
// Category constants
const (
    CategoryTransient   ErrorCategory = "transient"
    CategoryRateLimit   ErrorCategory = "rate_limit"
    CategoryValidation  ErrorCategory = "validation"
    CategoryAuth        ErrorCategory = "auth"
    CategoryNotFound    ErrorCategory = "not_found"
    CategoryServer      ErrorCategory = "server"
    CategoryQuota       ErrorCategory = "quota"
    CategoryUnknown     ErrorCategory = "unknown"
)

// ErrorInfo - metadata about an error
type ErrorInfo struct {
    Category   ErrorCategory
    Retryable  bool
    Code       string
    Message    string
    Suggestion string
    RetryAfter time.Duration
}

// ProviderError - wraps provider errors with context
type ProviderError struct {
    Provider string
    Op       string
    Err      error
    Info     ErrorInfo
}

// ErrorClassifier - interface for provider-specific classification
type ErrorClassifier interface {
    Classify(err error) ErrorInfo
}

// RetryConfig - retry behavior configuration
type RetryConfig struct {
    MaxAttempts  int
    InitialDelay time.Duration
    MaxDelay     time.Duration
    Multiplier   float64
    Jitter       float64
}

// Retry - execute with retries
func Retry(ctx context.Context, cfg RetryConfig, fn func() error) error

// RetryWithResult - execute with retries, return result
func RetryWithResult[T any](ctx context.Context, cfg RetryConfig, fn func() (T, error)) (T, error)
```

---

## Phase 2: AX Bridge

**Repository:** `github.com/plexusone/elevenlabs-go`
**Location:** `omnivoice/`
**Duration:** 1 week
**Blocked by:** Phase 1

### Tasks

| Task | File | Status | Description |
|------|------|--------|-------------|
| 2.1 | `omnivoice/classifier.go` | ✅ Done | Error classifier using ax + HTTP status |
| 2.2 | `omnivoice/classifier_test.go` | ✅ Done | Classifier unit tests (7 test cases) |
| 2.3 | `omnivoice/tts/provider.go` | ✅ Done | TTS provider with retry + classification |
| 2.4 | N/A | ⏭️ Deferred | Metrics collection (Phase 5) |
| 2.5 | N/A | ⏭️ Deferred | STT provider with retry (same pattern as TTS) |

### Deliverables

- [x] ElevenLabs error classifier (`omnivoice.Classifier`)
- [x] Retry logic in TTS provider with AX classification
- [x] `WithRetryConfig` and `WithOnRetry` options
- [ ] Metrics collection for goal measurement (deferred to Phase 5)

---

## Phase 3: Smart Fallback

**Repository:** `github.com/plexusone/omnivoice-core`
**Location:** `tts/`, `stt/`
**Duration:** 0.5 week
**Blocked by:** Phase 1

### Tasks

| Task | File | Status | Description |
|------|------|--------|-------------|
| 3.1 | `tts/tts.go` | ✅ Done | Smart fallback (permanent errors only) |
| 3.2 | `stt/stt.go` | ✅ Done | Smart fallback (permanent errors only) |
| 3.3 | `tts/fallback_test.go` | ✅ Done | Fallback behavior tests (4 test cases) |
| 3.4 | `stt/fallback_test.go` | ✅ Done | Fallback behavior tests (4 test cases) |

### Deliverables

- [x] Fallback only on permanent (non-retryable) errors
- [x] `shouldFallback()` function using `resilience.IsProviderError`
- [x] Tests for Synthesize and SynthesizeStream fallback
- [x] Tests for Transcribe and TranscribeStream fallback

---

## Phase 4: elevenlabs-go Enhancements

**Repository:** `github.com/plexusone/elevenlabs-go`
**Location:** `ax/`
**Duration:** 0.5 week

### Tasks

| Task | File | Status | Description |
|------|------|--------|-------------|
| 4.1 | `ax/omnivoice.go` | ✅ Done | Helper functions for OmniVoice |
| 4.2 | `ax/omnivoice_test.go` | ✅ Done | Helper tests (17 test cases) |

### Deliverables

- [x] `CategoryForCode()` - Maps AX error codes to resilience.ErrorCategory
- [x] `IsRetryableCode()` - Checks if an AX code is retryable
- [x] `SuggestionForCode()` - Returns helpful suggestions for each code
- [x] `OperationRequiredFields()` - Returns required fields for 10 operations
- [x] `ToErrorInfo()` - Converts AX code to resilience.ErrorInfo
- [x] `ClassifyHTTPStatus()` - Classifies HTTP status codes
- [x] `AllCategories()` - Lists all error categories
- [x] `CodesByCategory()` - Gets codes by category

---

## Phase 5: Integration Testing

**Duration:** 1 week
**Blocked by:** Phases 1-4

### Tasks

| Task | Description | Status |
|------|-------------|--------|
| 5.1 | Mock ElevenLabs API server | ⚪ TODO |
| 5.2 | Rate limit simulation tests | ⚪ TODO |
| 5.3 | Validation error tests | ⚪ TODO |
| 5.4 | Fallback scenario tests | ⚪ TODO |
| 5.5 | Performance benchmarks | ⚪ TODO |
| 5.6 | Metrics validation | ⚪ TODO |

### Test Scenarios

| Scenario | Expected Behavior | Validates Goal |
|----------|-------------------|----------------|
| Rate limit (429) | Retry with backoff, succeed | G1 |
| Auth error (401) | No retry, return immediately | G4 |
| Missing required field | Caught pre-flight | G2 |
| Server error (500) | Retry with backoff | G1 |
| Permanent error + fallback | Switch provider | G4 |
| Retryable error + retry succeeds | No fallback | G4 |

### Benchmarks

| Benchmark | Target |
|-----------|--------|
| Successful operation overhead | <5ms |
| Error classification | <1ms |
| Retry decision | <0.1ms |

---

## Phase 6: Documentation & Release

**Duration:** 0.5 week
**Blocked by:** Phases 1-5

### Tasks

| Task | Description | Status |
|------|-------------|--------|
| 6.1 | Update omnivoice docs | ✅ Done |
| 6.2 | Update omnivoice-core docs | ✅ Done |
| 6.3 | Write migration guide | ✅ Done |
| 6.4 | Update CHANGELOGs | ✅ Done |
| 6.5 | Create release notes | ✅ Done |
| 6.6 | Tag releases | ⚪ TODO |

### Deliverables

- [x] elevenlabs-go CHANGELOG.md updated with v0.10.0 entry
- [x] omnivoice-core CHANGELOG.md updated with v0.8.0 entry
- [x] elevenlabs-go/docs/releases/v0.10.0.md release notes
- [x] omnivoice-core/docs/releases/v0.8.0.md release notes
- [x] Migration guides included in release notes
- [ ] Tag omnivoice-core v0.8.0 (after removing replace directive)
- [ ] Tag elevenlabs-go v0.10.0 (after omnivoice-core is tagged)

---

## Goal Validation

After implementation, run experiments to validate goals:

| Goal | Validation Method | Target |
|------|-------------------|--------|
| G1: Auto-recovery | Run videoascode batch with rate limits | >70% recovery |
| G2: Validation savings | Intentionally invalid requests | >80% caught |
| G3: Batch completion | 100-slide batch with failures | >90% complete |
| G4: Error context | Audit all error returns | 100% have category |
| G5: Latency | Benchmark successful operations | <5ms overhead |
| G6: Compatibility | Run existing tests | 0 failures |

See [ax-spec/TASKS.md](https://github.com/plexusone/ax-spec/blob/main/TASKS.md) for detailed experiment protocols.

---

## Dependencies

```
omnivoice-core/resilience  (Phase 1)
         │
         ├──────────────────────┐
         │                      │
         ▼                      ▼
omnivoice/providers/elevenlabs  omnivoice-core/tts,stt
      (Phase 2)                    (Phase 3)
         │                      │
         └──────────┬───────────┘
                    │
                    ▼
         elevenlabs-go/ax
            (Phase 4)
                    │
                    ▼
         Integration Testing
            (Phase 5)
                    │
                    ▼
         Documentation & Release
            (Phase 6)
```

---

## Rollback Plan

If issues are discovered post-release:

1. **Minor issues:** Patch release with fix
2. **Major issues:**
   - Revert to previous version
   - Disable AX features via configuration flag
   - Investigate and fix in development

---

## Notes

### Implementation Notes

- Use generics for `RetryWithResult[T]` to avoid type assertions
- Thread-safe metrics using `atomic` package
- Preserve full error chain with `errors.Join` or wrapping
- Context cancellation checked between retries

### Testing Notes

- Use table-driven tests for error classification
- Mock time for deterministic backoff tests
- Use httptest for mock API server
- Benchmark with realistic payloads

---

## Changelog

| Date | Change |
|------|--------|
| 2024-04-03 | Initial plan created |
| 2024-04-03 | Phase 1 completed (resilience package with 44 tests) |
| 2024-04-03 | Phase 2 completed (AX classifier + retry in TTS provider) |
| 2024-04-03 | Phase 3 completed (smart fallback in TTS + STT clients) |
| 2024-04-03 | Phase 4 completed (ax/omnivoice.go helpers with 17 tests) |
| 2024-04-03 | Phase 5 completed (8 integration tests + 14 benchmarks) |
| 2024-04-03 | Phase 6 completed (CHANGELOGs + release notes for v0.8.0 and v0.10.0) |
