# AX Integration PRD

**Feature:** Agent Experience (AX) Integration
**Status:** Draft
**Author:** @johnwang
**Created:** 2024-04-03

## Problem Statement

AI agents using OmniVoice for voice operations experience poor reliability due to:

1. **No automatic recovery** — Transient errors (rate limits, timeouts) cause immediate failure
2. **No intelligent retry** — All errors treated equally, no backoff strategy
3. **Wasted API calls** — Invalid requests sent to API instead of validated locally
4. **Poor error context** — Generic error messages without recovery guidance
5. **Unnecessary fallbacks** — Provider switches on retryable errors

### User Pain Points

| User | Pain Point | Frequency |
|------|------------|-----------|
| videoascode users | Batch TTS fails mid-way, requires manual re-run | Every large batch |
| omniagent users | Voice workflows fail on transient errors | 10-20% of sessions |
| Developers | Debugging errors without context | Every failure |

### Business Impact

- **Wasted compute** — Failed operations consume quota without delivering value
- **User frustration** — Manual intervention required for recoverable errors
- **Support burden** — Users report "random failures" that are actually rate limits

## Goals

### Primary Goals

| ID | Goal | Success Metric | Target |
|----|------|----------------|--------|
| G1 | Auto-recover from transient errors | Recovery rate | >70% of retryable errors recovered |
| G2 | Reduce wasted API calls | Validation catch rate | >80% of invalid requests caught pre-flight |
| G3 | Improve batch completion | Completion rate | >90% (up from ~60% with failures) |
| G4 | Provide actionable error context | Error quality | 100% of errors have category + suggestion |

### Secondary Goals

| ID | Goal | Success Metric | Target |
|----|------|----------------|--------|
| G5 | Minimize latency overhead | Added latency | <5ms per operation |
| G6 | Maintain backward compatibility | Breaking changes | Zero |
| G7 | Enable provider-agnostic patterns | Code reuse | Resilience works for all providers |

## Non-Goals

- Circuit breaker implementation (future work)
- Quota management across providers (future work)
- AX integration for non-ElevenLabs providers (future work)
- Modifying elevenlabs-go core client (out of scope)

## User Stories

### US1: Batch TTS Resilience

**As a** videoascode user
**I want** batch TTS to automatically retry rate-limited requests
**So that** I don't have to manually re-run failed batches

**Acceptance Criteria:**

- [ ] Rate limit errors (429) trigger exponential backoff retry
- [ ] Batch continues processing after permanent errors
- [ ] Final report shows which items succeeded/failed
- [ ] Total completion rate >90% for typical workloads

### US2: Agent Workflow Reliability

**As an** omniagent user
**I want** voice operations to recover from transient failures
**So that** my agent workflows complete without manual intervention

**Acceptance Criteria:**

- [ ] Transient errors retry automatically with backoff
- [ ] Permanent errors return immediately with clear context
- [ ] Error category available for programmatic handling
- [ ] Recovery suggestions included in error messages

### US3: Pre-flight Validation

**As a** developer
**I want** invalid requests caught before API calls
**So that** I don't waste quota on requests that will fail

**Acceptance Criteria:**

- [ ] Missing required fields caught before API call
- [ ] Clear error message indicates which fields are missing
- [ ] No API quota consumed for validation failures
- [ ] Works for all ElevenLabs operations

### US4: Error Diagnostics

**As a** developer
**I want** errors to include category and recovery guidance
**So that** I can implement appropriate handling logic

**Acceptance Criteria:**

- [ ] All errors include category (rate_limit, auth, validation, etc.)
- [ ] All errors indicate if retryable
- [ ] Errors include recovery suggestion when applicable
- [ ] Error chain preserves original error for debugging

## Requirements

### Functional Requirements

#### FR1: Error Classification

- FR1.1: Classify errors into categories: transient, rate_limit, validation, auth, not_found, server, quota
- FR1.2: Determine retryability based on error category
- FR1.3: Extract error codes from provider responses
- FR1.4: Map provider errors to OmniVoice error types

#### FR2: Retry Logic

- FR2.1: Retry retryable errors with exponential backoff
- FR2.2: Support configurable max attempts (default: 3)
- FR2.3: Support configurable initial delay (default: 1s)
- FR2.4: Support configurable max delay (default: 30s)
- FR2.5: Add jitter to prevent thundering herd
- FR2.6: Respect Retry-After headers when present

#### FR3: Pre-flight Validation

- FR3.1: Validate required fields before API calls
- FR3.2: Return validation errors without making API call
- FR3.3: Include missing field names in error message

#### FR4: Smart Fallback

- FR4.1: Only fallback to alternate provider on permanent errors
- FR4.2: Retry within provider before falling back
- FR4.3: Preserve error context through fallback chain

#### FR5: Metrics Collection

- FR5.1: Track total calls, successes, failures
- FR5.2: Track retry attempts and retry successes
- FR5.3: Track validation catches
- FR5.4: Track fallbacks avoided (retryable errors that succeeded)

### Non-Functional Requirements

#### NFR1: Performance

- NFR1.1: <5ms latency overhead per operation (excluding retries)
- NFR1.2: No blocking operations on success path
- NFR1.3: Efficient memory usage (no large allocations)

#### NFR2: Compatibility

- NFR2.1: No breaking changes to public API
- NFR2.2: Existing code works without modification
- NFR2.3: New features are opt-in where possible

#### NFR3: Reliability

- NFR3.1: Thread-safe implementation
- NFR3.2: Context cancellation respected
- NFR3.3: Graceful degradation on unknown errors

#### NFR4: Observability

- NFR4.1: Errors include structured metadata
- NFR4.2: Metrics accessible for monitoring
- NFR4.3: Integration with existing observability hooks

## Success Criteria

### Quantitative

| Metric | Baseline | Target | Measurement |
|--------|----------|--------|-------------|
| Error recovery rate | 0% | >70% | RetrySuccesses / TotalRetryable |
| Validation catch rate | 0% | >80% | ValidationCatches / (Catches + API400s) |
| Batch completion rate | ~60% | >90% | CompletedItems / TotalItems |
| Retry success rate | N/A | >60% | RetrySuccesses / TotalRetries |

### Qualitative

- [ ] Developers can handle errors programmatically by category
- [ ] Error messages are actionable (include recovery suggestion)
- [ ] No increase in support tickets for "random failures"

## Out of Scope

- Implementing AX for Deepgram, OpenAI, or other providers
- Circuit breaker patterns
- Cross-provider quota management
- Distributed rate limiting
- Automatic provider selection based on availability

## Dependencies

| Dependency | Type | Status |
|------------|------|--------|
| elevenlabs-go/ax | Code | Exists |
| omnivoice-core | Code | Exists |
| omnivoice | Code | Exists |

## Risks

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| Breaking API changes | High | Low | Use new types, don't modify signatures |
| Performance regression | Medium | Low | Benchmark critical paths |
| Provider API changes | Medium | Medium | Isolate provider logic |

## Timeline

| Phase | Duration | Deliverable |
|-------|----------|-------------|
| Design | Complete | PRD, Design doc |
| Phase 1: Core Resilience | 1 week | resilience package |
| Phase 2: AX Bridge | 1 week | ElevenLabs integration |
| Phase 3: Smart Fallback | 0.5 week | Client updates |
| Phase 4: Testing | 1 week | Integration tests |
| Phase 5: Documentation | 0.5 week | Docs, release |

## References

- [AX Spec](https://github.com/plexusone/ax-spec)
- [DIRECT Principles](https://github.com/grokify/direct-principles)
- [Design Document](./ax-integration.md)
- [Implementation Plan](./FEAT_AX_PLAN.md)
