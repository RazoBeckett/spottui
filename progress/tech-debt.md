# SpotTUI Technical Debt

## Active Debt

### Test Coverage
- **Priority**: Medium
- **Description**: Current test coverage is basic. Need more comprehensive tests for edge cases.
- **Files affected**: All packages
- **Effort**: Medium

### Error Handling
- **Priority**: Medium
- **Description**: API errors show generic messages. Could benefit from more user-friendly error messages and retry logic.
- **Files affected**: `internal/tui/commands.go`
- **Effort**: Low

### Hardcoded Dimensions
- **Priority**: Low
- **Description**: Some UI calculations use magic numbers for heights/widths.
- **Files affected**: `internal/tui/render.go`, `internal/tui/views/`
- **Effort**: Low

---

## Resolved

_None yet_
