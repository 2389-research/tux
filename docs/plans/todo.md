# Tux Implementation Gaps

Gaps identified when comparing tux to hex for integration readiness.

## Completed

- [x] **Neo-Terminal Theme** - Added as `theme/neoterminal.go`
- [x] **Huh Form Integration** - Added as `form/` package wrapping huh with tux themes
- [x] **Help System** - Added as `help/` package with categories, mode filtering, and HelpModal
- [x] **Streaming Display** - Added `StreamingController` (via `shell.Streaming()`) and `StreamingContent` wrapper with typewriter effect

## Needs Design/Brainstorming

### 1. View Modes (vs Tabs)

**Priority:** Low
**Status:** Discussion needed

Hex has shell-level view modes:
- `ViewModeIntro` - Welcome/splash screen
- `ViewModeChat` - Main conversation
- `ViewModeHistory` - History browser
- `ViewModeTools` - Tool inspector

These are NOT tabs - they're full UI state changes. Currently tux has tabs only.

Questions:
- Are these just fullscreen modals?
- Should shell support view modes natively?
- Or let apps implement via tab content switching?

## Out of Scope (App Layer)

### Glamour Markdown Rendering

Hex uses glamour for rendering markdown in chat. This is content rendering that belongs at the app layer, not in tux. Apps can:
- Use glamour directly for their content
- Create custom content types that wrap glamour

## Summary

| Gap | Priority | Action |
|-----|----------|--------|
| Neo-Terminal Theme | High | ✅ Done |
| Huh Form Integration | High | ✅ Done |
| Help System | Medium | ✅ Done |
| Streaming Display | Medium | ✅ Done |
| View Modes | Low | Discuss |
| Markdown | N/A | Out of scope |
