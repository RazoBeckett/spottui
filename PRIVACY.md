# Privacy Policy

**Last Updated: January 21, 2026**

## Data Collection

SpotTUI only collects the following data:

### What We Collect
- **OAuth Access Tokens**: Stored locally at `~/.config/spottui/token.json` to authenticate with Spotify's API
- **Configuration Settings**: Your preferences stored at `~/.config/spottui/config.json`

### What We Do NOT Collect
- ❌ No telemetry or analytics
- ❌ No tracking or monitoring
- ❌ No personal information beyond what's required for Spotify authentication
- ❌ No data is sent to any third-party servers other than Spotify's official API

## Data Storage

All data collected by SpotTUI is stored **locally on your device only**:

| Data | Location | Purpose |
|------|----------|---------|
| Access tokens | `~/.config/spottui/token.json` | Authenticate with Spotify API |
| Refresh tokens | `~/.config/spottui/token.json` | Refresh access tokens |
| Configuration | `~/.config/spottui/config.json` | Store app settings |
| Cached API responses | In-memory (temporary) | Improve performance (TTL-based) |

**Security**: Token files are created with restricted permissions (0600 on Unix systems).

## Third-Party Services

SpotTUI uses the following third-party services:

| Service | Provider | Purpose |
|---------|----------|---------|
| Spotify Web API | Spotify AB | Music streaming, playback control, and metadata |

**Data Shared**: Only OAuth tokens are shared with Spotify AB's API endpoints. No data is shared with any other third parties.

## User Rights

You have the following rights regarding your data:

### Revocation of Access
- You can revoke SpotTUI's access to your Spotify account at any time through Spotify's account settings at [https://www.spotify.com/account](https://www.spotify.com/account)
- Navigate to "Apps" and remove SpotTUI from authorized applications

### Data Deletion
To delete all SpotTUI data from your device:

```bash
rm -rf ~/.config/spottui/
```

This will remove:
- Stored tokens
- Configuration files
- All cached data

### Token Management
- Tokens are automatically refreshed using the OAuth 2.0 PKCE flow
- Tokens persist between sessions for convenience
- You can manually delete `~/.config/spottui/token.json` to force re-authentication

## Cookies

SpotTUI does not use cookies.

## Children's Privacy

SpotTUI is not directed to children under the age of 13. We do not knowingly collect personal information from children.

## Changes to This Policy

We may update this privacy policy from time to time. We will notify users of any material changes by updating the "Last Updated" date.

## Contact

For questions about this privacy policy or your data:
- Open an issue on GitHub: [https://github.com/razobeckett/spottui/issues](https://github.com/razobeckett/spottui/issues)

---

## Spotify Compliance

This privacy policy is provided in compliance with Section V.12 of the [Spotify Developer Terms](https://developer.spotify.com/terms).
