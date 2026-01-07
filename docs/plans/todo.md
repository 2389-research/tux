# Tux Implementation Gaps

Gaps identified when comparing tux to hex for integration readiness.

## Completed

- [x] **Neo-Terminal Theme** - Added as `theme/neoterminal.go`
- [x] **Huh Form Integration** - Added as `form/` package wrapping huh with tux themes
- [x] **Help System** - Added as `help/` package with categories, mode filtering, and HelpModal

## Needs Design/Brainstorming

### 1. Streaming Display

**Priority:** Medium
**Status:** Enhancement needed

Our `Spinner` has token rate, but hex has fuller `StreamingDisplay`:
- Thinking animation with custom spinner frames
- Tool call progress tracking inline
- Typewriter effect (optional)
- Waiting state display

Current gap: Our spinner is for loading indication, not inline streaming display.

Options:
- Extend Spinner with thinking/streaming modes
- Create separate `StreamingDisplay` content type
- Leave to app layer (hex already has this)

### 4. View Modes (vs Tabs)

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
| Streaming Display | Medium | Brainstorm |
| View Modes | Low | Discuss |
| Markdown | N/A | Out of scope |
