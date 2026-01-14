# SpotTUI Technical Debt

## Active Debt

### Test Coverage
- **Priority**: Medium
- **Description**: Current test coverage improved but still below 80% target
- **Current Coverage** (Jan 14, 2026):
  - `tui`: 43.9%
  - `views`: 56.7%
  - `auth`: 23.6%
  - `styles`: 100%
- **Target**: 80%+ coverage
- **Files affected**: All packages
- **Effort**: Medium
- **Notes**: oauth.go coverage limited by external service dependencies

---

## Resolved

### Test Coverage Improvement (Jan 14, 2026)
- **Priority**: Medium
- **Description**: Initial test coverage was very low (tui: 4.8%, views: 9.1%)
- **Resolution**: Added comprehensive tests for model.go, commands.go, handlers.go, and auth/token.go
- **Coverage improvements**:
  - `tui`: 4.8% → 43.9%
  - `views`: 9.1% → 56.7%
  - `auth`: 20.8% → 23.6%
- **Files added**:
  - `internal/tui/model_test.go`
  - `internal/tui/commands_test.go`
  - `internal/tui/handlers_test.go`
  - `internal/tui/retry_test.go`
  - `internal/tui/views/player_test.go`

### Error Handling (Jan 14, 2026)
- **Priority**: Medium
- **Description**: API errors showed generic messages. Needed user-friendly error messages and retry logic.
- **Resolution**: Created `retry.go` with `withRetry[T]` generic function (exponential backoff) and `friendlyError()` for user-friendly messages. Wrapped all data fetching operations.
- **Files**: `internal/tui/retry.go`, `internal/tui/commands.go`

### Hardcoded Dimensions (Jan 13, 2026)
- **Priority**: Low
- **Description**: Some UI calculations used magic numbers for heights/widths.
- **Resolution**: Added layout constants (MinWidth, MinHeight, HorizontalPad, ListHeightSub, MinListHeight) and helper methods (listWidth(), listHeight())
- **Commit**: `ddb907d`
