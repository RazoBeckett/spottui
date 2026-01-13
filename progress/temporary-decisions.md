# SpotTUI Temporary Decisions

## Active Workarounds

### Token Storage Location
- **Decision**: Store OAuth tokens in `~/.config/spottui/token.json`
- **Reason**: Simple file-based storage for MVP
- **Future**: Consider system keychain integration (keyring, libsecret)
- **Date**: Project inception

### Lyrics API
- **Decision**: Using external lyrics API (lrclib or similar)
- **Reason**: Spotify doesn't provide lyrics via public API
- **Future**: Monitor for official Spotify lyrics API
- **Date**: Project inception

---

## Recently Resolved

_None yet_
