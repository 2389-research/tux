# Tux Implementation Gaps

Gaps identified when comparing tux to hex for integration readiness.

## Completed

- [x] **Neo-Terminal Theme** - Added as `theme/neoterminal.go`

## Needs Design/Brainstorming

### 1. Huh Form Integration

**Priority:** High
**Status:** Needs brainstorming session

Currently tux has hand-rolled modals (approval, wizard, confirm, list) but hex uses the `huh` library for forms. We should determine:

- Should tux wrap huh forms?
- Or provide theme integration for huh?
- Or leave form building to the app layer?

Hex usage:
- `forms/approval.go` - tool approval with risk assessment
- `forms/onboarding.go` - onboarding flow
- `forms/settings.go` - settings management
- `forms/quickactions.go` - quick actions menu

### 2. Help System

**Priority:** Medium
**Status:** Needs design

Hex has a context-aware help component (`components/help.go`) with:
- `HelpMode` enum: Chat, History, Tools, Approval, Search, QuickActions
- `KeyBindingCategory` with grouped bindings per mode
- Wraps bubbles `help.Model` with theme styling

Questions to resolve:
- Should this be a shell component or standalone?
- How to define keybindings per mode?
- Integration with config keybindings?

### 3. Streaming Display

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
| Huh Form Integration | High | Brainstorm |
| Help System | Medium | Design |
| Streaming Display | Medium | Enhance or skip |
| View Modes | Low | Discuss |
| Markdown | N/A | Out of scope |
