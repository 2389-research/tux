# Tux Implementation Gaps

Gaps identified when comparing tux to hex for integration readiness.

## Completed

- [x] **Neo-Terminal Theme** - Added as `theme/neoterminal.go`
- [x] **Huh Form Integration** - Added as `form/` package wrapping huh with tux themes
- [x] **Help System** - Added as `help/` package with categories, mode filtering, and HelpModal
- [x] **Streaming Display** - Added `StreamingController` (via `shell.Streaming()`) and `StreamingContent` wrapper with typewriter effect

## Designed (Needs Implementation)

### 1. Tabs and Modals

**Priority:** Medium
**Status:** Design complete, needs implementation
**Design:** `docs/plans/2026-01-07-tabs-and-modals-design.md`

Opinionated default for view management:
- **Tabs** for persistent content views (Chat, History, Tools)
- **Modals** for temporary overlays (Help, Approval, Forms)

Implementation needed:
- Tab keyboard shortcuts (Ctrl+Tab, Ctrl+1/2/3)
- Hidden tabs with shortcuts
- Tab content lifecycle hooks (OnActivate/OnDeactivate)

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
| Tabs and Modals | Medium | Designed, needs impl |
| Markdown | N/A | Out of scope |
