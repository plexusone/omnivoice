# Changelog

See [CHANGELOG.md](https://github.com/plexusone/omnivoice/blob/main/CHANGELOG.md) for the full changelog.

## Releases

| Version | Date | Highlights |
|---------|------|------------|
| [v0.14.0](releases/v0.14.0.md) | 2026-10-28 | OpenAI, ElevenLabs, and Deepgram TTS return 8 kHz mu-law or A-law phone-call audio from the same config, with the correct sample rate |
| [v0.13.0](releases/v0.13.0.md) | 2026-07-05 | Deepgram Voice Agent realtime provider support |
| v0.12.1 | 2026-06-29 | Patch release for dependency updates |
| [v0.12.0](releases/v0.12.0.md) | 2026-06-27 | Local TTS/STT providers for Apple Silicon (F5-TTS MLX, Whisper MLX) |
| [v0.11.0](releases/v0.11.0.md) | 2026-06-15 | Realtime factory API for getting factories by name |
| [v0.10.0](releases/v0.10.0.md) | 2026-06-15 | Unified registry via delegation to omnivoice-core |
| v0.9.1 | 2026-06-15 | Dependency updates |
| [v0.9.0](releases/v0.9.0.md) | 2026-06-15 | RealtimeProvider registry for native voice-to-voice APIs |
| [v0.8.0](releases/v0.8.0.md) | 2026-05-02 | CLI tool for speech-to-text transcription |
| [v0.7.1](releases/v0.7.1.md) | 2026-04-05 | AX-enabled dependencies with smart fallback and error classification |
| [v0.7.0](releases/v0.7.0.md) | 2026-03-22 | CallSystem provider registry with Twilio and Telnyx support |
| [v0.6.0](releases/v0.6.0.md) | 2026-03-07 | Re-export callsystem types from omnivoice-core for unified imports |
| [v0.5.0](releases/v0.5.0.md) | 2026-03-01 | Batteries-included voice pipeline framework for Go with unified STT and TTS interfaces |

## Latest Release: v0.14.0

- **Phone-call speech from any TTS provider**: OpenAI, ElevenLabs, and Deepgram return 8 kHz mu-law or A-law audio from the same config
- **Raw call audio**: Twilio and Telnyx voice gateways gain an audio mode
- **Keypad digits**: answered Twilio and Telnyx outbound calls can send DTMF tones
- **Docs**: user guides moved under `docs/guides/`

See [v0.14.0 Release Notes](releases/v0.14.0.md) for details.

## Versioning

OmniVoice follows [Semantic Versioning](https://semver.org/):

- **MAJOR**: Breaking API changes
- **MINOR**: New features, backward compatible
- **PATCH**: Bug fixes, backward compatible
