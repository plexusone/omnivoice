# Changelog

See [CHANGELOG.md](https://github.com/plexusone/omnivoice/blob/main/CHANGELOG.md) for the full changelog.

## Releases

| Version | Date | Highlights |
|---------|------|------------|
| [v0.7.0](releases/v0.7.0.md) | 2026-03-22 | CallSystem provider registry, Twilio & Telnyx support |
| [v0.6.0](releases/v0.6.0.md) | 2026-03-07 | CallSystem type re-exports, call options |
| [v0.5.0](releases/v0.5.0.md) | 2026-03-01 | Initial release, provider registry pattern |

## Latest Release: v0.7.0

- **feat(callsystem)**: Export `CallSystemClient` type and constructor
- **feat(registry)**: Add CallSystem provider registry with `GetCallSystemProvider`, `ListCallSystemProviders`
- **feat(providers)**: Register Twilio and Telnyx CallSystem providers
- **chore(deps)**: Update all provider dependencies

See [v0.7.0 Release Notes](releases/v0.7.0.md) for details.

## Versioning

OmniVoice follows [Semantic Versioning](https://semver.org/):

- **MAJOR**: Breaking API changes
- **MINOR**: New features, backward compatible
- **PATCH**: Bug fixes, backward compatible
