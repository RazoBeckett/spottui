# SpotTUI Deployment Log

## Releases

_No releases yet - project in development_

---

## Build Notes

### Local Development
```bash
# Build
go build -o spottui .

# Run
./spottui

# Test
go test ./...
```

### Environment Variables
| Variable | Required | Description |
|----------|----------|-------------|
| `SPOTIFY_CLIENT` | Yes | Spotify app client ID |
| `SPOTIFY_CALLBACK` | No | OAuth redirect URI (default: `http://localhost:8080/callback`) |
