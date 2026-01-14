# SpotTUI Technical Debt

## Active Debt

### Test Coverage
- **Priority**: Medium
- **Description**: Current test coverage is basic. Need more comprehensive tests for edge cases.
- **Files affected**: All packages
- **Effort**: Medium

---

## Resolved

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
