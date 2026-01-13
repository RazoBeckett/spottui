# SpotTUI Bugs & Issues

## Open Issues

_No known bugs currently_

---

## Recently Fixed (Jan 14, 2026)

- [x] ~~**Panic on startup due to uninitialized lists**~~ - @razobeckett
  - **Symptom**: `runtime error: invalid memory address or nil pointer dereference` when calling `SetSize` on list models
  - **Cause**: Window resize handler called `SetSize` on all lists, including those not yet initialized
  - **Fix**: Only resize lists that have items (`len(l.Items()) > 0`)
  - **Commit**: `7858fe4`

---

## Known Limitations

- Requires Spotify Premium for playback control
- No offline mode support
- Album art not displayed (terminal graphics not implemented)
- Queue management not available
