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

---

## Resolved

### Hardcoded Dimensions (Jan 13, 2026)
- **Priority**: Low
- **Description**: Some UI calculations used magic numbers for heights/widths.
- **Resolution**: Added layout constants (MinWidth, MinHeight, HorizontalPad, ListHeightSub, MinListHeight) and helper methods (listWidth(), listHeight())
- **Commit**: `ddb907d`
